package base

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
)

func EncryptWithMd5(plaintext string) (ciphertext string) {
	h := md5.New()
	io.WriteString(h, plaintext)
	ciphertext = fmt.Sprintf("%x", h.Sum(nil))
	return
}

func EncryptWithSha1(plaintext string) (ciphertext string) {
	h := sha1.New()
	io.WriteString(h, plaintext)
	ciphertext = fmt.Sprintf("%x", h.Sum(nil))
	return
}
