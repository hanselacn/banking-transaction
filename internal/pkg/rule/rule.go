package rule

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

var (
	AlphabetNumericSpaceChar = regexp.MustCompile(`^[\w\s\(\)\-\+\,\.\!\?\@\/]+$`)
	Amount                   = regexp.MustCompile(`^(1000000000000(\.0+)?|([1-9]\d{0,11}(\.\d+)?|0(\.\d+)?))$`)
	FullName                 = regexp.MustCompile(`^[A-Z][a-z]+(?: [A-Z][a-z]+)*$`)
	UserName                 = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]{3,100}$`)
	SpecialCharRegex         = regexp.MustCompile(`[@$!%*?&]`)
	DigitRegex               = regexp.MustCompile("[0-9]")
	LowercaseRegex           = regexp.MustCompile("[a-z]")
	UppercaseRegex           = regexp.MustCompile("[A-Z]")
	LengthRegex              = regexp.MustCompile(`^.{8,20}$`)
	InterestRate             = regexp.MustCompile(`^(0(\.\d+)?|1(\.0+)?)$`)
)

var (
	AlphabetNumericSpaceCharRule = validation.Match(AlphabetNumericSpaceChar).Error(`must be among or combination these characters (a-z, A-Z, 0-9, space, enter, tab, comma(,), dot(.), slash (/), question mark(?), exclamation mark(!), underscore(_), plus and minus(-+))`)
	AmountRule                   = validation.Match(Amount).Error(`invalid amount, must be between 0-1000000000000`)
	FullNameRule                 = validation.Match(FullName).Error(`invalid full name`)
	UserNameRule                 = validation.Match(UserName).Error(`must be among this combination (a-z,A-Z,0-9,dash(-),underscore(_)) with length between 3-100 characters`)
	InterestRateRule             = validation.Match(InterestRate).Error(`interest rate must be between 0-1`)
	SpecialCharRegexRule         = validation.Match(SpecialCharRegex).Error("password must be a combination of alphanumeric + symbols with length between 8-20 characters")
	DigitRegexRule               = validation.Match(SpecialCharRegex).Error("password must be a combination of alphanumeric + symbols with length between 8-20 characters")
	LowercaseRegexRule           = validation.Match(SpecialCharRegex).Error("password must be a combination of alphanumeric + symbols with length between 8-20 characters")
	UppercaseRegexRule           = validation.Match(UppercaseRegex).Error("password must be a combination of alphanumeric + symbols with length between 8-20 characters")
	LengthRegexRule              = validation.Match(LengthRegex).Error("password must be a combination of alphanumeric + symbols with length between 8-20 characters")
)
