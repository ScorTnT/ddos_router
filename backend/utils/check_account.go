package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings" // for MD4 support

	"github.com/GehirnInc/crypt/md5_crypt"
)

func CheckAccount(username, password string) bool {
	authenticated, err := authenticateUser(username, password)
	if err != nil {
		fmt.Printf("Authentication error: %v\n", err)
		return false
	}
	return authenticated
}

func authenticateUser(username, password string) (bool, error) {
	shadowFile, err := os.Open("/etc/shadow")
	if err != nil {
		return false, fmt.Errorf("error opening shadow file: %v", err)
	}
	defer shadowFile.Close()

	scanner := bufio.NewScanner(shadowFile)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if len(fields) < 2 || fields[0] != username {
			continue
		}

		if strings.HasPrefix(fields[1], "!") || strings.HasPrefix(fields[1], "*") {
			return false, fmt.Errorf("account is locked")
		}

		hash := fields[1]
		crypt := md5_crypt.New()
		err := crypt.Verify(hash, []byte(password))
		result := (err == nil)

		return result, nil
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("error reading shadow file: %v", err)
	}

	return false, fmt.Errorf("user not found")
}
