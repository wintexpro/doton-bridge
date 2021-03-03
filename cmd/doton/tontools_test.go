package main

import (
	"errors"
	"os"
	"testing"

	connection "github.com/ChainSafe/ChainBridge/connections/ton"
	log "github.com/ChainSafe/log15"
	"github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
	"github.com/wintexpro/chainbridge-utils/crypto/ed25519"
	"github.com/wintexpro/chainbridge-utils/keystore"
)

func TestDecodeSS58AddressToPublicKey(t *testing.T) {
	if decodeSS58AddressToPublicKey("5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY") != "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d" {
		t.Fatal(errors.New("Public key is not valid"))
	}
}

func TestSendGramsCmd(t *testing.T) {
	os.Setenv(keystore.EnvPassword, "123456")

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	logger := log.Root().New("test", "alice")

	ks := dir + "/../../keys"
	insecure := false

	kpI, err := keystore.KeypairFromAddress("0:2089148264fb4b40dbb9ed7ba7a862403a715abf50a5730637da33d4b6453dd2", keystore.TonChain, ks, insecure)
	if err != nil {
		t.Fatal(err)
	}
	kp, _ := kpI.(*ed25519.Keypair)
	conn := connection.NewConnection("http://localhost", false, "2", logger)

	workchainID := null.Int32From(int32(0))

	signer := client.Signer{
		Type: client.KeysSignerType,
		Keys: client.KeyPair{
			Public: kp.PublicKey(),
			Secret: kp.SecretKey(),
		},
	}

	sendGrams(conn, workchainID, &signer, logger)
}
