package ui

import (
	"testing"
)

// func TestNameValidator(t *testing.T) {
// 	// note that we're just testing a single case here since presumably the underlying implementation is already tested in k8s
// 	err := ValidatePrefix("some-valid-name")
// 	if err == nil {
// 		t.Errorf("name validator should have accepted name, but got: %v instead", err)
// 	}

// 	err = ValidatePrefix("abcdefghijklmnopqrstuvwxyzabcderfgighshshhshhshshshsshhshshsshshhshshshshsshhshshshshshshshshshshshshshshhshhshshs")
// 	if err == nil {
// 		t.Error("name validator should only attempt to validate non-nil strings")
// 	}
// }

func TestValidateSecretLength(t *testing.T) {
	// note that we're just testing a single case here since presumably the underlying implementation is already tested in k8s
	err := validateSecretLength("123")
	if err == nil {
		t.Errorf("name validator should have accepted name, but got: %v instead", err)
	}

	err = ValidatePrefix("abcdefghijklmnopqrstuvwxyzabcderfgighshshhshhshshshsshhshshsshshhshshshshsshhshshshshshshshshshshshshshshhshhshshs")
	if err == nil {
		t.Errorf("The secret length should 16 or more ")
	}
}
