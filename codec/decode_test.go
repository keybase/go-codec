// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bytes"
	"io"
	"testing"
)

func doReadx(r decReader, n int) (i interface{}) {
	defer func() {
		if x := recover(); x != nil {
			i = x
			return
		}
	}()

	r.readx(n)
	return nil
}

func testBufioDecReaderReadx(t *testing.T, n int) {
	var r bufioDecReader
	r.buf = make([]byte, 0, 10)
	r.reset(bytes.NewReader(nil))
	i := doReadx(&r, n)
	if i != io.EOF {
		t.Fatalf("(n=%d) expected EOF, got %v", n, i)
	}
}

func TestBufioDecReaderReadx(t *testing.T) {
	for n := 1; n < (1 << 10); n <<= 1 {
		t.Logf("n=%d", n)
		testBufioDecReaderReadx(t, n)
	}
}
