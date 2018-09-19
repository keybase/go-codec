// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"io"
	"strings"
	"testing"
)

type circularReader struct {
	b []byte
	n int
}

func (r *circularReader) Read(p []byte) (int, error) {
	for i := 0; i < len(p); i++ {
		p[i] = r.b[r.n]
		r.n = (r.n + 1) % len(r.b)
	}

	return len(p), nil
}

func assertMaxDepthError(t *testing.T, err error) {
	if err == nil || !strings.HasSuffix(err.Error(), "max depth exceeded") {
		t.Fatalf("Expected 'max depth exceeded', got %v", err)
	}
}

func testPattern(t *testing.T, b []byte) {
	r := circularReader{b: b}

	var h MsgpackHandle
	d := NewDecoder(&r, &h)

	var v interface{}
	err := d.Decode(&v)
	assertMaxDepthError(t, err)
}

func TestMsgpackDecodeInfiniteDepth(t *testing.T) {
	// [[[...
	testPattern(t, []byte{0x91})
	// {{{...
	testPattern(t, []byte{0x81})
	// [{[{...
	testPattern(t, []byte{0x91, 0x81})
	// [0x3f, {0x3f: [0x3f, {...
	testPattern(t, []byte{0x92, 0x3f, 0x81, 0x4e})
}

type selfer struct{}

func (s *selfer) CodecEncodeSelf(e *Encoder) {
	panic("CodecEncodeSelf unexpectedly called")
}

func (s *selfer) CodecDecodeSelf(d *Decoder) {
	d.MustDecode(&s)
}

// Make sure we're robust against reentrant calls.
func TestMsgpackDecodeSelfSelfer(t *testing.T) {
	var h MsgpackHandle
	d := NewDecoderBytes([]byte{0x00}, &h)

	var s selfer
	err := d.Decode(&s)
	assertMaxDepthError(t, err)
}

func TestMsgpackDecodeMaxDepthOption(t *testing.T) {
	// [[0x01]]
	b := []byte{0x91, 0x91, 0x01}

	var h MsgpackHandle
	d := NewDecoderBytes(b, &h)

	var v interface{}
	err := d.Decode(&v)
	if err != nil {
		t.Fatal(err)
	}

	h.MaxDepth = 1
	d = NewDecoderBytes(b, &h)

	err = d.Decode(&v)
	assertMaxDepthError(t, err)
}

func assertEOF(t *testing.T, err error) {
	if err != io.EOF && err != io.ErrUnexpectedEOF {
		t.Fatalf("expected EOF or ErrUnexpectedEOF, got %v", err)
	}
}

func testMsgpackDecodeMapSizeMismatch(t *testing.T, v interface{}) {
	// A map claiming to have 0x10eeeeee KV pairs, but only has 1.
	b := []byte{0xdf, 0x10, 0xee, 0xee, 0xee, 0x1, 0xa1, 0x1}

	var h MsgpackHandle
	d := NewDecoderBytes(b, &h)

	err := d.Decode(&v)
	assertEOF(t, err)
}

func TestMsgpackDecodeMapSizeMismatchFastPathNil(t *testing.T) {
	var m map[int]string
	testMsgpackDecodeMapSizeMismatch(t, &m)
}

func TestMsgpackDecodeMapSizeMismatchSlowPathNil(t *testing.T) {
	var m map[int][]byte
	testMsgpackDecodeMapSizeMismatch(t, &m)
}

func TestMsgpackDecodeSliceSizeMismatchFastPathNil(t *testing.T) {
	// An array claiming to have 0x10eeeeee elements, but only has 1.
	b := []byte{0xdd, 0x10, 0xee, 0xee, 0xee, 0x1}

	var h MsgpackHandle
	d := NewDecoderBytes(b, &h)

	var a []byte
	err := d.Decode(&a)
	assertEOF(t, err)
}

func TestMsgpackDecodeSliceSizeMismatchSlowPathNil(t *testing.T) {
	// An array claiming to have 0x10eeeeee elements, but only has 1.
	b := []byte{0xdd, 0x10, 0xee, 0xee, 0xee, 0x91, 0x1}

	var h MsgpackHandle
	d := NewDecoderBytes(b, &h)

	var a [][]byte
	err := d.Decode(&a)
	assertEOF(t, err)
}
