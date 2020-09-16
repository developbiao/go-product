package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// Advanced encryption Standard AES
var PwdKey = []byte("biaoge@*golang##")

// PKCS7 padding mode
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// pKCS7 unpadding
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("Encrypt string error!")
	} else {
		// Get padding string length
		unpadding := int(origData[length-1])
		// delete padding string
		return origData[:(length - unpadding)], nil
	}
}

// encrypter
func AesEncrypt(origData []byte, key []byte) ([]byte, error) {
	// create encrypt  algorithm instance
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Get block size
	blockSize := block.BlockSize()
	// data fill
	origData = PKCS7Padding(origData, blockSize)
	// AES CBC
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))

	// execute encrypt
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// Decrypter
func AesDeCrypt(cypted []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Get block size
	blockSize := block.BlockSize()
	// Create Client decrypt instance
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	// decrypt
	blockMode.CryptBlocks(origData, cypted)
	// remove padding string
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, err
}

// Base64 encoding
func EnPwdCode(pwd []byte) (string, error) {
	result, err := AesEncrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

// Base64 decode
func DePwdCode(pwd string) ([]byte, error) {
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}
	// AES decrypt
	return AesDeCrypt(pwdByte, PwdKey)
}
