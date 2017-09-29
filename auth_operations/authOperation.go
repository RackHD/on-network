package auth_operations

import (
//	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (c * Claims) SetToken(username string, password string ) string {
	// Expires the token and cookie in 1 hour
	expireToken := time.Now().Add(time.Hour * 1).Unix()

	// We'll manually assign the claims but in production you'd insert values from a database
	claims := Claims {
		username,
		jwt.StandardClaims {
			ExpiresAt: expireToken,
			Issuer:    "localhost:9000",
		},
	}

	// Create the token using your claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signs the token with a secret.
	signedToken, _ := token.SignedString([]byte("secret"))


	fmt.Println("Token created : -----> %+v" , signedToken)

	return signedToken

}

//// middleware to protect private pages
//func validate(page http.HandlerFunc) http.HandlerFunc {
//
//}


