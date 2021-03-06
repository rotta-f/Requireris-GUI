package main

import (
  "log"
  "fmt"
  "io/ioutil"
  "crypto/hmac"
  "crypto/sha256"
  "bytes"
  "crypto/rand"

  "github.com/rotta-f/Requireris"
  "golang.org/x/crypto/ssh/terminal"
  "os"
)

var PasswordUser []byte
var UserPhone string

const
(
  authByPhone = true
)

func generateNewCode() string {
  base32key := make([]byte, 12)
  _, err := rand.Read(base32key)
  if err != nil {
    log.Fatal(err)
  }
  encodeBase32 := "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
  for i := 0; i < 12; i++ {
    base32key[i] = encodeBase32[base32key[i] % 32]
  }
  otp := Requireris.Init(string(base32key), 6)
  return otp.TOTP()
}

func confirmByPhone(phoneNumber string) bool {
  if authByPhone == false {
    return true
  }
  code := generateNewCode()
  SendAuthCode(code, phoneNumber)

  var codeOnPhone string
  fmt.Print("\nConfirm with the code received by phone\nCode: ")
  _, err := fmt.Scanln(&codeOnPhone)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println()

  return codeOnPhone == code
}

func getHash() ([]byte, bool) {
  content, err := ioutil.ReadFile(".passwd")
  if err != nil {
    return []byte(""), false
  }
  return content, true
}

func initFirstConnection() ([]byte, string) {
  fmt.Print("New Password: ")
  passwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
  fmt.Println("")
  if err != nil {
    log.Fatal(err)
  }
  fmt.Print("Confirm Password: ")
  passwdConfirmed, err := terminal.ReadPassword(int(os.Stdin.Fd()))
  fmt.Println("")
  if err != nil {
    log.Fatal(err)
  }
  if bytes.Compare(passwdConfirmed, passwd) != 0 {
    log.Fatalln("Passwords differ")
  }
  var phoneNumber string
  fmt.Print("Phone number (ex: +33XXXXXXXXX): ")
  _, err = fmt.Scanln(&phoneNumber)
  fmt.Println()
  if err != nil {
    log.Fatal(err)
  }

  if confirmByPhone(phoneNumber) != true {
    log.Fatal("Error, bad code")
  }

  mac := hmac.New(sha256.New, passwd)
  mac.Write([]byte(phoneNumber))
  err = ioutil.WriteFile(".passwd", mac.Sum(nil), 0644)
  return passwd, phoneNumber
}

func newConnection(oldHash []byte) ([]byte, string) {
  fmt.Print("Password: ")
  passwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
  fmt.Println("")
  if err != nil {
    log.Fatal(err)
  }
  var phoneNumber string
  fmt.Print("Phone number (ex: +33XXXXXXXXX): ")
  _, err = fmt.Scanln(&phoneNumber)
  fmt.Println()
  if err != nil {
    log.Fatal(err)
  }
  mac := hmac.New(sha256.New, passwd)
  mac.Write([]byte(phoneNumber))
  if bytes.Compare(mac.Sum(nil), oldHash) != 0 {
    log.Fatal("Authentication failure")
  }
  if confirmByPhone(phoneNumber) != true {
    log.Fatal("Authentication failure")
  }
  return passwd, phoneNumber
}

func connect() ([]byte, bool) {
  oldHash, ok := getHash()
  if ok != true {
    PasswordUser, UserPhone = initFirstConnection()
  } else {
    PasswordUser, UserPhone = newConnection(oldHash)
  }
  return []byte(""), true
}
