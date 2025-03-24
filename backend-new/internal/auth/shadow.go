//go:build linux && cgo
// +build linux,cgo

package auth

/*
#include <stdlib.h>
#include <unistd.h>
#include <crypt.h>
*/
import "C"
import (
	"bufio"
	"os"
	"strings"
	"unsafe"
)

// Verify(username,password) returns true if /etc/shadow 에 저장된 hash 와 일치
func Verify(username, password string) (bool, error) {
	f, err := os.Open("/etc/shadow")
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), ":")
		if cols[0] != username {
			continue
		}
		hash := cols[1]
		cs := C.CString(password)
		ss := C.CString(hash)
		defer C.free(unsafe.Pointer(cs))
		defer C.free(unsafe.Pointer(ss))
		return C.GoString(C.crypt(cs, ss)) == hash, nil
	}
	return false, nil
}
