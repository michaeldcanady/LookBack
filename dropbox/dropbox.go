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

//Validation
//  q := &survey.Question{
//      Prompt: &survey.Input{Message: "Hello world validation"},
//      Validate: func (val interface{}) error {
//          // since we are validating an Input, the assertion will always succeed
//          if str, ok := val.(string) ; !ok || len(str) > 10 {
//              return errors.New("This response cannot be longer than 10 characters.")
//          }
//  	return nil
//      },
//  }

func main () {

//user string variable
a := ""
b := ""

prompts(&a, &b)

fmt.Println(a)
fmt.Println(b)

}
