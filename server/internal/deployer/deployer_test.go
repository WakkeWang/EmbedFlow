package deployer

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

// --- Render (requirement 3.1.2: the three variable families) ---

func TestRender_ThreeFamilies(t *testing.T) {
	v := Variables{
		Artifact: map[string]string{"name": "bootfw.bin", "checksum": "abc123"},
		Device:   map[string]string{"ip": "10.0.0.8"},
		Params:   map[string]string{"dir": "/tmp/x"},
	}
	got, missing := Render("cp ${artifact.name} ${params.dir}/ && ssh ${device.ip}", v)
	if got != "cp bootfw.bin /tmp/x/ && ssh 10.0.0.8" {
		t.Fatalf("render = %q", got)
	}
	if len(missing) != 0 {
		t.Fatalf("missing = %v", missing)
	}
}

func TestRender_UnknownVariableReported(t *testing.T) {
	got, missing := Render("deploy ${artifact.nosuch}", Variables{Artifact: map[string]string{"name": "a"}})
	if got != "deploy " {
		t.Fatalf("render = %q", got)
	}
	if len(missing) != 1 || missing[0] != "artifact.nosuch" {
		t.Fatalf("missing = %v", missing)
	}
}

// --- SSH integration: a loopback SSH server (decision 7A analog: no
// external device needed to test the executor path) ---

type sshServer struct {
	listener net.Listener
	config   *ssh.ServerConfig
	commands *[]string // observed commands (pointer for mutation via closure)
}

// newSSHServer starts a loopback SSH server accepting one password
// ("secret") and executing commands through a fake shell.
func newSSHServer(t *testing.T, failOn string) (addr string, commands *[]string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("host key: %v", err)
	}
	raw, _ := x509.MarshalECPrivateKey(key)
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: raw})
	signer, err := ssh.ParsePrivateKey(pemKey)
	if err != nil {
		t.Fatalf("signer: %v", err)
	}
	config := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			if conn.User() == "root" && string(password) == "secret" {
				return nil, nil
			}
			return nil, fmt.Errorf("auth rejected for %s", conn.User())
		},
	}
	config.AddHostKey(signer)

	commands = new([]string)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	s := &sshServer{listener: listener, config: config, commands: commands}
	go s.serve(t)
	t.Cleanup(func() { listener.Close() })
	return listener.Addr().String(), commands
}

func (s *sshServer) serve(t *testing.T) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handle(t, conn)
	}
}

func (s *sshServer) handle(t *testing.T, conn net.Conn) {
	serverConn, chans, reqs, err := ssh.NewServerConn(conn, s.config)
	if err != nil {
		return
	}
	defer serverConn.Close()
	go ssh.DiscardRequests(reqs)
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unsupported")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			return
		}
		go func() {
			for req := range requests {
				switch req.Type {
				case "exec":
					var payload struct{ Command string }
					ssh.Unmarshal(req.Payload, &payload)
					*s.commands = append(*s.commands, payload.Command)
					// Fake shell: echo back the command; a FAIL marker exits 3.
					// The channel closes after exit-status -- Session.Run waits
					// for it.
					if strings.Contains(payload.Command, "FAIL") {
						channel.Write([]byte("boom\n"))
						channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Code uint32 }{3}))
						req.Reply(true, nil)
					} else {
						channel.Write([]byte("ran: " + payload.Command + "\n"))
						channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Code uint32 }{0}))
						req.Reply(true, nil)
					}
					channel.Close()
					return
				default:
					if req.WantReply {
						req.Reply(true, nil)
					}
				}
			}
		}()
	}
}

func TestSSHRunner_CommandsExecute(t *testing.T) {
	addr, commands := newSSHServer(t, "")
	host, portStr, _ := strings.Cut(addr, ":")
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	var lines []string
	runner := &SSHRunner{
		Host: host, Port: port, User: "root", Password: "secret",
		Notify: func(_ int64, phase, line string) { lines = append(lines, phase+"|"+line) },
		Record: 1,
	}
	res := runner.Run(t.Context(), []string{"cmd-one", "cmd-two"})
	if !res.OK {
		t.Fatalf("run failed: %s", res.Detail)
	}
	if len(*commands) != 2 || (*commands)[0] != "cmd-one" {
		t.Fatalf("commands = %v", *commands)
	}
	if !strings.Contains(res.Detail, "2 command") {
		t.Fatalf("detail = %q", res.Detail)
	}
}

func TestSSHRunner_FailureStopsBatch(t *testing.T) {
	addr, _ := newSSHServer(t, "")
	host, portStr, _ := strings.Cut(addr, ":")
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	runner := &SSHRunner{Host: host, Port: port, User: "root", Password: "secret", Record: 1}
	res := runner.Run(t.Context(), []string{"ok-cmd", "FAIL-cmd", "never-runs"})
	if res.OK {
		t.Fatal("failing command must fail the batch")
	}
	if !strings.Contains(res.Detail, "failed") {
		t.Fatalf("detail = %q", res.Detail)
	}
}

func TestSSHRunner_BadCredentials(t *testing.T) {
	addr, _ := newSSHServer(t, "")
	host, portStr, _ := strings.Cut(addr, ":")
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	runner := &SSHRunner{Host: host, Port: port, User: "root", Password: "wrong", Record: 1}
	res := runner.Run(t.Context(), []string{"x"})
	if res.OK {
		t.Fatal("bad credentials must fail")
	}
	if !strings.Contains(res.Detail, "dial") {
		t.Fatalf("detail = %q", res.Detail)
	}
}

func TestWriteExecLog(t *testing.T) {
	dir := t.TempDir()
	path := t.TempDir() + "/log"
	_ = path
	// WriteExecLog creates <dataDir>/execs/<id>/log.txt under dataDir.
	if err := WriteExecLog(dir, 42, "ssh", "hello"); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := WriteExecLog(dir, 42, "ssh", "again"); err != nil {
		t.Fatalf("append: %v", err)
	}
	// Two lines accumulated.
	got := readTail(t, dir+"/execs/42/log.txt")
	if !strings.Contains(got, "hello") || !strings.Contains(got, "again") {
		t.Fatalf("log = %q", got)
	}
	_ = time.Now
}
