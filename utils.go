package capka

import (
	"bytes"

	"github.com/jamesruan/sodium"
)

type GetKey func(user string) (key []byte, err error)

type LoginData struct {
	User   string
	Nonce  []byte
	EphKey []byte
}

type LoginRequest struct {
	User      string `json:"user"`
	Nonce     string `json:"nonce"`
	EphKey    string `json:"ephkey"`
	Signature string `json:"signature"`
}

func NewLoginData(user string, nonce []byte) (*LoginData, sodium.BoxKP) {
	eph := sodium.MakeBoxKP()
	return &LoginData{
		User:   user,
		Nonce:  nonce,
		EphKey: eph.PublicKey.Bytes,
	}, eph
}

func (data *LoginData) MakeSigInput() sodium.Bytes {
	return bytes.Join([][]byte{
		[]byte(data.User),
		data.Nonce,
		data.EphKey,
	}, nil)
}
