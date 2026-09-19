// Package protocol defines the wire protocol shared by the EmbedFlow server
// and the Windows serial client (decision 5A).
//
// WebSocket dual-frame scheme (decision 4A): control frames are JSON text
// frames, data frames are binary frames carrying raw serial bytes.
package protocol

import (
	"encoding/json"
	"fmt"
)

// ProtocolVersion is the current wire-protocol version. A client whose version
// differs is refused at the handshake (CEO-11A).
const ProtocolVersion = 1

// Frame type numbers for control frames.
const (
	FrameAuth           uint16 = 1  // client -> server: authenticate
	FrameAuthOK         uint16 = 2  // server -> client: auth accepted
	FrameAuthFail       uint16 = 3  // server -> client: auth refused (reason in body)
	FrameHeartbeat      uint16 = 4  // both ways: liveness ping
	FrameShareRequest   uint16 = 5  // client -> server: share a serial port
	FrameDeviceStatus   uint16 = 6  // client -> server: port open/close/error state
	FrameSessionCtrl    uint16 = 7  // web -> server: session lifecycle commands
	FrameSessionState   uint16 = 8  // server -> clients: session state broadcast
	FrameExpectProgress uint16 = 9  // server -> web: expect step progress (DS-3A)
	FrameConfirm        uint16 = 10 // both ways: human confirmation card (issue #9)
	FrameBuildCtrl      uint16 = 11 // web -> server: build batch subscribe/unsubscribe
	FrameBuildEvent     uint16 = 12 // server -> web: batch/record progress events (M2)
)

// Envelope is the JSON shape of every control frame on the wire.
type Envelope struct {
	Type uint16          `json:"type"`
	Body json.RawMessage `json:"body,omitempty"`
}

// Frame is a parsed control frame: the type plus a decoded body.
// Unknown types parse with Body left as *UnknownBody so callers can
// ignore them with a warning instead of dropping the connection (CEO-3A).
type Frame struct {
	Type uint16
	Body any
}

// UnknownBody preserves the raw JSON of an unrecognized frame type.
type UnknownBody struct {
	Raw json.RawMessage `json:"-"`
}

// AuthFrame is the first frame a client or browser sends after connecting.
type AuthFrame struct {
	Token           string `json:"token"`
	ProtocolVersion uint16 `json:"protocol_version"`
	ClientVersion   string `json:"client_version,omitempty"`
}

// AuthOKFrame confirms authentication. ServerInfo echoes the server version
// so the client can warn on mismatch (issue #10).
type AuthOKFrame struct {
	ServerVersion string `json:"server_version,omitempty"`
	Role          string `json:"role,omitempty"` // admin | member (requirement 1.5)
}

// AuthFailFrame rejects authentication or a version mismatch.
type AuthFailFrame struct {
	Reason string `json:"reason"`
}

// HeartbeatFrame is a liveness ping; Seq echoes back for RTT inspection.
type HeartbeatFrame struct {
	Seq uint64 `json:"seq"`
}

// ShareRequestFrame asks the server to bind this client connection to a
// registered device as its serial transport.
type ShareRequestFrame struct {
	DeviceID int64  `json:"device_id"`
	Port     string `json:"port,omitempty"` // COM port name, for display
}

// DeviceStatusFrame reports the local serial port state on the client side.
type DeviceStatusFrame struct {
	Port  string `json:"port"`
	Open  bool   `json:"open"`
	Error string `json:"error,omitempty"`
}

// SessionCtrlFrame carries browser-side session commands.
// Commands: open (device_id), close (session_id), follow (device_id,
// read-only second viewer), abort (session_id, expect run), run
// (device_id, rule_id -- creates a task session, CEO-17A).
type SessionCtrlFrame struct {
	Command   string `json:"command"`
	SessionID int64  `json:"session_id,omitempty"`
	DeviceID  int64  `json:"device_id,omitempty"`
	RuleID    int64  `json:"rule_id,omitempty"`
}

// SessionStateFrame broadcasts session lifecycle state to subscribed browsers.
// Occupier fields are set on busy-rejection for the DS-2A popup.
type SessionStateFrame struct {
	SessionID int64  `json:"session_id"`
	DeviceID  int64  `json:"device_id"`
	State     string `json:"state"` // active | closed | failed | rejected | busy | idle_warning | readonly
	Detail    string `json:"detail,omitempty"`
	// Occupier identifies the current holder on busy-rejection (DS-2A).
	OccupierUser  string `json:"occupier_user,omitempty"`
	OccupierKind  string `json:"occupier_kind,omitempty"`
	OccupierSince string `json:"occupier_since,omitempty"`
	OwnerIsYou    bool   `json:"owner_is_you,omitempty"`
}

