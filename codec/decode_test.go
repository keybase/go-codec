// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bytes"
	"testing"
)

const maxInt = int(^uint(0) >> 1)

func TestBufioDecReaderReadx(t *testing.T) {
	var r bufioDecReader
	r.buf = make([]byte, 0, 10)
	r.reset(bytes.NewReader(nil))
	r.readx(maxInt >> 24)
}
