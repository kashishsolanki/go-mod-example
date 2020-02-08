package apiutils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// GetUserID to get user_id from token
func GetUserID(r *http.Request) (string, error) {

	hmacSecretString := os.Getenv("SECRET_HRMS_KEY") // "SECRET_HRMS_KEY" // Value
	hmacSecret := []byte(hmacSecretString)
	fmt.Println(strings.Split(r.Header["Authorization"][0], " "))

	token, err := jwt.Parse(strings.Split(r.Header["Authorization"][0], " ")[1], func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	fmt.Println("Token :: ", token)
	if err != nil {
		fmt.Println("Error while get data from JWT token")
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims, claims["id"])
		for key, val := range claims {
			fmt.Printf("Key: %v, value: %v\n", key, val)
		}
		return fmt.Sprint(claims["id"].(float64)), nil
	}
	log.Printf("Invalid JWT Token")
	return "", errors.New("Invalid JWT Token")
}
