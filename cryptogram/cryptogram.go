package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// PKCS7Padding 明文填充
func PKCS7Padding(text []byte, blockSize int) []byte {
	padding := blockSize - len(text)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(text, padtext...)
}

// PKS7UnPadding 明文截断
func PKS7UnPadding(text []byte) []byte {
	length := len(text)
	uppadding := int(text[length-1])
	return text[:(length - uppadding)]
}

// AesEncrypt 对称加密
func AesEncrypt(origin, key string) (string, error) {
	origBytes := []byte(origin)
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origBytes = PKCS7Padding(origBytes, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, keyBytes[:blockSize])
	cryted := make([]byte, len(origBytes))
	blockMode.CryptBlocks(cryted, origBytes)
	return base64.StdEncoding.EncodeToString(cryted), nil
}

// AesDecrypt 对称解密
func AesDecrypt(text, key string) (string, error) {
	crytedByte, _ := base64.StdEncoding.DecodeString(text)
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	orig := make([]byte, len(crytedByte))
	blockMode.CryptBlocks(orig, crytedByte)
	orig = PKS7UnPadding(orig)
	return string(orig), nil
}
func main() {
	orig := "http://c.biancheng.net/golang/"
	key := "123456781234567812345678"
	fmt.Println("原文：", orig)
	encryptCode, _ := AesEncrypt(orig, key)
	fmt.Println("密文：", encryptCode)
	decryptCode, _ := AesDecrypt(encryptCode, key)
	fmt.Println("解密结果：", decryptCode)
}
