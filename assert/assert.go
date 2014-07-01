package assert

import "reflect"
import "fmt"

func Equal(expected, actual interface{}) {
    if !ObjectAreEqual(expected, actual) {
        panic(fmt.Sprintf("Equal(%#v, %#v) fail", expected, actual))
    }
}

func NotEqual(expected, actual interface{}) {
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

func True(value bool) {
    if !value {
        panic(fmt.Sprintf("True(value=%#v) fail", value))
    }
}

func False(value bool) {
    if value {
        panic(fmt.Sprintf("False(value=%#v) fail", value))
    }
}
