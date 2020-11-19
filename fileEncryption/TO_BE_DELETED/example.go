package encryption

import(
  
)

func test(){
  domain, err := GetDomain()
  if err != nil{
    fmt.Println(err)
  }
  fmt.Printf("'%s'\n",domain)
  valid, err := ValidateDomain(domain,"Personal")
  if err != nil{
    panic(err)
  }
  fmt.Println(valid)
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
