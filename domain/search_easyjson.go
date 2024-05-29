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

func easyjsonD4176298DecodeProjectMessengerDomain(in *jlexer.Lexer, out *SearchRequest) {
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
		case "word":
			out.Word = string(in.String())
		case "user_id":
			out.UserID = uint(in.Uint())
		case "search_type":
			out.Type = string(in.String())
		case "chat_id":
			out.ChatID = uint(in.Uint())
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
func easyjsonD4176298EncodeProjectMessengerDomain(out *jwriter.Writer, in SearchRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"word\":"
		out.RawString(prefix[1:])
		out.String(string(in.Word))
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.Uint(uint(in.UserID))
	}
	{
		const prefix string = ",\"search_type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"chat_id\":"
		out.RawString(prefix)
		out.Uint(uint(in.ChatID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SearchRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD4176298EncodeProjectMessengerDomain(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SearchRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD4176298EncodeProjectMessengerDomain(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SearchRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD4176298DecodeProjectMessengerDomain(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SearchRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD4176298DecodeProjectMessengerDomain(l, v)
}
func easyjsonD4176298DecodeProjectMessengerDomain1(in *jlexer.Lexer, out *MessagesSearchResponse) {
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
		case "messages":
			if in.IsNull() {
				in.Skip()
				out.Messages = nil
			} else {
				in.Delim('[')
				if out.Messages == nil {
					if !in.IsDelim(']') {
						out.Messages = make([]Message, 0, 0)
					} else {
						out.Messages = []Message{}
					}
				} else {
					out.Messages = (out.Messages)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Message
					easyjsonD4176298DecodeProjectMessengerDomain2(in, &v1)
					out.Messages = append(out.Messages, v1)
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
func easyjsonD4176298EncodeProjectMessengerDomain1(out *jwriter.Writer, in MessagesSearchResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"messages\":"
		out.RawString(prefix[1:])
		if in.Messages == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Messages {
				if v2 > 0 {
					out.RawByte(',')
				}
				easyjsonD4176298EncodeProjectMessengerDomain2(out, v3)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessagesSearchResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD4176298EncodeProjectMessengerDomain1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessagesSearchResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD4176298EncodeProjectMessengerDomain1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessagesSearchResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD4176298DecodeProjectMessengerDomain1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessagesSearchResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD4176298DecodeProjectMessengerDomain1(l, v)
}
func easyjsonD4176298DecodeProjectMessengerDomain2(in *jlexer.Lexer, out *Message) {
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
		case "id":
			out.ID = uint(in.Uint())
		case "chat_id":
			out.ChatID = uint(in.Uint())
		case "user_id":
			out.UserID = uint(in.Uint())
		case "message_text":
			out.Message = string(in.String())
		case "edited":
			out.Edited = bool(in.Bool())
		case "edited_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.EditedAt).UnmarshalJSON(data))
			}
		case "sent_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "username":
			out.SenderUsername = string(in.String())
		case "file":
			if in.IsNull() {
				in.Skip()
				out.File = nil
			} else {
				if out.File == nil {
					out.File = new(FileInMessage)
				}
				easyjsonD4176298DecodeProjectMessengerDomain3(in, out.File)
			}
		case "sticker_path":
			out.StickerPath = string(in.String())
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
func easyjsonD4176298EncodeProjectMessengerDomain2(out *jwriter.Writer, in Message) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"chat_id\":"
		out.RawString(prefix)
		out.Uint(uint(in.ChatID))
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.Uint(uint(in.UserID))
	}
	{
		const prefix string = ",\"message_text\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"edited\":"
		out.RawString(prefix)
		out.Bool(bool(in.Edited))
	}
	{
		const prefix string = ",\"edited_at\":"
		out.RawString(prefix)
		out.Raw((in.EditedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"sent_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.SenderUsername))
	}
	{
		const prefix string = ",\"file\":"
		out.RawString(prefix)
		if in.File == nil {
			out.RawString("null")
		} else {
			easyjsonD4176298EncodeProjectMessengerDomain3(out, *in.File)
		}
	}
	{
		const prefix string = ",\"sticker_path\":"
		out.RawString(prefix)
		out.String(string(in.StickerPath))
	}
	out.RawByte('}')
}
func easyjsonD4176298DecodeProjectMessengerDomain3(in *jlexer.Lexer, out *FileInMessage) {
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
		case "original_name":
			out.OriginalName = string(in.String())
		case "path":
			out.Path = string(in.String())
		case "type":
			out.Type = string(in.String())
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
func easyjsonD4176298EncodeProjectMessengerDomain3(out *jwriter.Writer, in FileInMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"original_name\":"
		out.RawString(prefix[1:])
		out.String(string(in.OriginalName))
	}
	{
		const prefix string = ",\"path\":"
		out.RawString(prefix)
		out.String(string(in.Path))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	out.RawByte('}')
}
func easyjsonD4176298DecodeProjectMessengerDomain4(in *jlexer.Lexer, out *MessagesSearchRequest) {
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
		case "word":
			out.Word = string(in.String())
		case "user_id":
			out.UserID = uint(in.Uint())
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
func easyjsonD4176298EncodeProjectMessengerDomain4(out *jwriter.Writer, in MessagesSearchRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"word\":"
		out.RawString(prefix[1:])
		out.String(string(in.Word))
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.Uint(uint(in.UserID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessagesSearchRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD4176298EncodeProjectMessengerDomain4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessagesSearchRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD4176298EncodeProjectMessengerDomain4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessagesSearchRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD4176298DecodeProjectMessengerDomain4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessagesSearchRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD4176298DecodeProjectMessengerDomain4(l, v)
}
func easyjsonD4176298DecodeProjectMessengerDomain5(in *jlexer.Lexer, out *ContactsSearchResponse) {
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
		case "contacts":
			if in.IsNull() {
				in.Skip()
				out.Contacts = nil
			} else {
				in.Delim('[')
				if out.Contacts == nil {
					if !in.IsDelim(']') {
						out.Contacts = make([]Person, 0, 0)
					} else {
						out.Contacts = []Person{}
					}
				} else {
					out.Contacts = (out.Contacts)[:0]
				}
				for !in.IsDelim(']') {
					var v4 Person
					(v4).UnmarshalEasyJSON(in)
					out.Contacts = append(out.Contacts, v4)
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
func easyjsonD4176298EncodeProjectMessengerDomain5(out *jwriter.Writer, in ContactsSearchResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"contacts\":"
		out.RawString(prefix[1:])
		if in.Contacts == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Contacts {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ContactsSearchResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD4176298EncodeProjectMessengerDomain5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ContactsSearchResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD4176298EncodeProjectMessengerDomain5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ContactsSearchResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD4176298DecodeProjectMessengerDomain5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ContactsSearchResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD4176298DecodeProjectMessengerDomain5(l, v)
}
func easyjsonD4176298DecodeProjectMessengerDomain6(in *jlexer.Lexer, out *ContactsSearchRequest) {
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
		case "word":
			out.Word = string(in.String())
		case "user_id":
			out.UserID = uint(in.Uint())
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
func easyjsonD4176298EncodeProjectMessengerDomain6(out *jwriter.Writer, in ContactsSearchRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"word\":"
		out.RawString(prefix[1:])
		out.String(string(in.Word))
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.Uint(uint(in.UserID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ContactsSearchRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD4176298EncodeProjectMessengerDomain6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ContactsSearchRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD4176298EncodeProjectMessengerDomain6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ContactsSearchRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD4176298DecodeProjectMessengerDomain6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ContactsSearchRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD4176298DecodeProjectMessengerDomain6(l, v)
}
func easyjsonD4176298DecodeProjectMessengerDomain7(in *jlexer.Lexer, out *ChatSearchResponse) {
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
		case "chats":
			if in.IsNull() {
				in.Skip()
				out.Chats = nil
			} else {
				in.Delim('[')
				if out.Chats == nil {
					if !in.IsDelim(']') {
						out.Chats = make([]Chat, 0, 0)
					} else {
						out.Chats = []Chat{}
					}
				} else {
					out.Chats = (out.Chats)[:0]
				}
				for !in.IsDelim(']') {
					var v7 Chat
					easyjsonD4176298DecodeProjectMessengerDomain8(in, &v7)
					out.Chats = append(out.Chats, v7)
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
func easyjsonD4176298EncodeProjectMessengerDomain7(out *jwriter.Writer, in ChatSearchResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"chats\":"
		out.RawString(prefix[1:])
		if in.Chats == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Chats {
				if v8 > 0 {
					out.RawByte(',')
				}
				easyjsonD4176298EncodeProjectMessengerDomain8(out, v9)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ChatSearchResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD4176298EncodeProjectMessengerDomain7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ChatSearchResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD4176298EncodeProjectMessengerDomain7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ChatSearchResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD4176298DecodeProjectMessengerDomain7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ChatSearchResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD4176298DecodeProjectMessengerDomain7(l, v)
}
func easyjsonD4176298DecodeProjectMessengerDomain8(in *jlexer.Lexer, out *Chat) {
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
		case "id":
			out.ID = uint(in.Uint())
		case "type":
			out.Type = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "avatar":
			out.AvatarPath = string(in.String())
		case "creator":
			out.CreatorID = uint(in.Uint())
		case "messages":
			if in.IsNull() {
				in.Skip()
				out.Messages = nil
			} else {
				in.Delim('[')
				if out.Messages == nil {
					if !in.IsDelim(']') {
						out.Messages = make([]*Message, 0, 8)
					} else {
						out.Messages = []*Message{}
					}
				} else {
					out.Messages = (out.Messages)[:0]
				}
				for !in.IsDelim(']') {
					var v10 *Message
					if in.IsNull() {
						in.Skip()
						v10 = nil
					} else {
						if v10 == nil {
							v10 = new(Message)
						}
						easyjsonD4176298DecodeProjectMessengerDomain2(in, v10)
					}
					out.Messages = append(out.Messages, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "users":
			if in.IsNull() {
				in.Skip()
				out.Users = nil
			} else {
				in.Delim('[')
				if out.Users == nil {
					if !in.IsDelim(']') {
						out.Users = make([]*ChatUser, 0, 8)
					} else {
						out.Users = []*ChatUser{}
					}
				} else {
					out.Users = (out.Users)[:0]
				}
				for !in.IsDelim(']') {
					var v11 *ChatUser
					if in.IsNull() {
						in.Skip()
						v11 = nil
					} else {
						if v11 == nil {
							v11 = new(ChatUser)
						}
						easyjsonD4176298DecodeProjectMessengerDomain9(in, v11)
					}
					out.Users = append(out.Users, v11)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "created_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "edited_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.EditedAt).UnmarshalJSON(data))
			}
		case "last_action_date_time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.LastActionDateTime).UnmarshalJSON(data))
			}
		case "last_message":
			easyjsonD4176298DecodeProjectMessengerDomain2(in, &out.LastMessage)
		case "last_seen_message_id":
			out.LastSeenMessageID = int(in.Int())
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
func easyjsonD4176298EncodeProjectMessengerDomain8(out *jwriter.Writer, in Chat) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.AvatarPath))
	}
	{
		const prefix string = ",\"creator\":"
		out.RawString(prefix)
		out.Uint(uint(in.CreatorID))
	}
	{
		const prefix string = ",\"messages\":"
		out.RawString(prefix)
		if in.Messages == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v12, v13 := range in.Messages {
				if v12 > 0 {
					out.RawByte(',')
				}
				if v13 == nil {
					out.RawString("null")
				} else {
					easyjsonD4176298EncodeProjectMessengerDomain2(out, *v13)
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"users\":"
		out.RawString(prefix)
		if in.Users == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v14, v15 := range in.Users {
				if v14 > 0 {
					out.RawByte(',')
				}
				if v15 == nil {
					out.RawString("null")
				} else {
					easyjsonD4176298EncodeProjectMessengerDomain9(out, *v15)
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"edited_at\":"
		out.RawString(prefix)
		out.Raw((in.EditedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"last_action_date_time\":"
		out.RawString(prefix)
		out.Raw((in.LastActionDateTime).MarshalJSON())
	}
	{
		const prefix string = ",\"last_message\":"
		out.RawString(prefix)
		easyjsonD4176298EncodeProjectMessengerDomain2(out, in.LastMessage)
	}
	{
		const prefix string = ",\"last_seen_message_id\":"
		out.RawString(prefix)
		out.Int(int(in.LastSeenMessageID))
	}
	out.RawByte('}')
}
func easyjsonD4176298DecodeProjectMessengerDomain9(in *jlexer.Lexer, out *ChatUser) {
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
		case "chat_id":
			out.ChatID = int(in.Int())
		case "user_id":
			out.UserID = uint(in.Uint())
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
func easyjsonD4176298EncodeProjectMessengerDomain9(out *jwriter.Writer, in ChatUser) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"chat_id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ChatID))
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.Uint(uint(in.UserID))
	}
	out.RawByte('}')
}
func easyjsonD4176298DecodeProjectMessengerDomain10(in *jlexer.Lexer, out *ChatSearchRequest) {
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
		case "word":
			out.Word = string(in.String())
		case "user_id":
			out.UserID = uint(in.Uint())
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
func easyjsonD4176298EncodeProjectMessengerDomain10(out *jwriter.Writer, in ChatSearchRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"word\":"
		out.RawString(prefix[1:])
		out.String(string(in.Word))
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.Uint(uint(in.UserID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ChatSearchRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD4176298EncodeProjectMessengerDomain10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ChatSearchRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD4176298EncodeProjectMessengerDomain10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ChatSearchRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD4176298DecodeProjectMessengerDomain10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ChatSearchRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD4176298DecodeProjectMessengerDomain10(l, v)
}
