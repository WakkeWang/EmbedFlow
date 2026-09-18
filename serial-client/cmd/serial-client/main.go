// Command serial-client is the EmbedFlow serial tunnel client: login, share
// a real COM port to the server, pure passthrough. It doubles as the
// developer harness and the M1 latency probe (issue #2).
package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log/slog"
	"math"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/WakkeWang/EmbedFlow/serial-client"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "probe" {
		probeMain(os.Args[2:])
		return
	}

	var (
		server   = flag.String("server", "http://127.0.0.1:8420", "server URL")
		username = flag.String("user", "admin", "username")
		password = flag.String("pass", "", "password (required)")
		device   = flag.Int64("device", 0, "device id to share as (required)")
		port     = flag.String("port", "COM3", "local COM port name")
		baud     = flag.Int("baud", 115200, "baud rate")
		list     = flag.Bool("list-ports", false, "scan and print COM ports, then exit")
	)
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	if *list {
		ports, err := client.ListPorts()
		if err != nil {
			slog.Error("port scan failed", "err", err)
			os.Exit(1)
		}
		for _, p := range ports {
			fmt.Println(p)
		}
		if len(ports) == 0 {
			fmt.Println("(no serial ports found)")
		}
		return
	}

	if *password == "" || *device == 0 {
		fmt.Fprintln(os.Stderr, "usage: serial-client -server URL -user U -pass P -device N [-port COM3] [-baud 115200]")
		fmt.Fprintln(os.Stderr, "       serial-client probe -server URL [-n 200] [-mode ws|full] [-port COM3]")
		os.Exit(2)
	}

	ctx := signalContext()

	tok, err := client.Login(*server, *username, *password)
	if err != nil {
		slog.Error("login failed", "err", err)
		os.Exit(1)
	}
	slog.Info("logged in")

	cl := client.New(*server, tok)
	if err := cl.Connect(ctx); err != nil {
		slog.Error("connect failed", "err", err)
		os.Exit(1)
	}
	defer cl.Close()

	// Real COM port passthrough (requirement 3.4): open the port, wire it to
	// the connection, share it under the registered device.
	p, err := client.OpenPort(*port, *baud)
	if err != nil {
		slog.Error("open serial port failed", "port", *port, "err", err)
		os.Exit(1)
	}
	cl.SetPort(p)
	if err := cl.Share(*device, *port); err != nil {
		slog.Error("share failed", "err", err)
		os.Exit(1)
	}
	slog.Info("sharing", "device", *device, "port", *port, "baud", *baud)

	<-ctx.Done()
	slog.Info("bye")
}

// probeMain implements `serial-client probe`: the M1 latency spike runner
// (issue #2, CEO-12A). Two modes:
//
//   - ws (default): round trips against the server's -probe-echo endpoint
//     over pure WebSocket -- isolates the WS leg.
//   - full: 8-byte payloads go out the real COM port (loopback plug or live
//     console), through the tunnel, back -- the full-path number that drives
//     the server-side-matcher go/no-go.
func probeMain(args []string) {
	fs := flag.NewFlagSet("probe", flag.ExitOnError)
	var (
		server  = fs.String("server", "http://127.0.0.1:8420", "server URL (probe-echo mode)")
		n       = fs.Int("n", 200, "number of samples")
		mode    = fs.String("mode", "ws", "ws | full")
		port    = fs.String("port", "", "COM port (full mode only)")
		baud    = fs.Int("baud", 115200, "baud rate (full mode only)")
		device  = fs.Int64("device", 0, "device id to share as (full mode only)")
		timeout = fs.Duration("timeout", 2*time.Second, "per-sample timeout")
	)
	_ = fs.Parse(args)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	ctx := signalContext()

	switch *mode {
	case "ws":
		runWSProbe(ctx, *server, *n, *timeout)
	case "full":
		if *port == "" || *device == 0 {
			fmt.Fprintln(os.Stderr, "full mode needs -port COMx and -device N")
			os.Exit(2)
		}
		runFullProbe(ctx, *server, *n, *timeout, *port, *baud, *device)
	default:
		fmt.Fprintln(os.Stderr, "mode must be ws or full")
		os.Exit(2)
	}
}

