package main

import (
	"errors"
	"fmt"
	"github.com/anantadwi13/errorwrap"
)

var (
	ErrorCommonNotFound = errorwrap.New("error not found")

	ErrorInfraMysql = errorwrap.New("error infra mysql")
	ErrorInfraRedis = errorwrap.New("error infra redis")

	ErrorInfra   = errorwrap.New("error infra layer")
	ErrorDomain  = errorwrap.New("error domain layer")
	ErrorUseCase = errorwrap.New("error usecase layer")
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
	case errorwrap.Is(err, ErrorInfraRedis):
		fmt.Println("is ErrorInfraRedis")
	case errorwrap.Is(err, ErrorInfraMysql):
		fmt.Println("is ErrorInfraMysql") // go to this section
	}

	fmt.Println("is ErrorCommonNotFound", errorwrap.Is(err, ErrorCommonNotFound)) //true

	if err := errorwrap.Wrapper(err, ErrorCommonNotFound); err != nil {
		fmt.Printf("%+v\n", err)
		fmt.Println("Print file stack trace")
		fmt.Printf("%s\n", err.StackTrace())
		fmt.Println("Print file+line stack trace")
		fmt.Printf("%v\n", err.StackTrace())
		fmt.Println("Print function+file+line stack trace")
		fmt.Printf("%+v\n", err.StackTrace())
	}

	fmt.Println("is ErrorDomain", errorwrap.Is(err, ErrorDomain)) // true

	fmt.Println("is ErrorUseCase", errorwrap.Is(err, ErrorUseCase)) // true

	fmt.Printf("\nprint usecase stack trace%+v\n", errorwrap.Wrapper(err, ErrorUseCase).StackTrace()) // print stack trace until usecase layer only
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
			return errorwrap.Wrap(err, ErrorCommonNotFound, ErrorInfra, ErrorInfraRedis)
		}
	case 3:
		err := errors.New("another error")
		if err != nil {
			return errorwrap.NewErrorWithMessage("using wrapper message", err)
		}
	case 4:
		err := errors.New("error again")
		if err != nil {
			return errorwrap.WrapWithMessage(err, "unable to find resource in database", ErrorInfra, ErrorCommonNotFound, ErrorInfraMysql)
		}
	case 5:
		err := errors.New("just return it")
		if err != nil {
			return err
		}
	case 6:
		return errorwrap.NewError(ErrorInfra)
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
