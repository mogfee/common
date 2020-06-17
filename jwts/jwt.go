package jwts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"math/big"
)

var (
	//ES256 keys
	ECDSAKeyD = "7A429E82FF619D38CC8071111988FFA75625DD83B22E9EBEC29F17BFA7BF3A03"
	ECDSAKeyX = "76E93569AB21A614BCD581858D0066C8ED611DEFEEA2821CC43EC9E08948A151"
	ECDSAKeyY = "61BB8B7EF5333E2E87CDE6DF522BE6BF253C637768F9FA8D9EDCAB270E09B43C"
)

//获取token数据
func JWTGetMapString(map1 jwt.Claims) (string, error) {

	keyD := new(big.Int)
	keyX := new(big.Int)
	keyY := new(big.Int)
	keyD.SetString(ECDSAKeyD, 16)
	keyX.SetString(ECDSAKeyX, 16)
	keyY.SetString(ECDSAKeyY, 16)

	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     keyX,
		Y:     keyY,
	}
	privateKey := ecdsa.PrivateKey{D: keyD, PublicKey: publicKey}

	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodES256, map1)

	if tokenString, err := tokenJwt.SignedString(&privateKey); err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}

//解析jwt数据
func JWTGetStringMap(jwtString string) (map[string]interface{}, error) {
	keyX := new(big.Int)
	keyY := new(big.Int)

	keyX.SetString(ECDSAKeyX, 16)
	keyY.SetString(ECDSAKeyY, 16)
	publickKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     keyX,
		Y:     keyY,
	}
	jwtToken, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New(fmt.Sprintf("json 解析失败:%+v", token))
		}
		return &publickKey, nil
	})

	if err == nil {
		if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
			return claims, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("json 解析失败:%v", jwtString))
}

//
//func todogetkey() (map[string]string, error) {
//	randKey := rand.Reader
//	var err error
//	var prk *ecdsa.PrivateKey
//	prk, err = ecdsa.GenerateKey(elliptic.P256(), randKey)
//	if err != nil {
//		return nil, err
//	}
//	//puk := prk.PublicKey
//	map1 := make(map[string]string)
//	map1["prkD"] = fmt.Sprintf("%X", prk.D)
//	map1["prkX"] = fmt.Sprintf("%X", prk.X)
//	map1["prkY"] = fmt.Sprintf("%X", prk.Y)
//	return map1, nil
//}
