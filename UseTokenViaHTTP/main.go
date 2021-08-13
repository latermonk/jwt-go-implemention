package main

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	signKey    *rsa.PrivateKey
	serverPort int
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Define some custom types were going to use within our tokens
type CustomerInfo struct {
	Name string
	Kind string
}

type CustomClaimsExample struct {
	*jwt.StandardClaims
	TokenType string
	CustomerInfo
}

func main() {

	// Make a sample token
	// In a real world situation, this token will have been acquired from
	// some other API call (see Example_getTokenViaHTTP)
	token, err := createToken("foo")
	fatal(err)

	// Make request.  See func restrictedHandler for example request processor
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/restricted", serverPort), nil)
	fatal(err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	res, err := http.DefaultClient.Do(req)
	fatal(err)

	// Read the response body
	buf := new(bytes.Buffer)
	io.Copy(buf, res.Body)
	res.Body.Close()
	fmt.Println(buf.String())

}

func createToken(user string) (string, error) {

	t := jwt.New(jwt.GetSigningMethod("RS256"))

	t.Claims = &CustomClaimsExample{
		&jwt.StandardClaims{

			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		},
		"level1",
		CustomerInfo{user, "human"},
	}

	return t.SignedString(signKey)
}