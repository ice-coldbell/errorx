package errorx

import (
	"errors"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	runtimeCallers = func(skip int, pc []uintptr) int {
		pc[0] = globalTestPC
		return 1
	}

	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want CustomError
	}{
		{
			name: "common case",
			args: args{
				message: "test error",
			},
			want: &customError{
				err:   errors.New("test error"),
				stack: []frame{globalTestFrame},
				data:  make(map[string]any),
			},
		},
		{
			name: "empty error message",
			args: args{
				message: "",
			},
			want: &customError{
				err:   errors.New(""),
				stack: []frame{globalTestFrame},
				data:  make(map[string]any),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_customError_Error(t *testing.T) {
	runtimeCallers = func(skip int, pc []uintptr) int {
		pc[0] = globalTestPC
		return 1
	}

	type fields struct {
		err   error
		stack stack
		data  map[string]any
	}
	tests := []struct {
		name       string
		fields     fields
		want       string
		isNilError bool
	}{
		{
			name: "common case",
			fields: fields{
				err:   errors.New("common error"),
				stack: []frame{globalTestFrame},
				data:  make(map[string]any),
			},
			isNilError: false,
			want:       "common error",
		},
		{
			name:       "internal nil error",
			fields:     fields{err: nil, stack: nil, data: make(map[string]any)},
			isNilError: false,
			want:       "",
		},
		{
			name:       "nil error",
			isNilError: true,
			want:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &customError{
				err:   tt.fields.err,
				stack: tt.fields.stack,
				data:  tt.fields.data,
			}
			if tt.isNilError {
				e = nil
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("customError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_customError_Cause(t *testing.T) {
	runtimeCallers = func(skip int, pc []uintptr) int {
		pc[0] = globalTestPC
		return 1
	}
	commonErr := errors.New("common error")

	type fields struct {
		err   error
		stack stack
		data  map[string]any
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    error
		isNilError bool
	}{
		{
			name: "common case",
			fields: fields{
				err:   commonErr,
				stack: []frame{globalTestFrame},
				data:  make(map[string]any),
			},
			wantErr:    commonErr,
			isNilError: false,
		},
		{
			name:       "internal nil error",
			fields:     fields{err: nil, stack: nil, data: make(map[string]any)},
			wantErr:    nil,
			isNilError: false,
		},
		{
			name:       "nil error",
			wantErr:    nil,
			isNilError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &customError{
				err:   tt.fields.err,
				stack: tt.fields.stack,
				data:  tt.fields.data,
			}
			if tt.isNilError {
				e = nil
			}
			if err := e.Cause(); !errors.Is(err, tt.wantErr) {
				t.Errorf("customError.Cause() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_customError_With(t *testing.T) {
	runtimeCallers = func(skip int, pc []uintptr) int {
		pc[0] = globalTestPC
		return 1
	}
	commonErr := errors.New("common error")
	wantErr := New("common error").With("foo", "bar")

	type fields struct {
		err   error
		stack stack
		data  map[string]any
	}
	type args struct {
		key  string
		data any
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       CustomError
		isNilError bool
	}{
		{
			name: "common case",
			fields: fields{
				err:   commonErr,
				stack: []frame{globalTestFrame},
				data:  make(map[string]any),
			},
			args:       args{key: "foo", data: "bar"},
			want:       wantErr,
			isNilError: false,
		},
		{
			name:       "internal nil error",
			fields:     fields{err: nil, stack: nil, data: make(map[string]any)},
			args:       args{key: "foo", data: "bar"},
			want:       nil,
			isNilError: false,
		},
		{
			name:       "nil error",
			want:       nil,
			isNilError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &customError{
				err:   tt.fields.err,
				stack: tt.fields.stack,
				data:  tt.fields.data,
			}
			if tt.isNilError {
				e = nil
			}
			if got := e.With(tt.args.key, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("customError.With() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_customError_WithData(t *testing.T) {
	runtimeCallers = func(skip int, pc []uintptr) int {
		pc[0] = globalTestPC
		return 1
	}

	type fields struct {
		err   error
		stack stack
		data  map[string]any
	}
	type args struct {
		data map[string]any
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       CustomError
		isNilError bool
	}{
		{
			name: "common case",
			fields: fields{
				err:   errors.New("common error"),
				stack: []frame{globalTestFrame},
				data:  make(map[string]any),
			},
			args: args{data: map[string]any{"foo": "bar"}},
			want: New("common error").WithData(map[string]any{"foo": "bar"}),
		},
		{
			name:       "internal nil error",
			fields:     fields{err: nil, stack: nil, data: nil},
			args:       args{data: map[string]any{"foo": "bar"}},
			want:       nil,
			isNilError: false,
		},
		{
			name:       "nil error",
			args:       args{data: map[string]any{"foo": "bar"}},
			want:       nil,
			isNilError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &customError{
				err:   tt.fields.err,
				stack: tt.fields.stack,
				data:  tt.fields.data,
			}
			if tt.isNilError {
				e = nil
			}
			if got := e.WithData(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("customError.WithData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_customError_StackTrace(t *testing.T) {
	type fields struct {
		err   error
		stack stack
		data  map[string]any
	}
	tests := []struct {
		name       string
		fields     fields
		want       stack
		isNilError bool
	}{
		{
			name:       "common case",
			fields:     fields{err: errors.New("common error"), stack: []frame{globalTestFrame}},
			want:       []frame{globalTestFrame},
			isNilError: false,
		},
		{
			name:       "internal nil error",
			fields:     fields{err: nil, stack: nil, data: nil},
			want:       nil,
			isNilError: false,
		},
		{
			name:       "nil error",
			want:       nil,
			isNilError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &customError{
				err:   tt.fields.err,
				stack: tt.fields.stack,
				data:  tt.fields.data,
			}
			if tt.isNilError {
				e = nil
			}
			if got := e.StackTrace(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("customError.StackTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_customError_Unwrap(t *testing.T) {
	type fields struct {
		err   error
		stack stack
		data  map[string]any
	}
	commonErr := errors.New("common error")

	tests := []struct {
		name       string
		fields     fields
		want       error
		isNilError bool
	}{
		{
			name:   "common case",
			fields: fields{err: commonErr},
			want:   commonErr,
		},
		{
			name:       "internal nil error",
			fields:     fields{err: nil, stack: nil, data: nil},
			want:       nil,
			isNilError: false,
		},
		{
			name:       "nil error",
			want:       nil,
			isNilError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &customError{
				err:   tt.fields.err,
				stack: tt.fields.stack,
				data:  tt.fields.data,
			}
			if tt.isNilError {
				e = nil
			}
			if err := e.Unwrap(); !errors.Is(err, tt.want) {
				t.Errorf("customError.Unwrap() error = %v, want %v", err, tt.want)
			}
		})
	}
}
