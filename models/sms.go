package models

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func SendRecommendationText(PhoneTo string, name string) error {

	//Removed for GitHub
	accountSid := ""
	authToken := ""
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	v := url.Values{}
	v.Set("To", PhoneTo)
	//Removed from GitHub
	v.Set("From", "")
	v.Set("Body", "Hej "+name+"! Tack för din aktierekommendation. Utvecklingen följer du på: www.bit.do/linvestor")
	rb := *strings.NewReader(v.Encode())

	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	fmt.Println(resp.Status)

	return nil

}