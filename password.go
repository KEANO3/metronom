package password //import "metronom/password"

import (
	"errors"
	"math/rand"
	"time"
)

var (
	//ErrInvalidParametersMaxLength indicates an invalid and non processable request due to mismatching parameters
	ErrInvalidParametersMaxLength = errors.New("sum of requested special chars and numbers exceeds maximum length")

	//ErrParametersExceedMinLength indicates a faulty but processable request
	ErrParametersExceedMinLength = errors.New("sum of requested special chars and numbers exceeds minimum length")

	//DefaultMin is used as MinLength if a Request doesn't determine it
	DefaultMin = 8

	//DefaultMax is used as MaxLength if a Request doesn't determine it
	DefaultMax = 12

	//DefaultRequest can be used if no parameters were given
	DefaultRequest = Parameter{
		MinLength:    DefaultMin,
		MaxLength:    DefaultMax,
		Numbers:      -1,
		SpecialChars: -1,
	}

	//PoolNumbers contains all numbers which may be used in Generate
	PoolNumbers = []rune("0123456789")

	//PoolSpecial contains all special characters which may be used in Generate
	PoolSpecial = []rune(`!\"§$%&/()=?` + "´`" + `*+'#<>|-_.,;:²³{[]}~^°`)

	//PoolUpper contains all capitals which may be used in Generate
	PoolUpper = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	//PoolLower contains all minuscules which may be used in Generate
	PoolLower = []rune("abcdefghijklmnopqrstuvwxyz")

	rnd *rand.Rand
)

func init() {
	Seed(time.Now().UnixNano())
}

//Seed (re-)seeds the used PRNG
func Seed(i int64) {
	rnd = rand.New(rand.NewSource(i))
}

//Parameter contains all user parameters to generate a password
type Parameter struct {
	MinLength    int
	MaxLength    int
	Numbers      int
	SpecialChars int
	filled       bool
}

//Generate returns a password according to the given parameter. Parameters should be validated beforehand.
func Generate(p Parameter) string {
	if !p.filled {
		p = AutoComplete(p)
	}

	length := p.MinLength

	if p.MaxLength > p.MinLength {
		length = length + rnd.Intn(p.MaxLength-p.MinLength)
	}

	if length < 0 {
		length = 1
	}

	var runes []rune

	for s := 0; s < p.SpecialChars; s++ {
		runes = append(runes, randomPick(PoolSpecial))
	}

	for n := 0; n < p.Numbers; n++ {
		runes = append(runes, randomPick(PoolNumbers))
	}

	for i := 0; i < (length - (p.SpecialChars + p.Numbers)); i++ {
		if rnd.Int()%2 == 0 {
			runes = append(runes, randomPick(PoolUpper))
		} else {
			runes = append(runes, randomPick(PoolLower))
		}
	}

	rnd.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	return string(runes)
}

func randomPick(p []rune) rune {
	lenp := len(p)

	if lenp < 1 {
		return '\U0001f4a9'
	}

	if lenp == 1 {
		return p[0]
	}

	return p[rnd.Intn(lenp-1)]
}

//Validate returns an error if given noncoherent parameters
func Validate(p Parameter) error {
	if p.Numbers < 0 {
		p.Numbers = 0
	}

	if p.SpecialChars < 0 {
		p.SpecialChars = 0
	}

	switch {
	case p.SpecialChars+p.Numbers > p.MaxLength:
		return ErrInvalidParametersMaxLength
	case p.SpecialChars+p.Numbers > p.MinLength:
		return ErrParametersExceedMinLength
	}

	return nil
}

//AutoComplete fills up undetermined parameters
func AutoComplete(p Parameter) Parameter {
	if p.MaxLength < 0 {
		p.MaxLength = DefaultMax
	}

	if p.MinLength < 1 {
		p.MinLength = DefaultMin
	}

	if p.Numbers < 0 {
		n := 0
		if p.SpecialChars < 0 {
			n = p.MinLength / 3
		} else {
			n = p.MinLength - p.SpecialChars
		}

		if n < 1 {
			n = 1
		}

		p.Numbers = rnd.Intn(n)
	}

	if p.SpecialChars < 0 {
		n := (p.MinLength - p.Numbers) / 2

		if n < 1 {
			n = 1
		}

		p.SpecialChars = rnd.Intn(n)
	}

	p.filled = true

	return p
}
