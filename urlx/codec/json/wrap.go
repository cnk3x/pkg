package json

import "github.com/goccy/go-json"

type (
	Encoder          = json.Encoder
	EncodeOption     = json.EncodeOption
	EncodeOptionFunc = json.EncodeOptionFunc

	Decoder          = json.Decoder
	DecodeOption     = json.DecodeOption
	DecodeOptionFunc = json.DecodeOptionFunc

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
