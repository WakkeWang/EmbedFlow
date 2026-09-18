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
	FrameAuth         uint16 = 1 // client -> server: authenticate
	FrameAuthOK       uint16 = 2 // server -> client: auth accepted
	FrameAuthFail     uint16 = 3 // server -> client: auth refused (reason in body)
	FrameHeartbeat    uint16 = 4 // both ways: liveness ping
	FrameShareRequest uint16 = 5 // client -> server: share a serial port
	FrameDeviceStatus uint16 = 6 // client -> server: port open/close/error state
	FrameSessionCtrl  uint16 = 7 // web -> server: session lifecycle commands
	FrameSessionState uint16 = 8 // server -> clients: session state broadcast
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
	DeviceID int64 `json:"device_id"`
}

// DeviceStatusFrame reports the local serial port state on the client side.
type DeviceStatusFrame struct {
	Port  string `json:"port"`
	Open  bool   `json:"open"`
	Error string `json:"error,omitempty"`
}

// SessionCtrlFrame carries browser-side session commands (open/close/abort...).
type SessionCtrlFrame struct {
	Command   string `json:"command"`
	SessionID int64  `json:"session_id,omitempty"`
	DeviceID  int64  `json:"device_id,omitempty"`
}

// SessionStateFrame broadcasts session lifecycle state to subscribed browsers.
type SessionStateFrame struct {
	SessionID int64  `json:"session_id"`
	DeviceID  int64  `json:"device_id"`
	State     string `json:"state"`
	Detail    string `json:"detail,omitempty"`
}

// bodyTypes maps frame type numbers to concrete body structs.
var bodyTypes = map[uint16]func() any{
	FrameAuth:         func() any { return new(AuthFrame) },
	FrameAuthOK:       func() any { return new(AuthOKFrame) },
	FrameAuthFail:     func() any { return new(AuthFailFrame) },
	FrameHeartbeat:    func() any { return new(HeartbeatFrame) },
	FrameShareRequest: func() any { return new(ShareRequestFrame) },
	FrameDeviceStatus: func() any { return new(DeviceStatusFrame) },
	FrameSessionCtrl:  func() any { return new(SessionCtrlFrame) },
	FrameSessionState: func() any { return new(SessionStateFrame) },
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
	makeBody, known := bodyTypes[env.Type]
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
	var typ uint16
	switch body.(type) {
	case *AuthFrame:
		typ = FrameAuth
	case *AuthOKFrame:
		typ = FrameAuthOK
	case *AuthFailFrame:
		typ = FrameAuthFail
	case *HeartbeatFrame:
		typ = FrameHeartbeat
	case *ShareRequestFrame:
		typ = FrameShareRequest
	case *DeviceStatusFrame:
		typ = FrameDeviceStatus
	case *SessionCtrlFrame:
		typ = FrameSessionCtrl
	case *SessionStateFrame:
		typ = FrameSessionState
	default:
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
