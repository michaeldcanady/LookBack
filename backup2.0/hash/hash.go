package hash

import(
  "io"
  "log"
  "crypto/sha256"
  "os"
  "fmt"
)

// Hashes file for verification of a successful transfer
func HashFile(fileName string,UNIT int64)(string,error){
  file, err := os.Open(fileName)
	if err != nil {
		return "",err
	}
	defer file.Close()

	buf := make([]byte, 30*UNIT)
	sha256 := sha256.New()
	for {
		n, err := file.Read(buf)
		if n > 0 {
			_, err := sha256.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			//log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}

	sum := sha256.Sum(nil)
	return fmt.Sprintf("%x",sum),nil
}

//Compares two hash to see if they are the same
func CompareHash(hash1,hash2 string)bool{
  return hash1 == hash2
}
