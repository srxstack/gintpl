// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package httpx

import (
	"context"
	"testing"
)

const baidu = "https://www.baidu.com"

func TestNewHttpX(t *testing.T) {
	h := NewHTTPX()
	t.Log(h.POST(context.Background(), baidu, nil))
}
