package hsm

import "reflect"
import "fmt"

func AssertEqual(expected, actual interface{}) {
    if !ObjectAreEqual(expected, actual) {
        panic(fmt.Sprintf("Equal(%#v, %#v) fail", expected, actual))
    }
}

func AssertNotEqual(expected, actual interface{}) {
    if ObjectAreEqual(expected, actual) {
        panic(fmt.Sprintf("NotEqual(%#v, %#v) fail", expected, actual))
    }
}

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

func AssertTrue(value bool) {
    if !value {
        panic(fmt.Sprintf("True(value=%#v) fail", value))
    }
}

func AssertFalse(value bool) {
    if value {
        panic(fmt.Sprintf("False(value=%#v) fail", value))
    }
}

func AssertNil(value interface{}) {
    AssertEqual(nil, value)
}

func AssertNotNil(value interface{}) {
    AssertNotEqual(nil, value)
}
