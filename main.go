package main

import (
  "net/http"
  "html/template"
  "log"
  "time"
  "io/ioutil"
  "encoding/json"
)

type odpInfoJSON struct {
  Protocol string
  Name string
  Key string
  Counter int
}

type odpInfoWeb struct {
  Protocol string
  Name string
  Digits string
  Counter int
}

func getOdpInfoJSONFromFile(fileName string) []odpInfoJSON {
  content, err := ioutil.ReadFile(fileName)
  if err != nil {
    log.Print(err)
    content = []byte("[]")
  }
  var odpInfo []odpInfoJSON
  err = json.Unmarshal(content, &odpInfo)
  if err != nil {
    log.Print(err)
  }
  return odpInfo
}

func generateOtp(odpInfo []odpInfoJSON) []odpInfoWeb {
  oIfo := make([]odpInfoWeb, len(odpInfo))
  for i := 0; i < len(odpInfo) ; i++ {
    oIfo[i].Protocol = odpInfo[i].Protocol
    oIfo[i].Name = odpInfo[i].Name
    oIfo[i].Digits = "035294"
    oIfo[i].Counter = odpInfo[i].Counter
    if odpInfo[i].Protocol == "TOTP" {
      secs := time.Now().Unix()
      oIfo[i].Counter = int((30 - (secs % 30)) * 100 / 30)
    }
  }
  return oIfo
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("Templates/tableOdp.html", "Templates/index.html")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  oIfo := generateOtp(getOdpInfoJSONFromFile("info.json"))
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
