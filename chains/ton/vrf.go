// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"crypto/ed25519"
	"encoding/hex"

	ksEd25519 "github.com/wintexpro/chainbridge-utils/crypto/ed25519"
)

type Keypair struct {
	Public []byte
	Secret []byte
}

func ed25519KeypairToVrfKeypair(kp *ksEd25519.Keypair) (*Keypair, error) {
	decodedSk, err := hex.DecodeString(kp.SecretKey())
	if err != nil {
		return nil, err
	}

	privKey := ed25519.NewKeyFromSeed(decodedSk)
	pubKey := privKey.Public()

	return &Keypair{
		Public: pubKey.(ed25519.PublicKey),
		Secret: privKey,
	}, nil
}

func VrfGenerateKeypairFromSeed(seed string) (*Keypair, error) {
	kp, err := ksEd25519.NewKeypairFromSeed(seed)
	if err != nil {
		return nil, err
	}

	return ed25519KeypairToVrfKeypair(kp)
}

func VrfGenerateKeypair() (*Keypair, error) {
	kp, err := ksEd25519.GenerateKeypair()
	if err != nil {
		return nil, err
	}

	return ed25519KeypairToVrfKeypair(kp)
}

func (kp *Keypair) Sign(publicRandomness []byte) ([]byte, error) {
	return ed25519.Sign(kp.Secret, publicRandomness), nil
}

func (kp *Keypair) Verify(publicRandomness []byte, signature []byte) bool {
	return ed25519.Verify(kp.Public, publicRandomness, signature)
}
