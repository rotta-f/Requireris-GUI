package main

import (
  "net/http"
  "html/template"
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
  protocol := r.PostFormValue("protocol")
  service := r.PostFormValue("service")
  key := r.PostFormValue("key")
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
  TabOtpInfo = append(TabOtpInfo, otpInfoJSON{Protocol : protocol, Service : service, Key : key, Counter : 0})
  setOtpInfoJSONFromFile(FileBDD, TabOtpInfo)
  http.Redirect(w, r, "/", http.StatusFound)
}
