package hsm

import "reflect"
import "fmt"

// AssertEqual asserts the equality of actual and expected.
func AssertEqual(expected, actual interface{}) {
    if !ObjectAreEqual(expected, actual) {
        panic(fmt.Sprintf("Equal(%#v, %#v) fail", expected, actual))
    }
}

// AssertEqual asserts the inequality of actual and expected.
func AssertNotEqual(expected, actual interface{}) {
    if ObjectAreEqual(expected, actual) {
        panic(fmt.Sprintf("NotEqual(%#v, %#v) fail", expected, actual))
    }
}

// ObjectAreEqual test whether actual is equal to expected.
// It returns true when equal, otherwise returns false.
func ObjectAreEqual(expected, actual interface{}) bool {
    if expected == nil || actual == nil {
        return expected == actual
    }

    if reflect.DeepEqual(expected, actual) {
        return true
    }

    expectedValue := reflect.ValueOf(expected)
    actualValue := reflect.ValueOf(actual)
    if expectedValue == actualValue {
        return true
    }

    if actualValue.Type().ConvertibleTo(expectedValue.Type()) &&
        expectedValue == actualValue.Convert(expectedValue.Type()) {
        return true
    }

    if fmt.Sprintf("%#v", expected) == fmt.Sprintf("%#v", actual) {
        return true
    }

    return false
}

// AssertTrue asserts on truth of value.
func AssertTrue(value bool) {
    if !value {
        panic(fmt.Sprintf("True(value=%#v) fail", value))
    }
}

// AssertFalse asserts on falsehood of value.
func AssertFalse(value bool) {
    if value {
        panic(fmt.Sprintf("False(value=%#v) fail", value))
    }
}

// AssertNil asserts on nullability of value.
func AssertNil(value interface{}) {
    AssertEqual(nil, value)
}

// AssertNotNil is opposite to AssertNil.
func AssertNotNil(value interface{}) {
    AssertNotEqual(nil, value)
}
