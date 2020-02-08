package leaves

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/kashishsolanki/dt-hrms-golang/db"
	users "github.com/kashishsolanki/dt-hrms-golang/userAPI"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Email holds email name and address info
type Email struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"email,omitempty"`
}

// Claims use
type Claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.StandardClaims
}

// GetEmployeeLeaveInfo to get employee leave info from db
func GetEmployeeLeaveInfo(leaveID string, keyName string) (EmployeeLeave, error) {
	db, err := db.Open()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	q := `SELECT 
			id,
			reason,
			from_date,
			to_date,
			leave_type,
			from_date_day,
			to_date_day,
			total_days,
			leave_status,
			applied_by
		FROM employee_leave WHERE `
	q += keyName + " = ? "

	rows, err := db.Query(q, leaveID)

	if err != nil {
		fmt.Println(err)
		// http.Error(w, "Error querying users table", http.StatusInternalServerError)
		return EmployeeLeave{}, err
	}

	defer rows.Close()

	employeeLeave := EmployeeLeave{}
	for rows.Next() {
		err = rows.Scan(
			&employeeLeave.ID,
			&employeeLeave.Reason,
			&employeeLeave.FromDate,
			&employeeLeave.ToDate,
			&employeeLeave.LeaveType,
			&employeeLeave.FromDateDay,
			&employeeLeave.ToDateDay,
			&employeeLeave.TotalDays,
			&employeeLeave.LeaveStatus,
			&employeeLeave.AppliedBy)

		if err != nil {
			// handle this error
			fmt.Println(err)
			return EmployeeLeave{}, err
		}
	}

	return employeeLeave, nil
}

func getUserEmail(r *http.Request) (string, error) {

	hmacSecretString := os.Getenv("SECRET_HRMS_KEY") // "SECRET_HRMS_KEY" // Value
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims, claims["username"])
		for key, val := range claims {
			fmt.Printf("Key: %v, value: %v\n", key, val)
		}
		return claims["username"].(string), nil
	}
	log.Printf("Invalid JWT Token")
	return "", errors.New("Invalid JWT Token")
}

// getUserId to get user_id from token
func getUserID(r *http.Request) (string, error) {

	hmacSecretString := os.Getenv("SECRET_HRMS_KEY") // "SECRET_HRMS_KEY" // Value
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
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

// isUserAdmin to check user is admin or not
func isUserAdmin(r *http.Request) (bool, error) {

	hmacSecretString := os.Getenv("SECRET_HRMS_KEY") // "SECRET_HRMS_KEY" // Value
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims, claims["is_admin"])
		for key, val := range claims {
			fmt.Printf("Key: %v, value: %v\n", key, val)
		}
		return claims["is_admin"].(bool), nil
	}
	log.Printf("Invalid JWT Token")
	return false, errors.New("Invalid JWT Token")
}

func sendLeaveMail(userInfo users.Employee, employeeLeave EmployeeLeave) {

	fmt.Println("Sending email to HR")
	fmt.Println(userInfo.FirstName + " " + userInfo.LastName)
	fmt.Println(userInfo.Email)
	fmt.Println(os.Getenv("SENDGRID_MAIL_KEY"))
	fmt.Println(os.Getenv("HR_MAIL"))
	from := mail.NewEmail(userInfo.FirstName+" "+userInfo.LastName, userInfo.Email)
	to := mail.NewEmail("HR", os.Getenv("HR_MAIL"))

	subject := "Leave Application"
	switch employeeLeave.LeaveType {
	case 10:
		subject = "Casual Leave Application"
	case 20:
		subject = "Sick Leave Application"
	default:
		subject = "Leave Application"
	}

	plainTextContentHeader := "Hello, "
	plainTextContent := "I want leave from " + employeeLeave.FromDate + " to " + employeeLeave.ToDate + " because of " + employeeLeave.Reason + "."

	if employeeLeave.FromDateDay != 0 {
		halfDay := "second half"
		if employeeLeave.FromDateDay == 1 {
			halfDay = "first half"
		}
		plainTextContent += " I will be available on " + employeeLeave.FromDate + " for " + halfDay + "."
	}
	if employeeLeave.ToDateDay != 0 {
		halfDay := "second half"
		if employeeLeave.ToDateDay == 1 {
			halfDay = "first half"
		}
		plainTextContent += " I will be available on " + employeeLeave.ToDate + " for " + halfDay + "."
	}

	htmlContentString := "<p>Hello, <br> " + plainTextContent + "<br><br> Thanks & Regards,<br>" + userInfo.FirstName + " " + userInfo.LastName + "</p>"
	htmlContent := htmlContentString
	message := mail.NewSingleEmail(from, subject, to, plainTextContentHeader+plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_MAIL_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

func sendManageLeaveMail(userInfo users.Employee, leaveStatus int, employeeLeave EmployeeLeave) {

	fmt.Println("Sending manage leave mail to employee")

	var fromName, fromEmail, toName, toEmail, plainTextContentHeader, plainTextContent string
	subject := "Leave Application Status"

	if leaveStatus == 1 || leaveStatus == 2 {
		fromName = "HR"
		fromEmail = os.Getenv("HR_MAIL")
		toName = userInfo.FirstName + " " + userInfo.LastName
		toEmail = userInfo.Email

		plainTextContentHeader = "Hello " + userInfo.FirstName + " " + userInfo.LastName + " ,"
		plainTextContent = "Your leave has been "
		switch leaveStatus {
		case 1:
			plainTextContent += " approved."
		case 2:
			plainTextContent += " rejected."
		}

	} else if leaveStatus == 3 {
		fromName = userInfo.FirstName + " " + userInfo.LastName
		fromEmail = userInfo.Email
		toName = "HR"
		toEmail = os.Getenv("HR_MAIL")

		plainTextContentHeader = "Hello,"
		plainTextContent = "I want to cancel my leave that I have applied from " + employeeLeave.FromDate + " to " + employeeLeave.ToDate + ". Please consider this."
	}
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)

	htmlContentString := "<p>Hello, <br> " + plainTextContent + "<br><br> Thanks & Regards,<br>HR Executive<br>Digital Trons<br></p>"
	htmlContent := htmlContentString
	message := mail.NewSingleEmail(from, subject, to, plainTextContentHeader+plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_MAIL_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
