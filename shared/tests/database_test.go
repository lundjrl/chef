package tests

import "testing"

func TestFun(t *testing.T) {
	msg := "test message"
	err := "some error"
	t.Logf(`Hello("") = %q, %v, want "", error`, msg, err)
}
