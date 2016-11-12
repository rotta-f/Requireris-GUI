package main

import (
  "net/http"
  "html/template"
  "strconv"
)

func handleAdd(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("Templates/utils.html", "Templates/add.html")
  err = t.ExecuteTemplate(w, "content", &otpInfoJSON{Protocol : "", Service : "", Key : ""})
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func handleAddKey(w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()
  if err != nil {
      // Handle error here via logging and then return
  }
  protocolsAvailable := map[string]bool{"HOTP" : true, "TOTP" : true}
  digitsAvailable := map[int]bool{6: true, 7: true, 8: true, 9: true, 10: true}
  protocol := r.PostFormValue("protocol")
  service := r.PostFormValue("service")
  key := r.PostFormValue("key")
  digits, _ := strconv.Atoi(r.PostFormValue("digits"))
  if !digitsAvailable[digits] {
    // digits is not in digitsAvailable
    digits = 6
  }
  goodProtocol := protocolsAvailable[protocol]
  goodService := isServiceAvailable(service)
  if !goodProtocol || !goodService {
    t, err := template.ParseFiles("Templates/utils.html", "Templates/add.html")
    if goodProtocol {
      protocol = ""
    }
    if goodService {
      service = ""
    }
    err = t.ExecuteTemplate(w, "content", &otpInfoJSON{Protocol : protocol, Service : service, Key : key})
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    return
  }
  TabOtpInfo = append(TabOtpInfo, otpInfoJSON{Protocol : protocol, Service : service, Key : key, Counter : 0, Digits: digits})
  setOtpInfoJSONFromFile(FileBDD, TabOtpInfo)
  http.Redirect(w, r, "/", http.StatusFound)
}
