package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

//#;AES加解密使用 长度 16, 24, 32
//const aesKey = "sfe023f_9fd&fwfl"

////加密字符串
//func AesEncrypt(hstring string) (string, error) {
//	if result, err := aesEncrypt([]byte(hstring), []byte(aesKey)); err != nil {
//		return "", err
//	} else {
//		return base64.StdEncoding.EncodeToString(result), nil
//	}
//}
//
////解密字符串
//func AesDecrypt(hstring string) (string, error) {
//	if bye, err := base64.StdEncoding.DecodeString(hstring); err != nil {
//		return "", err
//	} else {
//		if result, err := aesDecrypt(bye, []byte(aesKey)); err != nil {
//			return "", err
//		} else {
//			return string(result), nil
//		}
//	}
//}
//func AESTestCode() {
//	// AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
//	key := []byte("sfe023f_9fd&fwfl")
//	result, err := aesEncrypt([]byte("polaris@studygolang"), key)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(result)
//	a, _ := AesEncrypt("polaris@studygolang")
//	aa, _ := base64.StdEncoding.DecodeString(a)
//	origData, err := aesDecrypt(aa, key)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(string(origData))
//}

func Encrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pKCS5Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func Decrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted

	blockMode.CryptBlocks(origData, crypted)

	origData = pKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

//func zeroPadding(ciphertext []byte, blockSize int) []byte {
//	padding := blockSize - len(ciphertext)%blockSize
//	padtext := bytes.Repeat([]byte{0}, padding)
//	return append(ciphertext, padtext...)
//}
//
//func zeroUnPadding(origData []byte) []byte {
//	length := len(origData)
//	unpadding := int(origData[length-1])
//	return origData[:(length - unpadding)]
//}
//
func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//
func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
