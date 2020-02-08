package leaves

import (
	r "github.com/kashishsolanki/dt-hrms-golang/routes"
)

var routes = r.Routes{
	// ApplyLeaves
	r.Route{
		Name:         "ApplyLeaves",
		Method:       "POST",
		Pattern:      "/v1/api/leaves/apply",
		HandlerFunc:  ApplyLeaves,
		VerifyJWT:    false,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
	// GetAllLeaves
	r.Route{
		Name:         "GetAllLeaves",
		Method:       "GET",
		Pattern:      "/v1/api/leaves",
		HandlerFunc:  GetAllLeaves,
		VerifyJWT:    false,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
	// UpdateLeave
	r.Route{
		Name:         "UpdateLeave",
		Method:       "PUT",
		Pattern:      "/v1/api/leaves/{id}",
		HandlerFunc:  UpdateLeave,
		VerifyJWT:    false,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
	// ManageLeave
	r.Route{
		Name:         "ManageLeave",
		Method:       "PUT",
		Pattern:      "/v1/api/leaves/manage/{id}",
		HandlerFunc:  ManageLeave,
		VerifyJWT:    false,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
}

// GetRoutes returns local variable routes which contain all methods for the API
func GetRoutes() r.Routes {
	return routes
}
