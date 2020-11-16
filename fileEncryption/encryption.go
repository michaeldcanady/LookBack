package main

import(
  "fmt"
  "io/ioutil"
  "os"
  "crypto/rsa"
)

func StorePublic(publicKey *rsa.PublicKey,file string)error{
  f, err := os.Create(file)
    if err != nil {
        return err
    }
    pempub, err := ExportRsaPublicKeyAsPemStr(publicKey); if err != nil{
      return err
    }
    l, err := f.WriteString(pempub)
    if err != nil {
        f.Close()
        return err
    }
    fmt.Println(l, "bytes written successfully")
    err = f.Close()
    if err != nil {
        return err
    }
    return nil
}

func RetrievePublic(file string)(*rsa.PublicKey,error){
  var pempub *rsa.PublicKey
  data, err := ioutil.ReadFile(file)
    if err != nil {
        return pempub,err
    }
  pempub, err = ParseRsaPublicKeyFromPemStr(string(data)); if err != nil{
    return pempub,err
  }
  return pempub,nil
}

func StorePrivate(PrivateKey *rsa.PrivateKey,file string)error{
  f, err := os.Create(file)
    if err != nil {
        return err
    }
    pempub := ExportRsaPrivateKeyAsPemStr(PrivateKey)
    _, err = f.WriteString(pempub)
    if err != nil {
        f.Close()
        return err
    }
    err = f.Close()
    if err != nil {
        return err
    }
    return nil
}

func RetrievePrivate(file string)(*rsa.PrivateKey,error){
  var pempub *rsa.PrivateKey
  data, err := ioutil.ReadFile(file)
    if err != nil {
        return pempub,err
    }
  pempub, err = ParseRsaPrivateKeyFromPemStr(string(data)); if err != nil{
    return pempub,err
  }
  return pempub,nil
}

func main(){
  file := "C:\\go\\src\\github.com\\michaeldcanady\\Project01\\.ssh\\id_rsa.public"
  privateKey, publicKey := GenerateRsaKeyPair()
  err := StorePublic(publicKey,file); if err != nil{
    fmt.Println("Storing error:",err)
  }
  _, err = RetrievePublic(file); if err != nil{
    fmt.Println("Retrieval err:", err)
  }

  fileBytes, err := EncryptFile("C:\\Users\\micha\\OneDrive\\Desktop\\New folder (3)\\backup1_FILE.log",publicKey)
  if err != nil {
    fmt.Println("Encryption Error:",err)
  }

  f, err := os.Create("C:\\Users\\micha\\OneDrive\\Desktop\\New folder (3)\\backup1_FILE (1).log")
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(len(fileBytes))
  f.WriteString(string(fileBytes))
  f.Close()

  fileBytes, err = DecryptFile("C:\\Users\\micha\\OneDrive\\Desktop\\New folder (3)\\backup1_FILE (1).log",privateKey)
  if err != nil {
    fmt.Println("Decryption Error:",err)
  }

  f, err = os.Create("C:\\Users\\micha\\OneDrive\\Desktop\\New folder (3)\\backup1_FILE (2).log")
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(len(fileBytes))
  f.Write(fileBytes)
  f.Close()

}
