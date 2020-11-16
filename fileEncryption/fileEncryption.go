package main

import(
  "io/ioutil"
  "crypto/rsa"
  "crypto/sha256"
  "crypto/rand"
  "fmt"
)

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

func EncryptFile(file string, PublicKey *rsa.PublicKey)([]byte, error){
  var encryptedBytes []byte

  // read the file into bytes
  data, err := ioutil.ReadFile(file)
  if err != nil {
    return encryptedBytes,err
  }
  fmt.Println(len(data))
  // Encrypts the file
  //fmt.Println(PublicKey.N.BitLen())
  //_,_ = strconv.Atoi((PublicKey.N).String())
  //ByteSlice := split(data,200)
  //var encryptedByte []byte
  rng := rand.Reader
  //for _,bytes := range ByteSlice{
    encryptedBytes, err = rsa.EncryptOAEP(
	     sha256.New(),
	     rng,
	     PublicKey,
	     data,
	     nil)
    if err != nil {
      return encryptedBytes,err
    }
    //encryptedBytes = append(encryptedBytes, encryptedByte...)
  //}

  // Returns file encrypted
  return encryptedBytes,nil
}

func DecryptFile(file string, privateKey *rsa.PrivateKey)([]byte, error){
  var decryptedBytes []byte
  // read the file into bytes
  data, err := ioutil.ReadFile(file)
  if err != nil {
    return decryptedBytes,err
  }
  fmt.Println(len(data))
  var decryptedByte []byte
  ByteSlice := split(data,200)
  rng := rand.Reader
  for _,bytes := range ByteSlice{
    decryptedBytes,err = rsa.DecryptOAEP(
      sha256.New(),
      rng,
      privateKey,
      bytes,
      nil)
    if err != nil {
  	 return decryptedBytes,err
    }
    decryptedBytes = append(decryptedBytes,decryptedByte...)
  }

  return decryptedBytes,nil
}
