package auth_operations

import (
//	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
	"os"
	"strings"
)

var signedToken= ""
const (
	bearer	string = "bearer"
	secret 	string = "secret"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}


func (c* Claims) ValidateLogin(username string, password string) (bool, error) {

	envUser,envChkUser := os.LookupEnv("SERVICE_USERNAME")
	envPass, envChkPass := os.LookupEnv("SERVICE_PASSWORD")

	if envChkUser == false || envChkPass == false{
		err := fmt.Errorf("Service Username or Password not set")
		return false, err
	}

	if username==envUser && password==envPass{
		return  true,nil
	}
	return false, nil
}

func (c * Claims) SetToken(username string ) string {

	// Expires the token and cookie in 1 hour
	expireToken := time.Now().Add(time.Hour * 2).Unix()

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
	signedToken, _ = token.SignedString([]byte(secret))


	fmt.Println("Token created : -----> %+v" , signedToken)

	return signedToken

}

// middleware to protect private pages
func  ValidateToken(tokenHeader string) bool {

	tokenString, ok := extractTokenFromAuthHeader(tokenHeader)
	if !ok {
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return  []byte(secret), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true
	} else {
		fmt.Println(err)
		return false
	}

}

func extractTokenFromAuthHeader(val string) (token string, ok bool) {
	authHeaderParts := strings.Split(val, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != bearer {
		return "", false
	}

	return authHeaderParts[1], true
}





