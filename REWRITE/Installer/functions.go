package main

import(
  "fmt"
  "os"
  "path/filepath"

  "github.com/go-ole/go-ole"
  "github.com/go-ole/go-ole/oleutil"
  "github.com/emersion/go-autostart"
  "github.com/mitchellh/go-homedir"
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

func CreateStructure(program string,files ...string)string{
  base := filepath.Join("C:\\Program Files",program)
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
  return filepath.Join(base,program+".exe")
}

func CreateShortcut(src, name string)error{
  fmt.Println(src)
  dir, err := homedir.Dir()
  if err != nil{
    return err
  }
  dir, err = homedir.Expand(dir)
  if err != nil{
    return err
  }
  dst := filepath.Join(dir,"Desktop",name+".lnk")
  fmt.Println(dst)


  ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
  oleShellObject, err := oleutil.CreateObject("WScript.Shell")
  if err != nil {
    return err
  }
  defer oleShellObject.Release()
  wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
  if err != nil {
    return err
  }
  defer wshell.Release()
  cs, err := oleutil.CallMethod(wshell, "CreateShortcut", dst)
  if err != nil {
    return err
  }
  idispatch := cs.ToIDispatch()
  oleutil.PutProperty(idispatch, "TargetPath", src)
  oleutil.CallMethod(idispatch, "Save")
  return nil
}
