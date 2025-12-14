package json

import (
	"encoding/json"
)

type (
	Encoder = json.Encoder
	Decoder = json.Decoder

	RawMessage = json.RawMessage
	Number     = json.Number
	Token      = json.Token
	Delim      = json.Delim
)

var (
	NewDecoder = json.NewDecoder
	Unmarshal  = json.Unmarshal

	NewEncoder    = json.NewEncoder
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent

	Compact    = json.Compact
	HTMLEscape = json.HTMLEscape
	Valid      = json.Valid
)
