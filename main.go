package main

import (
  "net/http"
  "html/template"
  "log"
  "time"
  "io/ioutil"
  "encoding/json"
)

type otpInfoJSON struct {
  Protocol string
  Name string
  Key string
  Counter int
}

type otpInfoWeb struct {
  Protocol string
  Name string
  Digits string
  Counter int
}

func getOtpInfoJSONFromFile(fileName string) []otpInfoJSON {
  content, err := ioutil.ReadFile(fileName)
  if err != nil {
    log.Print(err)
    content = []byte("[]")
  }
  var otpInfo []otpInfoJSON
  err = json.Unmarshal(content, &otpInfo)
  if err != nil {
    log.Print(err)
  }
  return otpInfo
}

func generateOtp(otpInfo []otpInfoJSON) []otpInfoWeb {
  oIfo := make([]otpInfoWeb, len(otpInfo))
  for i := 0; i < len(otpInfo) ; i++ {
    oIfo[i].Protocol = otpInfo[i].Protocol
    oIfo[i].Name = otpInfo[i].Name
    oIfo[i].Digits = "035294"
    oIfo[i].Counter = otpInfo[i].Counter
    if otpInfo[i].Protocol == "TOTP" {
      secs := time.Now().Unix()
      oIfo[i].Counter = int((30 - (secs % 30)) * 100 / 30)
    }
  }
  return oIfo
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("Templates/tableOtp.html", "Templates/index.html")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  oIfo := generateOtp(getOtpInfoJSONFromFile("info.json"))
  err = t.ExecuteTemplate(w, "content", oIfo)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func main() {
  var err error

  http.HandleFunc("/", handleIndex)
  log.Println("Ready to listen and serve.")
  err = http.ListenAndServe(":8080", nil)
  if err != nil {
    log.Fatal(err)
  }
}
