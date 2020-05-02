package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
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
	crytedByte, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
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

// RsaEncrypt 非对称加密
func RsaEncrypt(text, publicKey string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)
	encryped, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(text))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryped), nil
}

// RsaDecrypt 非对称解密
func RsaDecrypt(text, privateKey string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, priv, bytes)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// BcryptHash 适合做密码哈希的函数
func BcryptHash(text string, level int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(text), level)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckBcryptHash 检查给定的hash与字符串是否匹配
func CheckBcryptHash(text, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
	return err == nil
}

// BcryptHashLevel 返回哈希bcrypt的复杂等级
func BcryptHashLevel(hash string) (int, error) {
	bytes := []byte(hash)
	level, err := bcrypt.Cost(bytes)
	if err != nil {
		return 0, err
	}
	return level, nil
}

// MD5Hash md5哈希
func MD5Hash(text string) string {
	bytes := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", bytes)
}

func main() {
	// orig := ""
	// key := "123456781234567812345678"
	// fmt.Println("原文：", orig)
	// encryptCode, _ := AesEncrypt(orig, key)
	// fmt.Println("密文：", encryptCode)
	// decryptCode, _ := AesDecrypt(encryptCode, key)
	// fmt.Println("解密结果：", decryptCode)
	///////////////////////////////////////////////////////////////////////////////////////
	// 	privateKey := `
	// -----BEGIN RSA PRIVATE KEY-----
	// MIICXAIBAAKBgQCq6jTJ8FBBmZtU266AoaVEbZOX9bI2BQeSGy08vJ9MBsIhsuSf
	// 3EXiF8Tbod1t5FI7VYxK+aOdCJQEAIocxN+wWaT6w+HYHS+TbM9ld6kTplZHhxXm
	// Lmoy5pM8KGnO+XDxslp+RyFEijqjFUApcgnoPHZawreX5R12+WvhB7zVGQIDAQAB
	// AoGAFMxRRi2XJpRxoX7mGkMRYGVRK5G7cPVluXLTK2vsmozL/NvPFGDypBOuYeet
	// kVzSYOQFIW6FaWuMz3vdUqjma1U6m1oPoP1c4E1zo49PCUfuVX/Tr/RoK7Yt4oA3
	// oU2zR2ZO7SCrfLuorc3uuZyFXHDDLLf/QfBlzyZQYXHzXnECQQDiregnnhEZh7JC
	// E8k1LSWWHtzrXmJ/N1vtWAiaQKLG9IQv+6Li27NGWK3DcMc/AtGSkEiOJpRZs3hm
	// BMwcbNpFAkEAwQXBiSIy51j82+V4BSdDMcYOcW9LssowXJg9BYmmQ5viA7luGmOh
	// 66hjNh9kaAa0XljQdVh8u/1p4hh5zMJGxQJAGfuFGMVtgFelfx11errqDAMNn4z6
	// DUrzJoloswWloth0LmvRh99W9UFE6ZoLvJU+IGaZanHRn1aK/maeXsO7YQJAHjBS
	// Rqpsn7ZRResrzFKJtd0uIRuxETOCQNkmohkXQfe9b6uuxr3Sc1MQwURDpel2iFMH
	// t3Kxf41w5AXD9vxjqQJBAMRg2eKNNrag2xBEkCKOwJaHKP7xTpybSyADvVcCwZW3
	// Z37SlNRfOdd8q1bpHcBn41KXMXPJ61j3yQ+VkHq7OZQ=
	// -----END RSA PRIVATE KEY-----
	// `
	// 	publicKey := `
	// -----BEGIN PUBLIC KEY-----
	// MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCq6jTJ8FBBmZtU266AoaVEbZOX
	// 9bI2BQeSGy08vJ9MBsIhsuSf3EXiF8Tbod1t5FI7VYxK+aOdCJQEAIocxN+wWaT6
	// w+HYHS+TbM9ld6kTplZHhxXmLmoy5pM8KGnO+XDxslp+RyFEijqjFUApcgnoPHZa
	// wreX5R12+WvhB7zVGQIDAQAB
	// -----END PUBLIC KEY-----
	// `
	// 	data, err := RsaEncrypt("http://liuqiqi/golang/", publicKey)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// 	fmt.Println(data)
	// 	origData, err := RsaDecrypt(data, privateKey)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// 	fmt.Println(origData)
	//////////////////////////////////////////////////////////////////////////////////////////////
	// h, _ := BcryptHash("123456789", 14)
	// fmt.Println("密码哈希字符串： ", h)
	// t := CheckBcryptHash("123456789", h)
	// fmt.Printf("密码与哈希是否匹配： %v\n", t)
	// l, _ := BcryptHashLevel(h)
	// fmt.Println("哈希算法的等级：", l)
	//////////////////////////////////////////////////////////////////////////////////////////////
	fmt.Println(MD5Hash("FHDLFHDL"))
}
