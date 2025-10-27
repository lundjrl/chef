package tests

import "testing"

func TestFun(t *testing.T) {
	msg := "test message"
	err := "some error"
	t.Errorf(`Hello("") = %q, %v, want "", error`, msg, err)
}
