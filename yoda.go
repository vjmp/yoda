/*
Package yoda provides DSL for having terse conversation with your code.
Story goes like this:

Luke wanted to pair with master Yoda.

Luke: Shall we write this factorial function together, master Yoda?

Yoda: Equal 1 and sut.Fact(1) must be.

And there they go, Yoda is testing Luke again. Returning 1 Luke is.

Yoda: True sut.Fact(1) == sut.Fact(2) wont be.

...

*/
package yoda

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"testing/quick"
)

// Internal boolean type to bind "Must" and "Wont" functinality.
type Truth struct {
	Value bool
	Dump string
}

// Function wrapper type used for Panic testing.
type Callable func()

// Compares expected and actual values, so that they are deeply equal.
// Uses reflect.DeepEqual to do that.
//
// Result is just boolean truth, not verified yet. Anything goes.
// Actual verification should be done using Must or Wont calls.
func Equal(expected, actual interface{}) Truth {
	mustBeCleanStart()
	return Truth{
		Value:reflect.DeepEqual(expected, actual),
		Dump:fmt.Sprintf("%#v vs. %#v", expected, actual),
	}
}

// Determines if actual is really nil value or not.
// Uses recover, reflect.TypeOf and reflect.ValueOf to do that.
//
// Result is just boolean truth, not verified yet. Anything goes.
// Actual verification should be done using Must or Wont calls.
func Nil(actual interface{}) (result Truth) {
	mustBeCleanStart()
	defer func() {
		if recover() != nil {
			result = Truth{false, fmt.Sprintf("%#v", actual)}
		}
	}()
	return Truth{
		Value:reflect.TypeOf(actual) == nil || reflect.ValueOf(actual).IsNil(),
		Dump:fmt.Sprintf("%#v", actual),
	}
}

// Determines if function panics or not. Uses recover to do that.
//
// Result is just boolean truth, not verified yet. Anything goes.
// Actual verification should be done using Must or Wont calls.
func Panic(function Callable) (result Truth) {
	mustBeCleanStart()
	defer func() {
		if recover() != nil {
			result = Truth{true, ""}
		}
	}()
	function()
	return Truth{false, ""}
}

// Compares that expected and actual values seem to be same.
// Uses reflect.DeepEqual and fmt.Sprintf with %p and %#v to do that.
//
// Result is just boolean truth, not verified yet. Anything goes.
// Actual verification should be done using Must or Wont calls.
func Same(expected, actual interface{}) Truth {
	mustBeCleanStart()
	return Truth{
		Value: nice(expected) == nice(actual) && reflect.DeepEqual(actual, expected),
		Dump:fmt.Sprintf("%#v", actual),
	}
}

// Compares that expected and actual values converts to equal string
// representation. Uses fmt.Sprintf with %v to do that.
//
// Result is just boolean truth, not verified yet. Anything goes.
// Actual verification should be done using Must or Wont calls.
func Text(expected, actual interface{}) Truth {
	mustBeCleanStart()
	return Truth{
		Value:fmt.Sprintf("%v", expected) == fmt.Sprintf("%v", actual),
		Dump:fmt.Sprintf("%#v vs. %#v", expected, actual),
	}
}

// Determines if actual is really true or false.
// Uses type assertion to do that.
//
// Result is just boolean truth, not verified yet. Anything goes.
// Actual verification should be done using Must or Wont calls.
func True(actual interface{}) Truth {
	mustBeCleanStart()
	return Truth{actual.(bool), fmt.Sprintf("%#v", actual)}
}

// Determines if type of actual is represented as expected string.
//
// Result is just boolean truth, not verified yet. Anything goes.
// Actual verification should be done using Must or Wont calls.
func Type(expected string, actual interface{}) Truth {
	mustBeCleanStart()
	return Truth{
		Value: expected == reflect.TypeOf(actual).String(),
		Dump: fmt.Sprintf("%#v vs. %#v", expected, reflect.TypeOf(actual).String()),
	}
}

// Property is some function that accepts some random inputs and
// returns bool value indicating if property actually holds or not.
// Uses quick.Check internally.
//
// Result is just boolean truth, not verified yet. Anything goes.
// Actual verification should be done using Must or Wont calls.
func All(property interface{}) Truth {
	mustBeCleanStart()
	err := quick.Check(property, nil)
	return Truth{err == nil, fmt.Sprintf("try %v", err)}
}

// Verifies that given truth is actually true. It must be, or fatal
// failure follows.
//
// On case of failure, you get message like:
//
//        Must failed at darkside_test.go:42 !!!
//
// Use the source, Luke!
func (the Truth) Must(be *testing.T) {
	clean_start = true
	if !the.Value {
		failedTo(be, "Must", the.Dump)
	}
}

// Verifies that given truth is actually false. If it wont be, or failure follows.
//
// On case of failure, you get message like:
//
//        Wont failed at darkside_test.go:99 !!!
//
// Use the source, Luke!
func (the Truth) Wont(be *testing.T) {
	clean_start = true
	if the.Value {
		failedTo(be, "Wont", the.Dump)
	}
}

func failedTo(be *testing.T, what, dump string) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file, line = "<unknown>", 0
	}
	be.Fatalf("\n%s failed at %s:%d; Reason ~ %v!!!", what, relative(file), line, dump)
}

func nice(value interface{}) string {
	return fmt.Sprintf("%p %#v", value, value)
}

func relative(filename string) string {
	dir, err := os.Getwd()
	if err != nil {
		return filename
	}
	if strings.HasPrefix(filename, dir) {
		return filename[len(dir)+1:]
	}
	return filename
}

var (
	clean_start = true
)

func mustBeCleanStart() {
	if !clean_start {
		clean_start = true
		panic("Yoda: darkside you have, missing Must or Wont must be!")
	}
	clean_start = false
}
