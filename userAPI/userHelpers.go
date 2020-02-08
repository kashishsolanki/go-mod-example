package users

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/kashishsolanki/dt-hrms-golang/db"
)

// Claims use
type Claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("SECRET_HRMS_KEY")) // []byte("SECRET_HRMS_KEY")

func checkJWT(r *http.Request) bool {
	fmt.Println("Verifying JWT token")
	if r.Header["Token"] != nil || len(r.Header["Token"]) > 0 {
		fmt.Println("Token given")

		// Initialize a new instance of `Claims`
		claims := &Claims{}

		_, err := jwt.ParseWithClaims(r.Header["Token"][0], claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			return false
		}

		return true
	}
	return false
}

// GetUserInfo to get user detail from db
func GetUserInfo(userValue, keyName string) (Employee, error) {
	db, err := db.Open()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	q := `SELECT 
			id,
			first_name,
			last_name,
			email,
			gender,
			dob,
			address,
			phone_number,
			sick_leave,
			casual_leave,
			is_admin,
			is_active
		FROM employee and `
	q += keyName + " = ? "

	// fmt.Println("Get UserEmail query :: ", q)
	// fmt.Println("Key :: ", keyName, ", value :: ", userValue)

	rows, err := db.Query(q, userValue)

	if err != nil {
		fmt.Println(err)
		// http.Error(w, "Error querying users table", http.StatusInternalServerError)
		return Employee{}, err
	}

	defer rows.Close()

	employee := Employee{}
	for rows.Next() {
		err = rows.Scan(
			&employee.ID,
			&employee.FirstName,
			&employee.LastName,
			&employee.Email,
			&employee.Gender,
			&employee.DOB,
			&employee.Address,
			&employee.PhoneNumber,
			&employee.SickLeave,
			&employee.CasualLeave,
			&employee.IsAdmin,
			&employee.IsActive)

		if err != nil {
			// handle this error
			fmt.Println(err)
			return Employee{}, err
		}
	}

	return employee, nil
}

// UpdateUserLeave will update user leaves count
func UpdateUserLeave(userID int, sickLeave int, casualLeave int) error {

	db, err := db.Open()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	userUpdate, err := db.Prepare(`
		UPDATE employee
			SET
				sick_leave = ?,
				casual_leave = ?
			WHERE
				id = ? `)

	if err != nil {
		panic(err.Error())
	}

	_, userErr := userUpdate.Exec(
		sickLeave,
		casualLeave,
		userID,
	)

	return userErr
}

// GetEmployeeAllInfoByID to get all values of employee from db
func GetEmployeeAllInfoByID(userValue string) (*EmployeeProfile, error) {
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
		FROM employee e JOIN employee_profile ep ON e.id = ep.emp_id WHERE e.id = ?`)

	tx := db.MustBegin()
	var employeeData []*EmployeeProfile
	employeeProfileErr := tx.Select(&employeeData, q, userValue)

	if employeeProfileErr != nil {
		fmt.Println(employeeProfileErr)
		// http.Error(w, "Error querying users table", http.StatusInternalServerError)
		return nil, employeeProfileErr
	}

	tx.Commit()

	if len(employeeData) == 0 {
		return &EmployeeProfile{}, nil
	}
	return employeeData[0], nil
}
