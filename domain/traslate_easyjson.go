// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package domain

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

func easyjsonE4972506DecodeProjectMessengerDomain(in *jlexer.Lexer, out *Translations) {
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
		case "text":
			out.Text = string(in.String())
		case "detectedLanguageCode":
			out.DetectedLanguageCode = string(in.String())
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
func easyjsonE4972506EncodeProjectMessengerDomain(out *jwriter.Writer, in Translations) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix[1:])
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"detectedLanguageCode\":"
		out.RawString(prefix)
		out.String(string(in.DetectedLanguageCode))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Translations) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE4972506EncodeProjectMessengerDomain(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Translations) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE4972506EncodeProjectMessengerDomain(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Translations) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE4972506DecodeProjectMessengerDomain(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Translations) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE4972506DecodeProjectMessengerDomain(l, v)
}
func easyjsonE4972506DecodeProjectMessengerDomain1(in *jlexer.Lexer, out *TranslateResponse) {
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
		case "translations":
			if in.IsNull() {
				in.Skip()
				out.Translations = nil
			} else {
				in.Delim('[')
				if out.Translations == nil {
					if !in.IsDelim(']') {
						out.Translations = make([]Translations, 0, 2)
					} else {
						out.Translations = []Translations{}
					}
				} else {
					out.Translations = (out.Translations)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Translations
					(v1).UnmarshalEasyJSON(in)
					out.Translations = append(out.Translations, v1)
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
func easyjsonE4972506EncodeProjectMessengerDomain1(out *jwriter.Writer, in TranslateResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"translations\":"
		out.RawString(prefix[1:])
		if in.Translations == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Translations {
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

// MarshalJSON supports json.Marshaler interface
func (v TranslateResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE4972506EncodeProjectMessengerDomain1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v TranslateResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE4972506EncodeProjectMessengerDomain1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TranslateResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE4972506DecodeProjectMessengerDomain1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *TranslateResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE4972506DecodeProjectMessengerDomain1(l, v)
}
func easyjsonE4972506DecodeProjectMessengerDomain2(in *jlexer.Lexer, out *TranslateRequest) {
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
		case "texts":
			if in.IsNull() {
				in.Skip()
				out.Messages = nil
			} else {
				in.Delim('[')
				if out.Messages == nil {
					if !in.IsDelim(']') {
						out.Messages = make([]string, 0, 4)
					} else {
						out.Messages = []string{}
					}
				} else {
					out.Messages = (out.Messages)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.Messages = append(out.Messages, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "folderId":
			out.FolderID = string(in.String())
		case "targetLanguageCode":
			out.TargetLanguageCode = string(in.String())
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
func easyjsonE4972506EncodeProjectMessengerDomain2(out *jwriter.Writer, in TranslateRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"texts\":"
		out.RawString(prefix[1:])
		if in.Messages == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Messages {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"folderId\":"
		out.RawString(prefix)
		out.String(string(in.FolderID))
	}
	{
		const prefix string = ",\"targetLanguageCode\":"
		out.RawString(prefix)
		out.String(string(in.TargetLanguageCode))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v TranslateRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE4972506EncodeProjectMessengerDomain2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v TranslateRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE4972506EncodeProjectMessengerDomain2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TranslateRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE4972506DecodeProjectMessengerDomain2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *TranslateRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE4972506DecodeProjectMessengerDomain2(l, v)
}
