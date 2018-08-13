package encode

import (
	"fmt"
	"io/ioutil"

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
	Cmd     string
	Options interface{}
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
		return Keys{}
	}
	keys.Privkey, err = ioutil.ReadFile(privkey)
	if err != nil {
		return Keys{}
	}
	keys.Othpubkey, err = ioutil.ReadFile(othpubkey)
	if err != nil {
		return Keys{}
	}
	return keys
}

//Sign a message to authentify the author
func (keys Keys) Sign(msg []byte) []byte {
	entity, _ := pgp.GetEntity(keys.Pubkey, keys.Privkey)
	signature, _ := pgp.Sign(entity, msg)

	return signature
}

//Encrypt takes a message and encrypt it
func (keys Keys) Encrypt(msg string) []byte {
	pubEntity, err := pgp.GetEntity(keys.Othpubkey, []byte{})
	if err != nil {
		println(fmt.Errorf("Error getting entity: %v", err))
	}
	encrypted, err := pgp.Encrypt(pubEntity, []byte(msg))
	if err != nil {
		println(err)
	}

	return encrypted
}

//Uncrypt unseal a message
func (keys Keys) Uncrypt(msg string) string {
	privEntity, err := pgp.GetEntity(keys.Pubkey, keys.Privkey)
	if err != nil {
		println(fmt.Errorf("Error getting entity: %v", err))
	}
	decrypted, err := pgp.Decrypt(privEntity, []byte(msg))
	if err != nil {
		println(err)
	}
	decryptedMessage := string(decrypted)

	return decryptedMessage
}

//Verify that the message really come from the author
func (keys Keys) Verify(msg, signature []byte) bool {
	pubEntity, err := pgp.GetEntity(keys.Othpubkey, []byte{})
	if err != nil {
		println(fmt.Errorf("Error getting entity: %v", err))
	}
	err = pgp.Verify(pubEntity, []byte(msg), signature)
	if err != nil {
		return false
	}
	return true
}

/*
func main() {
	pubkey, err := ioutil.ReadFile("/home/florian/.gnupg/clefpub.asc")
	if err != nil {
		log.Fatal(err)
	}
	privkey, err := ioutil.ReadFile("/home/florian/.gnupg/clefpriv.asc")
	if err != nil {
		log.Fatal(err)
	}
	msg := "Salut à tous c benoit Magimel"

	msgchiffrer := encrypt(pubkey, msg)
	println("msg chiffrer : ", string(msgchiffrer))
	signature := sign(pubkey, privkey, []byte(msgchiffrer))
	println("Signature : ", string(signature))
	verify(pubkey, []byte(msgchiffrer), signature)
	msgclear := uncrypt(pubkey, privkey, string(msgchiffrer))
	println("message clear : ", msgclear)
}
*/
