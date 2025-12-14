package yaml

import (
	"go.yaml.in/yaml/v3"
)

type (
	Decoder = yaml.Decoder
	Encoder = yaml.Encoder
)

var (
	Unmarshal  = yaml.Unmarshal
	Marshal    = yaml.Marshal
	NewDecoder = yaml.NewDecoder
	NewEncoder = yaml.NewEncoder
)

// type (
// 	MapItem         = yaml.MapItem
// 	MapSlice        = yaml.MapSlice
// 	EncodeOption    = yaml.EncodeOption
// 	DecodeOption    = yaml.DecodeOption
// 	CommentMap      = yaml.CommentMap
// 	Comment         = yaml.Comment
// 	CommentPosition = yaml.CommentPosition
// 	Path            = yaml.Path
// 	PathBuilder     = yaml.PathBuilder
// 	StructField     = yaml.StructField
// 	StructFieldMap  = yaml.StructFieldMap
// )

// var (
// 	ToJSON      = yaml.YAMLToJSON
// 	FromJSON    = yaml.JSONToYAML
// 	HeadComment = yaml.HeadComment
// 	LineComment = yaml.LineComment
// 	PathString  = yaml.PathString
// )

// // EncodeOption
// var (
// 	JSON                       = yaml.JSON
// 	UseJSONMarshaler           = yaml.UseJSONMarshaler
// 	Indent                     = yaml.Indent
// 	IndentSequence             = yaml.IndentSequence
// 	Flow                       = yaml.Flow
// 	UseLiteralStyleIfMultiline = yaml.UseLiteralStyleIfMultiline
// 	MarshalAnchor              = yaml.MarshalAnchor
// 	WithComment                = yaml.WithComment
// )

// // DecodeOption
// var (
// 	ReferenceReaders     = yaml.ReferenceReaders
// 	ReferenceFiles       = yaml.ReferenceFiles
// 	ReferenceDirs        = yaml.ReferenceDirs
// 	RecursiveDir         = yaml.RecursiveDir
// 	Validator            = yaml.Validator
// 	Strict               = yaml.Strict
// 	DisallowUnknownField = yaml.DisallowUnknownField
// 	DisallowDuplicateKey = yaml.DisallowDuplicateKey
// 	UseOrderedMap        = yaml.UseOrderedMap
// 	UseJSONUnmarshaler   = yaml.UseJSONUnmarshaler
// 	CommentToMap         = yaml.CommentToMap
// )
