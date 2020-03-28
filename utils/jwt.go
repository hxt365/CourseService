package utils

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"path/filepath"
	"time"
)

func PublicKey() *rsa.PublicKey {
	publicPath, err := filepath.Abs("/code/public.pem")
	if err != nil {
		panic(err.Error())
	}
	publicBytes, err := ioutil.ReadFile(publicPath)
	if err != nil {
		panic(err.Error())
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		panic(err.Error())
	}
	return publicKey
}

func PrivateKey() *rsa.PrivateKey {
	privatePath, err := filepath.Abs("/code/private.pem")
	if err != nil {
		panic(err.Error())
	}
	privateBytes, err := ioutil.ReadFile(privatePath)
	if err != nil {
		panic(err.Error())
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		panic(err.Error())
	}
	return privateKey
}

func generateJWTToken(id uint, role string) string {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = id
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	t, _ := token.SignedString(PrivateKey())
	return t
}
