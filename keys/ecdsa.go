// Copyright 2017 Stratumn SAS. All rights reserved.
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

package keys

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"

	"github.com/pkg/errors"
	"github.com/stratumn/go-crypto/encoding"
)

const (
	// ECDSASecretPEMLabel is the label of a PEM-encoded ECDSA secret key.
	ECDSASecretPEMLabel = "EC PRIVATE KEY"

	// ECDSAPublicPEMLabel is the label of a PEM-encoded ECDSA public key.
	ECDSAPublicPEMLabel = "EC PUBLIC KEY"
)

// NewECDSAKeyPair generates a new ECDSA key pair using the P-256 curve.
func NewECDSAKeyPair() (crypto.PublicKey, *ecdsa.PrivateKey, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return priv.Public(), priv, nil
}

// EncodeECDSASecretKey encodes an ECDSA secret key in ASN.1 DER format within a PEM block
// embedded in PKCS#8.
func EncodeECDSASecretKey(sk *ecdsa.PrivateKey) ([]byte, error) {
	skBytes, err := x509.MarshalPKCS8PrivateKey(sk)
	if err != nil {
		return nil, err
	}
	return encoding.EncodePEM(skBytes, ECDSASecretPEMLabel)
}

// ParseECDSAPKCS8Key decodes a PEM block containing an ASN1. DER encoded
// secret key of type ECDSA embedded in PKCS#8.
func ParseECDSAPKCS8Key(sk []byte) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	DERBytes, err := encoding.DecodePEM(sk, ECDSASecretPEMLabel)
	if err != nil {
		return nil, nil, err
	}

	data, err := x509.ParsePKCS8PrivateKey(DERBytes)
	if err != nil {
		return nil, nil, err
	}

	key, ok := data.(*ecdsa.PrivateKey)
	if !ok {
		return nil, nil, errors.New("failed to parse ECDSA private key embedded in PKCS#8")
	}

	return key, key.Public().(*ecdsa.PublicKey), nil
}
