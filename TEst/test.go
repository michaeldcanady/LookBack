package main

import(
	"crypto/rand"
	"crypto/sha256"
	//"encoding/hex"
	"fmt"
	"os"
  "crypto/rsa"
)

var(
  test2048Key,_ = rsa.GenerateKey(rand.Reader, 4096)
  ciphertext []byte
)

func ExampleEncryptOAEP() {
	secretMessage := []byte("send reinforcements, we're going to advance")
	label := []byte("orders")

	// crypto/rand.Reader is a good source of entropy for randomizing the
	// encryption function.
	rng := rand.Reader

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &test2048Key.PublicKey, secretMessage, label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		return
	}

	// Since encryption is a randomized function, ciphertext will be
	// different each time.
	fmt.Printf("Ciphertext: %x\n", ciphertext)
}

func ExampleDecryptOAEP() {
	//ciphertext, _ := hex.DecodeString(ciphertext)
	label := []byte("orders")

	// crypto/rand.Reader is a good source of entropy for blinding the RSA
	// operation.
	rng := rand.Reader

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, test2048Key, ciphertext, label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from decryption: %s\n", err)
		return
	}

	fmt.Printf("Plaintext: %s\n", string(plaintext))

	// Remember that encryption only provides confidentiality. The
	// ciphertext should be signed before authenticity is assumed and, even
	// then, consider that messages might be reordered.
}

func main(){
  ExampleEncryptOAEP()
  ExampleDecryptOAEP()
}
