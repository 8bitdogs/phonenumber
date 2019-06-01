package phonenumber

import (
	"testing"
)

func TestPhoneNumberString(t *testing.T) {
	ph, err := Parse("+1-541-754-0301")
	if err != nil {
		t.Error("failed to parse.", err)
	}
	t.Logf("result %s", ph.String())
	if ph.String() != "+15417540301" {
		t.Error("Phone doesn't equal.")
	}
	t.Log(ph.Local())
}

func TestUnsupportableFormat(t *testing.T) {
	_, err := Parse("+12345678901234")
	if err == nil {

	}
}
