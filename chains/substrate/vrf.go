// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/gtank/merlin"
)

type Keypair struct {
	MiniSecretKey schnorrkel.MiniSecretKey
}

func miniSecretKeyToVrfKeypair(sk *schnorrkel.MiniSecretKey) (*Keypair, error) {
	var pubKey [32]byte
	err := sk.Public().Decode(pubKey)
	if err != nil {
		return nil, err
	}

	var privKey [32]byte
	err = sk.ExpandEd25519().Decode(privKey)
	if err != nil {
		return nil, err
	}

	return &Keypair{
		MiniSecretKey: *sk,
	}, nil
}

func newSigningContext(context, msg []byte) *merlin.Transcript {
	t := merlin.NewTranscript("SigningContext")
	t.AppendMessage([]byte(""), context)
	t.AppendMessage([]byte("sign-bytes"), msg)
	return t
}

func VrfGenerateKeypair() (*Keypair, error) {
	msk, err := schnorrkel.GenerateMiniSecretKey()
	if err != nil {
		return nil, err
	}

	return miniSecretKeyToVrfKeypair(msk)
}

func (kp *Keypair) Public() []byte {
	p := kp.MiniSecretKey.Public().Encode()
	return p[:]
}

func (kp *Keypair) Sign(publicRandomness []byte) ([]byte, []byte, error) {
	signTranscript := newSigningContext([]byte(""), publicRandomness)
	inout, proof, err := kp.MiniSecretKey.ExpandEd25519().VrfSign(signTranscript)
	if err != nil {
		return nil, nil, err
	}

	inoutb := inout.Output().Encode()
	proofb := proof.Encode()

	return inoutb[:], proofb[:], nil
}

func (kp *Keypair) Verify(publicRandomness []byte, output *schnorrkel.VrfOutput, proof *schnorrkel.VrfProof) (bool, error) {
	signTranscript := newSigningContext([]byte(""), publicRandomness)
	result, err := kp.MiniSecretKey.Public().VrfVerify(signTranscript, output, proof)
	if err != nil {
		return false, err
	}

	return result, nil
}
