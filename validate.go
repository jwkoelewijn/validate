package validate

import (
	"fmt"
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
	ValidateEmail(target interface{}, field string, allowEmptyInput bool) bool
	ValidateInclusion(target interface{}, field string, collection []interface{}, allowEmptyInput bool) bool
	ValidateWithFunction(target interface{}, field, message string, allowEmptyInput bool, function func(string) bool) bool
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
		v.appendViolation(field, "should be present")
		return false
	}
	return true
}

func (v *BasicValidator) ValidateEmail(target interface{}, field string, allowEmpty bool) bool {
	value, err := v.getValueForTargetField(target, field)
	if err != nil {
		v.appendViolation(field, fmt.Sprintf("could not find field '%s'", field))
		return false
	}

	if !v.mustBeEmail(value, allowEmpty) {
		v.appendViolation(field, "is not a valid email address")
		return false
	}
	return true
}

func (v *BasicValidator) ValidateInclusion(target interface{}, field string, collection []interface{}, allowEmpty bool) bool {
	value, err := v.getValueForTargetField(target, field)
	if err != nil {
		v.appendViolation(field, fmt.Sprintf("could not find field '%s'", field))
		return false
	}

	if !v.mustBeIn(value, collection, allowEmpty) {
		v.appendViolation(field, fmt.Sprintf("should be in '%+v'", collection))
		return false
	}
	return true
}

func (v *BasicValidator) ValidateWithFunction(target interface{}, field, message string, allowEmpty bool, function func(string) bool) bool {
	value, err := v.getValueForTargetField(target, field)
	if err != nil {
		v.appendViolation(field, fmt.Sprintf("could not find field '%s'", field))
		return false
	}

	if !v.validateValueWithFunc(value, allowEmpty, function) {
		v.appendViolation(field, message)
		return false
	}
	return true
}
func (v *BasicValidator) mustBePresent(value string) bool {
	return value != ""
}

func (v *BasicValidator) mustBeEmail(value string, allowEmpty bool) bool {
	if allowEmpty && value == "" {
		return true
	}
	return rxEmail.MatchString(value)
}

func (v *BasicValidator) mustBeIn(value string, collection []interface{}, allowEmpty bool) bool {
	if allowEmpty && value == "" {
		return true
	}

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

func (v *BasicValidator) validateValueWithFunc(value string, allowEmpty bool, function func(string) bool) bool {
	if allowEmpty && value == "" {
		return true
	}
	return function(value)
}

func (v *BasicValidator) appendViolation(field, message string) {
	var found bool
	if _, found = v.violations[field]; found == false {
		var violationList []string
		v.violations[field] = violationList
	}

	v.violations[field] = append(v.violations[field], message)
}

func (v *BasicValidator) getValueForTargetField(target interface{}, field string) (string, error) {
	value := reflect.ValueOf(target).Elem()
	fieldValue := value.FieldByName(field)

	if isZeroOfUnderlyingType(fieldValue) {
		return "", fmt.Errorf("Could not get a value for field '%s'", field)
	}

	res := fieldValueToString(fieldValue)
	return res, nil
}

func fieldValueToString(fieldValue reflect.Value) string {
	if isNil(fieldValue) {
		return ""
	} else if fieldValue.Kind() == reflect.Ptr {
		return fieldValueToString(fieldValue.Elem())
	} else if fieldValue.Kind() == reflect.Int {
		return strconv.FormatInt(fieldValue.Int(), 10)
	} else {
		return fieldValue.String()
	}
}

func isNil(a reflect.Value) bool {
	defer func() {
		recover()
	}()
	return a.IsNil()
}

func isZeroOfUnderlyingType(x interface{}) bool {
	return x == nil || x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
