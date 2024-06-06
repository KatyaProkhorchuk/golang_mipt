//go:build !solution

package testequal

import (
	"fmt"
)

func ErrorMsg(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	size := len(msgAndArgs)
	var message string
	if size == 1 {
		message = msgAndArgs[0].(string)
	} else if size == 0 {
		message = ""
	} else {
		message = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	t.Errorf(`not equal:'
	expected: %v
	actual: %v
	message: %s 
	`, expected, actual, message)
}

func Equal(t T, expected, actual interface{}) bool{
	t.Helper()
	if expected == nil || actual == nil {
		return expected == actual
	}
	switch expectedValue := expected.(type) {
	case uint, uint8, uint16, uint64, int, int8, int16, int64, int32, uint32:
		return expectedValue == actual
	case string:
		if actualValue, ok := actual.(string); ok {
			return actualValue == expectedValue
		} else {
			return false
		}
	case map[string]string:
		if actualValue, ok := actual.(map[string]string); ok {
			if len(actualValue) != len(expectedValue) {
				return false
			}
			if (expectedValue == nil && actualValue != nil) {
				return false
			}
			if (expectedValue != nil && actualValue == nil) {
				return false
			}
			for key := range expectedValue {
				if value, ok := actualValue[key]; ok {
					if value != expectedValue[key] {
						return false
					}
				} else {
					return false
				}
			}
			return true
		} else {
			return false
		}
	case []int:
		if actualValue, ok := actual.([]int); ok {
			if len(actualValue) != len(expectedValue) {
				return false
			}
			if (expectedValue == nil && actualValue != nil) {
				return false
			}
			if (expectedValue != nil && actualValue == nil) {
				return false
			}
			for i := range expectedValue {
				if expectedValue[i] != actualValue[i] {
					return false
				}
			}
			return true
		} else {
			return false
		}
	case []byte:
		if actualValue, ok := actual.([]byte); ok {
			if len(actualValue) != len(expectedValue) {
				return false
			}
			if (expectedValue == nil && actualValue != nil) {
				return false
			}
			if (expectedValue != nil && actualValue == nil) {
				return false
			}
			for i := range expectedValue {
				if expectedValue[i] != actualValue[i] {
					return false
				}
			}
			return true
		} else {
			return false
		}
	}
	return false
}

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if Equal(t, expected, actual) {
		return true
	}
	ErrorMsg(t, expected, actual, msgAndArgs...)
	return false
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if !Equal(t, expected, actual) {
		return true
	}
	ErrorMsg(t,expected, actual, msgAndArgs...)
	return false
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if Equal(t, expected, actual) {
		return
	}
	ErrorMsg(t,expected, actual, msgAndArgs...)
	t.FailNow()
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !Equal(t, expected, actual) {
		return
	}
	ErrorMsg(t, expected, actual, msgAndArgs...)
	t.FailNow()
}
