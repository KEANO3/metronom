package password_test

import (
	"fmt"
	"metronom/password"
)

func ExampleFill() {
	password.Seed(1337)
	fmt.Print(password.AutoComplete(password.Parameter{
		Numbers:      -1,
		SpecialChars: -1,
	}))
	// Output: {8 0 0 2 true}
}

func ExampleValidate() {
	r := password.DefaultRequest
	fmt.Println(password.Validate(r))
	r.SpecialChars = 13
	fmt.Println(password.Validate(r))
	r.SpecialChars = 9
	fmt.Println(password.Validate(r))

	// Output:
	// <nil>
	// sum of requested special chars and numbers exceeds maximum length
	// sum of requested special chars and numbers exceeds minimum length
}

func ExampleGenerate() {
	password.Seed(1337)
	r := password.DefaultRequest
	r.MaxLength = 0

	fmt.Println(password.Generate(r))
	r.MaxLength = 20
	fmt.Println(password.Generate(r))
	r.SpecialChars = 5
	fmt.Println(password.Generate(r))
	password.PoolSpecial = []rune{'a'}
	fmt.Println(password.Generate(r))

	// Output:
	// XVY$AQ-J
	// fFDEXHqgamI=
	// O}M²m[Rf5eY/;
	// LaaaINWaRWpQaEV2n
}
