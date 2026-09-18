package serialport

import (
	"io"
	"testing"
	"time"
)

func TestPipe_RoundTrip(t *testing.T) {
	client, device := NewPipe()

	go func() {
		n, err := client.Write([]byte("hello"))
		if err != nil || n != 5 {
			t.Errorf("client write: n=%d err=%v", n, err)
		}
	}()

	buf := make([]byte, 16)
	device.SetReadDeadline(time.Now().Add(time.Second))
	n, err := device.Read(buf)
	if err != nil {
		t.Fatalf("device read: %v", err)
	}
	if string(buf[:n]) != "hello" {
		t.Fatalf("device got %q, want %q", buf[:n], "hello")
	}
}

func TestPipe_BothDirections(t *testing.T) {
	client, device := NewPipe()

	go func() {
		device.Write([]byte{0x01, 0xFF}) // device -> client (binary safe)
	}()

	buf := make([]byte, 8)
	client.SetReadDeadline(time.Now().Add(time.Second))
	n, err := client.Read(buf)
	if err != nil {
		t.Fatalf("client read: %v", err)
	}
	if buf[0] != 0x01 || buf[1] != 0xFF {
		t.Fatalf("client got % x", buf[:n])
	}
}

func TestPipe_CloseUnblocksRead(t *testing.T) {
	client, device := NewPipe()

	done := make(chan error, 1)
	go func() {
		buf := make([]byte, 8)
		_, err := client.Read(buf)
		done <- err
	}()

	time.Sleep(50 * time.Millisecond)
	device.Close()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("read after close should error")
		}
	case <-time.After(time.Second):
		t.Fatal("read not unblocked by close")
	}
}

func TestPipe_ImplementsPort(t *testing.T) {
	var _ Port = (*pipeEnd)(nil)
	var _ io.ReadWriteCloser = (*pipeEnd)(nil)
}
