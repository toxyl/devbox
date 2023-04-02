package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/toxyl/gutils"
)

func StringToSha256(str string) string {
	return gutils.StringToSha256(str)
}

func FileToSha256(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return "0"
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "0"
	}
	return hex.EncodeToString(hash.Sum(nil))
}
