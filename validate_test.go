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
		if res := val.MustBePresent(test.input); res != test.expected {
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
		if res := val.MustBeEmail(test.input); res != test.expected {
			t.Errorf("Test '%s' failed: %s (Expected: %v, Found: %v)", test.title, test.message, test.expected, res)
		}
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
		if res := val.MustBeIn(test.input, test.collection); res != test.expected {
			t.Errorf("Test '%s' failed: %s (Expected: %v, Found: %v)", test.title, test.message, test.expected, res)
		}
	}
}
