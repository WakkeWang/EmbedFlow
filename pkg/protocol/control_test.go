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

func TestParseControl_ExpectProgressRoundTrip(t *testing.T) {
	wire, err := Encode(&ExpectProgressFrame{
		SessionID: 3, StepIndex: 1, StepTotal: 5, StepDesc: "await Login:",
		Phase: "running",
	})
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	f, err := ParseControl(wire)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	p, ok := f.Body.(*ExpectProgressFrame)
	if !ok {
		t.Fatalf("body = %T, want *ExpectProgressFrame", f.Body)
	}
	if p.SessionID != 3 || p.StepIndex != 1 || p.StepTotal != 5 || p.Phase != "running" {
		t.Fatalf("fields = %+v", p)
	}
}

func TestParseControl_ConfirmRoundTrip(t *testing.T) {
	wire, err := Encode(&ConfirmFrame{
		SessionID: 3, ConfirmID: 8, Prompt: "LED on?", State: "pending",
	})
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	f, err := ParseControl(wire)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	c, ok := f.Body.(*ConfirmFrame)
	if !ok {
		t.Fatalf("body = %T, want *ConfirmFrame", f.Body)
	}
	if c.ConfirmID != 8 || c.Prompt != "LED on?" || c.State != "pending" {
		t.Fatalf("fields = %+v", c)
	}
}

func TestParseControl_SessionStateOccupier(t *testing.T) {
	// DS-2A: busy rejection carries occupier identity.
	wire, err := Encode(&SessionStateFrame{
		DeviceID: 2, State: "busy", Detail: "device busy",
		OccupierUser: "bob", OccupierKind: "task", OccupierSince: "2026-09-18T10:00:00Z",
	})
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	f, err := ParseControl(wire)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	s, ok := f.Body.(*SessionStateFrame)
	if !ok {
		t.Fatalf("body = %T", f.Body)
	}
	if s.OccupierUser != "bob" || s.OccupierKind != "task" && s.OccupierKind != "" {
		t.Fatalf("fields = %+v", s)
	}
}
