// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package signature

import (
	"crypto/ecdsa"
	"crypto/rand"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type KeyBundle struct {
	onChainSigning *onChainPrivateKey
}

var curve = secp256k1.S256()

func NewKeyBundle() (*KeyBundle, error) {
	ecdsaKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	return &KeyBundle{
		onChainSigning: (*onChainPrivateKey)(ecdsaKey),
	}, nil
}

func (kb *KeyBundle) SignOnChain(msg []byte) ([]byte, error) {
	return kb.onChainSigning.Sign(msg)
}

// PublicKeyAddressOnChain returns public component of the keypair used in
// SignOnChain
func (kb *KeyBundle) PublicKeyAddressOnChain() OnChainSigningAddress {
	return kb.onChainSigning.Address()
}

type onChainPrivateKey ecdsa.PrivateKey

// Sign returns the signature on msgHash with k
func (k *onChainPrivateKey) Sign(msg []byte) (signature []byte, err error) {
	sig, err := crypto.Sign(onChainHash(msg), (*ecdsa.PrivateKey)(k))
	return sig, err
}

func (k onChainPrivateKey) Address() OnChainSigningAddress {
	return OnChainSigningAddress(crypto.PubkeyToAddress(k.PublicKey))
}
