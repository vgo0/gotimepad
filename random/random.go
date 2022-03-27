package random

import (
	"crypto/rand"
)

func GenerateRandomBytes(num uint) ([]byte, error) {
	ret := make([]byte, num)
	
	_, err := rand.Read(ret)

	return ret, err
}