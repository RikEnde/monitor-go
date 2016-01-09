package util

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"runtime"
)

func JsonResponse(value interface{}) string {
	jsonStr, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(jsonStr)
}

func Cert() string {
	return fmt.Sprintf("%s/certificates/cert.pem", runtime.GOROOT())
}

func Key() string {
	return fmt.Sprintf("%s/certificates/key.pem", runtime.GOROOT())
}

func Hash(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

func usersFile(name string) (string, error) {
	f, err := ioutil.ReadFile(fmt.Sprintf("%s/users/%s", runtime.GOROOT(), name))
	if err != nil {
		return "", errors.New("Not found: " + name)
	}
	return Trim(string(f)), nil
}

func credentials(name string) (map[string]string, error) {
	f, err := usersFile(name)
	if err != nil {
		return nil, err
	}
	var cred map[string]string
	if err := json.Unmarshal([]byte(f), &cred); err != nil {
		return nil, err
	}
	return cred, nil
}

func DbCredentials() map[string]string {
	cred, err := credentials("postgres")
	if err != nil {
		panic(err)
	}
	return cred
}

func Login(user, password string) bool {
	cred, err := credentials("users")
	if err != nil {
		/* Unknown user */
		return false
	}
	return cred[user] == Hash(password)
}
