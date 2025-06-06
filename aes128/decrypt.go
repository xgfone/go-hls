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
	"errors"
)

var (
	errInvalidEncryptedIV  = errors.New("invalid media segment encryption iv")
	errInvalidEncryptedKey = errors.New("invalid media segment encryption key")
)

func Decrypt(encryptedData, key, iv []byte, removePadding bool) (decryptedData []byte, err error) {
	switch {
	case len(key) != 16:
		return nil, errInvalidEncryptedKey

	case len(iv) != 16:
		return nil, errInvalidEncryptedIV
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	} else if len(encryptedData)%block.BlockSize() != 0 {
		return nil, errors.New("invalid encrypted data")
	}

	decryptedData = make([]byte, len(encryptedData))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(decryptedData, encryptedData)

	if removePadding {
		decryptedData, err = removePKCS7Padding(decryptedData)
	}

	return
}

// AutoDetectAndRemovePadding automatically detects the PKCS7 padding and removes it.
func AutoDetectAndRemovePadding(decryptedData []byte) []byte {
	decryptedData, _ = removePKCS7Padding(decryptedData)
	return decryptedData
}

func removePKCS7Padding(decryptedData []byte) ([]byte, error) {
	datalen := len(decryptedData)
	if datalen == 0 {
		return decryptedData, nil
	} else if datalen%aes.BlockSize != 0 {
		return decryptedData, errors.New("data length not a multiple of block size")
	}

	// 1. Get the padding length
	paddingLen := int(decryptedData[datalen-1])
	if paddingLen < 1 || paddingLen > aes.BlockSize {
		return decryptedData, errors.New("invalid padding length")
	}

	// 2. Check whether the padding bytes are consistent.
	for i := 0; i < paddingLen; i++ {
		if decryptedData[datalen-1-i] != byte(paddingLen) {
			return decryptedData, errors.New("invalid padding bytes")
		}
	}

	// 3. Remove the padding
	decryptedData = decryptedData[:datalen-paddingLen]
	return decryptedData, nil
}
