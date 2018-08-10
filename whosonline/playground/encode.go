package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/openpgp"
)

func main() {
	// Read armored public key into type EntityList
	// An EntityList contains one or more Entities.
	// This assumes there is only one Entity involved
	entitylist, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(publicKey))
	if err != nil {
		log.Fatal(err)
	}

	// Encrypt message using public key
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entitylist, nil, nil, nil)
	if err != nil {
	}

	const message = "Test message encrypted using public key"
	_, err = w.Write([]byte(message))
	if err != nil {
	}
	err = w.Close()
	if err != nil {
	}

	// Output as base64 encoded string
	bytes, err := ioutil.ReadAll(buf)
	str := base64.StdEncoding.EncodeToString(bytes)

	fmt.Println("Public key encrypted message (base64 encoded):", str)

}

// pub   1024R/7F98BBCE 2014-01-04
// uid                  Golang Test (Private key password is 'golang') <golangtest@test.com>
// sub   1024R/5F34A320 2014-01-04
const publicKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v1.4.15 (Darwin)

mI0EUsdthgEEAMKOAoeY+bAHEjdzcM9WhJ27T4QmX8SxLYcRo3rd2cuQawwCz7jf
bCzLCYyvMqoIvjSxuElVgFx97RyEv5yvLg7ngNfv6ADRlJXMVQ3YQahyeRofPJJ+
5S0F0JOahZlkAYWIHCUhLtoT/zpI7IeSwWjwtEL1b8YhZBLY9txp29TLABEBAAG0
REdvbGFuZyBUZXN0IChQcml2YXRlIGtleSBwYXNzd29yZCBpcyAnZ29sYW5nJykg
PGdvbGFuZ3Rlc3RAdGVzdC5jb20+iLgEEwECACIFAlLHbYYCGwMGCwkIBwMCBhUI
AgkKCwQWAgMBAh4BAheAAAoJEFVKIId/mLvOookD+wVQzZN8vZVkpYLsTU3XDBly
0H0F/vtJ4A9JWkYJnRyJRggV3DAajAq2OgOuxtiA+n5QY7JgwPq0bNYpomtBCgPJ
pCpVVGFs1cHsnPslPZqoocPW3tzHkV9TMMwE2i7dM5YeiYNfJAYMBQsBmeNo6Pz+
kN7qmjHIGW5KMwlTN8OmuI0EUsdthgEEAK1DA6pBp4PQqaZO91AVgXe44YW7ZNHm
kUIf4KFB4SiXq2eCzENtSCsiF/hkG7HA6XHKVzCOnk4V8ay/g/BuHDW+HsL09M3N
tPk/dc7YE/QP+FYn3BD0AhK06mP6GaYQM2TNaerEXp3NtnuNok9CIm3eYArNsJ0j
XlM8mw3LkIthABEBAAGInwQYAQIACQUCUsdthgIbDAAKCRBVSiCHf5i7zrpRA/9r
lIf6ozk+OvF6Cul7fN+8OOSUD6S6ohh/SiYKha1MSTMNWyBNhutOjmOoQoHhPmAv
Kp8tvYULV4SiKrlCP9ANait2gmYcKsqk/kI7xel4tIvx64EMAsgaKWN7hp3TG77Y
cVNCjtHerHjGZbRw6/GGlNSbw8DRQ0FbsPkasuexEw==
=jr2t
-----END PGP PUBLIC KEY BLOCK-----`
