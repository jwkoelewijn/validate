package validate

import (
	"regexp"
	"strconv"
)

const EmailPattern string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"

var rxEmail = regexp.MustCompile(EmailPattern)

type Validatable interface {
	Validate(v Validator) ValidationViolations
}

type Validator interface {
	MustBePresent(value string) bool
	MustBeEmail(value string) bool
	MustBeIn(value string, collection []interface{}) bool
	ValidateWithFunc(value string, function func(string) bool) bool
}

type ValidationResult struct {
	match   bool
	message string
}

type ValidationViolation struct {
	Field    string
	Messages []string
}

type ValidationViolations []ValidationViolation

type BasicValidator struct {
}

func (v *BasicValidator) MustBePresent(value string) bool {
	return value != ""
}

func (v *BasicValidator) MustBeEmail(value string) bool {
	return rxEmail.MatchString(value)
}

func (v *BasicValidator) MustBeIn(value string, collection []interface{}) bool {
	if len(collection) == 0 {
		return false
	}
	switch collection[0].(type) {
	case string:
		return stringSliceContainsString(collection, value)
	case int:
		val, err := strconv.Atoi(value)
		if err != nil {
			return false
		}
		return intSliceContainsInt(collection, val)
	default:
		return false
	}
}

func stringSliceContainsString(collection []interface{}, value string) bool {
	for _, elem := range collection {
		el := elem.(string)
		if el == value {
			return true
		}
	}
	return false
}

func intSliceContainsInt(collection []interface{}, value int) bool {
	for _, elem := range collection {
		el := elem.(int)
		if el == value {
			return true
		}
	}
	return false
}

func (v *BasicValidator) ValidateWithFunc(value string, function func(string) bool) bool {
	return function(value)
}
