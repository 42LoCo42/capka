package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/42LoCo42/capka"
	"github.com/go-faster/errors"
	"golang.org/x/term"
)

func main() {
	username := flag.String("u", "", "username")
	passwordFile := flag.String("p", "", "file with password, - for stdin")
	domain := flag.String("d", "", "domain")

	flag.Parse()

	if *username == "" {
		log.Fatal(errors.New("username must be set! (-u)"))
	}

	if *domain == "" {
		log.Fatal(errors.New("domain must be set! (-d)"))
	}

	var password []byte
	var err error

	if *passwordFile == "" {
		fmt.Fprint(os.Stderr, "Password: ")
		password, err = term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)
	} else {
		if *passwordFile == "-" {
			*passwordFile = "/dev/stdin"
		}

		password, err = os.ReadFile(*passwordFile)
		password = bytes.TrimSpace(password)
	}
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not read password"))
	}

	kp, err := capka.MakeKP(*username, string(password), *domain)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not create keypair"))
	}

	fmt.Println(base64.StdEncoding.EncodeToString(kp.PublicKey.Bytes))
}
