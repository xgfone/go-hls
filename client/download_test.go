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
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/xgfone/go-toolkit/httpx"
	"github.com/xgfone/go-toolkit/httpx/option"
)

func TestDownload(t *testing.T) {
	do := func(r *http.Request) (*http.Response, error) {
		ranger := r.Header.Get("Range")
		if ranger == "" {
			return nil, errors.New("missing the Range header")
		}

		ranger = strings.TrimPrefix(ranger, "bytes=")
		index := strings.IndexByte(ranger, '-')
		if index < 0 {
			return nil, errors.New("invalid the Range header")
		}

		start, err := strconv.ParseUint(ranger[:index], 10, 64)
		if err != nil {
			return nil, err
		}

		end, err := strconv.ParseUint(ranger[index+1:], 10, 64)
		if err != nil {
			return nil, err
		}

		header := make(http.Header, 2)
		header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/123456789", start, end))
		return &http.Response{
			StatusCode: 206,
			Header:     header,
			Body:       io.NopCloser(strings.NewReader(ranger)),
		}, nil
	}
	httpx.SetClient(httpx.DoFunc(do))

	const (
		baseurl = "http://localhost/dir"
		uri     = "file.txt"
	)

	url, err := ResolveURL(baseurl, uri)
	if err != nil {
		t.Fatal(err)
	}

	var body string
	bodydo := func(r *http.Response) (err error) {
		data, err := io.ReadAll(r.Body)
		if err == nil {
			body = string(data)
		}
		return
	}

	const expectbody = "0-99"
	err = Get(context.Background(), url, bodydo, option.ByteRange(0, 100))
	if err != nil {
		t.Fatal(err)
	} else if body != expectbody {
		t.Errorf("expect response body '%s', but got '%s'", expectbody, body)
	}
}
