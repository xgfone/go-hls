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

package playlist

import "fmt"

// ParseError represents an error when parsing a master/media playlist.
type ParseError struct {
	Line int    // Line Number
	Data string // Line Data
	Err  error
}

// Error implements the error interface.
func (e ParseError) Error() string {
	return fmt.Sprintf("line %d: %s: %v", e.Line, e.Data, e.Err)
}

// Unwrap is used to unwrap the inner error.
func (e ParseError) Unwrap() error {
	return e.Err
}
