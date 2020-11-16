package main

import(
  "fmt"
  "os"
  "math/rand"
  "time"
  "encoding/hex"
  "io/ioutil"
  "encoding/pem"
  "os/exec"
)

func StorePublic(key []byte, file string)error{
  f, err := os.Create(file)
  if err != nil {
    return err
  }

  privkey_pem := pem.EncodeToMemory(
          &pem.Block{
                  Type:  "RSA PRIVATE KEY",
                  Bytes: key,
          },
  )

  f.WriteString(string(privkey_pem))

  return nil
}

func RetrievePublic(file string)(string,error){
  data, err := ioutil.ReadFile(file)
  if err != nil {
    return "",err
  }
  key := hex.EncodeToString(data)
  return key,err
}

func GenerateKey()[]byte{
  key := make([]byte, 32)
  rand.Seed(time.Now().UnixNano())
  rand.Read(key)
  return key
}

// Currently only works for windows
func GetDomain()(string,error){
  commandString := "(Get-WmiObject Win32_ComputerSystem).Domain"
  output, err := exec.Command("Powershell", "-Command", commandString ).CombinedOutput()
  if err != nil{
    return "",err
  }
  return string(output),nil
}

func main(){
  domain, err := GetDomain()
  if err != nil{
    fmt.Println(err)
  }
  fmt.Println(domain)
  panic("TEST")
  file := "C:\\go\\src\\github.com\\michaeldcanady\\Project01\\.ssh\\id_rsa.public"
  publicKey := GenerateKey()
  err = StorePublic(publicKey,file); if err != nil{
    fmt.Println("Storing error:",err)
  }
  publicKey1, err := RetrievePublic(file); if err != nil{
    fmt.Println("Retrieval err:", err)
  }

  // how to encrypt a file returns encrypted bytes to be written
  fileBytes := Encrypt("C:\\Users\\dmcanady\\Desktop\\New folder (3)\\CS6545678_FILE.log",publicKey1)

  f, err := os.Create("C:\\Users\\dmcanady\\Desktop\\New folder (3)\\CS6545678_FILE (1).log")
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(len(fileBytes))
  f.WriteString(string(fileBytes))
  f.Close()

  fileBytes = Decrypt("C:\\Users\\dmcanady\\Desktop\\New folder (3)\\CS6545678_FILE (1).log",publicKey1)

  f, err = os.Create("C:\\Users\\dmcanady\\Desktop\\New folder (3)\\CS6545678_FILE (2).log")
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(len(fileBytes))
  f.Write(fileBytes)
  f.Close()

}
