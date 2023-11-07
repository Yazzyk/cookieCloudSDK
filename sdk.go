package cookieCloudSDK

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
)

type CookieCloudSDK struct {
	Url      string // 服务器地址
	UUID     string // UUID 用户KEY
	Password string // 密码 端对端加密密码
}

func NewCookieCloudSDK(host, uuid, password string) (*CookieCloudSDK, error) {
	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	return &CookieCloudSDK{
		Url:      url.String(),
		UUID:     uuid,
		Password: password,
	}, nil
}

// GetCookie 获取解密后的数据
func (this *CookieCloudSDK) GetCookie() (CookieCloudDecodeData, error) {
	resp, err := resty.New().R().Get(this.Url + "/get/" + this.UUID)
	if err != nil {
		return CookieCloudDecodeData{}, err
	}
	var rsJson CookieCloudRs
	if err = json.Unmarshal(resp.Body(), &rsJson); err != nil {
		return CookieCloudDecodeData{}, err
	}
	key := this.getKey()
	data, err := this.Decrypt(key, rsJson.Encrypted)
	if err != nil {
		return CookieCloudDecodeData{}, err
	}
	var cookieData CookieCloudDecodeData
	if err = json.Unmarshal(data, &cookieData); err != nil {
		return CookieCloudDecodeData{}, err
	}
	return cookieData, nil
}

// Decrypt 解密
func (this *CookieCloudSDK) Decrypt(passphrase, encryptedData string) ([]byte, error) {
	// base64 decode encrypted data
	encrypted, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("base64 decode err: %w", err)
	}
	salt := encrypted[8:16]
	key_iv, _ := this.bytesToKey([]byte(passphrase), salt, 48)
	key := key_iv[:32]
	iv := key_iv[32:]
	ciphertext := encrypted[16:]

	return this.AesDecrypt(ciphertext, key, iv)
}

func (this *CookieCloudSDK) bytesToKey(data []byte, salt []byte, output int) ([]byte, error) {
	if len(salt) != 8 {
		return nil, fmt.Errorf("expected salt of length 8, got %d", len(salt))
	}
	data = append(data, salt...)
	hash := md5.Sum(data)
	key := hash[:]
	finalKey := append([]byte(nil), key...)
	for len(finalKey) < output {
		hash = md5.Sum(append(key, data...))
		key = hash[:]
		finalKey = append(finalKey, key...)
	}
	return finalKey[:output], nil
}

func (this *CookieCloudSDK) AesEncrypt(data, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	encryptBytes := this.pkcs7Padding(data, blockSize)
	crypted := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

func (this *CookieCloudSDK) AesDecrypt(data, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)
	crypted, err = this.pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

func (this *CookieCloudSDK) pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func (this *CookieCloudSDK) pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("empty data to unpad")
	}
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func (this *CookieCloudSDK) getKey() string {
	hash := md5.Sum([]byte(this.UUID + "-" + this.Password))
	key := hex.EncodeToString(hash[:])[:16]
	return key
}
