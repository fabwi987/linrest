package models

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendRecommendationMail(toadress string, name string) error {

	from := mail.NewEmail("Linvestor Messaging", "recommendations@linvestor.se")
	subject := "Rekommendation noterad"
	to := mail.NewEmail(name, toadress)
	content := mail.NewContent("text/plain", "Hej "+name+"! Tack för din aktierekommendation. Utvecklingen följer du på: www.bit.do/linvestor")
	m := mail.NewV3MailInit(from, subject, to, content)

	//Removed for GitHub
	request := sendgrid.GetRequest("", "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return nil
}
