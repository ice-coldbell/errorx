package errorx

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	err1 := New("1")
	erra := wrapped{"wrap 2", err1}
	errb := wrapped{"wrap 3", erra}

	err3 := New("3")

	poser := &poser{"either 1 or 3", func(err error) bool {
		return err == err1 || err == err3
	}}

	testCases := []struct {
		err    error
		target error
		match  bool
	}{
		{nil, nil, true},
		{err1, nil, false},
		{err1, err1, true},
		{erra, err1, true},
		{errb, err1, true},
		{err1, err3, false},
		{erra, err3, false},
		{errb, err3, false},
		{poser, err1, true},
		{poser, err3, true},
		{poser, erra, false},
		{poser, errb, false},
		{errorUncomparable{}, errorUncomparable{}, true},
		{errorUncomparable{}, &errorUncomparable{}, false},
		{&errorUncomparable{}, errorUncomparable{}, true},
		{&errorUncomparable{}, &errorUncomparable{}, false},
		{errorUncomparable{}, err1, false},
		{&errorUncomparable{}, err1, false},
		{multiErr{}, err1, false},
		{multiErr{err1, err3}, err1, true},
		{multiErr{err3, err1}, err1, true},
		{multiErr{err1, err3}, New("x"), false},
		{multiErr{err3, errb}, errb, true},
		{multiErr{err3, errb}, erra, true},
		{multiErr{err3, errb}, err1, true},
		{multiErr{errb, err3}, err1, true},
		{multiErr{poser}, err1, true},
		{multiErr{poser}, err3, true},
		{multiErr{nil}, nil, false},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			if got := Is(tc.err, tc.target); got != tc.match {
				t.Errorf("Is(%v, %v) = %v, want %v", tc.err, tc.target, got, tc.match)
			}
		})
	}
}

type poser struct {
	msg string
	f   func(error) bool
}

var poserPathErr = &fs.PathError{Op: "poser"}

func (p *poser) Error() string     { return p.msg }
func (p *poser) Is(err error) bool { return p.f(err) }
func (p *poser) As(err any) bool {
	switch x := err.(type) {
	case **poser:
		*x = p
	case *errorT:
		*x = errorT{"poser"}
	case **fs.PathError:
		*x = poserPathErr
	default:
		return false
	}
	return true
}

