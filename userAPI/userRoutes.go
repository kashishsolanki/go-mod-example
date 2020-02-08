package users

import (
	r "github.com/kashishsolanki/dt-hrms-golang/routes"
)

var routes = r.Routes{
	// SignIn
	r.Route{
		Name:         "SignIn",
		Method:       "POST",
		Pattern:      "/v1/api/employees/signin",
		HandlerFunc:  SignIn,
		VerifyJWT:    false,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
	// ChangePassword
	r.Route{
		Name:         "ChangePassword",
		Method:       "POST",
		Pattern:      "/v1/api/employees/password/change",
		HandlerFunc:  ChangePassword,
		VerifyJWT:    true,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
	// AddEmployee
	r.Route{
		Name:         "AddEmployee",
		Method:       "POST",
		Pattern:      "/v1/api/employees/create",
		HandlerFunc:  AddEmployee,
		VerifyJWT:    true,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
	// AddEmployeeProfile
	r.Route{
		Name:         "AddEmployeeProfile",
		Method:       "POST",
		Pattern:      "/v1/api/profile/insert",
		HandlerFunc:  AddEmployeeProfile,
		VerifyJWT:    true,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
	// GetAllEmployees
	r.Route{
		Name:         "GetAllEmployees",
		Method:       "GET",
		Pattern:      "/v1/api/employees",
		HandlerFunc:  GetAllEmployees,
		VerifyJWT:    true,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
	// GetEmployee
	r.Route{
		Name:         "GetEmployee",
		Method:       "GET",
		Pattern:      "/v1/api/employees/{id}",
		HandlerFunc:  GetEmployee,
		VerifyJWT:    true,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
	// UpdateEmployeeData
	r.Route{
		Name:         "UpdateEmployeeData",
		Method:       "PUT",
		Pattern:      "/v1/api/employees/{id}",
		HandlerFunc:  UpdateEmployeeData,
		VerifyJWT:    true,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
	// DeleteEmployee
	r.Route{
		Name:         "DeleteEmployee",
		Method:       "DELETE",
		Pattern:      "/v1/api/employees/{id}",
		HandlerFunc:  DeleteEmployee,
		VerifyJWT:    true,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
	// GetRoles
	r.Route{
		Name:         "GetRoles",
		Method:       "GET",
		Pattern:      "/v1/api/roles",
		HandlerFunc:  GetRoles,
		VerifyJWT:    true,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
	// GetTeams
	r.Route{
		Name:         "GetTeams",
		Method:       "GET",
		Pattern:      "/v1/api/teams",
		HandlerFunc:  GetTeams,
		VerifyJWT:    true,
		VerifyPerms:  false,
		VerifyAPIKey: false,
	},
}

// GetRoutes returns local variable routes which contain all methods for the API
func GetRoutes() r.Routes {
	return routes
}
