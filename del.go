package main

import (
  "net/http"
  "html/template"
  "log"
)

type argDel struct {
  Err otpInfoJSON
  ListOfService []string
}

func getArgDel(service string) argDel {
  listOfService := make([]string, len(TabOtpInfo))
  for i := 0; i < len(TabOtpInfo); i++ {
    listOfService[i] = TabOtpInfo[i].Service
  }
  return argDel{Err : otpInfoJSON{Service : ""}, ListOfService : listOfService}
}

func handleDelKey(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("Templates/utils.html", "Templates/del.html")
  err = t.ExecuteTemplate(w, "content", getArgDel(""))
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func handleCheckDelKey(w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()
  if err != nil {
    //Error
  }
  service := r.PostFormValue("service")
  goodService := isServiceAvailable(service)
  if goodService {
    t, err := template.ParseFiles("Templates/utils.html", "Templates/del.html")
    err = t.ExecuteTemplate(w, "content", getArgDel(service))
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    return
  }
  log.Println("DelService", service)
  idx := getServiceIndex(service)
  TabOtpInfo = append(TabOtpInfo[:idx], TabOtpInfo[idx + 1:]...)
  setOtpInfoJSONFromFile(FileBDD, TabOtpInfo)
  http.Redirect(w, r, "/", http.StatusFound)
}
