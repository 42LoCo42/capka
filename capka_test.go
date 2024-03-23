package capka

import (
	"encoding/base64"
	"testing"
)

func TestCAPKA(t *testing.T) {
	username := "alice"
	password := "hunter2"
	domain := "example.org"

	kp, err := MakeKP(username, password, domain)
	if err != nil {
		t.Fatal(err)
	}

	pkWant := "2ajG+pbyrnXGTjQB8TEpdERkOHcCea9xGSj0tnw/ogM="
	pkHave := base64.StdEncoding.EncodeToString(kp.PublicKey.Bytes)
	if pkWant != pkHave {
		t.Fatalf("public key mismatch: want %s, have %s", pkWant, pkHave)
	}
}
