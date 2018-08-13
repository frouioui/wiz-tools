package encode

import (
	"fmt"

	"github.com/jchavannes/go-pgp/pgp"
)

func sign(pubkey, privkey, msg []byte) []byte {
	entity, _ := pgp.GetEntity(pubkey, privkey)
	signature, _ := pgp.Sign(entity, msg)

	return signature
}

func encrypt(pubkey []byte, msg string) []byte {
	pubEntity, err := pgp.GetEntity([]byte(pubkey), []byte{})
	if err != nil {
		println(fmt.Errorf("Error getting entity: %v", err))
	}
	encrypted, err := pgp.Encrypt(pubEntity, []byte(msg))
	if err != nil {
		println(err)
	}

	return encrypted
}

func uncrypt(pubkey, privkey []byte, msg string) string {
	privEntity, err := pgp.GetEntity([]byte(pubkey), []byte(privkey))
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

func verify(pubkey, msg, signature []byte) bool {
	pubEntity, err := pgp.GetEntity(pubkey, []byte{})
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
