package validate

import "testing"

type validationTest struct {
	title    string
	input    string
	expected bool
	message  string
}

func TestMustBePresent(t *testing.T) {
	val := &BasicValidator{}
	tests := []validationTest{
		{
			title:    "Empty input is not present",
			input:    "",
			expected: false,
			message:  "Expected empty string to not be present",
		},
		{
			title:    "minimum input",
			input:    ".",
			expected: true,
			message:  "Expected '.' to be present",
		},
		{
			title:    "Some string is present indeed",
			input:    "some string",
			expected: true,
			message:  "Expected 'some string' to be present",
		},
	}

	for _, test := range tests {
		if res := val.mustBePresent(test.input); res != test.expected {
			t.Errorf("Test '%s' failed: %s (Expected: %v, Found: %v)", test.title, test.message, test.expected, res)
		}
	}
}

func TestMustBeEmail(t *testing.T) {
	val := &BasicValidator{}
	tests := []validationTest{
		{
			title:    "Empty input is not an email",
			input:    "",
			expected: false,
			message:  "Expected '' to not be an email address",
		},
		{
			title:    "Missing country extension",
			input:    "some_email@gmail",
			expected: false,
			message:  "Email without a country extension is not an email",
		},
		{
			title:    "Missing @",
			input:    "some_emailgmail.com",
			expected: false,
			message:  "Email without an @ is not an email address",
		},
		{
			title:    "Missing username before @",
			input:    "@gmail.com",
			expected: false,
			message:  "Email without username part is not valid",
		},
		{
			title:    "Correct email",
			input:    "username@gmail.com",
			expected: true,
			message:  "Email should be correct",
		},
	}

	for _, test := range tests {
		if res := val.mustBeEmail(test.input); res != test.expected {
			t.Errorf("Test '%s' failed: %s (Expected: %v, Found: %v)", test.title, test.message, test.expected, res)
		}
	}
}

type withFuncTest struct {
	title    string
	input    string
	function func(string) bool
	expected bool
	message  string
}

func TestValidateWithFunc(t *testing.T) {
	val := &BasicValidator{}

	alwaysTrueFunc := func(input string) bool { return true }
	alwaysFalseFunc := func(input string) bool { return false }

	marker := false

	markerFunc := func(input string) bool {
		marker = true
		return marker
	}

	tests := []withFuncTest{
		{
			title:    "Function always returning true will pass empty string",
			input:    "",
			function: alwaysTrueFunc,
			expected: true,
			message:  "Validating using a function always returning true should always pass, even with empty string",
		},
		{
			title:    "Function always returning true will pass random string",
			input:    "dfjshakjfdshkafds",
			function: alwaysTrueFunc,
			expected: true,
			message:  "Validating using a function always returning true should always pass, even with random string",
		},
		{
			title:    "Function always returning false will not pass empty string",
			input:    "",
			function: alwaysFalseFunc,
			expected: false,
			message:  "Validating using a function always returning true should never pass, even with empty string",
		},
		{
			title:    "Function always returning false will not pass random string",
			input:    "dfjshakjfdshkafds",
			function: alwaysFalseFunc,
			expected: false,
			message:  "Validating using a function always returning false should never pass, even with random string",
		},
		{
			title:    "Function is actually called",
			input:    "dummy",
			function: markerFunc,
			expected: true,
			message:  "Expected the function to be called",
		},
	}

	for _, test := range tests {
		if res := val.validateValueWithFunc(test.input, test.function); res != test.expected {
			t.Errorf("Test '%s' failed: %s (Expected: %v, Found: %v)", test.title, test.message, test.expected, res)
		}
	}

	if marker != true {
		t.Errorf("Expected the marker to be switched to false (it seems the function was not called)")
	}
}

type inValidationTest struct {
	title      string
	input      string
	collection []interface{}
	expected   bool
	message    string
}

func TestMustBeIn(t *testing.T) {
	val := &BasicValidator{}
	tests := []inValidationTest{
		{
			title:      "Empty input in empty collection",
			input:      "",
			expected:   false,
			message:    "Empty collection cannot contain anything",
			collection: []interface{}{},
		},
		{
			title:      "regular input in empty collection",
			input:      "word",
			expected:   false,
			message:    "Empty collection cannot contain anything",
			collection: []interface{}{},
		},
		{
			title:      "word not in int collection",
			input:      "word",
			expected:   false,
			message:    "Word should not be found in an integer collection",
			collection: []interface{}{1, 2, 3, 4},
		},
		{
			title:      "word not in string collection",
			input:      "contain",
			expected:   false,
			message:    "Word should not be found in a string collection that does not contain it",
			collection: []interface{}{"I", "just", "cannot", "this"},
		},
		{
			title:      "word actually in string collection",
			input:      "contain",
			expected:   true,
			message:    "Word should be found in a string collection that does not contain it",
			collection: []interface{}{"I", "just", "cannot", "contain", "this"},
		},
		{
			title:      "int not in int collection",
			input:      "0",
			expected:   false,
			message:    "0 should not be should not be found in an integer collection [1,2,3,4]",
			collection: []interface{}{1, 2, 3, 4},
		},
		{
			title:      "int actually in int collection",
			input:      "2",
			expected:   true,
			message:    "Word should not be found in an integer collection",
			collection: []interface{}{1, 2, 3, 4},
		},
		{
			title:      "something in a non-int, non-string collection",
			input:      "2",
			expected:   false,
			message:    "'2' should not be found in a float collection",
			collection: []interface{}{1.0, 2.0, 3.0, 4.0},
		},
	}

	for _, test := range tests {
		if res := val.mustBeIn(test.input, test.collection); res != test.expected {
			t.Errorf("Test '%s' failed: %s (Expected: %v, Found: %v)", test.title, test.message, test.expected, res)
		}
	}
}

type tester struct {
	field string
}

type testCase struct {
	title    string
	field    string
	value    string
	expected bool
	message  string
	errors   int
}

func TestValidatePresent(t *testing.T) {
	testCases := []testCase{
		{
			title:    "non existing field is not present",
			field:    "nofield",
			value:    "value",
			expected: false,
			message:  "Expected a non-existing field to not be present",
			errors:   1,
		},
		{
			title:    "empty field is not present",
			field:    "",
			value:    "value",
			expected: false,
			message:  "Expected a field without a name to not be present",
			errors:   1,
		},
		{
			title:    "empty string is not present",
			field:    "field",
			value:    "",
			expected: false,
			message:  "Expected an empty string as value to not be present",
			errors:   1,
		},
		{
			title:    "non empty string is present",
			field:    "field",
			value:    "value",
			expected: true,
			message:  "Expected a non empty string to be present",
			errors:   0,
		},
	}

	validator := BasicValidator{}
	for _, tc := range testCases {
		target := &tester{field: tc.value}

		validator.ClearViolations()
		res := validator.ValidatePresent(target, tc.field)
		if res != tc.expected || len(validator.Violations()[tc.field]) > tc.errors {
			t.Errorf("Test '%s' failed: %s (expected errors: %d, found %d) (Expected: %v, Found: %v). Errors: %+v", tc.title,
				tc.message,
				tc.errors,
				len(validator.Violations()[tc.field]),
				tc.expected,
				res,
				validator.Violations()[tc.field],
			)
		}
	}
}
