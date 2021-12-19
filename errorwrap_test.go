package errorwrap

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var (
	ErrorRootCause = errors.New("error root cause")
	ErrorInfra     = errors.New("error infra layer")
	ErrorDomain    = errors.New("error domain layer")
	ErrorUseCase   = errors.New("error usecase layer")
	ErrorApp       = errors.New("error app layer")
)

func infraLayer(intType ...int) error {
	it := 0
	if len(intType) > 0 {
		it = intType[0]
	}
	switch it {
	case 1:
		return NewErrorWithMessage(ErrorRootCause, "database not found")
	case 2:
		return NewErrorWithMessage(ErrorInfra, "redis not found")
	default:
		return Wrap(ErrorRootCause, ErrorInfra)
	}
}

func domainLayer(intType ...int) error {
	err := infraLayer(intType...)
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
			stackTrace := fmt.Sprintf("%+v", got)
			assert.NotEmpty(t, stackTrace)
			stackTrace = strings.Replace(stackTrace, tt.args.message+"\n", "", 1)
			assert.NotEmpty(t, stackTrace)
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
			name: "success ABC->A",
			args: args{
				wrapper: usecaseLayer(),
				target:  ErrorInfra,
			},
			want: true,
		},
		{
			name: "success ABC->B",
			args: args{
				wrapper: usecaseLayer(),
				target:  ErrorDomain,
			},
			want: true,
		},
		{
			name: "success ABC->C",
			args: args{
				wrapper: usecaseLayer(),
				target:  ErrorUseCase,
			},
			want: true,
		},
		{
			name: "error ABC->D",
			args: args{
				wrapper: usecaseLayer(),
				target:  ErrorApp,
			},
			want: false,
		},
		{
			name: "success AB->A",
			args: args{
				wrapper: domainLayer(),
				target:  ErrorInfra,
			},
			want: true,
		},
		{
			name: "success AB->B",
			args: args{
				wrapper: domainLayer(),
				target:  ErrorDomain,
			},
			want: true,
		},
		{
			name: "error AB->C",
			args: args{
				wrapper: domainLayer(),
				target:  ErrorUseCase,
			},
			want: false,
		},
		{
			name: "error AB->D",
			args: args{
				wrapper: domainLayer(),
				target:  ErrorApp,
			},
			want: false,
		},
		{
			name: "error A->AB",
			args: args{
				wrapper: ErrorInfra,
				target:  domainLayer(),
			},
			want: false,
		},
		{
			name: "success ABCD->A",
			args: args{
				wrapper: appLayer(),
				target:  ErrorInfra,
			},
			want: true,
		},
		{
			name: "success ABCD->B",
			args: args{
				wrapper: appLayer(),
				target:  ErrorDomain,
			},
			want: true,
		},
		{
			name: "success ABCD->C",
			args: args{
				wrapper: appLayer(),
				target:  ErrorUseCase,
			},
			want: true,
		},
		{
			name: "success ABCD->D",
			args: args{
				wrapper: appLayer(),
				target:  ErrorApp,
			},
			want: true,
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
			name: "get ErrorInfra from appLayer",
			args: args{
				err:    appLayer(),
				target: ErrorInfra,
			},
			wantErr: assert.Error,
		},
		{
			name: "invalid get ErrorInfra from appLayer",
			args: args{
				err:    ErrorInfra,
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
