package main

import (
	"errors"
	"fmt"
	"github.com/anantadwi13/errorwrap"
)

var (
	ErrorInfra         = errorwrap.New("error infra layer")
	ErrorInfraNotFound = errorwrap.New("error infra layer [not found]")
	ErrorDomain        = errorwrap.New("error domain layer")
	ErrorUseCase       = errorwrap.New("error usecase layer")
)

func main() {
	err := usecaseLayer(4)
	fmt.Println(err) // print usecase error message only
	fmt.Println()
	fmt.Printf("%+s\n", err) // print full error message
	fmt.Println()
	fmt.Printf("%+v\n", err) // print full error message with stack trace
	fmt.Println()

	switch {
	case errorwrap.Is(err, ErrorInfra):
		fmt.Println("is ErrorInfra")
	case errorwrap.Is(err, ErrorInfraNotFound):
		fmt.Println("is ErrorInfraNotFound") // go to this section
	}

	if err := errorwrap.Wrapper(err, ErrorInfraNotFound); err != nil {
		fmt.Printf("%+v\n", err)
		fmt.Println("Print file stack trace")
		fmt.Printf("%s\n", err.StackTrace())
		fmt.Println("Print file+line stack trace")
		fmt.Printf("%v\n", err.StackTrace())
		fmt.Println("Print function+file+line stack trace")
		fmt.Printf("%+v\n", err.StackTrace())
	}

	fmt.Println(errorwrap.Is(err, ErrorDomain)) // true

	fmt.Println(errorwrap.Is(err, ErrorUseCase)) // true

	fmt.Printf("%+v\n", errorwrap.Wrapper(err, ErrorDomain).StackTrace()) // print stack trace until domain layer only
}

func infraLayer(intType int) error {
	switch intType {
	case 1:
		err := errors.New("standard error")
		if err != nil {
			return errorwrap.Wrap(err, ErrorInfra)
		}
	case 2:
		err := errors.New("resource not found")
		if err != nil {
			return errorwrap.Wrap(err, ErrorInfraNotFound)
		}
	case 3:
		err := errors.New("another error")
		if err != nil {
			return errorwrap.NewErrorWithMessage(err, "using wrapper message")
		}
	case 4:
		err := errors.New("error again")
		if err != nil {
			return errorwrap.WrapWithMessage(err, ErrorInfraNotFound, "unable to find resource in database")
		}
	case 5:
		err := errors.New("just return it")
		if err != nil {
			return err
		}
	}
	return nil
}

func domainLayer(intType int) error {
	err := infraLayer(intType)
	if err != nil {
		return errorwrap.Wrap(err, ErrorDomain)
	}
	return nil
}

func usecaseLayer(intType int) error {
	err := domainLayer(intType)
	if err != nil {
		return errorwrap.Wrap(err, ErrorUseCase)
	}
	return nil
}
