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
	"net/url"
	"strings"
	"unsafe"

	"github.com/xgfone/go-toolkit/httpx"
)

func request(ctx context.Context, method, url string, body io.Reader,
	do func(*http.Response) error, options ...Option) (err error) {

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return
	}

	for i := range options {
		req = options[i](req)
	}

	resp, err := httpx.GetClient().Do(req)
	if err != nil {
		return
	}

	return do(resp)
}

// Get is a convenient function to download something by HTTP.
func Get(ctx context.Context, url string, do func(*http.Response) error, options ...Option) error {
	return request(ctx, http.MethodGet, url, nil, func(r *http.Response) (err error) {
		if r.StatusCode < 200 || r.StatusCode >= 300 {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				return fmt.Errorf("statuscode=%d, err=%w", r.StatusCode, err)
			}

			msg := unsafe.String(unsafe.SliceData(data), len(data))
			return fmt.Errorf("statuscode=%d, err=%s", r.StatusCode, msg)
		}
		return do(r)
	}, options...)
}

// ResolveURL tries to reslove the relative url based on baseurl
// if uri is relative, and returns it.
func ResolveURL(baseurl, uri string) (string, error) {
	switch {
	case uri == "":
		return "", errors.New("missing uri")

	case strings.HasPrefix(uri, "http://"),
		strings.HasPrefix(uri, "https://"):
		return uri, nil

	case baseurl == "":
		return "", errors.New("missing base url")
	}

	bu, err := url.Parse(baseurl)
	if err != nil {
		return "", fmt.Errorf("invalid baseurl: %w", err)
	}

	ru, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("invalid baseurl: %w", err)
	}

	uri = bu.ResolveReference(ru).String()
	return uri, nil
}
