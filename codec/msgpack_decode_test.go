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

func TestMsgpackDecodeInfinitelyNestedArray(t *testing.T) {
	r := circularReader{b: []byte{0x91}}

	var h MsgpackHandle
	d := NewDecoder(&r, &h)

	var v interface{}
	err := d.Decode(&v)
	if err == nil || !strings.HasSuffix(err.Error(), "max depth exceeded") {
		t.Fatalf("Expected 'max depth exceeded', got %v", err)
	}
}
