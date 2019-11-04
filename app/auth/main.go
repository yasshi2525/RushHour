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
	"log"

	"github.com/gomodule/oauth1/oauth"
	"golang.org/x/oauth2"

	"github.com/yasshi2525/RushHour/app/config"
)

// Auther enc/dec sensitive data
type Auther struct {
	baseURL       string
	state         string
	salt          string
	block         cipher.Block
	twitterClient *oauth.Client
	githubConf    *oauth2.Config
	googleConf    *oauth2.Config
}

// OAuthInfo represents infomation from OAuth App
type OAuthInfo struct {
	Handler     *Auther
	DisplayName string
	Image       string
	LoginID     string
	OAuthToken  string
	OAuthSecret string
	IsEnc       bool
}

// GetAuther initiates cipher, connection
func GetAuther(conf config.CnfAuth) (*Auther, error) {
	a := &Auther{
		baseURL: conf.BaseURL,
		state:   conf.State,
		salt:    conf.Salt,
	}
	key := "0123456789abcdef"
	if len(conf.Key) == 16 {
		key = conf.Key
	} else {
		log.Printf("auth.key %s must be 16 length. set to 0123456789abcdef", conf.Key)
	}

	var err error
	if a.block, err = aes.NewCipher([]byte(key)); err != nil {
		return nil, err
	}

	a.initTwitter(conf.Twitter)
	a.initGoogle(conf.Google)
	a.initGitHub(conf.GitHub)

	return a, nil
}

// IsValid returns whether OAuthToken exists or not
func (i *OAuthInfo) IsValid() bool {
	if i.IsEnc {
		return i.Handler.Decrypt(i.OAuthToken) != ""
	}
	return i.OAuthToken != ""
}

// Enc returns encrypted info
func (i *OAuthInfo) Enc() (*OAuthInfo, error) {
	if i.IsEnc {
		return nil, fmt.Errorf("%v is already encrypted", i)
	}
	return &OAuthInfo{
		Handler:     i.Handler,
		LoginID:     i.Handler.Encrypt(i.LoginID),
		OAuthToken:  i.Handler.Encrypt(i.OAuthToken),
		OAuthSecret: i.Handler.Encrypt(i.OAuthSecret),
		DisplayName: i.Handler.Encrypt(i.DisplayName),
		Image:       i.Handler.Encrypt(i.Image),
		IsEnc:       true,
	}, nil
}

// Dec returns decrypted info
func (i *OAuthInfo) Dec() (*OAuthInfo, error) {
	if !i.IsEnc {
		return nil, fmt.Errorf("%v is already decrypted", i)
	}
	return &OAuthInfo{
		Handler:     i.Handler,
		LoginID:     i.Handler.Decrypt(i.LoginID),
		OAuthToken:  i.Handler.Decrypt(i.OAuthToken),
		OAuthSecret: i.Handler.Decrypt(i.OAuthSecret),
		DisplayName: i.Handler.Decrypt(i.DisplayName),
		Image:       i.Handler.Decrypt(i.Image),
		IsEnc:       false,
	}, nil
}

// Digest returns hash value of plain
func (a *Auther) Digest(plain string) string {
	mac := hmac.New(sha512.New, []byte(a.salt))
	mac.Write([]byte(plain))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Encrypt returns AES and Base64 encoded string.
func (a *Auther) Encrypt(plain string) string {
	encrypter, iv := a.genEncrypter()
	padPlain := pad([]byte(plain))
	enc := make([]byte, len(padPlain))
	encrypter.CryptBlocks(enc, []byte(padPlain))
	// combine iv after enc
	enc = append(enc, iv...)
	return base64.StdEncoding.EncodeToString(enc)
}

// Decrypt returns plain text from AES and Base64 encoded string.
func (a *Auther) Decrypt(enc64 string) string {
	if enc64 == "" {
		return ""
	}
	if combined, err := base64.StdEncoding.DecodeString(enc64); err != nil {
		panic(err)
	} else {
		// trim iv after enc
		piv := len(combined) - a.block.BlockSize()
		enc := combined[:piv]
		iv := combined[piv:]
		plain := make([]byte, len(enc))
		cipher.NewCBCDecrypter(a.block, iv).CryptBlocks(plain, enc)
		return string(unpad(plain))
	}
}

func (a *Auther) genEncrypter() (cipher.BlockMode, []byte) {
	iv := make([]byte, a.block.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		panic(err.Error())
	}
	return cipher.NewCBCEncrypter(a.block, iv), iv
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
