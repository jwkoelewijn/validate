package validate

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
)

const EmailPattern string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"

var rxEmail = regexp.MustCompile(EmailPattern)

type Validatable interface {
	Validate(v Validator) ValidationViolations
}

type Validator interface {
	ClearViolations()
	Violations() ValidationViolations

	ValidatePresent(target interface{}, field string) bool
	ValidateEmail(target interface{}, field string) bool
	ValidateInclusion(target interface{}, field string, collection []interface{}) bool
	ValidateWithFunction(target interface{}, field string, function func(string) bool) bool
}

type ValidationViolations map[string][]string

type BasicValidator struct {
	violations ValidationViolations
}

func (v *BasicValidator) ClearViolations() {
	v.violations = make(map[string][]string)
}

func (v *BasicValidator) Violations() ValidationViolations {
	return v.violations
}

func (v *BasicValidator) ValidatePresent(target interface{}, field string) bool {
	value, err := v.getValueForTargetField(target, field)
	if err != nil {
		v.appendViolation(field, fmt.Sprintf("could not find field '%s'", field))
		return false
	}
	if !v.mustBePresent(value) {
		v.appendViolation(field, fmt.Sprintf("expected %s to be present", field))
		return false
	}
	return true
}

func (v *BasicValidator) ValidateEmail(target interface{}, field string) bool {
	value, err := v.getValueForTargetField(target, field)
	if err != nil {
		v.appendViolation(field, fmt.Sprintf("could not find field '%s'", field))
		return false
	}

	if !v.mustBeEmail(value) {
		v.appendViolation(field, fmt.Sprintf("expected '%s' to be an email address", value))
		return false
	}
	return true
}

func (v *BasicValidator) ValidateInclusion(target interface{}, field string, collection []interface{}) bool {
	value, err := v.getValueForTargetField(target, field)
	if err != nil {
		v.appendViolation(field, fmt.Sprintf("could not find field '%s'", field))
		return false
	}

	if !v.mustBeIn(value, collection) {
		v.appendViolation(field, fmt.Sprintf("expected '%+v' to include '%s'", collection, value))
		return false
	}
	return true
}

func (v *BasicValidator) ValidateWithFunction(target interface{}, field string, function func(string) bool) bool {
	value, err := v.getValueForTargetField(target, field)
	if err != nil {
		v.appendViolation(field, fmt.Sprintf("could not find field '%s'", field))
		return false
	}

	if !v.validateValueWithFunc(value, function) {
		v.appendViolation(field, fmt.Sprintf("expected provided function to evaluate to true when applied to input '%s'", value))
		return false
	}
	return true
}
func (v *BasicValidator) mustBePresent(value string) bool {
	return value != ""
}

func (v *BasicValidator) mustBeEmail(value string) bool {
	return rxEmail.MatchString(value)
}

func (v *BasicValidator) mustBeIn(value string, collection []interface{}) bool {
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

func (v *BasicValidator) validateValueWithFunc(value string, function func(string) bool) bool {
	return function(value)
}

func (v *BasicValidator) appendViolation(field, message string) {
	var found bool
	if _, found = v.violations[field]; found == false {
		v.violations[field] = make([]string, 1, 1)
	}

	v.violations[field] = append(v.violations[field], message)
	log.Println(fmt.Sprintf("Violations: %+v", v.violations))
}

func (v *BasicValidator) getValueForTargetField(target interface{}, field string) (string, error) {
	value := reflect.ValueOf(target).Elem()
	fieldValue := value.FieldByName(field)
	if isZeroOfUnderlyingType(fieldValue) {
		return "", fmt.Errorf("Could not get a value for field '%s'", field)
	}
	res := fieldValue.String()
	return res, nil
}

func isZeroOfUnderlyingType(x interface{}) bool {
	return x == nil || x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