func TestAs(t *testing.T) {
	var errT errorT
	var errP *fs.PathError
	var timeout interface{ Timeout() bool }
	var p *poser
	_, errF := os.Open("non-existing")
	poserErr := &poser{"oh no", nil}

	testCases := []struct {
		err    error
		target any
		match  bool
		want   any // value of target on match
	}{{
		nil,
		&errP,
		false,
		nil,
	}, {
		wrapped{"pitied the fool", errorT{"T"}},
		&errT,
		true,
		errorT{"T"},
	}, {
		errF,
		&errP,
		true,
		errF,
	}, {
		errorT{},
		&errP,
		false,
		nil,
	}, {
		wrapped{"wrapped", nil},
		&errT,
		false,
		nil,
	}, {
		&poser{"error", nil},
		&errT,
		true,
		errorT{"poser"},
	}, {
		&poser{"path", nil},
		&errP,
		true,
		poserPathErr,
	}, {
		poserErr,
		&p,
		true,
		poserErr,
	}, {
		New("err"),
		&timeout,
		false,
		nil,
	}, {
		errF,
		&timeout,
		true,
		errF,
	}, {
		wrapped{"path error", errF},
		&timeout,
		true,
		errF,
	}, {
		multiErr{},
		&errT,
		false,
		nil,
	}, {
		multiErr{New("a"), errorT{"T"}},
		&errT,
		true,
		errorT{"T"},
	}, {
		multiErr{errorT{"T"}, New("a")},
		&errT,
		true,
		errorT{"T"},
	}, {
		multiErr{errorT{"a"}, errorT{"b"}},
		&errT,
		true,
		errorT{"a"},
	}, {
		multiErr{multiErr{New("a"), errorT{"a"}}, errorT{"b"}},
		&errT,
		true,
		errorT{"a"},
	}, {
		multiErr{wrapped{"path error", errF}},
		&timeout,
		true,
		errF,
	}, {
		multiErr{nil},
		&errT,
		false,
		nil,
	}}
	for i, tc := range testCases {
		name := fmt.Sprintf("%d:As(Errorf(..., %v), %v)", i, tc.err, tc.target)
		// Clear the target pointer, in case it was set in a previous test.
		rtarget := reflect.ValueOf(tc.target)
		rtarget.Elem().Set(reflect.Zero(reflect.TypeOf(tc.target).Elem()))
		t.Run(name, func(t *testing.T) {
			match := As(tc.err, tc.target)
			if match != tc.match {
				t.Fatalf("match: got %v; want %v", match, tc.match)
			}
			if !match {
				return
			}
			if got := rtarget.Elem().Interface(); got != tc.want {
				t.Fatalf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestAsValidation(t *testing.T) {
	var s string
	testCases := []any{
		nil,
		(*int)(nil),
		"error",
		&s,
	}
	err := New("error")
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T(%v)", tc, tc), func(t *testing.T) {
			defer func() {
				recover()
			}()
			if As(err, tc) {
				t.Errorf("As(err, %T(%v)) = true, want false", tc, tc)
				return
			}
			t.Errorf("As(err, %T(%v)) did not panic", tc, tc)
		})
	}
}

func TestUnwrap(t *testing.T) {
	err1 := New("1")
	erra := wrapped{"wrap 2", err1}

	testCases := []struct {
		err  error
		want error
	}{
		{nil, nil},
		{wrapped{"wrapped", nil}, nil},
		// {err1, err1.err},
		{erra, err1},
		{wrapped{"wrap 3", erra}, erra},
	}
	for _, tc := range testCases {
		if got := Unwrap(tc.err); got != tc.want {
			t.Errorf("Unwrap(%v) = %v, want %v", tc.err, got, tc.want)
		}
	}
}

type errorT struct{ s string }

func (e errorT) Error() string { return fmt.Sprintf("errorT(%s)", e.s) }

type wrapped struct {
	msg string
	err error
}

func (e wrapped) Error() string { return e.msg }
func (e wrapped) Unwrap() error { return e.err }

type multiErr []error

func (m multiErr) Error() string   { return "multiError" }
func (m multiErr) Unwrap() []error { return []error(m) }

type errorUncomparable struct {
	f []string
}

func (errorUncomparable) Error() string {
	return "uncomparable error"
}

func (errorUncomparable) Is(target error) bool {
	_, ok := target.(errorUncomparable)
	return ok
}

func ExampleIs() {
	if _, err := os.Open("non-existing"); err != nil {
		if Is(err, fs.ErrNotExist) {
			fmt.Println("file does not exist")
		} else {
			fmt.Println(err)
		}
	}

	// Output:
	// file does not exist
}

func ExampleAs() {
	if _, err := os.Open("non-existing"); err != nil {
		var pathError *fs.PathError
		if As(err, &pathError) {
			fmt.Println("Failed at path:", pathError.Path)
		} else {
			fmt.Println(err)
		}
	}

	// Output:
	// Failed at path: non-existing
}

func ExampleUnwrap() {
	err1 := New("error1")
	err2 := fmt.Errorf("error2: [%w]", err1)
	fmt.Println(err2)
	fmt.Println(Unwrap(err2))
	// Output
	// error2: [error1]
	// error1
}

func TestWrap(t *testing.T) {
	const errMessage = "test wrap error"

	err := errors.New(errMessage)
	wrapErr := Wrap(err)
	assert.Equal(t, errMessage, wrapErr.Error())
	assert.Equal(t, err, wrapErr.Unwrap())

	const (
		key  = "test_key"
		data = "test_data"
	)
	cErr := wrapErr.With(key, data)
	val, ok := cErr.(*customError)
	assert.True(t, ok)
	assert.Equal(t, data, val.data[key])
}

func TestWrapWithCustomError(t *testing.T) {
	const errMessage = "test wrap error"

	err := New(errMessage)
	wrapErr := Wrap(err)
	assert.Equal(t, errMessage, wrapErr.Error())
	assert.Equal(t, err, wrapErr)

	const (
		key  = "test_key"
		data = "test_data"
	)
	cErr := wrapErr.With(key, data)
	val, ok := cErr.(*customError)
	assert.True(t, ok)
	assert.Equal(t, data, val.data[key])
}

func TestWrapWithData(t *testing.T) {
	const errMessage = "test wrap error"

	err := errors.New(errMessage)

	wrapErr := Wrap(err).With("foo", "bar")
	wrapDataErr := WrapWithData(err, map[string]any{"foo": "bar"})

	assert.Equal(t, wrapErr, wrapDataErr)
}

// Nil Error Case
func TestWraprWithData2(t *testing.T) {
	err := WrapWithData(nil, map[string]any{"foo": "bar"})
	assert.Equal(t, err, nil)
}

// Custom Error Case
func TestWrapWithData3(t *testing.T) {
	const errMessage = "test wrap error"

	err := New(errMessage).With("foo", "bar")
	wrapErr := WrapWithData(err, map[string]any{"foo": "bar"})
	assert.Equal(t, err, wrapErr)
}

func TestWrapReturnsNil(t *testing.T) {
	if err := Wrap(nil); err != nil {
		t.Errorf("Wrap(nil) = %v, want nil", err)
	}
}

func TestNil(t *testing.T) {
	nilErr := Wrap(nil)
	assert.Nil(t, nilErr)
}
