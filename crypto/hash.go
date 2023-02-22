package crypto

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(text string) string {
	ctx := sha256.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func Md5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func Md5Raw(text string) []byte {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return ctx.Sum(nil)
}