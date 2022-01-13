package core

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	DateFormat = "2006-01-02 15:02"
)

func init() {

}

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

