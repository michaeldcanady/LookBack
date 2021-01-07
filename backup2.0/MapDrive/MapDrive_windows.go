package MapDrive

import (
	"syscall"
	"unsafe"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	NO_ERROR				= 0x00000000
	ERROR_SESSION_CREDENTIAL_CONFLICT	= 0x000004C3
	RESOURCETYPE_DISK			= 0x00000001
	CONNECT_TEMPORARY			= 0x00000004
)

var(
  mpr,_                    = syscall.LoadDLL("mpr.dll")
  procWNetAddConnection2,_ = mpr.FindProc("WNetAddConnection2W")
)

func test() {
	ip := "\\\\server\\path\\to\\loc"
	user := ""
	pass := ""

	err := WNetAddConnection2(ip, user, pass)
	if err != nil {
		panic("Error mapping shared resource: " + err.Error())
	}
  err = os.Mkdir(filepath.Join(ip,"LookBackTest"),os.ModePerm)
  if err != nil{
    panic(err)
  }
}

func WNetAddConnection2(rpath, user, pass string) error {

	//https://msdn.microsoft.com/library/windows/desktop/aa385353(v=vs.85).aspx
	type NETRESOURCE struct {
		dwScope uint32
		dwType uint32
		dwDisplayType uint32
		dwUsage  uint32
		lpLocalName uintptr
		lpRemoteName uintptr
		lpComment uintptr
		lpProvider uintptr
	}

	ret, _, _ := procWNetAddConnection2.Call(
		uintptr(unsafe.Pointer(&NETRESOURCE{
			dwType: RESOURCETYPE_DISK,
			lpRemoteName: uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(rpath))),
		})),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(pass))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(user))),
		CONNECT_TEMPORARY)

	switch ret {
		case NO_ERROR, ERROR_SESSION_CREDENTIAL_CONFLICT:	//succesfully authenticate or already did
			return nil
		default:
			return fmt.Errorf("Error has appeared, return value %x", ret)
	}
}

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil
}

func run(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Done")
}
