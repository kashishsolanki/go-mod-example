package users

import "gopkg.in/guregu/null.v3"

// EmployeeRole struct for define fields
type EmployeeRole struct {
	ID          int    `json:"id"`
	RoleName    string `json:"role_name"`
	Description string `json:"role_description"`
	Position    string `json:"employee_position"`
}

// Team struct for define fields
type Team struct {
	ID          int    `json:"id"`
	TeamName    string `json:"team_name"`
	Description string `json:"team_description"`
}

// EmployeeProfile struct for define fields
type EmployeeProfile struct {
	ID                          int    `db:"id" json:"id"`
	IFSCCode                    string `db:"ifsc_code" json:"ifsc_code"`
	AccountNo                   string `db:"account_number" json:"account_number"`
	BranchAddress               string `db:"branch_address" json:"branch_address"`
	OrigDocName                 string `db:"origin_doc_name" json:"origin_doc_name"`
	OrigDocStatus               string `db:"origin_doc_status" json:"origin_doc_status"` // Submitted/Returned/None/Pending
	EmergencyContactName        string `db:"emergency_contact_name" json:"emergency_contact_name"`
	EmergencyContactPhoneNumber string `db:"emergency_contact_phone_number" json:"emergency_contact_phone_number"`
	EmergencyContactRelation    string `db:"emergency_contact_relation" json:"emergency_contact_relation"`
	EmpAdharCard                string `db:"adhar_card_number" json:"adhar_card_number"`
	EmpPanCard                  string `db:"pan_card_number" json:"pan_card_number"`
	ProfilePhoto                string `db:"profile_pic" json:"profile_pic"`
	MaritalStatus               string `db:"marital_status" json:"marital_status"`   // Married/Unmarried
	EmployeeStatus              string `db:"employee_status" json:"employee_status"` // Probation/Permanent
	TeamID                      int    `db:"team_id" json:"team_id"`                 // Flourish/React/Management
	RoleID                      int    `db:"role_id" json:"role_id"`                 // Flourish/React/Management
	EmployeeID                  int    `db:"emp_id" json:"emp_id"`
	Employee                    `db:"employee"`
}

// Employee struct for define fields
type Employee struct {
	ID            int         `db:"id" json:"id"`
	FirstName     string      `db:"first_name" json:"first_name"`
	LastName      string      `db:"last_name" json:"last_name"`
	Password      []byte      `db:"password" json:"password"`
	Email         string      `db:"email" json:"email"`
	PersonalEmail string      `db:"personal_email" json:"personal_email"`
	Gender        string      `db:"gender" json:"gender"`
	DOB           string      `db:"dob" json:"dob"`
	Address       string      `db:"address" json:"address"`
	PhoneNumber   string      `db:"phone_number" json:"phone_number"`
	BloodGroup    string      `db:"blood_group" json:"blood_group"`
	JoiningDate   null.String `db:"joining_date" json:"joining_date"`
	SickLeave     int         `db:"sick_leave" json:"sick_leave"`
	CasualLeave   int         `db:"casual_leave" json:"casual_leave"`
	FloatingLeave int         `db:"floating_leave" json:"floating_leave"`
	IsActive      bool        `db:"is_active" json:"is_active"`
	IsAdmin       bool        `db:"is_admin" json:"is_admin"`
	// OriginalDoc     OriginalDoc            `json:"original_doc"`
	// BankDetail      BankDetail             `json:"bank_details"`
	// DocumentDetails int                    `json:"document_details"`
}
