package apiutils

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmployeePassword to send random generated password to employee
func SendEmployeePassword(empEmail string, empPwd string) {

	fmt.Println("Sending temporary password employee")

	var fromName, fromEmail, toName, toEmail, plainTextContentHeader, plainTextContent string
	subject := "DT HRMS temporary password"

	fromName = "HR"
	fromEmail = os.Getenv("HR_MAIL")
	toName = empEmail
	toEmail = empEmail

	plainTextContentHeader = "Hello " + empEmail + " ,"
	plainTextContent = "Welcome to DT HRMS. Here is a temporary password with you can login to DT HRMS system." +
		" Please change your password once you login. <br><b>Password</b> : " + empPwd

	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)

	htmlContentString := "<p>Hello, <br> " + plainTextContent + "<br><br> Thanks & Regards,<br>DT HRMS<br>Digital Trons<br></p>"
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
