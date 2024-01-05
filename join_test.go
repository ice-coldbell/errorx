package errorx

import (
	"errors"
	"reflect"
	"testing"
)

func TestJoinReturnsNil(t *testing.T) {
	if err := Join(); err != nil {
		t.Errorf("Join() = %v, want nil", err)
	}
	if err := Join(nil); err != nil {
		t.Errorf("Join(nil) = %v, want nil", err)
	}
	if err := Join(nil, nil); err != nil {
		t.Errorf("Join(nil, nil) = %v, want nil", err)
	}
}

func TestJoin(t *testing.T) {
	err1 := New("err1")
	err2 := New("err2")

	for _, test := range []struct {
		errs []error
		want []error
	}{{
		errs: []error{err1},
		want: []error{err1},
	}, {
		errs: []error{err1, err2},
		want: []error{err1, err2},
	}, {
		errs: []error{err1, nil, err2},
		want: []error{err1, err2},
	}} {
		got := Join(test.errs...).Unwrap()
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Join(%v) = %v; want %v", test.errs, got, test.want)
		}
		if len(got) != cap(got) {
			t.Errorf("Join(%v) returns errors with len=%v, cap=%v; want len==cap", test.errs, len(got), cap(got))
		}
	}
}

func TestJoinExtendError(t *testing.T) {
	errSTD := errors.New("err std")
	err1 := New("err1")
	err2 := New("err2")
	errs := Join(err1, err2)

	for _, test := range []struct {
		errs []error
		want []error
	}{{
		errs: []error{errSTD},
		want: []error{errSTD},
	}, {
		errs: []error{errs},
		want: []error{err1, err2},
	}, {
		errs: []error{errSTD, errs},
		want: []error{errSTD, err1, err2},
	}} {
		got := Join(test.errs...).Error()
		want := Join(test.want...).Error()
		if got != want {
			t.Errorf("Join(%v).Error() = %q; want %q", test.errs, got, test.want)
		}
	}
}

func TestJoinErrorMethod(t *testing.T) {
	err1 := New("err1")
	err2 := New("err2")
	for _, test := range []struct {
		errs []error
		want string
	}{{
		errs: []error{err1},
		want: "err1",
	}, {
		errs: []error{err1, err2},
		want: "err1\nerr2",
	}, {
		errs: []error{err1, nil, err2},
		want: "err1\nerr2",
	}} {
		got := Join(test.errs...).Error()
		if got != test.want {
			t.Errorf("Join(%v).Error() = %q; want %q", test.errs, got, test.want)
		}
	}
}
