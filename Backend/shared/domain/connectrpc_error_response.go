package domain

import "encoding/json"

type ConnectError struct {
	Code    string        `json:"code"`              // e.g., "invalid_argument"
	Message string        `json:"message"`           // Human-readable description
	Details []ErrorDetail `json:"details,omitempty"` // Optional structured data
}

type ErrorDetail struct {
	Type  string          `json:"type"`            // Protobuf message type URL
	Value string          `json:"value"`           // Base64-encoded binary data
	Debug json.RawMessage `json:"debug,omitempty"` // Human-readable JSON representation
}
