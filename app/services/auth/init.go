package auth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

// Init initiates cipher, connection
func Init(config Config) {
	baseURL = config.BaseURL
	state = config.State
	key := config.Key
	salt = config.Salt

	var err error
	if block, err = aes.NewCipher([]byte(key)); err != nil {
		panic(err.Error())
	}

	initTwitter(config.Twitter)
	initGoogle(config.Google)
	initGitHub(config.GitHub)
}

// IsValid returns whether OAuthToken exists or not
func (i *OAuthInfo) IsValid() bool {
	if i.IsEnc {
		return Decrypt(i.OAuthToken) != ""
	} else {
		return i.OAuthToken != ""
	}
}

// Enc returns encrypted info
func (i *OAuthInfo) Enc() *OAuthInfo {
	if i.IsEnc {
		panic(fmt.Errorf("%s is already encrypted", i.LoginID))
	}
	return &OAuthInfo{
		LoginID:     Encrypt(i.LoginID),
		OAuthToken:  Encrypt(i.OAuthToken),
		OAuthSecret: Encrypt(i.OAuthSecret),
		DisplayName: Encrypt(i.DisplayName),
		Image:       Encrypt(i.Image),
		IsEnc:       true,
	}
}

// Dec returns decrypted info
func (i *OAuthInfo) Dec() *OAuthInfo {
	if !i.IsEnc {
		panic(fmt.Errorf("%s is already decrypted", i.LoginID))
	}
	return &OAuthInfo{
		LoginID:     Decrypt(i.LoginID),
		OAuthToken:  Decrypt(i.OAuthToken),
		OAuthSecret: Decrypt(i.OAuthSecret),
		DisplayName: Decrypt(i.DisplayName),
		Image:       Decrypt(i.Image),
		IsEnc:       false,
	}
}

// Token returns digest of access token.
func (i *OAuthInfo) Token() string {
	if !i.IsValid() {
		panic(fmt.Errorf("%s's access token is empty", i.LoginID))
	}
	var key string
	if i.IsEnc {
		key = Decrypt(i.OAuthToken)
	} else {
		key = i.OAuthToken
	}

	return Digest(key)
}

// Digest returns hash value of plain
func Digest(plain string) string {
	mac := hmac.New(sha512.New, []byte(salt))
	mac.Write([]byte(plain))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Encrypt returns AES and Base64 encoded string.
func Encrypt(plain string) string {
	encrypter, iv := genEncrypter()
	padPlain := pad([]byte(plain))
	enc := make([]byte, len(padPlain))
	encrypter.CryptBlocks(enc, []byte(padPlain))
	// combine iv after enc
	enc = append(enc, iv...)
	return base64.StdEncoding.EncodeToString(enc)
}

// Decrypt returns plain text from AES and Base64 encoded string.
func Decrypt(enc64 string) string {
	if enc64 == "" {
		return ""
	}
	if combined, err := base64.StdEncoding.DecodeString(enc64); err != nil {
		panic(err)
	} else {
		// trim iv after enc
		piv := len(combined) - block.BlockSize()
		enc := combined[:piv]
		iv := combined[piv:]
		plain := make([]byte, len(enc))
		cipher.NewCBCDecrypter(block, iv).CryptBlocks(plain, enc)
		return string(unpad(plain))
	}
}

func genEncrypter() (cipher.BlockMode, []byte) {
	iv := make([]byte, block.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		panic(err.Error())
	}
	return cipher.NewCBCEncrypter(block, iv), iv
}

func pad(b []byte) []byte {
	padSize := aes.BlockSize - (len(b) % aes.BlockSize)
	pad := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(b, pad...)
}

func unpad(b []byte) []byte {
	padSize := int(b[len(b)-1])
	return b[:len(b)-padSize]
}
