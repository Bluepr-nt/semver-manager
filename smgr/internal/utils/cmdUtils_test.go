package utils

import (
	"reflect"
	"testing"
)

type Config struct {
	Field1 string `san:"trim"`
	Field2 int
}

func TestSanitizeInputs(t *testing.T) {
	tests := []struct {
		name        string
		inputs      interface{}
		expected    interface{}
		expectedErr bool
	}{
		{
			name: "sanitizes string field",
			inputs: &Config{
				Field1: "test  ",
			},
			expected: &Config{
				Field1: "test",
			},
			expectedErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := SanitizeInputs(test.inputs)
			if test.expectedErr && err == nil {
				t.Errorf("expected error but got none")
			} else if !test.expectedErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !test.expectedErr && !reflect.DeepEqual(test.inputs, test.expected) {
				t.Errorf("unexpected result: expected %v but got %v", test.expected, test.inputs)
			}
		})
	}
}
