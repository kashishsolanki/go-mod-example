package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/kashishsolanki/dt-hrms-golang/apiutils"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

		fmt.Println(secret)
		if len(secret) == 0 {
			return nil, errors.New("Auth0 Client Secret Not Set")
		}

		return []byte(secret), nil
	},

	Extractor: extractTokenFromCookie,
	Debug:     false,
})

func extractTokenFromCookie(r *http.Request) (string, error) {

	// get Authorization Header
	auth := r.Header.Get("Authorization")

	// Authorization token:
	// Bearer XXXXX
	//
	// Get actual token
	splitAuth := strings.Split(auth, " ")
	if len(splitAuth) < 2 {
		return "", errors.New("Invalid Auth: length of split authorization token is less than 2")
	}

	// Assume full token passed in header
	headerToken := splitAuth[1]
	headerTokenParts := strings.Split(headerToken, ".")

	if len(headerTokenParts) == 3 {
		return headerToken, nil
	}

	headerTokenPayload := headerToken

	// Base64 decode
	decodedAuth, err := base64.RawStdEncoding.DecodeString(headerTokenPayload)
	if err != nil {
		return "", fmt.Errorf("Invalid Auth. Could not decode from base64: %s", err)
	}

	// Unmarshal into HrmsJwtClaims
	hrmsAuth := &apiutils.HrmsJwtClaims{}
	err = json.Unmarshal(decodedAuth, hrmsAuth)
	if err != nil {
		return "", fmt.Errorf("Invalid Auth: Could not unmarshal Auth header into HrmsJwtClaims: %s", err)
	}

	// Get header and key from permission_lookup object
	// used to find the cookie access token
	header := hrmsAuth.PermissionLookup.Header
	key := hrmsAuth.PermissionLookup.Key

	// Get cookie header
	cookieHeaderSlice, ok := r.Header[header]
	if !ok {
		return "", fmt.Errorf("Header %s specified in permission lookup object not found", header)
	}

	// Ensure Cookie is populated
	if len(cookieHeaderSlice) < 1 {
		return "", fmt.Errorf("Header %s specified in permission lookup object is empty", header)
	}

	// Cookie string currently looks like:
	// 		foo=bar; hello=goodbye
	// Convert to map:
	// 		cookieMap[foo] = bar
	//		cookieMap[hello] = goodbye
	cookieMap, err := cookieToMap(cookieHeaderSlice[0])
	if err != nil {
		return "", fmt.Errorf("Error converting cookie to map: %s", err)
	}

	// From the map converted above, get the actual key
	// specified in permission_lookup object
	cookie, ok := cookieMap[key]
	if !ok {
		return "", fmt.Errorf("Key %s specified in permission lookup object not found in header %s", key, header)
	}

	// Split the token by periods
	cSplit := strings.Split(cookie, ".")
	if len(cSplit) < 2 {
		return "", errors.New("Invalid access token in cookie")
	}

	// Validate token has not been tampered with
	if splitAuth[1] != cSplit[1] {
		return "", errors.New("Token revoked, showed evidence of tampering")
	}

	return cookie, nil
}

// Input:
// 		"foo=bar; hello=goodbye==;"
// Output:
// 		map["foo"] = "bar"
//		map["hello"] = "goodbye=="
func cookieToMap(cookie string) (map[string]string, error) {
	ret := make(map[string]string)

	cookies := strings.Split(cookie, ";")

	for _, c := range cookies {
		if len(c) < 1 {
			continue
		}

		trimmedCookie := strings.TrimSpace(c)
		cookieSplit := strings.SplitN(trimmedCookie, "=", 2)
		if len(cookieSplit) < 2 {
			return nil, errors.New("invalid cookie map")
		}

		ret[cookieSplit[0]] = cookieSplit[1]
	}

	return ret, nil
}
