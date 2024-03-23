package capka

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/go-faster/errors"
	"github.com/jamesruan/sodium"
)

var Nonces = map[string]any{}

func RandomBytes(length int) []byte {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		log.Fatal(errors.Wrap(err, "generating random bytes should never fail"))
	}

	return buf
}

func RandomString(length int) string {
	return base64.StdEncoding.EncodeToString(RandomBytes(length))
}

func GetNonce(length int, age time.Duration) string {
	nonce := RandomString(length)
	Nonces[nonce] = nil

	go func() {
		time.Sleep(age)
		delete(Nonces, nonce)
	}()

	return nonce
}

func DecodeLoginRequestJSON(from io.Reader) (*LoginRequest, error) {
	req := &LoginRequest{}
	if err := json.NewDecoder(from).Decode(req); err != nil {
		return nil, errors.Wrap(err, "could not decode login request JSON")
	}

	return req, nil
}

func (req *LoginRequest) Decode(key sodium.SignPublicKey) (*LoginData, error) {
	nonce, err := base64.StdEncoding.DecodeString(req.Nonce)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode nonce")
	}

	ephKey, err := base64.StdEncoding.DecodeString(req.EphKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode ephkey")
	}

	signature, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signature")
	}

	data := &LoginData{
		User:   req.User,
		Nonce:  nonce,
		EphKey: ephKey,
	}

	if err := data.MakeSigInput().SignVerifyDetached(
		sodium.Signature{Bytes: signature},
		key,
	); err != nil {
		return nil, errors.Wrap(err, "signature verification failed")
	}

	return data, nil
}

func (req *LoginData) Encrypt(data sodium.Bytes) sodium.Bytes {
	return data.SealedBox(sodium.BoxPublicKey{Bytes: req.EphKey})
}
