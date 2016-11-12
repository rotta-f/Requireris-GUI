package main

import (
  "net/http"
  "html/template"
  "log"
  "time"
  "io/ioutil"
  "encoding/json"

  "github.com/rotta-f/Requireris"
  "net/url"
)

type otpInfoJSON struct {
  Protocol string
  Service string
  Key string
  Counter uint64
  Digits int
}

type otpInfoWeb struct {
  Protocol string
  Service string
  Digits string
  Counter uint64
}

var TabOtpInfo []otpInfoJSON
const FileBDD string = "info.json"

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

func setOtpInfoJSONFromFile(fileName string, otpInfo []otpInfoJSON) {
  content, err := json.Marshal(otpInfo)
  if err != nil {
    log.Print(err)
    return
  }
  err = ioutil.WriteFile(fileName, content, 0644)
  if err != nil {
    log.Print(err)
  }
}

func generateOtp(otpInfo []otpInfoJSON) []otpInfoWeb {
  oIfo := make([]otpInfoWeb, len(otpInfo))
  for i := 0; i < len(otpInfo) ; i++ {
    oIfo[i].Protocol = otpInfo[i].Protocol
    oIfo[i].Service = otpInfo[i].Service
    oIfo[i].Counter = otpInfo[i].Counter
    oIfo[i].Digits = "Not Defined"
    otp := Requireris.Init(otpInfo[i].Key, otpInfo[i].Digits)
    switch otpInfo[i].Protocol {
    case "TOTP":
      secs := time.Now().Unix()
      oIfo[i].Digits = otp.TOTP()
      oIfo[i].Counter = uint64(30 - (secs % 30)) * 100 / 30
    case "HOTP":
      oIfo[i].Digits = otp.HOTP(otpInfo[i].Counter)
      oIfo[i].Counter = otpInfo[i].Counter
    }
  }
  return oIfo
}

func isServiceAvailable(service string) bool {
  for i := 0; i < len(TabOtpInfo); i++ {
    if service == TabOtpInfo[i].Service {
      return false
    }
  }
  return true
}

func getServiceIndex(service string) int {
  for i := 0; i < len(TabOtpInfo); i++ {
    if service == TabOtpInfo[i].Service {
      return i
    }
  }
  return -1
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("Templates/utils.html", "Templates/index.html")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  err = t.ExecuteTemplate(w, "content", nil)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func handleGetOtp(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("Templates/tableOtp.html")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  oIfo := generateOtp(TabOtpInfo)
  err = t.ExecuteTemplate(w, "tableOtp", oIfo)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func handleIncrementHOTP(w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()
  if err != nil {
      // Handle error here via logging and then return
  }
  service := r.PostFormValue("Service")
  for i := 0; i < len(TabOtpInfo); i++ {
    if service == TabOtpInfo[i].Service {
      TabOtpInfo[i].Counter += 1
    }
  }
  setOtpInfoJSONFromFile(FileBDD, TabOtpInfo)
}
func handleSmsSend(w http.ResponseWriter, r *http.Request) {
  m, _ :=url.ParseQuery(r.URL.RawQuery)
  SendCode(m["code"][0], m["phone"][0])
}

func main() {
  var err error

  TabOtpInfo = getOtpInfoJSONFromFile(FileBDD)
  http.HandleFunc("/", handleIndex)
  http.HandleFunc("/add", handleAdd)
  http.HandleFunc("/addKey", handleAddKey)
  http.HandleFunc("/getOtp", handleGetOtp)
  http.HandleFunc("/incrementHOTP", handleIncrementHOTP)
  http.HandleFunc("/del", handleDelKey)
  http.HandleFunc("/checkDel", handleCheckDelKey)
  http.HandleFunc("/sms", handleSmsSend)
  log.Println("Ready to listen and serve.")
  err = http.ListenAndServe(":8080", nil)
  if err != nil {
    log.Fatal(err)
  }
}
