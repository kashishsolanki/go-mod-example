package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/kashishsolanki/dt-hrms-golang/apiutils"
	"github.com/kashishsolanki/dt-hrms-golang/db"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

// PasswordChangeStruct struct for login
type PasswordChangeStruct struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// User struct for login
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SigninResponse to return json reponse in below format
type SigninResponse struct {
	Token   string `json:"token"`
	EmpID   int    `json:"emp_id"`
	IsAdmin bool   `json:"is_admin"`
}

// SignIn will allow to signin in system
//
// Request Type: POST
//
// URL - /v1/api/users/signin
func SignIn(w http.ResponseWriter, r *http.Request) {
	// Set content type returned to JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var user User
	// var employee Employee
	json.NewDecoder(r.Body).Decode(&user)

	if user.Username == "" {
		http.Error(w, "Please enter username", http.StatusBadRequest)
		return
	} else if user.Password == "" {
		http.Error(w, "Please enter password", http.StatusBadRequest)
		return
	}

	db, err := db.OpenSqlx()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	tx := db.MustBegin()

	q := fmt.Sprintf(`
		SELECT id, email, password, is_admin, is_active FROM employee WHERE email = ?
	`)

	empLogin := []struct {
		ID       int    `db:"id" json:"id"`
		Email    string `db:"email" json:"email"`
		Password string `db:"password" json:"password"`
		IsAdmin  bool   `db:"is_admin" json:"is_admin"`
		IsActive int    `db:"is_active" json:"is_active"`
	}{}

	empDataErr := tx.Select(&empLogin, q, user.Username)
	if empDataErr != nil {
		fmt.Println("empData error :: ", empDataErr.Error())
		http.Error(w, "Sql query error while get employee data with email", http.StatusInternalServerError)
		return
	}

	if len(empLogin) == 0 {
		http.Error(w, "There is no records with provided email", http.StatusBadRequest)
		return
	}

	if empLogin[0].IsActive == 0 {
		http.Error(w, "You are now inactivate to the system", http.StatusBadRequest)
		return
	}

	tx.Commit()

	// Compare entered and stored password
	pwdCheckErr := bcrypt.CompareHashAndPassword([]byte(empLogin[0].Password), []byte(user.Password))
	if pwdCheckErr != nil {
		fmt.Println("Password is not same :: " + pwdCheckErr.Error())
		http.Error(w, "Password is wrong", http.StatusBadRequest)
		return
	}

	// setting expiration time
	expirationTime := time.Now().Add(5000000 * time.Minute)
	claims := &Claims{
		ID:       empLogin[0].ID,
		Username: user.Username,
		IsAdmin:  empLogin[0].IsAdmin,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_HRMS_KEY"))) // os.Getenv("SECRET_HRMS_KEY")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error while creating jwt token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SigninResponse{Token: tokenString, EmpID: empLogin[0].ID, IsAdmin: empLogin[0].IsAdmin})

}

// ChangePassword for change employee password
//
// Request Type: PUT
//
// URL - /v1/api/employees/password/change
func ChangePassword(w http.ResponseWriter, r *http.Request) {

	// Set content type returned to JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var passwordChangeStruct PasswordChangeStruct
	_ = json.NewDecoder(r.Body).Decode(&passwordChangeStruct)

	if passwordChangeStruct.CurrentPassword == "" {
		http.Error(w, "Current password is required", http.StatusBadRequest)
		return
	}
	if passwordChangeStruct.NewPassword == "" {
		http.Error(w, "New password is required", http.StatusBadRequest)
		return
	}

	db, err := db.Open()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Get Employee ID from JWT token
	empID, empIDErr := apiutils.GetUserID(r)
	if empIDErr != nil {
		fmt.Println("Get Employee ID error" + empIDErr.Error())
		http.Error(w, "Unable to extract employee id from token", http.StatusInternalServerError)
		return
	}

	var storedPassword []byte
	fmt.Println("Stored Password :: ", storedPassword)
	// get already stored password from database
	row := db.QueryRow(`SELECT 
		password
	FROM employee WHERE id=?`, empID)
	employeeErr := row.Scan(
		&storedPassword)

	if employeeErr != nil {
		fmt.Println("Error in getting stored password from database :: " + employeeErr.Error())
		http.Error(w, "Error in retrieve stored password information", http.StatusInternalServerError)
		return
	}

	pwdCheckErr := bcrypt.CompareHashAndPassword(storedPassword, []byte(passwordChangeStruct.CurrentPassword))
	if pwdCheckErr != nil {
		fmt.Println("Password not match :: " + pwdCheckErr.Error())
		http.Error(w, "Your entered current password is wrong", http.StatusBadRequest)
		return
	}

	empPasswordUpdate, udpateEmpPasswordErr := db.Prepare(`
		UPDATE employee
			SET
				password = ?
			WHERE
				id = ? `)

	if udpateEmpPasswordErr != nil {
		panic(udpateEmpPasswordErr.Error())
	}

	hashPwd, hashPwdErr := bcrypt.GenerateFromPassword([]byte(passwordChangeStruct.NewPassword), bcrypt.MinCost)
	if hashPwdErr != nil {
		fmt.Println("Error in create hashPassword :: " + hashPwdErr.Error())
		http.Error(w, "Error in generate hash password", http.StatusInternalServerError)
		return
	}

	updateResponse, udpateErr := empPasswordUpdate.Exec(
		hashPwd,
		empID,
	)

	fmt.Println(updateResponse)
	if udpateErr != nil {
		fmt.Println("Udpate Employee password error : " + udpateErr.Error())
		http.Error(w, "Error in update password : ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// AddEmployee will insert user to db
//
// Request Type: POST
//
// URL - /v1/api/employees/create
func AddEmployee(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Inserting Employee information...")

	// Set content type returned to JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	employee := &Employee{}
	json.NewDecoder(r.Body).Decode(&employee)

	db, err := db.OpenSqlx()
	if err != nil {
		fmt.Println("Error in database connection")
		panic(err.Error())
	}

	defer db.Close()

	// Generate random password
	randomPwd, pwdGenErr := password.Generate(8, 2, 2, false, false)
	if pwdGenErr != nil {
		fmt.Println("Error in generating random password :: " + pwdGenErr.Error())
		http.Error(w, "Error in generate random password", http.StatusInternalServerError)
		return
	}

	// create encrypted string from random password
	hashPwd, hashPwdErr := bcrypt.GenerateFromPassword([]byte(randomPwd), bcrypt.MinCost)
	if hashPwdErr != nil {
		fmt.Println("Error in create hashPassword :: " + hashPwdErr.Error())
		http.Error(w, "Error in generate hash password", http.StatusInternalServerError)
		return
	}
	employee.Password = hashPwd

	tx := db.MustBegin()
	q := fmt.Sprintf(`
		INSERT INTO employee
			(
				first_name,
				last_name,
				password,
				email,
				personal_email,
				gender,
				dob,
				address,
				phone_number,
				blood_group,
				sick_leave,
				casual_leave,
				floating_leave,
				is_admin,
				is_active
			)
		VALUES (
				:first_name,
				:last_name,
				:password,
				:email,
				:personal_email,
				:gender,
				:dob,
				:address,
				:phone_number,
				:blood_group,
				:sick_leave,
				:casual_leave,
				:floating_leave,
				:is_admin,
				:is_active
			);
	`)

	empInsertRes, empInsertErr := tx.NamedExec(q, employee)
	if empInsertErr != nil {
		fmt.Println(empInsertErr.Error())
		http.Error(w, "Sql query error while insert data", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	fmt.Println("email : , ", employee.Email, ", random password :: ", randomPwd)

	// send hashpassword to employee
	apiutils.SendEmployeePassword(employee.Email, randomPwd)

	empIDValue, empIDErr := empInsertRes.LastInsertId()
	if empIDErr != nil {
		fmt.Println(empIDErr.Error())
		http.Error(w, "Sql query error in get employee id", http.StatusInternalServerError)
		return
	}
	employee.ID = int(empIDValue)
	employee.Password = nil

	// Write OK back to client
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(employee)
}

// AddEmployeeProfile will insert user to db
//
// Request Type: POST
//
// URL - /v1/api/employees/create
func AddEmployeeProfile(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Inserting Employee profile information...")

	// Set content type returned to JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	employeeProfile := &EmployeeProfile{}
	json.NewDecoder(r.Body).Decode(&employeeProfile)

	db, err := db.OpenSqlx()
	if err != nil {
		fmt.Println("Error in database connection")
		panic(err.Error())
	}
	defer db.Close()

	// empID, _ := apiutils.GetUserID(r)
	// employeeProfile.EmployeeID, _ = strconv.Atoi(empID)

	tx := db.MustBegin()
	q := fmt.Sprintf(`
		INSERT INTO employee_profile
			(
				ifsc_code,
				account_number,
				branch_address,
				origin_doc_name,
				origin_doc_status,
				emergency_contact_name,
				emergency_contact_phone_number,
				emergency_contact_relation,
				adhar_card_number,
				pan_card_number,
				profile_pic,
				marital_status,
				employee_status,
				team_id,
				role_id,
				emp_id
			)
		VALUES (
				:ifsc_code,
				:account_number,
				:branch_address,
				:origin_doc_name,
				:origin_doc_status,
				:emergency_contact_name,
				:emergency_contact_phone_number,
				:emergency_contact_relation,
				:adhar_card_number,
				:pan_card_number,
				:profile_pic,
				:marital_status,
				:employee_status,
				:team_id,
				:role_id,
				:emp_id
			);
	`)

	empProfileRes, empProfileErr := tx.NamedExec(q, employeeProfile)
	if empProfileErr != nil {
		fmt.Println(empProfileErr)
		http.Error(w, "Sql query error while insert data", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	empIDValue, empIDErr := empProfileRes.LastInsertId()
	if empIDErr != nil {
		fmt.Println(empIDErr)
		http.Error(w, "Sql query error in get employee id", http.StatusInternalServerError)
		return
	}

	employeeProfile.ID = int(empIDValue)

	// Write OK back to client
	w.WriteHeader(http.StatusOK)
	// w.Write(empIDjson)
	// json.NewEncoder(w).Encode(employeeProfile)
}

// GetAllEmployees will get list of users
//
// Request Type: GET
//
// URL - /v1/api/employees
func GetAllEmployees(w http.ResponseWriter, r *http.Request) {

	db, err := db.OpenSqlx()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	q := fmt.Sprintf(`
		SELECT 
			e.first_name "employee.first_name",
			e.last_name "employee.last_name",
			e.email "employee.email",
			e.personal_email "employee.personal_email",
			e.gender "employee.gender",
			e.dob "employee.dob",
			e.address "employee.address",
			e.phone_number "employee.phone_number",
			e.blood_group "employee.blood_group",
			e.joining_date "employee.joining_date",
			e.sick_leave "employee.sick_leave",
			e.casual_leave "employee.casual_leave",
			e.floating_leave "employee.floating_leave",
			e.is_admin "employee.is_admin",
			e.is_active "employee.is_active",
			ep.id,
			ep.emp_id,
			ep.ifsc_code,
			ep.account_number,
			ep.branch_address,
			ep.origin_doc_name,
			ep.origin_doc_status,
			ep.emergency_contact_name,
			ep.emergency_contact_phone_number,
			ep.emergency_contact_relation,
			ep.adhar_card_number,
			ep.pan_card_number,
			ep.profile_pic,
			ep.marital_status,
			ep.employee_status,
			ep.team_id,
			ep.role_id
		FROM employee e LEFT JOIN employee_profile ep ON e.id = ep.emp_id`)

	tx := db.MustBegin()
	employees := []EmployeeProfile{}
	employeeProfileErr := tx.Select(&employees, q)
	if employeeProfileErr != nil {
		http.Error(w, employeeProfileErr.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, jsonError := json.Marshal(employees)
	if jsonError != nil {
		fmt.Println(jsonError)
		http.Error(w, "Error while parsing employee data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// GetEmployee will get detail of user
//
// Request Type: GET
//
// URL - /v1/api/employees/id
func GetEmployee(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	empID := mux.Vars(r)["id"]

	if empID == "" {
		http.Error(w, "Employee id is missing", http.StatusBadRequest)
		return
	}

	employeeData, employeeDataErr := GetEmployeeAllInfoByID(empID)
	if employeeDataErr != nil {
		http.Error(w, employeeDataErr.Error(), http.StatusInternalServerError)
		return
	}

	if employeeData.ID == 0 {
		http.Error(w, "No records available with given employee id", http.StatusBadRequest)
		return
	}

	jsonResponse, jsonError := json.Marshal(employeeData)
	if jsonError != nil {
		http.Error(w, "Error while map employee data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// UpdateEmployeeData will update user detail
//
// Request Type: PUT
//
// URL - /v1/api/employees/id
func UpdateEmployeeData(w http.ResponseWriter, r *http.Request) {

	empID := mux.Vars(r)["id"]
	fmt.Println(empID)

	if empID == "" {
		http.Error(w, "Employee id is missing", http.StatusBadRequest)
		return
	}

	db, err := db.OpenSqlx()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var employeeProfile EmployeeProfile
	json.NewDecoder(r.Body).Decode(&employeeProfile)

	q := fmt.Sprintf(
		`
		UPDATE employee e JOIN employee_profile ep
		ON e.id = ep.emp_id
		SET 
			e.first_name = :employee.first_name,
			e.last_name = :employee.last_name,
			e.password = :employee.password,
			e.email = :employee.email,
			e.personal_email = :employee.personal_email,
			e.gender = :employee.gender,
			e.dob = :employee.dob,
			e.address = :employee.address,
			e.phone_number = :employee.phone_number,
			e.blood_group = :employee.blood_group,
			e.joining_date = :employee.joining_date,
			e.sick_leave = :employee.sick_leave,
			e.casual_leave = :employee.casual_leave,
			e.floating_leave = :employee.floating_leave,
			e.is_admin = :employee.is_admin,
			e.is_active = :employee.is_active,
			ep.ifsc_code = :ifsc_code,
			ep.account_number = :account_number,
			ep.branch_address = :branch_address,
			ep.origin_doc_name = :origin_doc_name,
			ep.origin_doc_status = :origin_doc_status,
			ep.emergency_contact_name = :emergency_contact_name,
			ep.emergency_contact_phone_number = :emergency_contact_phone_number,
			ep.emergency_contact_relation = :emergency_contact_relation,
			ep.adhar_card_number = :adhar_card_number,
			ep.pan_card_number = :pan_card_number,
			ep.profile_pic = :profile_pic,
			ep.marital_status = :marital_status,
			ep.employee_status = :employee_status,
			ep.team_id = :team_id,
			ep.role_id = :role_id
		WHERE e.id = %[1]s
		`, empID)

	tx := db.MustBegin()
	empUpdate, err := tx.NamedExec(q, employeeProfile)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error while update employee data", http.StatusInternalServerError)
	}

	fmt.Println(empUpdate)
	// employeeInfo, employeeErr := GetEmployeeAllInfoByID(empID)
	// if employeeErr != nil {
	// 	http.Error(w, "error in retriving data", http.StatusInternalServerError)
	// }

	tx.Commit()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(employeeInfo)
}

// DeleteEmployee for activate/deactivate employee
//
// Request Type: DELETE
//
// URL - /v1/api/employees/id
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	isActiveParam := r.URL.Query().Get("is_active")

	empID := mux.Vars(r)["id"]
	if empID == "" {
		http.Error(w, "Employee id is missing", http.StatusBadRequest)
		return
	}

	db, err := db.OpenSqlx()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	q := fmt.Sprintf(`UPDATE employee
			SET 
				is_active = ?
			WHERE
				id = ?`)

	tx := db.MustBegin()
	update, err := tx.Exec(q, isActiveParam, empID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error updateing employee status", http.StatusInternalServerError)
		return
	}
	fmt.Println(update)
	tx.Commit()

	w.WriteHeader(http.StatusOK)
	// w.Write(jsonResponse)
}

// GetRoles will get list of employee roles
//
// Request Type: GET
//
// URL - /v1/api/roles
func GetRoles(w http.ResponseWriter, r *http.Request) {

	var employeeRoles []EmployeeRole

	db, err := db.Open()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	rows, err := db.Query(`SELECT 
		id,
		name
	FROM employee_role`)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error querying employee role table", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	employeeRole := EmployeeRole{}
	for rows.Next() {
		err = rows.Scan(
			&employeeRole.ID,
			&employeeRole.RoleName)

		if err != nil {
			// handle this error
			fmt.Println(err)
			http.Error(w, "Error while retrieve data", http.StatusInternalServerError)
			return
		}
		employeeRoles = append(employeeRoles, employeeRole)
	}

	jsonResponse, jsonError := json.Marshal(employeeRoles)

	if jsonError != nil {
		fmt.Println(jsonError)
		http.Error(w, "Error in parsing retrieved data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// GetTeams will get list of teams list
//
// Request Type: GET
//
// URL - /v1/api/teams
func GetTeams(w http.ResponseWriter, r *http.Request) {

	var teams []Team

	db, err := db.Open()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	rows, err := db.Query(`SELECT 
		id,
		name,
		description
	FROM teams`)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error querying teams table", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	team := Team{}
	for rows.Next() {
		err = rows.Scan(
			&team.ID,
			&team.TeamName,
			&team.Description)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error while retrieve team list data", http.StatusInternalServerError)
			return
		}
		teams = append(teams, team)
	}

	jsonResponse, jsonError := json.Marshal(teams)

	if jsonError != nil {
		fmt.Println(jsonError)
		http.Error(w, "Error in parsing retrieved team list data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
