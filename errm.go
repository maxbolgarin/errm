// Package errm is a library for convinient and easy use of errors in a golang code.
package errm

import (
	"fmt"
	"io"
	"strings"

	"github.com/rotisserie/eris"
)

type errorImpl struct {
	err error
}

func newError(err error) errorImpl {
	return errorImpl{err: err}
}

// Error implements error interface, it just returns error message with applied fields in field=val format.
func (e errorImpl) Error() string {
	return e.err.Error()
}

// String is a wrapper of Error method.
func (e errorImpl) String() string {
	return e.Error()
}

// StackForLogger returns slice ["stack", "[...]"] that can be used as fields for logger if you want to log stack trace.
func (e errorImpl) StackForLogger() []any {
	jsonErr := ToJSON(e.err)
	root, ok := jsonErr["root"].(map[string]any)
	if !ok {
		return nil
	}
	return []any{"stack", root["stack"]}
}

// Format is used to handle %+v in formatted print, that will print stack trace.
func (e errorImpl) Format(s fmt.State, verb rune) {
	var withTrace bool
	switch verb {
	case 'v':
		if s.Flag('+') {
			withTrace = true
		}
	default:
		break
	}
	str := eris.ToString(e.err, withTrace)
	_, _ = io.WriteString(s, str)
}

// New creates a new error with a static message and pairs of fields in a field=val format.
func New(msg string, fields ...any) error {
	return newError(eris.New(buildErrorMessage(msg, fields)))
}

// Errorf creates a new error with a formatted message and pairs of fields in a field=val format.
func Errorf(msg string, args ...any) error {
	args, fields := separateArgsAndFields(msg, args)
	if len(args) == 0 {
		return New(msg, fields...)
	}
	return newError(eris.Errorf(buildErrorMessage(msg, fields), args...))
}

// Wrap adds additional context to all error types while maintaining the type of the original error;
// It also adds pairs of fields in a field=val format to message.
func Wrap(err error, msg string, fields ...any) error {
	if err == nil {
		return New(msg, fields...)
	}
	return newError(eris.Wrap(unwrap(err), buildErrorMessage(msg, fields)))
}

// Wrapf adds additional context to all error types while maintaining the type of the original error;
// It waits for formatted input and also adds pairs of fields in a field=val format to message.
func Wrapf(err error, msg string, args ...any) error {
	if err == nil {
		return Errorf(msg, args...)
	}
	args, fields := separateArgsAndFields(msg, args)
	if len(args) == 0 {
		return Wrap(err, msg, args...)
	}
	return newError(eris.Wrapf(unwrap(err), buildErrorMessage(msg, fields), args...))
}

// Is reports whether any error in err's chain matches target.
func Is(err, target error, targets ...error) bool {
	var set setError
	if eris.As(err, &set) {
		return set.Has(target, targets...)
	}
	var list listError
	if eris.As(err, &list) {
		return list.Has(target, targets...)
	}

	res := eris.Is(unwrap(err), unwrap(target))
	if !res && len(targets) > 0 {
		for _, t := range targets {
			if eris.Is(unwrap(err), unwrap(t)) {
				return true
			}
		}
	}
	return res
}

// Contains reports whether any error in err's chain contains target string.
func Contains(err error, target string) bool {
	return err != nil && strings.Contains(eris.ToString(unwrap(err), false), target)
}

// ContainsErr reports whether any error in err's chain contains target string representation.
func ContainsErr(err, target error) bool {
	return target != nil && Contains(err, eris.ToString(unwrap(target), false))
}

// ToJSON returns a JSON formatted map for a given error.
func ToJSON(err error) map[string]any {
	return eris.ToJSON(unwrap(err), true)
}

// StackForLogger returns slice ["stack", "[...]"] that can be used as fields for logger if you want to log stack trace.
func StackForLogger(err error) []any {
	jsonErr := ToJSON(err)
	root, ok := jsonErr["root"].(map[string]any)
	if !ok {
		return nil
	}
	return []any{"stack", root["stack"]}
}

// Check returns true if the provided error is the one that was created using methods from this package.
func Check(err error) bool {
	return eris.As(err, &errorImpl{})
}

func unwrap(err error) error {
	var errObject errorImpl
	if eris.As(err, &errObject) {
		err = errObject.err
	}
	return err
}

var fieldAverageLength = 8

func buildErrorMessage(baseErr string, fields []any) string {
	if len(fields) < 2 {
		return baseErr
	}
	out := strings.Builder{}
	out.Grow(len(baseErr) + len(fields)*(fieldAverageLength+1))
	out.WriteString(baseErr)

	for i := 0; i < len(fields); i += 2 {
		if len(fields) <= i+1 {
			break
		}
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		out.WriteRune(' ')
		out.WriteString(key)
		out.WriteRune('=')
		out.WriteString(fmt.Sprint(fields[i+1]))
	}

	return out.String()
}

func separateArgsAndFields(msg string, args []any) ([]any, []any) {
	var fields []any
	numberOfFormats := strings.Count(msg, "%")
	if numberOfFormats == 0 {
		return nil, args
	}
	if numberOfFormats <= len(args) {
		fields = args[numberOfFormats:]
		args = args[:numberOfFormats]
	}
	return args, fields
}

// JoinErrors joins error messages using '; ' as separator (instead of '\n' like errors.Join() does).
//
//	a := errm.New("first error")
//	b := errm.New("second error")
//	JoinErrors(a, b)  // "first error; second error"
func JoinErrors(errs ...error) error {
	var b []byte
	for i, err := range errs {
		if err == nil {
			continue
		}
		msg := err.Error()
		if msg == "" {
			continue
		}
		if i > 0 {
			b = append(b, ';', ' ')
		}
		b = append(b, msg...)
	}
	if len(b) > 0 {
		return New(string(b))
	}
	return nil
}
