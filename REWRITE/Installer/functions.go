package main

import(
  "github.com/emersion/go-autostart"
  "fmt"
  "os"
  "path/filepath"
)

// exists returns whether the given file or directory exists
func exists(path string) (bool) {
    _, err := os.Stat(path)
    if err == nil { return true}
    if os.IsNotExist(err) { return false}
    return false
}

func OnStart(name, displayName string, kwargs ...string) {
	app := &autostart.App{
		Name: name,
		DisplayName: displayName,
		Exec: kwargs,
	}
  var action string
	if app.IsEnabled() {
		fmt.Printf("[-] Removing %s from startup...\n",name)
    action = "removed"

		if err := app.Disable(); err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("[+] Adding %s to startup...\n",name)
    action = "added"

		if err := app.Enable(); err != nil {
			panic(err)
		}
	}
  var sign string
  if action == "removed"{
    sign = "[-]"
  }else{
    sign = "[+]"
  }
	fmt.Printf("%s %s has been %s...\n",sign,name,action)
}

func CreateStructure(program string,files ...string){
  base := fmt.Sprintf("C:\\Program Files\\%s",program)
  if exists(base){
    fmt.Printf("[-] Removing existing file structure.\n")
    os.RemoveAll(base)
  }else{
    fmt.Printf("[+] Creating needed file structure.\n")
    for _,file := range files{
      path := filepath.Join(base,file)
      fmt.Printf("[+] Creating %s.\n",path)
      err := os.MkdirAll(path, 0755)
      if err != nil {
          panic(err)
      }
    }
  }
}
