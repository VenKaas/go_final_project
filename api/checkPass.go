package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/VenKaas/go_final_project/env"
	"github.com/golang-jwt/jwt"
)

type AuthPass struct {
	Pass string `json:"password"`
}

type AuthPassError struct {
	MyToken string `json:"token,omitempty"`
	Err     string `json:"error,omitempty"`
}

var AuthResult AuthPassError
var buffer bytes.Buffer
var auth AuthPass

func (srv Server) CheckPass(rw http.ResponseWriter, rq *http.Request) {
	//читаем тело запроса
	_, err := buffer.ReadFrom(rq.Body)
	checkErr(err)

	//десериализируем тело запроса в структуру auth
	if err = json.Unmarshal(buffer.Bytes(), &auth); err != nil {
		fmt.Println("ошибка десериализации")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	//проверяем на совпадение TODO_PASSWORD и сам запрос
	if auth.Pass == os.Getenv("TODO_PASSWORD") {
		//формируем токен
		secret := []byte(auth.Pass)
		jwtToken := jwt.New(jwt.SigningMethodHS256)
		AuthResult.MyToken, err = jwtToken.SignedString(secret)
		checkErr(err)

	} else {
		AuthResult.Err = "Неверный пароль"
	}

	//возвращаем токен в поле token или ошибку
	srv.Server.Response(AuthResult, rw)

}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		var jwt string // переменная для хранения токена
		// получаем пароль
		pass := env.SetPass()
		if len(pass) > 0 {
			// получаем куки
			cookie, err := rq.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			if jwt != AuthResult.MyToken {
				// если не совпадает то возвращаем ошибку 401
				http.Error(rw, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(rw, rq)
	})
}
