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

package aes

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	msg           = []byte("coucou, tu veux voir mon message ?")
	ciphertext, _ = base64.StdEncoding.DecodeString("/D/5lTo/tWXZEbIAP43VfxlxYtVT+A8DVzrUkg78YYgkytXgAemz2vwsoh9YncVx+gIU8WM1NwnK9FthnvY=")
	key           = []byte("NV5oa0/G3hNbdqfBdDn2LKl6sf0JEOi4gj/EWXZdO5A=")
)

func TestEncrypt(t *testing.T) {
	ct, k, err := Encrypt(msg)

	require.NoError(t, err)

	pt, err := Decrypt(ct, k)
	require.NoError(t, err)

	assert.Equal(t, pt, msg)
}

func TestDecrypt(t *testing.T) {
	pt, err := Decrypt(ciphertext, key)
	require.NoError(t, err)
	assert.Equal(t, msg, pt)
}

func TestDecrypt_Fail(t *testing.T) {
	k := make([]byte, KeyLength)
	rand.Read(key)
	pt, err := Decrypt(ciphertext, k)
	assert.EqualError(t, err, ErrCouldNotDecrypt.Error())
	assert.Nil(t, pt)
}
