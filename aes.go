package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
)

type AesCryptor struct {
	key []byte
	iv  []byte
}

var aesCryptor *AesCryptor

func GetAesCryptor() *AesCryptor {
	if aesCryptor == nil {
		//TODO get from db
		aesCryptor = &AesCryptor{key: []byte("xxxxxxxxxxxx"), iv: []byte("xxxxxxxxxxxx")}
	}
	return aesCryptor
}

//加密数据
func (a *AesCryptor) Encrypt(data string) string {
	aesBlockEncrypter, err := aes.NewCipher(a.key)
	if err != nil {
		log.Println("ase util err", err.Error())
		return ""
	}
	content := PKCS5Padding([]byte(data), aesBlockEncrypter.BlockSize())
	encrypted := make([]byte, len(content))
	aesEncrypter := cipher.NewCBCEncrypter(aesBlockEncrypter, a.iv)
	aesEncrypter.CryptBlocks(encrypted, content)
	str := hex.EncodeToString(encrypted)
	return strings.ToUpper(str)
}

//解密数据
func (a *AesCryptor) Decrypt(src string) (string, error) {
	defer recoverAes(src)
	if src == "" {
		return "", nil
	}
	aesBlockDecrypter, err := aes.NewCipher(a.key)
	if err != nil {
		log.Println("ase util decrypt aesBlockDecrypter err", err.Error())
		return "", err
	}
	var encryptByt []byte
	encryptByt, err = hex.DecodeString(src)
	if err != nil {
		log.Println("ase util decrypt encryptByt err", err.Error())
		return "", err
	}
	decrypted := make([]byte, len(encryptByt))
	aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, a.iv)
	aesDecrypter.CryptBlocks(decrypted, encryptByt)
	return string(PKCS5Trimming(decrypted)), nil
}

/**
 * PKCS5包装
 */
func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

/*
 * 解包装
 */
func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func recoverAes(src string) {
	if r := recover(); r != nil {
		fmt.Printf("aes string : %s err, recovered from %v", src, r)
	}
}