// ExpectProgressFrame reports expect run progress (issue #8, DS-3A):
// current step index/total, the step description, and the phase.
type ExpectProgressFrame struct {
	SessionID int64  `json:"session_id"`
	StepIndex int    `json:"step_index"`
	StepTotal int    `json:"step_total"`
	StepDesc  string `json:"step_desc,omitempty"`
	Phase     string `json:"phase"` // running | completed | failed | aborted
	Detail    string `json:"detail,omitempty"`
}

// ConfirmFrame is a human confirmation card (issue #9): LED observations and
// cable unplugging recorded as first-class session steps. Cards never
// auto-dismiss; resolution is explicit PASS/FAIL with an optional note.
type ConfirmFrame struct {
	SessionID int64  `json:"session_id"`
	ConfirmID int64  `json:"confirm_id"`
	Prompt    string `json:"prompt"`
	State     string `json:"state"`            // pending | resolved
	Result    string `json:"result,omitempty"` // pass | fail
	Note      string `json:"note,omitempty"`
}

// BuildCtrlFrame carries browser-side build subscriptions: the batch view
// subscribes on entry and unsubscribes on leave, so live record/log events
// reach only the pages that want them.
type BuildCtrlFrame struct {
	Command string `json:"command"` // subscribe | unsubscribe
	BatchID int64  `json:"batch_id"`
}

// BuildEventFrame pushes build progress to subscribed browsers (M2):
// record state changes and streamed log lines. Phase vocabulary:
// started | log | succeeded | failed | skipped | canceled | batch_done.
type BuildEventFrame struct {
	BatchID  int64  `json:"batch_id"`
	RecordID int64  `json:"record_id,omitempty"`
	ItemName string `json:"item_name,omitempty"`
	Phase    string `json:"phase"`
	Detail   string `json:"detail,omitempty"`
	LogLine  string `json:"log_line,omitempty"`
}

// frameDef binds a frame type to its concrete body constructor.
type frameDef struct {
	typ uint16
	new func() any
}

// registry is the single source of truth for frame <-> type mapping; the
// decode map and the encode reverse map are both derived from it.
var registry = []frameDef{
	{FrameAuth, func() any { return new(AuthFrame) }},
	{FrameAuthOK, func() any { return new(AuthOKFrame) }},
	{FrameAuthFail, func() any { return new(AuthFailFrame) }},
	{FrameHeartbeat, func() any { return new(HeartbeatFrame) }},
	{FrameShareRequest, func() any { return new(ShareRequestFrame) }},
	{FrameDeviceStatus, func() any { return new(DeviceStatusFrame) }},
	{FrameSessionCtrl, func() any { return new(SessionCtrlFrame) }},
	{FrameSessionState, func() any { return new(SessionStateFrame) }},
	{FrameExpectProgress, func() any { return new(ExpectProgressFrame) }},
	{FrameConfirm, func() any { return new(ConfirmFrame) }},
	{FrameBuildCtrl, func() any { return new(BuildCtrlFrame) }},
	{FrameBuildEvent, func() any { return new(BuildEventFrame) }},
}

var decodeByType = func() map[uint16]func() any {
	m := make(map[uint16]func() any, len(registry))
	for _, d := range registry {
		m[d.typ] = d.new
	}
	return m
}()

var encodeByBody = func() map[string]uint16 {
	m := make(map[string]uint16, len(registry))
	for _, d := range registry {
		m[frameBodyName(d.new())] = d.typ
	}
	return m
}()

func frameBodyName(v any) string {
	// Pointer type name, e.g. "*protocol.AuthFrame".
	return fmt.Sprintf("%T", v)
}

// ParseControl decodes a JSON text frame into a Frame.
// Known types decode into their concrete body struct; unknown types keep the
// raw body so the caller can ignore-and-log (CEO-3A).
func ParseControl(wire []byte) (*Frame, error) {
	var env Envelope
	if err := json.Unmarshal(wire, &env); err != nil {
		return nil, fmt.Errorf("protocol: bad control frame: %w", err)
	}
	if env.Type == 0 {
		return nil, fmt.Errorf("protocol: control frame missing type")
	}
	makeBody, known := decodeByType[env.Type]
	if !known {
		return &Frame{Type: env.Type, Body: &UnknownBody{Raw: env.Body}}, nil
	}
	body := makeBody()
	if len(env.Body) > 0 {
		if err := json.Unmarshal(env.Body, body); err != nil {
			return nil, fmt.Errorf("protocol: bad body for frame type %d: %w", env.Type, err)
		}
	}
	return &Frame{Type: env.Type, Body: body}, nil
}

// Encode marshals a body struct into wire bytes under its frame type.
func Encode(body any) ([]byte, error) {
	typ, ok := encodeByBody[frameBodyName(body)]
	if !ok {
		return nil, fmt.Errorf("protocol: unencodable body %T", body)
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("protocol: marshal body: %w", err)
	}
	wire, err := json.Marshal(Envelope{Type: typ, Body: raw})
	if err != nil {
		return nil, fmt.Errorf("protocol: marshal envelope: %w", err)
	}
	return wire, nil
}
