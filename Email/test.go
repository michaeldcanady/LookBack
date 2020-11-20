package email

import(
  //"crypto/tls"
  "fmt"

  gomail"gopkg.in/mail.v2"
)

func SendEmail(email string, key string){
  var Sender = "SENDER EMAIL"
  var server = "SENDER SERVER"
  var password = "SENDER PASSWORD"


  m := gomail.NewMessage()

  // set E-mail sender
  m.SetHeader("From",Sender)

  // set Email receivers
  m.SetHeader("To", email)

  // set Email subject
  m.SetHeader("Subject","Encrypt test")

  // set email body
  m.SetBody("text/plain",key)

  // settings for smtp server
  d := gomail.NewDialer(server,587,Sender,password)

  d.StartTLSPolicy = gomail.MandatoryStartTLS

  if err := d.DialAndSend(m); err != nil{
    fmt.Println(err)
    panic(err)
  }

  return
}
