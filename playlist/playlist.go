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

// Package playlist is used to encode or decode the HLS playlist
// based on the M3U8 format.
package playlist

import "io"

// Define the playlist types.
const (
	PlayListTypeMaster = "Master"
	PlayListTypeMedia  = "Media"
)

// PlayList is a playlist interface.
type PlayList interface {
	Type() string // Master or Media
	MinVersion() uint64
}

// Parse reads the data from r and decodes it as a master or media playlist.
func Parse(r io.Reader, options ...Option) (PlayList, error) {
	var p _Parser
	p.configure(options...)
	return p.Parse(r)
}
