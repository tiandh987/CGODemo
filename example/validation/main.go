package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func main() {
	validate := validator.New()

	var boolTest bool
	err := validate.Var(boolTest, "required")
	if err != nil {
		fmt.Println(err)
	}
	var stringTest string = ""
	err = validate.Var(stringTest, "required")
	if err != nil {
		fmt.Println(err)
	}

	var emailTest string = "test@126.com"
	err = validate.Var(emailTest, "email")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success") // 输出： success。 说明验证成功
	}

	emailTest2 := "test.126.com"
	errs := validate.Var(emailTest2, "required,email")
	if errs != nil {
		fmt.Println(errs) // 输出: Key: "" Error:Field validation for "" failed on the "email" tag。验证失败
	}

	fmt.Println("\r\nEnd!!")

}
