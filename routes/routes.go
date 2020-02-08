package routes

import (
	"net/http"
)

type Route struct {
	Name         string
	Method       string
	Pattern      string
	HandlerFunc  http.HandlerFunc
	VerifyJWT    bool
	VerifyPerms  bool
	VerifyAPIKey bool
}

type Routes []Route

func (routes *Routes) AppendRoutes(newRoutes Routes) {
	for _, r := range newRoutes {
		*routes = append(*routes, r)
	}
}
