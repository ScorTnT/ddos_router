package router

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/GehirnInc/crypt/md5_crypt"
)

func CheckAccount(username, password string) (bool, error) {
	shadowFile, err := os.Open("/etc/shadow")
	if err != nil {
		return false, fmt.Errorf("failed to open /etc/shadow: %v", err)
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
			// Account is locked
			return false, nil
		}

		hash := fields[1]
		crypto := md5_crypt.New()
		err := crypto.Verify(hash, []byte(password))
		result := (err == nil)

		return result, nil
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("failed to read /etc/shadow: %v", err)
	}

	return false, fmt.Errorf("user not found")
}
