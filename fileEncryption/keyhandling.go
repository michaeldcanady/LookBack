package encryption

import(
  "os"
  "io/ioutil"
  "crypto/rand"
  mrand"math/rand"
  "io"
  "encoding/hex"
  "crypto/md5"
  "crypto/aes"
  "crypto/cipher"
  "time"
  "encoding/pem"
)


func StorePublic(key []byte, file string)error{
  f, err := os.Create(file)
  if err != nil {
    return err
  }
  privkey_pem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY",Bytes: key,},)
  f.WriteString(string(privkey_pem))
  return nil
}

func RetrievePublic(file string)(string,error){
  data, err := ioutil.ReadFile(file)
  if err != nil {
    return "",err
  }
  return hex.EncodeToString(data),err
}

func GenerateKey()[]byte{
  key := make([]byte, 32)
  mrand.Seed(time.Now().UnixNano())
  mrand.Read(key)
  return key
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(file string, passphrase string) []byte {
  var ciphertext []byte
  data, err := ioutil.ReadFile(file)
  if err != nil {
    return ciphertext
  }
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext = gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func Decrypt(file string, passphrase string) []byte {
  var plaintext []byte
  data, err := ioutil.ReadFile(file)
  if err != nil {
    return plaintext
  }
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}
