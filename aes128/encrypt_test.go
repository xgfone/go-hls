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

package aes128

import (
	"bytes"
	"testing"
)

func TestAddPKCS7Padding(t *testing.T) {
	data := []byte("0123456789")
	bs := addPKCS7Padding(data, 16)
	if !bytes.Equal(data, bs[:len(data)]) {
		t.Errorf("expect '%s', but got '%s'", string(data), string(bs[:len(data)]))
	} else if !bytes.Equal([]byte{6, 6, 6, 6, 6, 6}, bs[len(data):]) {
		t.Errorf("expect '%s', but got '%s'", string([]byte{6, 6, 6, 6, 6, 6}), string(bs[len(data):]))
	}
}

func TestEncrypt(t *testing.T) {
	var (
		data = []byte("0123456789")
		key  = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
		iv   = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	)

	encrypted, err := Encrypt(data, key, iv)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := Decrypt(encrypted, key, iv, true)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, decrypted) {
		t.Errorf("expect '%s', but got '%s'", string(data), decrypted)
	}
}
