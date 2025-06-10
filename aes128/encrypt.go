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
	"crypto/aes"
	"crypto/cipher"
)

// Encrypt encrypts the data with key and iv by AES-128.
func Encrypt(data, key, iv []byte) (encrypted []byte, err error) {
	switch {
	case len(key) != 16:
		return nil, errInvalidEncryptedKey

	case len(iv) != 16:
		return nil, errInvalidEncryptedIV
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	paddata := addPKCS7Padding(data, block.BlockSize())
	encrypted = make([]byte, len(paddata))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(encrypted, paddata)
	return
}

func addPKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	paddata := make([]byte, len(data)+padding)
	copy(paddata, data)
	for i, _len := len(data), len(data)+padding; i < _len; i++ {
		paddata[i] = byte(padding)
	}
	return paddata
}
