package main

import (
  "net/http"
  "net/url"
  "fmt"
  "strings"
  "io/ioutil"
  "encoding/json"
)

func send(phone string, body string) {
  accountSid := ""
  authToken := ""
  urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

  v := url.Values{}
  v.Set("To", phone)
  v.Set("From", "")
  v.Set("Body", body)
  rb := *strings.NewReader(v.Encode())

  // Create client
  client := &http.Client{}

  req, _ := http.NewRequest("POST", urlStr, &rb)
  req.SetBasicAuth(accountSid, authToken)
  req.Header.Add("Accept", "application/json")
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

  // Make request
  resp, _ := client.Do(req)
  if ( resp.StatusCode >= 200 && resp.StatusCode < 300 ) {
    var data map[string]interface{}
    bodyBytes, _ := ioutil.ReadAll(resp.Body)
    err := json.Unmarshal(bodyBytes, &data)
    if ( err == nil ) {
      fmt.Println("Message sent to " + phone)
    }
  } else {
    fmt.Println(resp.Status);
  }
}

func SendCode(code string, phone string) {
  send(phone, "Votre code Requireris est : " + code)
}

func SendAuthCode(code string, phone string) {
  send(phone, "Votre code d'authentification Requireris est : " + code)
}