func runWSProbe(ctx context.Context, server string, n int, timeout time.Duration) {
	pc, err := client.ProbeDial(ctx, server)
	if err != nil {
		slog.Error("probe dial", "err", err)
		os.Exit(1)
	}
	defer pc.Close()
	if err := pc.Handshake(ctx); err != nil {
		slog.Error("probe handshake", "err", err)
		os.Exit(1)
	}

	samples := make([]time.Duration, 0, n)
	timedOut := 0
	for i := 0; i < n; i++ {
		if ctx.Err() != nil {
			break
		}
		payload := make([]byte, 8)
		binary.BigEndian.PutUint64(payload, uint64(i+1))
		rtt, err := pc.RTT(ctx, payload, timeout)
		if err != nil {
			timedOut++
			continue
		}
		samples = append(samples, rtt)
		time.Sleep(10 * time.Millisecond)
	}
	printSamples("ws-echo", samples, timedOut)
}

func runFullProbe(ctx context.Context, server string, n int, timeout time.Duration, portName string, baud int, device int64) {
	// Full mode runs against the production server: login + share, then push
	// probe bytes through the real COM port (loopback plug). The port is
	// bound WITHOUT the passthrough pump: the probe reads the port itself,
	// and a second concurrent reader would steal echo bytes (silently
	// skewed numbers).
	tok, err := client.Login(server, os.Getenv("EF_USER"), os.Getenv("EF_PASS"))
	if err != nil {
		slog.Error("login failed (set EF_USER/EF_PASS for full mode)", "err", err)
		os.Exit(1)
	}
	cl := client.New(server, tok)
	if err := cl.Connect(ctx); err != nil {
		slog.Error("connect", "err", err)
		os.Exit(1)
	}
	defer cl.Close()

	p, err := client.OpenPort(portName, baud)
	if err != nil {
		slog.Error("open port", "port", portName, "err", err)
		os.Exit(1)
	}
	defer p.Close()
	cl.SetLocalPortOnly(p)

	if err := cl.Share(device, portName); err != nil {
		slog.Error("share", "err", err)
		os.Exit(1)
	}
	time.Sleep(500 * time.Millisecond) // let the share settle

	res, err := client.MeasureLatency(ctx, cl, p, n, timeout)
	if err != nil {
		slog.Error("measure", "err", err)
		os.Exit(1)
	}
	printResult("full-path", res)
}

func printSamples(label string, samples []time.Duration, timedOut int) {
	if len(samples) == 0 {
		fmt.Printf("%s: no samples (all timed out)\n", label)
		return
	}
	sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
	w := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	fmt.Fprintf(w, "samples\t%d\n", len(samples))
	fmt.Fprintf(w, "timed out\t%d\n", timedOut)
	fmt.Fprintf(w, "p50\t%s\n", percentile(samples, 0.50))
	fmt.Fprintf(w, "p90\t%s\n", percentile(samples, 0.90))
	fmt.Fprintf(w, "p99\t%s\n", percentile(samples, 0.99))
	fmt.Fprintf(w, "max\t%s\n", samples[len(samples)-1])
	_ = w.Flush()
}

// percentile is the nearest-rank percentile (ceil indexing) -- the same
// formula the latency core uses, so CLI summaries and library callers
// never disagree.
func percentile(sorted []time.Duration, q float64) time.Duration {
	idx := int(math.Ceil(q*float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

func printResult(label string, res *client.LatencyResult) {
	w := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	fmt.Fprintf(w, "samples\t%d\n", res.N)
	fmt.Fprintf(w, "timed out\t%d\n", res.TimedOut)
	fmt.Fprintf(w, "p50\t%s\n", res.P50)
	fmt.Fprintf(w, "p90\t%s\n", res.P90)
	fmt.Fprintf(w, "p99\t%s\n", res.P99)
	fmt.Fprintf(w, "max\t%s\n", res.Max)
	fmt.Fprintf(w, "mean\t%s\n", res.Mean)
	_ = w.Flush()
}
