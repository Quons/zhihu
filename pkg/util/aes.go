package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type AesEncrypt struct {
}

func (this *AesEncrypt) getKey() []byte {
	strKey := "60f2fa84aed5a7a766c37e07817966dd"
	keyLen := len(strKey)
	if keyLen < 16 {
		panic("res key 长度不能小于16")
	}
	arrKey := []byte(strKey)
	if keyLen >= 32 {
		//取前32个字节
		return arrKey[:32]
	}
	if keyLen >= 24 {
		//取前24个字节
		return arrKey[:24]
	}
	//取前16个字节
	return arrKey[:16]
}

//加密字符串
func (this *AesEncrypt) AesEncrypt(strMesg string) string {
	key := this.getKey()
	var iv = []byte(key)[:aes.BlockSize]
	encrypted := make([]byte, len(strMesg))
	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, []byte(strMesg))
	return base64.URLEncoding.EncodeToString(encrypted)
}

//解密字符串
func (this *AesEncrypt) AesDecrypt(str string) string {
	key := this.getKey()
	var iv = []byte(key)[:aes.BlockSize]
	src, _ := base64.URLEncoding.DecodeString(str)
	decrypted := make([]byte, len(src))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, src)
	return string(decrypted)
}
