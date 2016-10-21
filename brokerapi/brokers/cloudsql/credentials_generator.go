package cloudsql

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const (
	maxUsernameLength       = 16 // Limit from http://dev.mysql.com/doc/refman/5.7/en/user-names.html
	generatedPasswordLength = 32
)

func GenerateUsername(instanceID, bindingID string) (string, error) {
	if len(instanceID)+len(bindingID) == 0 {
		return "", fmt.Errorf("empty instanceID and bindingID")
	}

	username := bindingID + instanceID
	if len(username) > maxUsernameLength {
		username = username[:maxUsernameLength]
	}

	return username, nil
}

func GeneratePassword() (string, error) {
	rb := make([]byte, generatedPasswordLength)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	rs := base64.URLEncoding.EncodeToString(rb)

	return rs, nil
}
