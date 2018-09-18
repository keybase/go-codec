// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
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
