package errorx

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"testing"
)

var (
	globalTestPC, _, _, _          = runtime.Caller(0)
	globalTestFunction             = runtime.FuncForPC(globalTestPC)
	globalTestFile, globalTestLine = globalTestFunction.FileLine(globalTestPC)
	globalTestFrame                = frameForPC(globalTestPC)
)

func Test_frameForPC(t *testing.T) {
	type args struct {
		pc uintptr
	}
	tests := []struct {
		name string
		args args
		want frame
	}{
		{
			name: "global test case",
			args: args{pc: globalTestPC},
			want: frame{
				pc:       globalTestPC,
				file:     globalTestFile,
				function: globalTestFunction.Name(),
				line:     globalTestLine,
			},
		},
		{
			name: "zero pc",
			args: args{pc: uintptr(0)},
			want: frame{
				pc:       uintptr(0),
				file:     "",
				function: "",
				line:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := frameForPC(tt.args.pc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("frameForPC() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_frame_Format(t *testing.T) {
	type fields struct {
		pc       uintptr
		file     string
		function string
		line     int
	}
	type args struct {
		s    fmt.State
		verb rune
	}
	tests := []struct {
		name   string
		fields fields
		format string
		want   string
	}{
		{
			name: "verb s", // Format file base path
			fields: fields{
				pc:       globalTestPC,
				file:     globalTestFile,
				function: globalTestFunction.Name(),
				line:     globalTestLine,
			},
			format: "%s",
			want:   path.Base(globalTestFile),
		},
		{
			name: "verb +s", // Format function and file path
			fields: fields{
				pc:       globalTestPC,
				file:     globalTestFile,
				function: globalTestFunction.Name(),
				line:     globalTestLine,
			},
			format: "%+s",
			want:   fmt.Sprint(globalTestFunction.Name(), "\n\t", globalTestFile),
		},
		{
			name: "verb d", // Format line number
			fields: fields{
				pc:       globalTestPC,
				file:     globalTestFile,
				function: globalTestFunction.Name(),
				line:     globalTestLine,
			},
			format: "%d",
			want:   strconv.Itoa(globalTestLine),
		},
		{
			name: "verb n", // Format function
			fields: fields{
				pc:       globalTestPC,
				file:     globalTestFile,
				function: globalTestFunction.Name(),
				line:     globalTestLine,
			},
			format: "%n",
			want:   globalTestFunction.Name(),
		},
		{
			name: "verb v", // Format function, file path and line number
			fields: fields{
				pc:       globalTestPC,
				file:     globalTestFile,
				function: globalTestFunction.Name(),
				line:     globalTestLine,
			},
			format: "%v",
			want:   path.Base(globalTestFile) + ":" + strconv.Itoa(globalTestLine),
		},
		{
			name: "verb #v", // Format Syntax
			fields: fields{
				pc:       globalTestPC,
				file:     globalTestFile,
				function: globalTestFunction.Name(),
				line:     globalTestLine,
			},
			format: "%#v",
			want: fmt.Sprintf(
				"errorx.frame{pc: %#x, file: %s, function:%s, line:%d}",
				globalTestPC,
				globalTestFile,
				globalTestFunction.Name(),
				globalTestLine,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := frame{
				pc:       tt.fields.pc,
				file:     tt.fields.file,
				function: tt.fields.function,
				line:     tt.fields.line,
			}
			if got := fmt.Sprintf(tt.format, f); got != tt.want {
				t.Errorf("fmt.Sprintf(\"%s\") = %#v, want %s", tt.format, got, tt.want)
			}
		})
	}
}

func Test_frame_MarshalText(t *testing.T) {
	type fields struct {
		pc       uintptr
		file     string
		function string
		line     int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "global test case",
			fields: fields{
				pc:       globalTestPC,
				file:     globalTestFile,
				function: globalTestFunction.Name(),
				line:     globalTestLine,
			},
			want: []byte(fmt.Sprintf(
				"%s %s:%d",
				globalTestFunction.Name(),
				globalTestFile,
				globalTestLine),
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := frame{
				pc:       tt.fields.pc,
				file:     tt.fields.file,
				function: tt.fields.function,
				line:     tt.fields.line,
			}
			got, err := f.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("frame.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("frame.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stack_Format(t *testing.T) {
	type args struct {
		st   fmt.State
		verb rune
	}
	tests := []struct {
		name   string
		s      stack
		format string
		want   string
	}{
		{
			name:   "verb s",
			s:      []frame{globalTestFrame},
			format: "%s",
			want:   fmt.Sprintf("[%s]", globalTestFrame),
		},
		{
			name:   "verb s, double frame",
			s:      []frame{globalTestFrame, globalTestFrame},
			format: "%s",
			want:   fmt.Sprintf("[%s %s]", globalTestFrame, globalTestFrame),
		},
		{
			name:   "verb v",
			s:      []frame{globalTestFrame},
			format: "%v",
			want:   fmt.Sprintf("[%s:%d]", globalTestFrame, globalTestLine),
		},
		{
			name:   "verb +v",
			s:      []frame{globalTestFrame},
			format: "%+v",
			want:   fmt.Sprintf("\n%+v", globalTestFrame),
		},
		{
			name:   "verb #v", // Format Syntax
			s:      []frame{globalTestFrame},
			format: "%#v",
			want:   fmt.Sprintf("[]errorx.frame{\n\t%#v,\n}", globalTestFrame),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmt.Sprintf(tt.format, tt.s); got != tt.want {
				t.Errorf("fmt.Sprintf(\"%s\") = %#v, want %s", tt.format, got, tt.want)
			}
		})
	}
}

func Test_callers(t *testing.T) {
	runtimeCallers = func(skip int, pc []uintptr) int {
		pc[0] = globalTestPC
		return 1
	}

	type args struct {
		skip int
	}
	tests := []struct {
		name string
		args args
		want stack
	}{
		{
			name: "global test case",
			args: args{
				skip: 0,
			},
			want: []frame{globalTestFrame},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := callers(tt.args.skip); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("callers() = %v, want %v", got, tt.want)
			}
		})
	}
}
