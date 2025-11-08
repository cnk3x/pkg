package yaml

import "github.com/goccy/go-yaml"

type (
	MapItem         = yaml.MapItem
	MapSlice        = yaml.MapSlice
	Decoder         = yaml.Decoder
	Encoder         = yaml.Encoder
	EncodeOption    = yaml.EncodeOption
	DecodeOption    = yaml.DecodeOption
	CommentMap      = yaml.CommentMap
	Comment         = yaml.Comment
	CommentPosition = yaml.CommentPosition
	Path            = yaml.Path
	PathBuilder     = yaml.PathBuilder
	StructField     = yaml.StructField
	StructFieldMap  = yaml.StructFieldMap
)

var (
	Unmarshal   = yaml.Unmarshal
	Marshal     = yaml.Marshal
	NewDecoder  = yaml.NewDecoder
	NewEncoder  = yaml.NewEncoder
	ToJSON      = yaml.YAMLToJSON
	FromJSON    = yaml.JSONToYAML
	HeadComment = yaml.HeadComment
	LineComment = yaml.LineComment
	PathString  = yaml.PathString
)

//EncodeOption
var (
	JSON                       = yaml.JSON
	UseJSONMarshaler           = yaml.UseJSONMarshaler
	Indent                     = yaml.Indent
	IndentSequence             = yaml.IndentSequence
	Flow                       = yaml.Flow
	UseLiteralStyleIfMultiline = yaml.UseLiteralStyleIfMultiline
	MarshalAnchor              = yaml.MarshalAnchor
	WithComment                = yaml.WithComment
)

//DecodeOption
var (
	ReferenceReaders     = yaml.ReferenceReaders
	ReferenceFiles       = yaml.ReferenceFiles
	ReferenceDirs        = yaml.ReferenceDirs
	RecursiveDir         = yaml.RecursiveDir
	Validator            = yaml.Validator
	Strict               = yaml.Strict
	DisallowUnknownField = yaml.DisallowUnknownField
	DisallowDuplicateKey = yaml.DisallowDuplicateKey
	UseOrderedMap        = yaml.UseOrderedMap
	UseJSONUnmarshaler   = yaml.UseJSONUnmarshaler
	CommentToMap         = yaml.CommentToMap
)
