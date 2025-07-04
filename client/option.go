// Copyright 2025 xgfone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"
	"net/http"
	"strconv"
)

// Option is a function that modifies the HTTP request.
type Option func(*http.Request) *http.Request

func noop(r *http.Request) *http.Request { return r }

// ByteRange returns an Option that sets the "Range" header for the HTTP request.
func ByteRange(offset, length uint64) Option {
	if offset == 0 && length == 0 {
		return noop
	}

	var end string
	if length > 0 {
		end = strconv.FormatUint(offset+length-1, 10)
	}

	return func(r *http.Request) *http.Request {
		r.Header.Set("Range", fmt.Sprintf("bytes=%d-%s", offset, end))
		return r
	}
}
