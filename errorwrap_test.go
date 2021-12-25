package errorwrap

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	ErrorCommonNotFound = New("error not found")

	ErrorMysqlDb = errors.New("error database mysql")
	ErrorRedisDb = New("error database redis")

	ErrorInfraDatabase = New("error infra layer")
	ErrorDomain        = New("error domain layer")
	ErrorUseCase       = New("error usecase layer")
	ErrorApp           = New("error app layer")

	ErrorTestA = New("error test a")
	ErrorTestB = errors.New("error test b")
)

const (
	MYSQL = iota + 1
	REDIS
)

func infraDbLayer(intType ...int) error {
	it := 0
	if len(intType) > 0 {
		it = intType[0]
	}
	switch it {
	case MYSQL:
		return NewError(ErrorInfraDatabase, ErrorCommonNotFound, ErrorMysqlDb)
	case REDIS:
		return NewErrorWithMessage("redis not found", ErrorInfraDatabase, ErrorCommonNotFound, ErrorRedisDb)
	default:
		return NewError(ErrorInfraDatabase)
	}
}

func domainLayer(intType ...int) error {
	err := infraDbLayer(intType...)
	if err != nil {
		return Wrap(err, ErrorDomain)
	}
	return nil
}

func usecaseLayer(intType ...int) error {
	err := domainLayer(intType...)
	if err != nil {
		return Wrap(err, ErrorUseCase)
	}
	return nil
}

func appLayer(intType ...int) error {
	err := usecaseLayer(intType...)
	if err != nil {
		return Wrap(err, ErrorApp)
	}
	return nil
}

func TestNew(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "success", args: args{message: "error message"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.message)
			assert.NotNil(t, got)
			assert.Equal(t, tt.args.message, got.Error())
		})
	}
}

func TestErrorIs(t *testing.T) {
	type args struct {
		wrapper error
		target  error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success usecaseLayer->ErrorInfraDatabase",
			args: args{
				wrapper: usecaseLayer(),
				target:  ErrorInfraDatabase,
			},
			want: true,
		},
		{
			name: "success usecaseLayer->ErrorDomain",
			args: args{
				wrapper: usecaseLayer(),
				target:  ErrorDomain,
			},
			want: true,
		},
		{
			name: "success usecaseLayer->ErrorUseCase",
			args: args{
				wrapper: usecaseLayer(),
				target:  ErrorUseCase,
			},
			want: true,
		},
		{
			name: "error usecaseLayer->ErrorApp",
			args: args{
				wrapper: usecaseLayer(),
				target:  ErrorApp,
			},
			want: false,
		},
		{
			name: "success domainLayer->ErrorInfraDatabase",
			args: args{
				wrapper: domainLayer(),
				target:  ErrorInfraDatabase,
			},
			want: true,
		},
		{
			name: "success domainLayer->ErrorDomain",
			args: args{
				wrapper: domainLayer(),
				target:  ErrorDomain,
			},
			want: true,
		},
		{
			name: "error domainLayer->ErrorUseCase",
			args: args{
				wrapper: domainLayer(),
				target:  ErrorUseCase,
			},
			want: false,
		},
		{
			name: "error domainLayer->ErrorApp",
			args: args{
				wrapper: domainLayer(),
				target:  ErrorApp,
			},
			want: false,
		},
		{
			name: "error ErrorInfraDatabase->domainLayer",
			args: args{
				wrapper: ErrorInfraDatabase,
				target:  domainLayer(),
			},
			want: false,
		},
		{
			name: "success appLayer->ErrorInfraDatabase",
			args: args{
				wrapper: appLayer(),
				target:  ErrorInfraDatabase,
			},
			want: true,
		},
		{
			name: "success appLayer->ErrorDomain",
			args: args{
				wrapper: appLayer(),
				target:  ErrorDomain,
			},
			want: true,
		},
		{
			name: "success appLayer->ErrorUseCase",
			args: args{
				wrapper: appLayer(),
				target:  ErrorUseCase,
			},
			want: true,
		},
		{
			name: "success appLayer->ErrorApp",
			args: args{
				wrapper: appLayer(),
				target:  ErrorApp,
			},
			want: true,
		},
		{
			name: "success appLayer->ErrorRedisDb",
			args: args{
				wrapper: appLayer(REDIS),
				target:  ErrorRedisDb,
			},
			want: true,
		},
		{
			name: "error appLayer->ErrorRedisDb",
			args: args{
				wrapper: appLayer(MYSQL),
				target:  ErrorRedisDb,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.args.wrapper, tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapper(t *testing.T) {
	type args struct {
		err    error
		target error
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "get ErrorInfraDatabase from appLayer",
			args: args{
				err:    appLayer(),
				target: ErrorInfraDatabase,
			},
			wantErr: assert.Error,
		},
		{
			name: "invalid get ErrorInfraDatabase from appLayer",
			args: args{
				err:    ErrorInfraDatabase,
				target: appLayer(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "get ErrorDomain from domainLayer",
			args: args{
				err:    domainLayer(),
				target: ErrorDomain,
			},
			wantErr: assert.Error,
		},
		{
			name: "invalid get ErrorDomain from domainLayer",
			args: args{
				err:    ErrorDomain,
				target: domainLayer(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "failed to get ErrorUseCase from domainLayer",
			args: args{
				err:    domainLayer(),
				target: ErrorUseCase,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Wrapper(tt.args.err, tt.args.target)
			if tt.wantErr(t, err) {
				assert.NotEmpty(t, fmt.Sprintf("%+v", err))
			}
		})
	}
}

func TestAppendInto(t *testing.T) {
	type args struct {
		errWrapper   error
		currentError []error
		err          []error
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{
			name: "expect nil",
			args: args{
				errWrapper:   nil,
				currentError: nil,
				err:          nil,
			},
			wantError: false,
		},
		{
			name: "success appLayer & ErrorTestA + ErrorTestB",
			args: args{
				errWrapper:   appLayer(),
				currentError: []error{ErrorApp},
				err:          []error{ErrorTestA, ErrorTestB},
			},
			wantError: true,
		},
		{
			name: "success appLayer & ErrorTestA",
			args: args{
				errWrapper:   appLayer(),
				currentError: []error{ErrorApp},
				err:          []error{ErrorTestA},
			},
			wantError: true,
		},
		{
			name: "success ErrorTestA only",
			args: args{
				errWrapper:   nil,
				currentError: nil,
				err:          []error{ErrorTestA},
			},
			wantError: true,
		},
		{
			name: "success infraDbLayer Redis & ErrorTestA + ErrorTestB",
			args: args{
				errWrapper:   infraDbLayer(REDIS),
				currentError: []error{ErrorInfraDatabase, ErrorCommonNotFound, ErrorRedisDb},
				err:          []error{ErrorTestA},
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AppendInto(tt.args.errWrapper, tt.args.err...)
			if !tt.wantError {
				assert.Nil(t, got)
				return
			}
			for _, err := range tt.args.currentError {
				assert.True(t, IsExact(got, err))
			}
			for _, err := range tt.args.err {
				assert.True(t, IsExact(got, err))
			}
		})
	}
}
