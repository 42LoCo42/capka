package capka

import (
	"encoding/base64"
	"encoding/json"
	"log"
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

	eph := sodium.MakeBoxKP()

	dataWant := &LoginData{
		User:   username,
		Nonce:  nonce,
		EphKey: eph.PublicKey.Bytes,
	}

	req := dataWant.Encode(kp.SecretKey)
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatal(errors.Wrap(err, "could not encode request as JSON"))
	}
	log.Print(string(raw))

	dataHave, err := req.Decode(kp.PublicKey)
	if err != nil {
		t.Fatal(errors.Wrap(err, "could not decode login request"))
	}
	if !reflect.DeepEqual(dataWant, dataHave) {
		t.Fatalf("original and decoded login data mismatch:\nwant: %+v\nhave: %+v", dataWant, dataHave)
	}
}
