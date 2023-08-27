// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD2b7633eDecodeGithubComDupreehkudaAvitoSegmentsInternalModels(in *jlexer.Lexer, out *UserSetRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "userID":
			out.UserID = string(in.String())
		case "segments":
			if in.IsNull() {
				in.Skip()
				out.Segments = nil
			} else {
				in.Delim('[')
				if out.Segments == nil {
					if !in.IsDelim(']') {
						out.Segments = make([]UserSegment, 0, 1)
					} else {
						out.Segments = []UserSegment{}
					}
				} else {
					out.Segments = (out.Segments)[:0]
				}
				for !in.IsDelim(']') {
					var v1 UserSegment
					(v1).UnmarshalEasyJSON(in)
					out.Segments = append(out.Segments, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComDupreehkudaAvitoSegmentsInternalModels(out *jwriter.Writer, in UserSetRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"userID\":"
		out.RawString(prefix[1:])
		out.String(string(in.UserID))
	}
	{
		const prefix string = ",\"segments\":"
		out.RawString(prefix)
		if in.Segments == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Segments {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserSetRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComDupreehkudaAvitoSegmentsInternalModels(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserSetRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComDupreehkudaAvitoSegmentsInternalModels(l, v)
}
func easyjsonD2b7633eDecodeGithubComDupreehkudaAvitoSegmentsInternalModels1(in *jlexer.Lexer, out *UserSegment) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "slug":
			out.Slug = string(in.String())
		case "expire":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Expire).UnmarshalJSON(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComDupreehkudaAvitoSegmentsInternalModels1(out *jwriter.Writer, in UserSegment) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"slug\":"
		out.RawString(prefix[1:])
		out.String(string(in.Slug))
	}
	if true {
		const prefix string = ",\"expire\":"
		out.RawString(prefix)
		out.Raw((in.Expire).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserSegment) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComDupreehkudaAvitoSegmentsInternalModels1(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserSegment) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComDupreehkudaAvitoSegmentsInternalModels1(l, v)
}
func easyjsonD2b7633eDecodeGithubComDupreehkudaAvitoSegmentsInternalModels2(in *jlexer.Lexer, out *UserResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "userID":
			out.UserID = string(in.String())
		case "slugs":
			if in.IsNull() {
				in.Skip()
				out.Slugs = nil
			} else {
				in.Delim('[')
				if out.Slugs == nil {
					if !in.IsDelim(']') {
						out.Slugs = make([]string, 0, 4)
					} else {
						out.Slugs = []string{}
					}
				} else {
					out.Slugs = (out.Slugs)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.Slugs = append(out.Slugs, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComDupreehkudaAvitoSegmentsInternalModels2(out *jwriter.Writer, in UserResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"userID\":"
		out.RawString(prefix[1:])
		out.String(string(in.UserID))
	}
	{
		const prefix string = ",\"slugs\":"
		out.RawString(prefix)
		if in.Slugs == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Slugs {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComDupreehkudaAvitoSegmentsInternalModels2(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComDupreehkudaAvitoSegmentsInternalModels2(l, v)
}
func easyjsonD2b7633eDecodeGithubComDupreehkudaAvitoSegmentsInternalModels3(in *jlexer.Lexer, out *UserDeleteRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "userID":
			out.UserID = string(in.String())
		case "slugs":
			if in.IsNull() {
				in.Skip()
				out.Slugs = nil
			} else {
				in.Delim('[')
				if out.Slugs == nil {
					if !in.IsDelim(']') {
						out.Slugs = make([]string, 0, 4)
					} else {
						out.Slugs = []string{}
					}
				} else {
					out.Slugs = (out.Slugs)[:0]
				}
				for !in.IsDelim(']') {
					var v7 string
					v7 = string(in.String())
					out.Slugs = append(out.Slugs, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComDupreehkudaAvitoSegmentsInternalModels3(out *jwriter.Writer, in UserDeleteRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"userID\":"
		out.RawString(prefix[1:])
		out.String(string(in.UserID))
	}
	{
		const prefix string = ",\"slugs\":"
		out.RawString(prefix)
		if in.Slugs == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Slugs {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.String(string(v9))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserDeleteRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComDupreehkudaAvitoSegmentsInternalModels3(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserDeleteRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComDupreehkudaAvitoSegmentsInternalModels3(l, v)
}
func easyjsonD2b7633eDecodeGithubComDupreehkudaAvitoSegmentsInternalModels4(in *jlexer.Lexer, out *Segment) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "slug":
			out.Slug = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "DeletedAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.DeletedAt).UnmarshalJSON(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComDupreehkudaAvitoSegmentsInternalModels4(out *jwriter.Writer, in Segment) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"slug\":"
		out.RawString(prefix[1:])
		out.String(string(in.Slug))
	}
	if in.Description != "" {
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"DeletedAt\":"
		out.RawString(prefix)
		out.Raw((in.DeletedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Segment) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComDupreehkudaAvitoSegmentsInternalModels4(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Segment) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComDupreehkudaAvitoSegmentsInternalModels4(l, v)
}
