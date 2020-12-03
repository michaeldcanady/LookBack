package encryption

import(
  "fmt"
  "os/exec"
  "errors"
  "strings"
)

// Currently only works for windows
func GetDomain()(string,error){
  commandString := "(Get-WmiObject Win32_ComputerSystem).Domain"
  output, err := exec.Command("Powershell", "-Command", commandString ).CombinedOutput()
  if err != nil{
    return "",err
  }
  return strings.TrimSpace(string(output)),nil
}

func ValidateDomain(domain string,dataType string)(bool,error){
  if domain == "Placeholder for domain" && dataType == "Work"{
    return true,nil
  }else if domain != "Placeholder for domain" && dataType == "Work"{
    return false,errors.New("Work data cannot be decrypted on non-work device")
  }else if domain == "Placeholder for domain" && dataType != "Work"{
    return false,errors.New("Personal data cannot be decrypted on work device")
  }else if domain != "Placeholder for domain" && dataType == "Personal"{
    return true,nil
  }
  panic(fmt.Sprintf("Unforseen exception met: Domain: %s, DataType: %s",domain,dataType))
}
