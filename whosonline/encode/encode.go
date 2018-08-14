package encode

import (
	"io/ioutil"
	"log"

	"github.com/jchavannes/go-pgp/pgp"
)

//Keys should contains all the info's about the keys
type Keys struct {
	//Pubkey is your Own public key
	Pubkey []byte
	//Privkey is your Own private key
	Privkey []byte
	//Othpubkey is the other person public key
	Othpubkey []byte
}

//Requete contiens la requete en json
type Requete struct {
	Cmd string
	Opt interface{}
}

//Options contiens les details (pas definitif à revoir)
type Options struct {
	Title string
	Text  string
	Image string
}

//Init initialize the struct
func Init(pubkey, privkey, othpubkey string) Keys {
	var keys Keys
	var err error
	keys.Pubkey, err = ioutil.ReadFile(pubkey)
	if err != nil {
		log.Println(err)
		return Keys{}
	}
	keys.Privkey, err = ioutil.ReadFile(privkey)
	if err != nil {
		log.Println(err)
		return Keys{}
	}
	keys.Othpubkey, err = ioutil.ReadFile(othpubkey)
	if err != nil {
		log.Println(err)
		return Keys{}
	}
	return keys
}

//Sign a message to authentify the author
func (keys Keys) Sign(msg []byte) []byte {
	entity, err := pgp.GetEntity(keys.Pubkey, keys.Privkey)
	if err != nil {
		log.Println(err)
		return nil
	}
	signature, err := pgp.Sign(entity, msg)
	if err != nil {
		log.Println(err)
		return nil
	}
	return signature
}

//Encrypt takes a message and encrypt it
func (keys Keys) Encrypt(msg []byte) []byte {
	pubEntity, err := pgp.GetEntity(keys.Othpubkey, []byte{})
	if err != nil {
		log.Println(err)
		return nil
	}
	encrypted, err := pgp.Encrypt(pubEntity, []byte(msg))
	if err != nil {
		log.Println(err)
		return nil
	}
	return encrypted
}

//Uncrypt unseal a message
func (keys Keys) Uncrypt(msg []byte) []byte {
	privEntity, err := pgp.GetEntity(keys.Pubkey, keys.Privkey)
	if err != nil {
		log.Println(err)
		return nil
	}
	decrypted, err := pgp.Decrypt(privEntity, msg)
	if err != nil {
		log.Println(err)
		return nil
	}
	return decrypted
}

//Verify that the message really come from the author
func (keys Keys) Verify(msg, signature []byte) bool {
	pubEntity, err := pgp.GetEntity(keys.Othpubkey, []byte{})
	if err != nil {
		log.Println(err)
		return false
	}
	err = pgp.Verify(pubEntity, []byte(msg), signature)
	if err != nil {
		return false
	}
	return true
}
