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

func easyjson3a862f94Decode20242FIGHTCLUBDomain(in *jlexer.Lexer, out *UpdateAdRequest) {
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
		case "cityName":
			out.CityName = string(in.String())
		case "address":
			out.Address = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "roomsNumber":
			out.RoomsNumber = int(in.Int())
		case "dateFrom":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.DateFrom).UnmarshalJSON(data))
			}
		case "dateTo":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.DateTo).UnmarshalJSON(data))
			}
		case "rooms":
			if in.IsNull() {
				in.Skip()
				out.Rooms = nil
			} else {
				in.Delim('[')
				if out.Rooms == nil {
					if !in.IsDelim(']') {
						out.Rooms = make([]AdRoomsResponse, 0, 2)
					} else {
						out.Rooms = []AdRoomsResponse{}
					}
				} else {
					out.Rooms = (out.Rooms)[:0]
				}
				for !in.IsDelim(']') {
					var v1 AdRoomsResponse
					easyjson3a862f94Decode20242FIGHTCLUBDomain1(in, &v1)
					out.Rooms = append(out.Rooms, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "squareMeters":
			out.SquareMeters = int(in.Int())
		case "floor":
			out.Floor = int(in.Int())
		case "buildingType":
			out.BuildingType = string(in.String())
		case "hasBalcony":
			out.HasBalcony = bool(in.Bool())
		case "hasElevator":
			out.HasElevator = bool(in.Bool())
		case "hasGas":
			out.HasGas = bool(in.Bool())
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
func easyjson3a862f94Encode20242FIGHTCLUBDomain(out *jwriter.Writer, in UpdateAdRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"cityName\":"
		out.RawString(prefix[1:])
		out.String(string(in.CityName))
	}
	{
		const prefix string = ",\"address\":"
		out.RawString(prefix)
		out.String(string(in.Address))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"roomsNumber\":"
		out.RawString(prefix)
		out.Int(int(in.RoomsNumber))
	}
	{
		const prefix string = ",\"dateFrom\":"
		out.RawString(prefix)
		out.Raw((in.DateFrom).MarshalJSON())
	}
	{
		const prefix string = ",\"dateTo\":"
		out.RawString(prefix)
		out.Raw((in.DateTo).MarshalJSON())
	}
	{
		const prefix string = ",\"rooms\":"
		out.RawString(prefix)
		if in.Rooms == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Rooms {
				if v2 > 0 {
					out.RawByte(',')
				}
				easyjson3a862f94Encode20242FIGHTCLUBDomain1(out, v3)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"squareMeters\":"
		out.RawString(prefix)
		out.Int(int(in.SquareMeters))
	}
	{
		const prefix string = ",\"floor\":"
		out.RawString(prefix)
		out.Int(int(in.Floor))
	}
	{
		const prefix string = ",\"buildingType\":"
		out.RawString(prefix)
		out.String(string(in.BuildingType))
	}
	{
		const prefix string = ",\"hasBalcony\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasBalcony))
	}
	{
		const prefix string = ",\"hasElevator\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasElevator))
	}
	{
		const prefix string = ",\"hasGas\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasGas))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UpdateAdRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3a862f94Encode20242FIGHTCLUBDomain(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UpdateAdRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3a862f94Encode20242FIGHTCLUBDomain(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UpdateAdRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3a862f94Decode20242FIGHTCLUBDomain(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UpdateAdRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3a862f94Decode20242FIGHTCLUBDomain(l, v)
}
func easyjson3a862f94Decode20242FIGHTCLUBDomain1(in *jlexer.Lexer, out *AdRoomsResponse) {
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
		case "type":
			out.Type = string(in.String())
		case "squareMeters":
			out.SquareMeters = int(in.Int())
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
func easyjson3a862f94Encode20242FIGHTCLUBDomain1(out *jwriter.Writer, in AdRoomsResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix[1:])
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"squareMeters\":"
		out.RawString(prefix)
		out.Int(int(in.SquareMeters))
	}
	out.RawByte('}')
}
func easyjson3a862f94Decode20242FIGHTCLUBDomain2(in *jlexer.Lexer, out *PlacesResponse) {
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
		case "places":
			(out.Places).UnmarshalEasyJSON(in)
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
func easyjson3a862f94Encode20242FIGHTCLUBDomain2(out *jwriter.Writer, in PlacesResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"places\":"
		out.RawString(prefix[1:])
		(in.Places).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PlacesResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3a862f94Encode20242FIGHTCLUBDomain2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PlacesResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3a862f94Encode20242FIGHTCLUBDomain2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PlacesResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3a862f94Decode20242FIGHTCLUBDomain2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PlacesResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3a862f94Decode20242FIGHTCLUBDomain2(l, v)
}
func easyjson3a862f94Decode20242FIGHTCLUBDomain3(in *jlexer.Lexer, out *PaymentInfo) {
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
		case "cardNumber":
			out.CardNumber = string(in.String())
		case "cardExpiry":
			out.CardExpiry = string(in.String())
		case "cardCVC":
			out.CardCvc = string(in.String())
		case "donationAmount":
			out.DonationAmount = string(in.String())
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
func easyjson3a862f94Encode20242FIGHTCLUBDomain3(out *jwriter.Writer, in PaymentInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"cardNumber\":"
		out.RawString(prefix[1:])
		out.String(string(in.CardNumber))
	}
	{
		const prefix string = ",\"cardExpiry\":"
		out.RawString(prefix)
		out.String(string(in.CardExpiry))
	}
	{
		const prefix string = ",\"cardCVC\":"
		out.RawString(prefix)
		out.String(string(in.CardCvc))
	}
	{
		const prefix string = ",\"donationAmount\":"
		out.RawString(prefix)
		out.String(string(in.DonationAmount))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PaymentInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3a862f94Encode20242FIGHTCLUBDomain3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PaymentInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3a862f94Encode20242FIGHTCLUBDomain3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PaymentInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3a862f94Decode20242FIGHTCLUBDomain3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PaymentInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3a862f94Decode20242FIGHTCLUBDomain3(l, v)
}
func easyjson3a862f94Decode20242FIGHTCLUBDomain4(in *jlexer.Lexer, out *GetOneAdResponse) {
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
		case "place":
			(out.Place).UnmarshalEasyJSON(in)
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
func easyjson3a862f94Encode20242FIGHTCLUBDomain4(out *jwriter.Writer, in GetOneAdResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"place\":"
		out.RawString(prefix[1:])
		(in.Place).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetOneAdResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3a862f94Encode20242FIGHTCLUBDomain4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetOneAdResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3a862f94Encode20242FIGHTCLUBDomain4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetOneAdResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3a862f94Decode20242FIGHTCLUBDomain4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetOneAdResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3a862f94Decode20242FIGHTCLUBDomain4(l, v)
}
func easyjson3a862f94Decode20242FIGHTCLUBDomain5(in *jlexer.Lexer, out *GetAllAdsResponse) {
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
			out.UUID = string(in.String())
		case "cityId":
			out.CityID = int(in.Int())
		case "authorUUID":
			out.AuthorUUID = string(in.String())
		case "address":
			out.Address = string(in.String())
		case "publicationDate":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.PublicationDate).UnmarshalJSON(data))
			}
		case "description":
			out.Description = string(in.String())
		case "roomsNumber":
			out.RoomsNumber = int(in.Int())
		case "viewsCount":
			out.ViewsCount = int(in.Int())
		case "squareMeters":
			out.SquareMeters = int(in.Int())
		case "floor":
			out.Floor = int(in.Int())
		case "buildingType":
			out.BuildingType = string(in.String())
		case "hasBalcony":
			out.HasBalcony = bool(in.Bool())
		case "hasElevator":
			out.HasElevator = bool(in.Bool())
		case "hasGas":
			out.HasGas = bool(in.Bool())
		case "likesCount":
			out.LikesCount = int(in.Int())
		case "priority":
			out.Priority = int(in.Int())
		case "endBoostDate":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.EndBoostDate).UnmarshalJSON(data))
			}
		case "cityName":
			out.CityName = string(in.String())
		case "adDateFrom":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.AdDateFrom).UnmarshalJSON(data))
			}
		case "adDateTo":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.AdDateTo).UnmarshalJSON(data))
			}
		case "isFavorite":
			out.IsFavorite = bool(in.Bool())
		case "author":
			easyjson3a862f94Decode20242FIGHTCLUBDomain6(in, &out.AdAuthor)
		case "images":
			if in.IsNull() {
				in.Skip()
				out.Images = nil
			} else {
				in.Delim('[')
				if out.Images == nil {
					if !in.IsDelim(']') {
						out.Images = make([]ImageResponse, 0, 2)
					} else {
						out.Images = []ImageResponse{}
					}
				} else {
					out.Images = (out.Images)[:0]
				}
				for !in.IsDelim(']') {
					var v4 ImageResponse
					easyjson3a862f94Decode20242FIGHTCLUBDomain7(in, &v4)
					out.Images = append(out.Images, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "rooms":
			if in.IsNull() {
				in.Skip()
				out.Rooms = nil
			} else {
				in.Delim('[')
				if out.Rooms == nil {
					if !in.IsDelim(']') {
						out.Rooms = make([]AdRoomsResponse, 0, 2)
					} else {
						out.Rooms = []AdRoomsResponse{}
					}
				} else {
					out.Rooms = (out.Rooms)[:0]
				}
				for !in.IsDelim(']') {
					var v5 AdRoomsResponse
					easyjson3a862f94Decode20242FIGHTCLUBDomain1(in, &v5)
					out.Rooms = append(out.Rooms, v5)
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
func easyjson3a862f94Encode20242FIGHTCLUBDomain5(out *jwriter.Writer, in GetAllAdsResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.UUID))
	}
	{
		const prefix string = ",\"cityId\":"
		out.RawString(prefix)
		out.Int(int(in.CityID))
	}
	{
		const prefix string = ",\"authorUUID\":"
		out.RawString(prefix)
		out.String(string(in.AuthorUUID))
	}
	{
		const prefix string = ",\"address\":"
		out.RawString(prefix)
		out.String(string(in.Address))
	}
	{
		const prefix string = ",\"publicationDate\":"
		out.RawString(prefix)
		out.Raw((in.PublicationDate).MarshalJSON())
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"roomsNumber\":"
		out.RawString(prefix)
		out.Int(int(in.RoomsNumber))
	}
	{
		const prefix string = ",\"viewsCount\":"
		out.RawString(prefix)
		out.Int(int(in.ViewsCount))
	}
	{
		const prefix string = ",\"squareMeters\":"
		out.RawString(prefix)
		out.Int(int(in.SquareMeters))
	}
	{
		const prefix string = ",\"floor\":"
		out.RawString(prefix)
		out.Int(int(in.Floor))
	}
	{
		const prefix string = ",\"buildingType\":"
		out.RawString(prefix)
		out.String(string(in.BuildingType))
	}
	{
		const prefix string = ",\"hasBalcony\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasBalcony))
	}
	{
		const prefix string = ",\"hasElevator\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasElevator))
	}
	{
		const prefix string = ",\"hasGas\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasGas))
	}
	{
		const prefix string = ",\"likesCount\":"
		out.RawString(prefix)
		out.Int(int(in.LikesCount))
	}
	{
		const prefix string = ",\"priority\":"
		out.RawString(prefix)
		out.Int(int(in.Priority))
	}
	{
		const prefix string = ",\"endBoostDate\":"
		out.RawString(prefix)
		out.Raw((in.EndBoostDate).MarshalJSON())
	}
	{
		const prefix string = ",\"cityName\":"
		out.RawString(prefix)
		out.String(string(in.CityName))
	}
	{
		const prefix string = ",\"adDateFrom\":"
		out.RawString(prefix)
		out.Raw((in.AdDateFrom).MarshalJSON())
	}
	{
		const prefix string = ",\"adDateTo\":"
		out.RawString(prefix)
		out.Raw((in.AdDateTo).MarshalJSON())
	}
	{
		const prefix string = ",\"isFavorite\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsFavorite))
	}
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix)
		easyjson3a862f94Encode20242FIGHTCLUBDomain6(out, in.AdAuthor)
	}
	{
		const prefix string = ",\"images\":"
		out.RawString(prefix)
		if in.Images == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v6, v7 := range in.Images {
				if v6 > 0 {
					out.RawByte(',')
				}
				easyjson3a862f94Encode20242FIGHTCLUBDomain7(out, v7)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"rooms\":"
		out.RawString(prefix)
		if in.Rooms == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Rooms {
				if v8 > 0 {
					out.RawByte(',')
				}
				easyjson3a862f94Encode20242FIGHTCLUBDomain1(out, v9)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetAllAdsResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3a862f94Encode20242FIGHTCLUBDomain5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetAllAdsResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3a862f94Encode20242FIGHTCLUBDomain5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetAllAdsResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3a862f94Decode20242FIGHTCLUBDomain5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetAllAdsResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3a862f94Decode20242FIGHTCLUBDomain5(l, v)
}
func easyjson3a862f94Decode20242FIGHTCLUBDomain7(in *jlexer.Lexer, out *ImageResponse) {
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
			out.ID = int(in.Int())
		case "path":
			out.ImagePath = string(in.String())
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
func easyjson3a862f94Encode20242FIGHTCLUBDomain7(out *jwriter.Writer, in ImageResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"path\":"
		out.RawString(prefix)
		out.String(string(in.ImagePath))
	}
	out.RawByte('}')
}
func easyjson3a862f94Decode20242FIGHTCLUBDomain6(in *jlexer.Lexer, out *UserResponce) {
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
		case "rating":
			out.Rating = float64(in.Float64())
		case "avatar":
			out.Avatar = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "sex":
			out.Sex = string(in.String())
		case "birthDate":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Birthdate).UnmarshalJSON(data))
			}
		case "guestCount":
			out.GuestCount = int(in.Int())
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
func easyjson3a862f94Encode20242FIGHTCLUBDomain6(out *jwriter.Writer, in UserResponce) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"rating\":"
		out.RawString(prefix[1:])
		out.Float64(float64(in.Rating))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"sex\":"
		out.RawString(prefix)
		out.String(string(in.Sex))
	}
	{
		const prefix string = ",\"birthDate\":"
		out.RawString(prefix)
		out.Raw((in.Birthdate).MarshalJSON())
	}
	{
		const prefix string = ",\"guestCount\":"
		out.RawString(prefix)
		out.Int(int(in.GuestCount))
	}
	out.RawByte('}')
}
func easyjson3a862f94Decode20242FIGHTCLUBDomain8(in *jlexer.Lexer, out *GetAllAdsListResponse) {
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
		case "housing":
			if in.IsNull() {
				in.Skip()
				out.Housing = nil
			} else {
				in.Delim('[')
				if out.Housing == nil {
					if !in.IsDelim(']') {
						out.Housing = make([]GetAllAdsResponse, 0, 0)
					} else {
						out.Housing = []GetAllAdsResponse{}
					}
				} else {
					out.Housing = (out.Housing)[:0]
				}
				for !in.IsDelim(']') {
					var v10 GetAllAdsResponse
					(v10).UnmarshalEasyJSON(in)
					out.Housing = append(out.Housing, v10)
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
func easyjson3a862f94Encode20242FIGHTCLUBDomain8(out *jwriter.Writer, in GetAllAdsListResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"housing\":"
		out.RawString(prefix[1:])
		if in.Housing == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v11, v12 := range in.Housing {
				if v11 > 0 {
					out.RawByte(',')
				}
				(v12).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetAllAdsListResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3a862f94Encode20242FIGHTCLUBDomain8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetAllAdsListResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3a862f94Encode20242FIGHTCLUBDomain8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetAllAdsListResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3a862f94Decode20242FIGHTCLUBDomain8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetAllAdsListResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3a862f94Decode20242FIGHTCLUBDomain8(l, v)
}
func easyjson3a862f94Decode20242FIGHTCLUBDomain9(in *jlexer.Lexer, out *Favorites) {
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
		case "adId":
			out.AdId = string(in.String())
		case "userId":
			out.UserId = string(in.String())
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
func easyjson3a862f94Encode20242FIGHTCLUBDomain9(out *jwriter.Writer, in Favorites) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"adId\":"
		out.RawString(prefix[1:])
		out.String(string(in.AdId))
	}
	{
		const prefix string = ",\"userId\":"
		out.RawString(prefix)
		out.String(string(in.UserId))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Favorites) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3a862f94Encode20242FIGHTCLUBDomain9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Favorites) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3a862f94Encode20242FIGHTCLUBDomain9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Favorites) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3a862f94Decode20242FIGHTCLUBDomain9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Favorites) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3a862f94Decode20242FIGHTCLUBDomain9(l, v)
}
func easyjson3a862f94Decode20242FIGHTCLUBDomain10(in *jlexer.Lexer, out *CreateAdRequest) {
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
		case "cityName":
			out.CityName = string(in.String())
		case "address":
			out.Address = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "roomsNumber":
			out.RoomsNumber = int(in.Int())
		case "dateFrom":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.DateFrom).UnmarshalJSON(data))
			}
		case "dateTo":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.DateTo).UnmarshalJSON(data))
			}
		case "rooms":
			if in.IsNull() {
				in.Skip()
				out.Rooms = nil
			} else {
				in.Delim('[')
				if out.Rooms == nil {
					if !in.IsDelim(']') {
						out.Rooms = make([]AdRoomsResponse, 0, 2)
					} else {
						out.Rooms = []AdRoomsResponse{}
					}
				} else {
					out.Rooms = (out.Rooms)[:0]
				}
				for !in.IsDelim(']') {
					var v13 AdRoomsResponse
					easyjson3a862f94Decode20242FIGHTCLUBDomain1(in, &v13)
					out.Rooms = append(out.Rooms, v13)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "squareMeters":
			out.SquareMeters = int(in.Int())
		case "floor":
			out.Floor = int(in.Int())
		case "buildingType":
			out.BuildingType = string(in.String())
		case "hasBalcony":
			out.HasBalcony = bool(in.Bool())
		case "hasElevator":
			out.HasElevator = bool(in.Bool())
		case "hasGas":
			out.HasGas = bool(in.Bool())
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
func easyjson3a862f94Encode20242FIGHTCLUBDomain10(out *jwriter.Writer, in CreateAdRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"cityName\":"
		out.RawString(prefix[1:])
		out.String(string(in.CityName))
	}
	{
		const prefix string = ",\"address\":"
		out.RawString(prefix)
		out.String(string(in.Address))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"roomsNumber\":"
		out.RawString(prefix)
		out.Int(int(in.RoomsNumber))
	}
	{
		const prefix string = ",\"dateFrom\":"
		out.RawString(prefix)
		out.Raw((in.DateFrom).MarshalJSON())
	}
	{
		const prefix string = ",\"dateTo\":"
		out.RawString(prefix)
		out.Raw((in.DateTo).MarshalJSON())
	}
	{
		const prefix string = ",\"rooms\":"
		out.RawString(prefix)
		if in.Rooms == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v14, v15 := range in.Rooms {
				if v14 > 0 {
					out.RawByte(',')
				}
				easyjson3a862f94Encode20242FIGHTCLUBDomain1(out, v15)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"squareMeters\":"
		out.RawString(prefix)
		out.Int(int(in.SquareMeters))
	}
	{
		const prefix string = ",\"floor\":"
		out.RawString(prefix)
		out.Int(int(in.Floor))
	}
	{
		const prefix string = ",\"buildingType\":"
		out.RawString(prefix)
		out.String(string(in.BuildingType))
	}
	{
		const prefix string = ",\"hasBalcony\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasBalcony))
	}
	{
		const prefix string = ",\"hasElevator\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasElevator))
	}
	{
		const prefix string = ",\"hasGas\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasGas))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateAdRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3a862f94Encode20242FIGHTCLUBDomain10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateAdRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3a862f94Encode20242FIGHTCLUBDomain10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateAdRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3a862f94Decode20242FIGHTCLUBDomain10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateAdRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3a862f94Decode20242FIGHTCLUBDomain10(l, v)
}
func easyjson3a862f94Decode20242FIGHTCLUBDomain11(in *jlexer.Lexer, out *AdFilter) {
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
		case "Location":
			out.Location = string(in.String())
		case "Rating":
			out.Rating = string(in.String())
		case "NewThisWeek":
			out.NewThisWeek = string(in.String())
		case "HostGender":
			out.HostGender = string(in.String())
		case "GuestCount":
			out.GuestCount = string(in.String())
		case "Limit":
			out.Limit = int(in.Int())
		case "Offset":
			out.Offset = int(in.Int())
		case "DateFrom":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.DateFrom).UnmarshalJSON(data))
			}
		case "DateTo":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.DateTo).UnmarshalJSON(data))
			}
		case "Favorites":
			out.Favorites = string(in.String())
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
func easyjson3a862f94Encode20242FIGHTCLUBDomain11(out *jwriter.Writer, in AdFilter) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Location\":"
		out.RawString(prefix[1:])
		out.String(string(in.Location))
	}
	{
		const prefix string = ",\"Rating\":"
		out.RawString(prefix)
		out.String(string(in.Rating))
	}
	{
		const prefix string = ",\"NewThisWeek\":"
		out.RawString(prefix)
		out.String(string(in.NewThisWeek))
	}
	{
		const prefix string = ",\"HostGender\":"
		out.RawString(prefix)
		out.String(string(in.HostGender))
	}
	{
		const prefix string = ",\"GuestCount\":"
		out.RawString(prefix)
		out.String(string(in.GuestCount))
	}
	{
		const prefix string = ",\"Limit\":"
		out.RawString(prefix)
		out.Int(int(in.Limit))
	}
	{
		const prefix string = ",\"Offset\":"
		out.RawString(prefix)
		out.Int(int(in.Offset))
	}
	{
		const prefix string = ",\"DateFrom\":"
		out.RawString(prefix)
		out.Raw((in.DateFrom).MarshalJSON())
	}
	{
		const prefix string = ",\"DateTo\":"
		out.RawString(prefix)
		out.Raw((in.DateTo).MarshalJSON())
	}
	{
		const prefix string = ",\"Favorites\":"
		out.RawString(prefix)
		out.String(string(in.Favorites))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AdFilter) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3a862f94Encode20242FIGHTCLUBDomain11(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AdFilter) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3a862f94Encode20242FIGHTCLUBDomain11(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AdFilter) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3a862f94Decode20242FIGHTCLUBDomain11(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AdFilter) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3a862f94Decode20242FIGHTCLUBDomain11(l, v)
}
func easyjson3a862f94Decode20242FIGHTCLUBDomain12(in *jlexer.Lexer, out *Ad) {
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
			out.UUID = string(in.String())
		case "cityId":
			out.CityID = int(in.Int())
		case "authorUUID":
			out.AuthorUUID = string(in.String())
		case "address":
			out.Address = string(in.String())
		case "publicationDate":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.PublicationDate).UnmarshalJSON(data))
			}
		case "description":
			out.Description = string(in.String())
		case "roomsNumber":
			out.RoomsNumber = int(in.Int())
		case "viewsCount":
			out.ViewsCount = int(in.Int())
		case "squareMeters":
			out.SquareMeters = int(in.Int())
		case "floor":
			out.Floor = int(in.Int())
		case "buildingType":
			out.BuildingType = string(in.String())
		case "hasBalcony":
			out.HasBalcony = bool(in.Bool())
		case "hasElevator":
			out.HasElevator = bool(in.Bool())
		case "hasGas":
			out.HasGas = bool(in.Bool())
		case "likesCount":
			out.LikesCount = int(in.Int())
		case "priority":
			out.Priority = int(in.Int())
		case "endBoostDate":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.EndBoostDate).UnmarshalJSON(data))
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
func easyjson3a862f94Encode20242FIGHTCLUBDomain12(out *jwriter.Writer, in Ad) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.UUID))
	}
	{
		const prefix string = ",\"cityId\":"
		out.RawString(prefix)
		out.Int(int(in.CityID))
	}
	{
		const prefix string = ",\"authorUUID\":"
		out.RawString(prefix)
		out.String(string(in.AuthorUUID))
	}
	{
		const prefix string = ",\"address\":"
		out.RawString(prefix)
		out.String(string(in.Address))
	}
	{
		const prefix string = ",\"publicationDate\":"
		out.RawString(prefix)
		out.Raw((in.PublicationDate).MarshalJSON())
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"roomsNumber\":"
		out.RawString(prefix)
		out.Int(int(in.RoomsNumber))
	}
	{
		const prefix string = ",\"viewsCount\":"
		out.RawString(prefix)
		out.Int(int(in.ViewsCount))
	}
	{
		const prefix string = ",\"squareMeters\":"
		out.RawString(prefix)
		out.Int(int(in.SquareMeters))
	}
	{
		const prefix string = ",\"floor\":"
		out.RawString(prefix)
		out.Int(int(in.Floor))
	}
	{
		const prefix string = ",\"buildingType\":"
		out.RawString(prefix)
		out.String(string(in.BuildingType))
	}
	{
		const prefix string = ",\"hasBalcony\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasBalcony))
	}
	{
		const prefix string = ",\"hasElevator\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasElevator))
	}
	{
		const prefix string = ",\"hasGas\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasGas))
	}
	{
		const prefix string = ",\"likesCount\":"
		out.RawString(prefix)
		out.Int(int(in.LikesCount))
	}
	{
		const prefix string = ",\"priority\":"
		out.RawString(prefix)
		out.Int(int(in.Priority))
	}
	{
		const prefix string = ",\"endBoostDate\":"
		out.RawString(prefix)
		out.Raw((in.EndBoostDate).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Ad) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3a862f94Encode20242FIGHTCLUBDomain12(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Ad) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3a862f94Encode20242FIGHTCLUBDomain12(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Ad) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3a862f94Decode20242FIGHTCLUBDomain12(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Ad) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3a862f94Decode20242FIGHTCLUBDomain12(l, v)
}