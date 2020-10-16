package apierrors

import "fmt"

func ExampleErrors_String() {
	fmt.Println(New())

	// Output: Errors(Code=200, SystemError='<nil>', UserError='<nil>')
}

func ExamlpleNew() {
	fmt.Println(New())
	fmt.Println(New().Occurred())

	// Output: Errors(Code=200, SystemError='<nil>', UserError='<nil>')
	// false
}

func ExampleErrors_Occurred() {
	err := New()
	fmt.Println(err.Occurred())

	err.SetSystemError("test")
	fmt.Println(err.Occurred())

	err.SetUserError("test")
	fmt.Println(err.Occurred())

	err.SystemError = nil
	fmt.Println(err.Occurred())

	// Output: false
	// true
	// true
	// true
}
