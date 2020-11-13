package main

import(
  "github.com/AlecAivazis/survey/v2"
  "fmt"
)

func prompts(a *string, b *string){
  //user prompt for username
  prompt1 := &survey.Input{
      Message: "Dropbox username",
  }
  survey.AskOne(prompt1, a)
  //user prompt for password
  prompt2 := &survey.Password{
      Message: "Dropbox password",
  }
  survey.AskOne(prompt2, b)

  }

func main () {

//user string variable
a := ""
b := ""

prompts(&a, &b)

fmt.Println(a)
fmt.Println(b)

}
