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
}

func ExampleNewError() {
	err := errorwrap.NewError(errors.New("error message"), errorwrap.New("another error message"))

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
	//	 -  error message
	//	    another error message
	//
	//	Second
	//	 -  error message
	//	    another error message
	//
	//	Third
	//	 -  error message
	//	    another error message
	//
	//	Fourth
	//	 -  error message
	//	    another error message
	//
	//	main.main
	//		/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:37
	//	runtime.main
	//		/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//		/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
}

func ExampleNewErrorWithMessage() {
	err := errorwrap.NewErrorWithMessage("context message", errors.New("error message"), errorwrap.New("another error message"))

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
	//	 -  error message
	//	    another error message
	//	    context: context message
	//
	//	Second
	//	 -  error message
	//	    another error message
	//	    context: context message
	//
	//	Third
	//	 -  error message
	//	    another error message
	//	    context: context message
	//
	//	Fourth
	//	 -  error message
	//	    another error message
	//	    context: context message
	//
	//	main.main
	//		/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:75
	//	runtime.main
	//		/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//		/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
}

func ExampleWrap() {
	ErrorStd := errors.New("standard error")
	ErrorA := errorwrap.New("Error A")
	ErrorB := errorwrap.New("Error B")
	err := errorwrap.Wrap(ErrorA, ErrorB, errorwrap.New("same level with B"))
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
	//	 -  Error B
	//	    same level with B
	//
	//	Second
	//	 -  Error B
	//	    same level with B
	//	 -  Error A
	//
	//	Third
	//	 -  Error B
	//	    same level with B
	//
	//	Fourth
	//	 -  Error B
	//	    same level with B
	//	 -  Error A
	//
	//	main.main
	//		/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:120
	//	runtime.main
	//		/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//		/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
	//
	//	Fifth
	//	 -  Error B
	//	 -  standard error
	//
	//	main.main
	//		/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:121
	//	runtime.main
	//		/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//		/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
}

func ExampleWrapWithMessage() {
	ErrorStd := errors.New("standard error")
	ErrorA := errorwrap.New("Error A")
	ErrorB := errorwrap.New("Error B")
	err := errorwrap.WrapWithMessage(ErrorA, "context message", ErrorB)
	err2 := errorwrap.WrapWithMessage(ErrorStd, "context message", ErrorB)

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
	//	 -  Error B
	//	    same level with B
	//	    context: context message
	//
	//	Second
	//	 -  Error B
	//	    same level with B
	//	    context: context message
	//	 -  Error A
	//
	//	Third
	//	 -  Error B
	//	    same level with B
	//	    context: context message
	//
	//	Fourth
	//	 -  Error B
	//	    same level with B
	//	    context: context message
	//	 -  Error A
	//
	//	main.main
	//		/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:177
	//	runtime.main
	//		/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//		/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
	//
	//	Fifth
	//	 -  Error B
	//	    context: context message
	//	 -  standard error
	//
	//	main.main
	//		/Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:178
	//	runtime.main
	//		/usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//		/usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
}

func ExampleWrapper() {
	var (
		ErrorA              = errorwrap.New("Error A")
		ErrorSameLevelWithA = errors.New("error same level with A")
		ErrorB              = errorwrap.New("Error B")
	)

	errA := errorwrap.NewError(ErrorA, ErrorSameLevelWithA)
	err := errorwrap.WrapWithMessage(errA, "context message", ErrorB)

	fmt.Println("wrapper of error definition (should be nil)")
	fmt.Println(errorwrap.Wrapper(ErrorA, ErrorA) == nil)
	fmt.Println()

	wrapErrorA := errorwrap.Wrapper(errA, ErrorA)
	fmt.Println("wrapErrorA")
	fmt.Println(wrapErrorA == nil)
	fmt.Printf("RootCause:\n%v\n", wrapErrorA.RootCause())
	fmt.Printf("ParentError:\n%v\n", wrapErrorA.ParentError())
	fmt.Printf("CurrentError:\n%q\n", wrapErrorA.CurrentError())
	fmt.Println("wrapErrorA message:")
	fmt.Printf("%v\n", wrapErrorA)
	fmt.Println("wrapErrorA message + stack trace:")
	fmt.Printf("%+v\n", wrapErrorA)
	fmt.Println()

	errWrapper1 := errorwrap.Wrapper(err, ErrorA)
	fmt.Println("errWrapper1")
	fmt.Println(errWrapper1 == nil)
	fmt.Printf("RootCause:\n%v\n", errWrapper1.RootCause())
	fmt.Printf("ParentError:\n%v\n", errWrapper1.ParentError())
	fmt.Printf("CurrentError:\n%q\n", errWrapper1.CurrentError())
	fmt.Println("errWrapper1 message:")
	fmt.Printf("%v\n", errWrapper1)
	fmt.Println("errWrapper1 message + stack trace:")
	fmt.Printf("%+v\n", errWrapper1)
	fmt.Println()

	errWrapper2 := errorwrap.Wrapper(err, ErrorB)
	fmt.Println("errWrapper2")
	fmt.Println(errWrapper2 == nil)
	fmt.Printf("RootCause:\n%v\n", errWrapper2.RootCause())
	fmt.Printf("ParentError:\n%v\n", errWrapper2.ParentError())
	fmt.Printf("CurrentError:\n%q\n", errWrapper2.CurrentError())
	fmt.Println("errWrapper2 message:")
	fmt.Printf("%v\n", errWrapper2)
	fmt.Println("errWrapper2 message + stack trace:")
	fmt.Printf("%+v\n", errWrapper2)
	fmt.Println()

	//	Example output:
	//
	//	wrapper of error definition (should be nil)
	//	true
	//
	//	wrapErrorA
	//	false
	//	RootCause:
	//	<nil>
	//	ParentError:
	//	<nil>
	//	CurrentError:
	//	["Error A" "error same level with A"]
	//	wrapErrorA message:
	//	 -  Error A
	//	    error same level with A
	//	wrapErrorA message + stack trace:
	//	 -  Error A
	//	    error same level with A
	//
	//	main.main
	//	        /Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:242
	//	runtime.main
	//	        /usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//	        /usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
	//
	//	errWrapper1
	//	false
	//	RootCause:
	//	<nil>
	//	ParentError:
	//	<nil>
	//	CurrentError:
	//	["Error A" "error same level with A"]
	//	errWrapper1 message:
	//	 -  Error A
	//	    error same level with A
	//	errWrapper1 message + stack trace:
	//	 -  Error A
	//	    error same level with A
	//
	//	main.main
	//	        /Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:242
	//	runtime.main
	//	        /usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//	        /usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
	//
	//	errWrapper2
	//	false
	//	RootCause:
	//	 -  Error A
	//	    error same level with A
	//	ParentError:
	//	 -  Error A
	//	    error same level with A
	//	CurrentError:
	//	["Error B"]
	//	errWrapper2 message:
	//	 -  Error B
	//	    context: context message
	//	errWrapper2 message + stack trace:
	//	 -  Error B
	//	    context: context message
	//	 -  Error A
	//	    error same level with A
	//
	//	main.main
	//	        /Users/go/src/github.com/anantadwi13/errorwrap/example_test.go:242
	//	runtime.main
	//	        /usr/local/Cellar/go/1.17.5/libexec/src/runtime/proc.go:255
	//	runtime.goexit
	//	        /usr/local/Cellar/go/1.17.5/libexec/src/runtime/asm_amd64.s:1581
}
