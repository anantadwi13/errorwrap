package errorwrap_test

import (
	"errors"
	"fmt"
	"github.com/anantadwi13/errorwrap"
)

func ExampleNew() {
	err := errorwrap.New("error message")

	fmt.Println("\nFirst")
	fmt.Printf("%s\n", err)
	fmt.Println("\nSecond")
	fmt.Printf("%+s\n", err)
	fmt.Println("\nThird")
	fmt.Printf("%v\n", err)
	fmt.Println("\nFourth")
	fmt.Printf("%+v\n", err)

	//	Example output:
	//
	//	First
	//	error message
	//
	//	Second
	//	error message
	//
	//	Third
	//	error message
	//
	//	Fourth
	//	error message
	//	main.main
	//	/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:10
	//	runtime.main
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
}

func ExampleNewError() {
	err := errorwrap.NewError(errors.New("error message"))

	fmt.Println("\nFirst")
	fmt.Printf("%s\n", err)
	fmt.Println("\nSecond")
	fmt.Printf("%+s\n", err)
	fmt.Println("\nThird")
	fmt.Printf("%v\n", err)
	fmt.Println("\nFourth")
	fmt.Printf("%+v\n", err)

	//	Example output:
	//
	//	First
	//	error message
	//
	//	Second
	//	error message
	//
	//	Third
	//	error message
	//
	//	Fourth
	//	error message
	//	main.main
	//	/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:43
	//	runtime.main
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
}

func ExampleNewErrorWithMessage() {
	err := errorwrap.NewErrorWithMessage(errors.New("error message"), "context message")

	fmt.Println("\nFirst")
	fmt.Printf("%s\n", err)
	fmt.Println("\nSecond")
	fmt.Printf("%+s\n", err)
	fmt.Println("\nThird")
	fmt.Printf("%v\n", err)
	fmt.Println("\nFourth")
	fmt.Printf("%+v\n", err)

	//	Example output:
	//
	//	First
	//	error message: context message
	//
	//	Second
	//	error message: context message
	//
	//	Third
	//	error message: context message
	//
	//	Fourth
	//	error message: context message
	//	main.main
	//	/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:76
	//	runtime.main
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
}

func ExampleWrap() {
	ErrorStd := errors.New("standard error")
	ErrorA := errorwrap.New("Error A")
	ErrorB := errorwrap.New("Error B")
	err := errorwrap.Wrap(ErrorA, ErrorB)
	err2 := errorwrap.Wrap(ErrorStd, ErrorB)

	fmt.Println("\nFirst")
	fmt.Printf("%s\n", err)
	fmt.Println("\nSecond")
	fmt.Printf("%+s\n", err)
	fmt.Println("\nThird")
	fmt.Printf("%v\n", err)
	fmt.Println("\nFourth")
	fmt.Printf("%+v\n", err)
	fmt.Println("\nFifth")
	fmt.Printf("%+v\n", err2)

	//	Example output:
	//
	//	First
	//	Error B
	//
	//	Second
	//	Error B
	//	Error A
	//
	//	Third
	//	Error B
	//
	//	Fourth
	//	Error B
	//	Error A
	//	main.main
	//	/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:110
	//	runtime.main
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
	//
	//	Fifth
	//	Error B
	//	standard error
	//	main.main
	//	/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:113
	//	runtime.main
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
}

func ExampleWrapWithMessage() {
	ErrorStd := errors.New("standard error")
	ErrorA := errorwrap.New("Error A")
	ErrorB := errorwrap.New("Error B")
	err := errorwrap.WrapWithMessage(ErrorA, ErrorB, "context message")
	err2 := errorwrap.WrapWithMessage(ErrorStd, ErrorB, "context message")

	fmt.Println("\nFirst")
	fmt.Printf("%s\n", err)
	fmt.Println("\nSecond")
	fmt.Printf("%+s\n", err)
	fmt.Println("\nThird")
	fmt.Printf("%v\n", err)
	fmt.Println("\nFourth")
	fmt.Printf("%+v\n", err)
	fmt.Println("\nFifth")
	fmt.Printf("%+v\n", err2)

	//	Example output:
	//
	//	First
	//	Error B: context message
	//
	//	Second
	//	Error B: context message
	//	Error A
	//
	//	Third
	//	Error B: context message
	//
	//	Fourth
	//	Error B: context message
	//	Error A
	//	main.main
	//	/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:161
	//	runtime.main
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
	//
	//	Fifth
	//	Error B: context message
	//	standard error
	//	main.main
	//	/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:164
	//	runtime.main
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
}

func ExampleWrapper() {
	ErrorA := errorwrap.New("Error A")
	ErrorB := errorwrap.New("Error B")
	err := errorwrap.WrapWithMessage(ErrorA, ErrorB, "context message")

	wrapErrorA := errorwrap.Wrapper(ErrorA, ErrorA)
	fmt.Println("wrapErrorA")
	fmt.Println(wrapErrorA == nil)
	fmt.Println(wrapErrorA.RootCause(), wrapErrorA.ParentError(), wrapErrorA.CurrentError())
	fmt.Println("wrapErrorA message:")
	fmt.Printf("%v\n", wrapErrorA)
	fmt.Println("wrapErrorA message + stack trace:")
	fmt.Printf("%+v\n", wrapErrorA)
	fmt.Println()

	errWrapper1 := errorwrap.Wrapper(err, ErrorA)
	fmt.Println("errWrapper1")
	fmt.Println(errWrapper1 == nil)
	fmt.Println(errWrapper1.RootCause(), errWrapper1.ParentError(), errWrapper1.CurrentError())
	fmt.Println("errWrapper1 message:")
	fmt.Printf("%v\n", errWrapper1)
	fmt.Println("errWrapper1 message + stack trace:")
	fmt.Printf("%+v\n", errWrapper1)
	fmt.Println()

	errWrapper2 := errorwrap.Wrapper(err, ErrorB)
	fmt.Println("errWrapper2")
	fmt.Println(errWrapper2 == nil)
	fmt.Println(errWrapper2.RootCause(), errWrapper2.ParentError(), errWrapper2.CurrentError())
	fmt.Println("errWrapper2 message:")
	fmt.Printf("%v\n", errWrapper2)
	fmt.Println("errWrapper2 message + stack trace:")
	fmt.Printf("%+v\n", errWrapper2)
	fmt.Println()

	//	Example output:
	//
	//	wrapErrorA
	//	false
	//	<nil> <nil> Error A
	//	wrapErrorA message:
	//	Error A
	//	wrapErrorA message + stack trace:
	//	Error A
	//	main.main
	//	/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:183
	//	runtime.main
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
	//
	//	errWrapper1
	//	false
	//	<nil> <nil> Error A
	//	errWrapper1 message:
	//	Error A
	//	errWrapper1 message + stack trace:
	//	Error A
	//	main.main
	//	/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:183
	//	runtime.main
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
	//
	//	errWrapper2
	//	false
	//	Error A Error A Error B
	//	errWrapper2 message:
	//	Error B: context message
	//	errWrapper2 message + stack trace:
	//	Error B: context message
	//	Error A
	//	main.main
	//	/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:183
	//	runtime.main
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//	/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
}
