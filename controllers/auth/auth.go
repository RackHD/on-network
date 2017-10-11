package auth

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	ao "github.com/RackHD/on-network/auth_operations"
	"github.com/RackHD/on-network/models"
)

// SwitchConfig is a struct for the http objects
type Auth struct {
	Request *http.Request
	Handler ao.Claims
	Login models.Login

}

// MiddleWare handles the route call
func MiddleWare(r *http.Request, body *models.Login) middleware.Responder {


	fmt.Println("Authentication")
	return &Auth{
		Request: r,
		Login: *body,
	}
}

// WriteResponse implements the CRUD logic behind the /credentials route
func (a *Auth) WriteResponse(rw http.ResponseWriter, rp runtime.Producer) {
	switch a.Request.Method {
	case http.MethodPost:
		a.postLogin(rw, rp)
	default:
		a.notSupported(rw, rp)
	}
}

func (a *Auth) notSupported(rw http.ResponseWriter, rp runtime.Producer) {
	rw.WriteHeader(http.StatusNotImplemented)
}

func (a *Auth) postLogin(rw http.ResponseWriter, rp runtime.Producer) {

	isValid, err := a.Handler.ValidateLogin(a.Login.Username, a.Login.Password)
	if err != nil {
		loginError := models.LoginError{
			Message: fmt.Sprintf("Invalid Credential, %+v" , err),
			Error: "401",
		}
		rw.WriteHeader(401)

		if err := rp.Produce(rw, loginError); err != nil {
			panic(err)
		}
	}else{
		if isValid {
			tokenValue :=  a.Handler.SetToken(a.Login.Username)
			token := models.Token{
				Token: tokenValue,
			}

			if err := rp.Produce(rw, token); err != nil {
				panic(err)
			}
		} else{
			loginError := models.LoginError{
				Message:"Invalid Credential" ,
				Error: "401",
			}
			rw.WriteHeader(401)

			if err := rp.Produce(rw, loginError); err != nil {
			panic(err)
			}
		}
	}
}


