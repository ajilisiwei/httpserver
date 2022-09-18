package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

var hmacSampleSecret = []byte("weiweiwei")

var USERS = map[string]string{"foo": "123", "bar": "456"}

type Claims struct {
	Acc string `json:"acc"`
	jwt.StandardClaims
}

func main() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		acc := r.PostForm["acc"][0]
		psw := r.PostForm["psw"][0]
		if acc != "" && psw != "" {
			if USERS[acc] == psw {
				// Create a new token object, specifying signing method and the claims
				// you would like it to contain.
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"acc": acc,
					"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
				})
				// Sign and get the complete encoded token as a string using the secret
				tokenString, _ := token.SignedString(hmacSampleSecret)
				// 设置响应头 Token 需要客户端配合读取
				w.Header().Add("Token", tokenString)
				w.WriteHeader(200)
				w.Write([]byte("hello," + acc))
				return
			}
		}
		w.WriteHeader(403)
		w.Write([]byte("Unauthorizate"))
	})

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		// 读取 Token
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(401)
			w.Write([]byte("Unauthorizate"))
		}
		// 解析 Token
		token := strings.Split(auth, " ")[1]
		tokenClaims, _ := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return hmacSampleSecret, nil
		})
		if claims, ok := tokenClaims.Claims.(*Claims); ok {
			// TODO 判断是否过期
			if claims != nil {
				fmt.Printf("------:%v\n", claims)
				fmt.Printf("ExpiresAt:%v\n", claims.ExpiresAt)
			}
			w.WriteHeader(200)
			w.Write([]byte("hello," + claims.Acc))
		}
	})

	fmt.Println("Listen on 8080....")
	http.ListenAndServe(":8080", nil)
}
