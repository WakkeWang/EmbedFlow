package protocol

import "testing"

func TestParseControl_AuthOK(t *testing.T) {
	wire, err := Encode(&AuthFrame{Token: "t1", ProtocolVersion: 1, ClientVersion: "0.1.0"})
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	f, err := ParseControl(wire)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if f.Type != FrameAuth {
		t.Fatalf("type = %d, want %d", f.Type, FrameAuth)
	}
	a, ok := f.Body.(*AuthFrame)
	if !ok {
		t.Fatalf("body type = %T, want *AuthFrame", f.Body)
	}
	if a.Token != "t1" || a.ProtocolVersion != 1 || a.ClientVersion != "0.1.0" {
		t.Fatalf("fields = %+v", a)
	}
}

func TestParseControl_UnknownTypeIgnorable(t *testing.T) {
	// Unknown frame type parses into a generic frame: caller ignores it and
	// logs a warning (CEO-3A), connection stays up.
	wire := []byte(`{"type":9999,"body":{"x":1}}`)
	f, err := ParseControl(wire)
	if err != nil {
		t.Fatalf("parse unknown type: %v", err)
	}
	if f.Type != 9999 {
		t.Fatalf("type = %d, want 9999", f.Type)
	}
}

func TestParseControl_MalformedJSONRejected(t *testing.T) {
	if _, err := ParseControl([]byte(`{not json`)); err == nil {
		t.Fatal("malformed JSON should error")
	}
}

func TestParseControl_MissingTypeRejected(t *testing.T) {
	if _, err := ParseControl([]byte(`{"body":{}}`)); err == nil {
		t.Fatal("missing type should error")
	}
}

func TestParseControl_HeartbeatRoundTrip(t *testing.T) {
	wire, err := Encode(&HeartbeatFrame{Seq: 7})
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	f, err := ParseControl(wire)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	h, ok := f.Body.(*HeartbeatFrame)
	if !ok {
		t.Fatalf("body type = %T, want *HeartbeatFrame", f.Body)
	}
	if h.Seq != 7 {
		t.Fatalf("seq = %d, want 7", h.Seq)
	}
}
