package apiutils

import (
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	secret = []byte(os.Getenv("SECRET_HRMS_KEY")) // []byte("SECRET_HRMS_KEY")
)

// HrmsJwtClaims contains information currently being sent in a JWT
type HrmsJwtClaims struct {
	ID               int                  `json:"id"`
	Username         string               `json:"username"`
	IsAdmin          bool                 `json:"is_admin"`
	PermissionLookup HrmsPermissionLookup `json:"permissions_lookup"`
	jwt.StandardClaims
}

// HrmsPermissionLookup is used to know which cookie to look at
// to get permissions
type HrmsPermissionLookup struct {
	Header string `json:"header"`
	Key    string `json:"key"`
}
