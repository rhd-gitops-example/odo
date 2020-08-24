package ui

import "testing"

func TestValidatePrefix(t *testing.T) {

	validator := MakePrefixValidator()

	if err := validator("tst"); err != nil {
		t.Fatalf("got %s error for long prefix", err)
	}
	if err := validator("testing"); err != nil {
		t.Fatalf("got %s error for long prefix", err)
	}
}
