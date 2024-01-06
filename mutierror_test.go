package errorx

import (
	"testing"
)

func Test_customErrors_Is(t *testing.T) {
	err := New("test error")
	err2 := New("test error 2")

	type args struct {
		target error
	}
	tests := []struct {
		name string
		e    CustomErrors
		args args
		want bool
	}{
		{
			name: "target error is not multiple errors",
			e:    Join(err),
			args: args{
				target: err,
			},
			want: false,
		},
		{
			name: "number of errors is different",
			e:    Join(err),
			args: args{
				target: Join(err, err2),
			},
			want: false,
		},
		{
			name: "some errors do not match",
			e:    Join(err, err2, err2),
			args: args{
				target: Join(err, err2, err),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Is(tt.args.target); got != tt.want {
				t.Errorf("customErrors.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
