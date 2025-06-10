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
)

// Option is used to configure the request.
type Option func(r *http.Request)

// ByteRange returns a request option to add the http request header "Range".
//
// start may be equal to 0, but length not.
func ByteRange(start, length uint64) Option {
	if length == 0 {
		panic("ByteRange: length must not be equal to 0")
	}

	end := start + length - 1
	return func(r *http.Request) {
		r.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	}
}

// Download is a convenient function to download something by HTTP.
func Download(ctx context.Context, url string, do func(http.Header, io.Reader) error, options ...Option) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	for _, o := range options {
		o(req)
	}

	resp, err := Get().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("statuscode=%d, err=%w", resp.StatusCode, err)
		}

		msg := unsafe.String(unsafe.SliceData(data), len(data))
		return fmt.Errorf("statuscode=%d, err=%s", resp.StatusCode, msg)
	}

	return do(resp.Header, resp.Body)
}

// ResolveURL tries to reslove the relative url based on baseurl
// if uri is relative, and returns it.
func ResolveURL(baseurl, uri string) (string, error) {
	if uri == "" {
		return "", errors.New("missing uri")
	}

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
