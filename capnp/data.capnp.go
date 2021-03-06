// Code generated by capnpc-go. DO NOT EDIT.

package data

import (
	capnp "zombiezen.com/go/capnproto2"
	text "zombiezen.com/go/capnproto2/encoding/text"
	schemas "zombiezen.com/go/capnproto2/schemas"
)

type Stream struct{ capnp.Struct }

// Stream_TypeID is the unique identifier for the type Stream.
const Stream_TypeID = 0xcf2858cf0a901da3

func NewStream(s *capnp.Segment) (Stream, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return Stream{st}, err
}

func NewRootStream(s *capnp.Segment) (Stream, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return Stream{st}, err
}

func ReadRootStream(msg *capnp.Message) (Stream, error) {
	root, err := msg.RootPtr()
	return Stream{root.Struct()}, err
}

func (s Stream) String() string {
	str, _ := text.Marshal(0xcf2858cf0a901da3, s.Struct)
	return str
}

func (s Stream) Times() (capnp.Int64List, error) {
	p, err := s.Struct.Ptr(0)
	return capnp.Int64List{List: p.List()}, err
}

func (s Stream) HasTimes() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Stream) SetTimes(v capnp.Int64List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewTimes sets the times field to a newly
// allocated capnp.Int64List, preferring placement in s's segment.
func (s Stream) NewTimes(n int32) (capnp.Int64List, error) {
	l, err := capnp.NewInt64List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.Int64List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Stream) Values() (capnp.Float64List, error) {
	p, err := s.Struct.Ptr(1)
	return capnp.Float64List{List: p.List()}, err
}

func (s Stream) HasValues() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Stream) SetValues(v capnp.Float64List) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewValues sets the values field to a newly
// allocated capnp.Float64List, preferring placement in s's segment.
func (s Stream) NewValues(n int32) (capnp.Float64List, error) {
	l, err := capnp.NewFloat64List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.Float64List{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
}

// Stream_List is a list of Stream.
type Stream_List struct{ capnp.List }

// NewStream creates a new list of Stream.
func NewStream_List(s *capnp.Segment, sz int32) (Stream_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	return Stream_List{l}, err
}

func (s Stream_List) At(i int) Stream { return Stream{s.List.Struct(i)} }

func (s Stream_List) Set(i int, v Stream) error { return s.List.SetStruct(i, v.Struct) }

func (s Stream_List) String() string {
	str, _ := text.MarshalList(0xcf2858cf0a901da3, s.List)
	return str
}

// Stream_Promise is a wrapper for a Stream promised by a client call.
type Stream_Promise struct{ *capnp.Pipeline }

func (p Stream_Promise) Struct() (Stream, error) {
	s, err := p.Pipeline.Struct()
	return Stream{s}, err
}

type StreamCollection struct{ capnp.Struct }

// StreamCollection_TypeID is the unique identifier for the type StreamCollection.
const StreamCollection_TypeID = 0xc1faf22fef67d15b

func NewStreamCollection(s *capnp.Segment) (StreamCollection, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return StreamCollection{st}, err
}

func NewRootStreamCollection(s *capnp.Segment) (StreamCollection, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return StreamCollection{st}, err
}

func ReadRootStreamCollection(msg *capnp.Message) (StreamCollection, error) {
	root, err := msg.RootPtr()
	return StreamCollection{root.Struct()}, err
}

func (s StreamCollection) String() string {
	str, _ := text.Marshal(0xc1faf22fef67d15b, s.Struct)
	return str
}

func (s StreamCollection) Times() (capnp.Int64List, error) {
	p, err := s.Struct.Ptr(0)
	return capnp.Int64List{List: p.List()}, err
}

func (s StreamCollection) HasTimes() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s StreamCollection) SetTimes(v capnp.Int64List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewTimes sets the times field to a newly
// allocated capnp.Int64List, preferring placement in s's segment.
func (s StreamCollection) NewTimes(n int32) (capnp.Int64List, error) {
	l, err := capnp.NewInt64List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.Int64List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s StreamCollection) Streams() (Stream_List, error) {
	p, err := s.Struct.Ptr(1)
	return Stream_List{List: p.List()}, err
}

func (s StreamCollection) HasStreams() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s StreamCollection) SetStreams(v Stream_List) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewStreams sets the streams field to a newly
// allocated Stream_List, preferring placement in s's segment.
func (s StreamCollection) NewStreams(n int32) (Stream_List, error) {
	l, err := NewStream_List(s.Struct.Segment(), n)
	if err != nil {
		return Stream_List{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
}

// StreamCollection_List is a list of StreamCollection.
type StreamCollection_List struct{ capnp.List }

// NewStreamCollection creates a new list of StreamCollection.
func NewStreamCollection_List(s *capnp.Segment, sz int32) (StreamCollection_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	return StreamCollection_List{l}, err
}

func (s StreamCollection_List) At(i int) StreamCollection { return StreamCollection{s.List.Struct(i)} }

func (s StreamCollection_List) Set(i int, v StreamCollection) error {
	return s.List.SetStruct(i, v.Struct)
}

func (s StreamCollection_List) String() string {
	str, _ := text.MarshalList(0xc1faf22fef67d15b, s.List)
	return str
}

// StreamCollection_Promise is a wrapper for a StreamCollection promised by a client call.
type StreamCollection_Promise struct{ *capnp.Pipeline }

func (p StreamCollection_Promise) Struct() (StreamCollection, error) {
	s, err := p.Pipeline.Struct()
	return StreamCollection{s}, err
}

const schema_d213afc571b92e25 = "x\xda\x12\x98\xeb\xc0d\xc8*\xce\xc4\xc0\x10(\xc1\xca" +
	"\xf6?\xfab\xfa{\xfdO\xbf\x0e2\x08r3\xfeW" +
	"\xd5\xdbYxt\xbd\xf0%\x06V&v\x06\x06\xc1\x8f" +
	"\x8f\x04\xff\x82\xe8\x9f\xe5\x0c\x8c\xff\x17\xcbN\xe0:\x1f" +
	"\xa1q\x1e\x8bB\xe1R\xc6I\xc2\xb5\x8c V%c" +
	"9\x83\xee\xff\x94\xc4\x92D\xbd\xe4\xc4\x02\xa6\xbc\x02\xab" +
	"\xe0\x92\xa2\xd4\xc4\\\xe7\xfc\x9c\x9c\xd4\xe4\x12\xf6\xcc\xfc" +
	"\xbc\x00F\xc6@\x0ef\x16\x06\x06\x16F\x06\x06AM" +
	"#\x06\x86@\x15f\xc6@\x07&FAFF\x11F" +
	"\x90\xa0\xad\x13\x03C\xa0\x053c`\x08\x13\xa3|I" +
	"fnj1#\x1f\x03c\x003##+\x03\x13\x88" +
	"Y_\x0c6\x15.,\x80p\x1e\x03#H\x10\xee\x04" +
	"F\x98\x13\x18s\x89\xb1\xd8\x8a\x80\xc5\xf6e\x899\xa5" +
	"\x08Qn\x88( \x00\x00\xff\xffL\x86Ju"

func init() {
	schemas.Register(schema_d213afc571b92e25,
		0xc1faf22fef67d15b,
		0xcf2858cf0a901da3)
}
