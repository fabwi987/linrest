package models

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func SendRecommendationText(PhoneTo string, name string) error {

	//Removed for GitHub
	accountSid := "ACda82d1d2837367826c0bd6e41d15f30a"
	authToken := "3fb8953dfe3c0665d74c41623aec0eb3"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	v := url.Values{}
	v.Set("To", PhoneTo)
	//Removed from GitHub
	v.Set("From", "+46769437171")
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
