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
	"net/http"
	"sync/atomic"
)

// Client represents a http client.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// DoFunc is a function to send a http request and get the http response.
type DoFunc func(*http.Request) (*http.Response, error)

// Do implements the Client interface.
func (f DoFunc) Do(r *http.Request) (*http.Response, error) {
	return f(r)
}

type _Client struct {
	Client
}

var _client atomic.Value

func init() {
	Set(http.DefaultClient)
}

// Get returns the http client.
func Get() Client {
	return _client.Load().(_Client).Client
}

// Set resets the http client.
//
// Default: http.DefaultClient
func Set(c Client) {
	if c == nil {
		panic("client.Set: client must not be nil")
	}
	_client.Store(_Client{Client: c})
}
