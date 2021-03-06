// Copyright 2019 Stratumn SAS. All rights reserved.
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

package encryption

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"

	"github.com/pkg/errors"
	"github.com/stratumn/go-crypto/aes"
	"github.com/stratumn/go-crypto/keys"
)

// Decrypt decrypts a long message with the private key.
// Only RSA keys are supported for now.
// The ciphertext is actually composed of:
// - a RSA-OAEP encrypted AES key
// - an AES-256-GCM encrypted message
// Returns the bytes of the plaintext.
func Decrypt(secretKey, data []byte) ([]byte, error) {
	sk, _, err := keys.ParseSecretKey(secretKey)
	if err != nil {
		return nil, err
	}

	var opts crypto.DecrypterOpts
	switch sk.(type) {
	case *rsa.PrivateKey:
		opts = &rsa.OAEPOptions{Hash: crypto.SHA1}
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrNotImplemented
	}

	if len(data) < aes.KeyLength*8 {
		return nil, ErrCouldNotDecrypt
	}

	encryptedSymKey := data[:aes.KeyLength*8]
	cipherText := data[aes.KeyLength*8:]

	decrypter, ok := sk.(crypto.Decrypter)
	if !ok {
		return nil, errors.New("private key does not implement crypto.Decrypter")
	}

	aesKey, err := decrypter.Decrypt(rand.Reader, encryptedSymKey, opts)
	if err != nil {
		return nil, ErrCouldNotDecrypt
	}

	return aes.Decrypt(cipherText, aesKey)
}

// DecryptShort decrypt a short message.
// for 2048-bit RSA keys, the max message size is 214 bytes.
// Only RSA keys are supported for now.
// The message is directly RSA-OAEP decrypted.
// Returns the bytes of the plaintext.
func DecryptShort(secretKey, data []byte) ([]byte, error) {
	sk, _, err := keys.ParseSecretKey(secretKey)
	if err != nil {
		return nil, err
	}

	var opts crypto.DecrypterOpts
	switch sk.(type) {
	case *rsa.PrivateKey:
		opts = &rsa.OAEPOptions{Hash: crypto.SHA1}
	default:
		return nil, ErrNotImplemented
	}

	decrypter, ok := sk.(crypto.Decrypter)
	if !ok {
		return nil, errors.New("private key does not implement crypto.Decrypter")
	}

	return decrypter.Decrypt(rand.Reader, data, opts)
}
