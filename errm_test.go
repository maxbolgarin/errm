package errm_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/maxbolgarin/errm"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		id  string
		err error
		exp string
	}{
		{
			id:  "empty",
			err: errm.New(""),
			exp: "",
		},
		{
			id:  "simple",
			err: errm.New("some-err"),
			exp: "some-err",
		},
		{
			id:  "field",
			err: errm.New("some-err", "field", "value"),
			exp: "some-err field=value",
		},
		{
			id:  "many_fields",
			err: errm.New("some-err", "field", "value", "field2", []any{123, 321}, "field3", 123, "field4"),
			exp: "some-err field=value field2=[123 321] field3=123",
		},
	}

	for _, test := range testCases {
		t.Run(test.id, func(t *testing.T) {
			if test.err.Error() != test.exp {
				t.Errorf("expected %s, got %s", test.exp, test.err)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	testCases := []struct {
		id  string
		err error
		exp string
	}{
		{
			id:  "empty",
			err: errm.Errorf(""),
			exp: "",
		},
		{
			id:  "format_simple",
			err: errm.Errorf("some-err"),
			exp: "some-err",
		},
		{
			id:  "format",
			err: errm.Errorf("some-err %s", "a"),
			exp: "some-err a",
		},
		{
			id:  "format_many_fields",
			err: errm.Errorf("some-err", "field", "value", "field2", []any{123, 321}, "field3", 123, "field4"),
			exp: "some-err field=value field2=[123 321] field3=123",
		},
		{
			id:  "format_many_fields_2",
			err: errm.Errorf("some-err %s %d", "a", 1, "field", "value", "field2", []any{123, 321}, "field3", 123, "field4"),
			exp: "some-err a 1 field=value field2=[123 321] field3=123",
		},
	}

	for _, test := range testCases {
		t.Run(test.id, func(t *testing.T) {
			if test.err.Error() != test.exp {
				t.Errorf("expected %s, got %s", test.exp, test.err)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	err := errm.Errorf("some-err %s %d", "a", 1, "field", "value", "field2", []any{123, 321}, "field3", 123, "field4")
	exp := "some-err a 1 field=value field2=[123 321] field3=123"
	if err.Error() != exp {
		t.Errorf("expected %s, got %s", exp, err)
	}

	err = errm.Wrap(err, "second error", "field", 123, "testtt", "122", "asd")
	exp = "second error field=123 testtt=122: " + exp
	if err.Error() != exp {
		t.Errorf("expected %s, got %s", exp, err)
	}

	err = errm.Wrapf(err, "third error with %s", "format", "v", "a", "testtt")
	exp = "third error with format v=a: " + exp
	if err.Error() != exp {
		t.Errorf("expected %s, got %s", exp, err)
	}
}

func TestIs(t *testing.T) {
	err := errm.Errorf("some-err %s %d", "a", 1, "field", "value", "field2", []any{123, 321}, "field3", 123, "field4")
	f1 := func() error {
		return err
	}
	f2 := func() error {
		return errm.Wrap(err, "second error", "field", 123, "testtt", "122", "asd")
	}
	if !errm.Is(f1(), err) {
		t.Errorf("expected true, got %t", errm.Is(f1(), err))
	}
	if !errm.Is(f2(), err) {
		t.Errorf("expected true, got %t", errm.Is(f2(), err))
	}

	base := errors.New("common-error")
	f3 := func() error {
		return errm.Wrap(base, "second error", "field", 123, "testtt", "122", "asd")
	}
	if !errm.Is(f3(), base) {
		t.Errorf("expected true, got %t", errm.Is(f3(), base))
	}

	f4 := func() error {
		return fmt.Errorf("second error: %w", base)
	}
	if !errm.Is(f4(), base) {
		t.Errorf("expected true, got %t", errm.Is(f4(), base))
	}

	// Errorf doesn't wrap in eris because it is using Sprintf
	f5 := func() error {
		return errm.Errorf("second error: %w", base)
	}
	if errm.Is(f5(), base) {
		t.Errorf("expected false, got %t", errm.Is(f5(), base))
	}

	wrapped := errm.Wrap(base, "some-error")
	f6 := func() error {
		return errm.Wrap(wrapped, "second error")
	}
	if !errm.Is(f6(), base) {
		t.Errorf("expected true, got %t", errm.Is(f6(), base))
	}

	f7 := func() error {
		return fmt.Errorf("second error: %w", wrapped)
	}
	if !errm.Is(f7(), base) {
		t.Errorf("expected true, got %t", errm.Is(f7(), base))
	}

	fmtWrap := fmt.Errorf("second error: %w", wrapped)
	f8 := func() error {
		return errm.Wrap(fmtWrap, "second error")
	}
	if !errm.Is(f8(), base) {
		t.Errorf("expected true, got %t", errm.Is(f8(), base))
	}

	fmtWrap2 := fmt.Errorf("second error: %w", fmtWrap)
	f9 := func() error {
		return fmt.Errorf("second error: %w", fmtWrap2)
	}
	if !errm.Is(f9(), fmtWrap) {
		t.Errorf("expected true, got %t", errm.Is(f9(), fmtWrap))
	}
	if !errm.Is(f9(), fmtWrap, errm.New("ABC")) {
		t.Errorf("expected true, got %t", errm.Is(f9(), fmtWrap, errm.New("ABC")))
	}
	if !errm.Is(f9(), errm.New("ABC"), fmtWrap) {
		t.Errorf("expected true, got %t", errm.Is(f9(), errm.New("ABC"), fmtWrap))
	}
	if !errm.Is(f9(), wrapped, fmtWrap) {
		t.Errorf("expected true, got %t", errm.Is(f9(), wrapped, fmtWrap))
	}
	if !errm.Is(f9(), fmtWrap, wrapped) {
		t.Errorf("expected true, got %t", errm.Is(f9(), fmtWrap, wrapped))
	}
	if !errm.Is(f9(), errm.New("ABC"), errm.New("ABC"), errm.New("ABC"), errm.New("ABC"),
		errm.New("ABC"), errm.New("ABC"), errm.New("ABC"), wrapped) {
		t.Errorf("expected true, got %t", errm.Is(f9(), errm.New("ABC"), errm.New("ABC"),
			errm.New("ABC"), errm.New("ABC"), errm.New("ABC"), errm.New("ABC"), errm.New("ABC"), wrapped))
	}
	if errm.Is(f9(), errm.New("ABC"), errm.New("ABC"), errm.New("ABC"), errm.New("ABC"),
		errm.New("ABC"), errm.New("ABC"), errm.New("ABC")) {
		t.Errorf("expected false, got %t", errm.Is(f9(), errm.New("ABC"), errm.New("ABC"),
			errm.New("ABC"), errm.New("ABC"), errm.New("ABC"), errm.New("ABC"), errm.New("ABC")))
	}

	if errm.Is(fmtWrap, nil) {
		t.Errorf("expected false, got %t", errm.Is(fmtWrap, nil))
	}
	if errm.Is(nil, fmtWrap) {
		t.Errorf("expected false, got %t", errm.Is(nil, fmtWrap))
	}
	if !errm.Is(nil, nil) {
		t.Errorf("expected true, got %t", errm.Is(nil, nil))
	}
}

func TestIsSet(t *testing.T) {
	err := errors.New("A")
	f1 := func() error {
		s := errm.NewSet()
		s.Add(err)
		s.Clear()
		s.Add(err)
		return s.Err()
	}
	if !errm.Is(f1(), err) {
		t.Errorf("expected true, got %t", errm.Is(f1(), err))
	}

	f2 := func() error {
		s := errm.NewSet()
		s.Add(err)
		s.Add(err)
		return s.Err()
	}
	if !errm.Is(f2(), err) {
		t.Errorf("expected true, got %t", errm.Is(f2(), err))
	}

	f3 := func() error {
		s := errm.NewSet()
		s.Add(errm.New("B"))
		return s.Err()
	}
	if errm.Is(f3(), err) {
		t.Errorf("expected false, got %t", errm.Is(f3(), err))
	}

	f4 := func() error {
		s := errm.NewSet()
		s.Wrap(err, "C")
		s.Add(errm.New("D"))
		return s.Err()
	}
	if !errm.Is(f4(), err) {
		t.Errorf("expected true, got %t", errm.Is(f4(), err))
	}

	err2 := errm.Wrap(err, "abc")
	f5 := func() error {
		s := errm.NewSet()
		s.Wrap(err2, "C")
		s.Add(errm.New("D"))
		return s.Err()
	}
	if !errm.Is(f5(), err) {
		t.Errorf("expected true, got %t", errm.Is(f5(), err))
	}
	if !errm.Is(f5(), err2) {
		t.Errorf("expected true, got %t", errm.Is(f5(), err2))
	}

	err3 := errm.Errorf("AAA")
	f6 := func() error {
		s := errm.NewSet()
		s.Wrap(err, "FF")
		s.Wrap(err3, "F")
		return s.Err()
	}
	if !errm.Is(f6(), err) {
		t.Errorf("expected true, got %t", errm.Is(f6(), err))
	}
	if !errm.Is(f6(), err3) {
		t.Errorf("expected true, got %t", errm.Is(f6(), err3))
	}
	if errm.Is(f6(), err2) {
		t.Errorf("expected false, got %t", errm.Is(f6(), err2))
	}

	if !errm.Is(f6(), err2, err3) {
		t.Errorf("expected true, got %t", errm.Is(f6(), err2, err3))
	}
	if !errm.Is(f6(), err3, err2) {
		t.Errorf("expected true, got %t", errm.Is(f6(), err3, err2))
	}
	if errm.Is(f6(), err2, err2, err2, err2, err2, err2, err2, err2) {
		t.Errorf("expected false, got %t", errm.Is(f6(), err2, err2, err2, err2, err2, err2, err2, err2))
	}
	if !errm.Is(f6(), err2, err2, err2, err2, err2, err2, err2, err2, err3) {
		t.Errorf("expected true, got %t", errm.Is(f6(), err2, err2, err2, err2, err2, err2, err2, err2, err3))
	}
}

func TestIsList(t *testing.T) {
	err := errors.New("A")
	f1 := func() error {
		s := errm.NewList()
		s.Add(err)
		return s.Err()
	}
	if !errm.Is(f1(), err) {
		t.Errorf("expected true, got %t", errm.Is(f1(), err))
	}

	f2 := func() error {
		s := errm.NewList()
		s.Add(err)
		s.Add(err)
		return s.Err()
	}
	if !errm.Is(f2(), err) {
		t.Errorf("expected true, got %t", errm.Is(f2(), err))
	}

	f3 := func() error {
		s := errm.NewList()
		s.Add(errm.New("B"))
		return s.Err()
	}
	if errm.Is(f3(), err) {
		t.Errorf("expected false, got %t", errm.Is(f3(), err))
	}

	f4 := func() error {
		s := errm.NewList()
		s.Wrap(err, "C")
		s.Add(errm.New("D"))
		return s.Err()
	}
	if !errm.Is(f4(), err) {
		t.Errorf("expected true, got %t", errm.Is(f4(), err))
	}

	err2 := errm.Wrap(err, "abc")
	f5 := func() error {
		s := errm.NewList()
		s.Wrap(err2, "C")
		s.Add(errm.New("D"))
		return s.Err()
	}
	if !errm.Is(f5(), err) {
		t.Errorf("expected true, got %t", errm.Is(f5(), err))
	}
	if !errm.Is(f5(), err2) {
		t.Errorf("expected true, got %t", errm.Is(f5(), err2))
	}

	err3 := errm.Errorf("AAA")
	f6 := func() error {
		s := errm.NewList()
		s.Wrap(err, "FF")
		s.Wrap(err3, "F")
		return s.Err()
	}
	if !errm.Is(f6(), err) {
		t.Errorf("expected true, got %t", errm.Is(f6(), err))
	}
	if !errm.Is(f6(), err3) {
		t.Errorf("expected true, got %t", errm.Is(f6(), err3))
	}
	if errm.Is(f6(), err2) {
		t.Errorf("expected false, got %t", errm.Is(f6(), err2))
	}

	if !errm.Is(f6(), err2, err3) {
		t.Errorf("expected true, got %t", errm.Is(f6(), err2, err3))
	}
	if !errm.Is(f6(), err3, err2) {
		t.Errorf("expected true, got %t", errm.Is(f6(), err3, err2))
	}
	if errm.Is(f6(), err2, err2, err2, err2, err2, err2, err2, err2) {
		t.Errorf("expected false, got %t", errm.Is(f6(), err2, err2, err2, err2, err2, err2, err2, err2))
	}
	if !errm.Is(f6(), err2, err2, err2, err2, err2, err2, err2, err2, err3) {
		t.Errorf("expected true, got %t", errm.Is(f6(), err2, err2, err2, err2, err2, err2, err2, err2, err3))
	}
}

func TestContains(t *testing.T) {
	err := errm.Errorf("some-err %s %d", "a", 1, "field", "value", "field2", []any{123, 321}, "field3", 123, "field4")

	if !errm.Contains(err, "some-err") {
		t.Errorf("expected true, got %t", errm.Contains(err, "some-err"))
	}
	if !errm.Contains(err, "field3") {
		t.Errorf("expected true, got %t", errm.Contains(err, "field3"))
	}
	if errm.Contains(err, "field4") {
		t.Errorf("expected false, got %t", errm.Contains(err, "field4"))
	}
	if errm.Contains(err, "another-err") {
		t.Errorf("expected false, got %t", errm.Contains(err, "another-err"))
	}

	anotherErr := errm.Wrap(err, "another-err")
	if !errm.Contains(anotherErr, "some-err") {
		t.Errorf("expected true, got %t", errm.Contains(anotherErr, "some-err"))
	}
	if !errm.Contains(anotherErr, "field3") {
		t.Errorf("expected true, got %t", errm.Contains(anotherErr, "field3"))
	}
	if errm.Contains(anotherErr, "field4") {
		t.Errorf("expected false, got %t", errm.Contains(anotherErr, "field4"))
	}
	if !errm.Contains(anotherErr, "another-err") {
		t.Errorf("expected true, got %t", errm.Contains(anotherErr, "another-err"))
	}
}

func TestContainsErr(t *testing.T) {
	err := errm.Errorf("some-err %s %d", "a", 1, "field", "value", "field2", []any{123, 321}, "field3", 123, "field4")

	if !errm.ContainsErr(err, err) {
		t.Errorf("expected true, got %t", errm.ContainsErr(err, err))
	}
	if !errm.ContainsErr(err, errors.New("field3")) {
		t.Errorf("expected true, got %t", errm.ContainsErr(err, errors.New("field3")))
	}
	if errm.ContainsErr(err, fmt.Errorf("field4")) {
		t.Errorf("expected false, got %t", errm.ContainsErr(err, fmt.Errorf("field4")))
	}

	anotherErr := errm.Wrap(err, "another-err")
	if errm.ContainsErr(err, anotherErr) {
		t.Errorf("expected false, got %t", errm.ContainsErr(err, anotherErr))
	}

	if !errm.ContainsErr(anotherErr, err) {
		t.Errorf("expected true, got %t", errm.ContainsErr(anotherErr, err))
	}
	if !errm.ContainsErr(anotherErr, errors.New("field3")) {
		t.Errorf("expected true, got %t", errm.ContainsErr(anotherErr, errors.New("field3")))
	}
	if errm.ContainsErr(anotherErr, fmt.Errorf("field4")) {
		t.Errorf("expected false, got %t", errm.ContainsErr(anotherErr, fmt.Errorf("field4")))
	}
	if !errm.ContainsErr(anotherErr, anotherErr) {
		t.Errorf("expected true, got %t", errm.ContainsErr(anotherErr, anotherErr))
	}
}

func TestJoinErrors(t *testing.T) {
	err1 := errm.New("first error")
	err2 := errm.New("second error")
	err3 := errm.New("third error")
	err4 := errm.New("")
	var err5 error = nil

	t.Run("TestNoErrors", func(t *testing.T) {
		err := errm.JoinErrors()
		if err != nil {
			t.Errorf("expected nil, got %s", err)
		}
	})

	t.Run("TestNoErrors2", func(t *testing.T) {
		err := errm.JoinErrors(err4, err5)
		if err != nil {
			t.Errorf("expected nil, got %s", err)
		}
	})

	t.Run("TestSingleError", func(t *testing.T) {
		err := errm.JoinErrors(err1)
		expectedErr := errm.New("first error")
		if err.Error() != expectedErr.Error() {
			t.Errorf("expected %s, got %s", expectedErr, err)
		}
	})

	t.Run("TestMultipleErrors", func(t *testing.T) {
		err := errm.JoinErrors(err1, err2)
		expectedErr := errm.New("first error; second error")
		if err.Error() != expectedErr.Error() {
			t.Errorf("expected %s, got %s", expectedErr, err)
		}

		err = errm.JoinErrors(err1, err2, err3)
		expectedErr = errm.New("first error; second error; third error")
		if err.Error() != expectedErr.Error() {
			t.Errorf("expected %s, got %s", expectedErr, err)
		}
	})
}
func TestSet(t *testing.T) {
	s := errm.NewSet()
	if s.Len() != 0 {
		t.Errorf("expected 0, got %d", s.Len())
	}

	err := errm.New("A")

	s.Add(err)
	if s.Len() != 1 {
		t.Errorf("expected 1, got %d", s.Len())
	}

	s.Add(err)
	if s.Len() != 1 {
		t.Errorf("expected 1, got %d", s.Len())
	}

	s.Add(errm.New("A"))
	if s.Len() != 1 {
		t.Errorf("expected 1, got %d", s.Len())
	}

	s.Add(errm.Wrap(err, "B"))
	if s.Len() != 2 {
		t.Errorf("expected 2, got %d", s.Len())
	}
}

func TestList(t *testing.T) {
	s := errm.NewList()
	if s.Len() != 0 {
		t.Errorf("expected 0, got %d", s.Len())
	}

	err := errm.New("A")

	s.Add(err)
	if s.Len() != 1 {
		t.Errorf("expected 1, got %d", s.Len())
	}

	s.Add(err)
	if s.Len() != 2 {
		t.Errorf("expected 2, got %d", s.Len())
	}

	s.Add(errm.Wrap(err, "B"))
	if s.Len() != 3 {
		t.Errorf("expected 3, got %d", s.Len())
	}
}
