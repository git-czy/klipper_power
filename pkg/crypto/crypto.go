package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"log"
)

type DeviceKey struct {
	key [16]byte
	iv  [16]byte
}

func DeviceKeyFromToken(token []byte) DeviceKey {
	key := Md5Byte(token[:])
	return DeviceKey{
		key: key,
		iv:  Md5Byte(key[:], token[:]),
	}
}

func Md5Byte(chunks ...[]byte) [16]byte {
	hash := md5.New()
	for _, chunk := range chunks {
		hash.Write(chunk)
	}
	var result [16]byte
	copy(result[:], hash.Sum(nil))
	return result
}

func (d DeviceKey) Encrypt(data []byte) []byte {
	block, err := aes.NewCipher(d.key[:])
	if err != nil {
		log.Println("New cipher failed at encrypt" + err.Error())
		return nil
	}
	mode := cipher.NewCBCEncrypter(block, d.iv[:])
	padData := pad(data, mode.BlockSize())
	res := make([]byte, len(padData))
	mode.CryptBlocks(res, padData)
	return res
}

func (d DeviceKey) Decrypt(data []byte) []byte {
	block, err := aes.NewCipher(d.key[:])
	if err != nil {
		log.Println("New cipher failed at decrypt" + err.Error())
		return nil
	}
	mode := cipher.NewCBCDecrypter(block, d.iv[:])
	res := make([]byte, len(data))
	mode.CryptBlocks(res, data)
	return res
}

func pad(data []byte, blockSize int) []byte {
	lenData := len(data)
	padding := blockSize - lenData%blockSize
	result := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, result...)
}
