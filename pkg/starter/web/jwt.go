// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/hidevopsio/hiboot/pkg/system"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
)

// JwtMap is the JWT map
type JwtMap map[string]interface{}

// Token is the token string
type Token string

const (
	privateKeyPath = "/config/ssl/app.rsa"
	pubKeyPath     = "/config/ssl/app.rsa.pub"
)

var (
	verifyKey  *rsa.PublicKey
	signKey    *rsa.PrivateKey
	jwtHandler *JwtMiddleware
	jwtEnabled bool
)


// InitJwt init JWT, it read config/ssl/rsa and config/ssl/rsa.pub
func InitJwt(wd string) error {

	// check if key exist
	if io.IsPathNotExist(wd + privateKeyPath) {
		return &system.NotFoundError{Name: wd + privateKeyPath}
	}

	signBytes, err := ioutil.ReadFile(wd + privateKeyPath)
	if err != nil {
		return err
	}

	signKey, err = jwtgo.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}

	verifyBytes, err := ioutil.ReadFile(wd + pubKeyPath)
	if err != nil {
		return err
	}

	verifyKey, err = jwtgo.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}

	jwtHandler = NewJwtMiddleware(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwtgo.Token) (interface{}, error) {
			//log.Debug(token)
			return verifyKey, nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwtgo.SigningMethodRS256,
	})

	jwtEnabled = true

	return nil
}

// GenerateJwtToken generates JWT token with specified exired time
func GenerateJwtToken(payload JwtMap, expired int64, unit time.Duration) (*Token, error) {
	if jwtEnabled {
		claim := jwtgo.MapClaims{
			"exp": time.Now().Add(unit * time.Duration(expired)).Unix(),
			"iat": time.Now().Unix(),
		}

		for k, v := range payload {
			claim[k] = v
		}

		token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS256, claim)

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(signKey)

		jwtToken := Token(tokenString)

		return &jwtToken, err
	}

	return nil, fmt.Errorf("JWT does not work")
}
