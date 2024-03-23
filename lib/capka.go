package lib

// #include <sodium.h>
// #include "capka.h"
import "C"
import (
	"unsafe"

	"github.com/go-faster/errors"
	"github.com/jamesruan/sodium"
)

func MakeKP(username, password, domain string) (sodium.SignKP, error) {
	return MakeKP_Ex(
		[]byte(username),
		[]byte(password),
		[]byte(domain),
		C.crypto_pwhash_OPSLIMIT_INTERACTIVE,
		C.crypto_pwhash_MEMLIMIT_INTERACTIVE,
	)
}

func MakeKP_Ex(username, password, domain []byte, ops, mem C.int) (sodium.SignKP, error) {
	kp := sodium.SignKP{}

	saltData := append(username, domain...)

	pk := make([]byte, C.crypto_sign_PUBLICKEYBYTES)
	sk := make([]byte, C.crypto_sign_SECRETKEYBYTES)

	if C.capka_makeKeypair(
		(*C.char)(unsafe.Pointer(&password[0])),
		(*C.char)(unsafe.Pointer(&saltData[0])),
		ops,
		mem,
		(*C.uchar)(unsafe.Pointer(&pk[0])),
		(*C.uchar)(unsafe.Pointer(&sk[0])),
	) != 0 {
		return kp, errors.New("CAPKA: internal failure")
	}

	kp.PublicKey = sodium.SignPublicKey{Bytes: pk}
	kp.SecretKey = sodium.SignSecretKey{Bytes: sk}
	return kp, nil
}
