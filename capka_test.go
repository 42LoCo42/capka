package capka

import (
	"bytes"
	"encoding/base64"
	"reflect"
	"testing"
	"time"

	"github.com/go-faster/errors"
	"github.com/jamesruan/sodium"
)

func TestCAPKA(t *testing.T) {
	username := "alice"
	password := "hunter2"
	domain := "example.org"

	kp, err := MakeKP(username, password, domain)
	if err != nil {
		t.Fatal(errors.Wrap(err, "could not create keypair"))
	}

	pkWant := "2ajG+pbyrnXGTjQB8TEpdERkOHcCea9xGSj0tnw/ogM="
	pkHave := base64.StdEncoding.EncodeToString(kp.PublicKey.Bytes)
	if pkWant != pkHave {
		t.Fatalf("public key mismatch: want %s, have %s", pkWant, pkHave)
	}

	nonce, err := base64.StdEncoding.DecodeString(GetNonce(32, time.Second*5))
	if err != nil {
		t.Fatal(errors.Wrap(err, "could not decode nonce"))
	}

	dataWant, eph := NewLoginData(username, nonce)
	raw, err := dataWant.EncodeJSON(kp.SecretKey)
	if err != nil {
		t.Fatal(errors.Wrap(err, "could not encode login data"))
	}

	req, err := DecodeLoginRequestJSON(bytes.NewReader(raw))
	if err != nil {
		t.Fatal(errors.Wrap(err, "could not decode login request JSON"))
	}

	dataHave, err := req.Decode(kp.PublicKey)
	if err != nil {
		t.Fatal(errors.Wrap(err, "could not decode login request"))
	}
	if !reflect.DeepEqual(dataWant, dataHave) {
		t.Fatalf("original and decoded login data mismatch:\nwant: %+v\nhave: %+v", dataWant, dataHave)
	}

	secureDataWant := sodium.Bytes(RandomBytes(32))
	secureDataHave, err := Decrypt(dataHave.Encrypt(secureDataWant), eph)
	if err != nil {
		t.Fatal(errors.Wrap(err, "could not decrypt secure data"))
	}
	if !reflect.DeepEqual(secureDataWant, secureDataHave) {
		t.Fatalf(
			"original and decrypted secure data mismatch:\nwant: %+v\nhave: %+v",
			secureDataWant, secureDataHave)
	}
}
