package errm

import (
	"sync"
)

// Set object is useful for collecting multiple unique errors into a single error,
// in which error messages are separated by a ";". This object is not safe for concurrent/parallel usage.
// It is not very optimal thing, because it is calling err.Error() to make a key for the map.
// So you have time-overhead caused by Error() and space-overhead because it stores an error twice (string key and value).
// But you can win with it versus [List] when you have a lot of similar errors.
type Set struct {
	errs map[string]error
}

// NewSet returns a new [Set] instance with an empty underlying map.
// Working with [Set] will cause allocations, use [NewSetWithCapacity] if you know the number of unique errors.
func NewSet() *Set {
	return &Set{errs: make(map[string]error)}
}

// NewSetWithCapacity returns a new [Set] instance with an initialized underlying map.
// It may be useful if you know the number of errors and you want to optimize code.
func NewSetWithCapacity(capacity int) *Set {
	return &Set{errs: make(map[string]error, capacity)}
}

// Add sets an error to the underlying map. It is noop if you provide a nil error.
// It will call err.Error() to make a key for the map.
func (e *Set) Add(err error) {
	if err == nil {
		return
	}
	e.errs[err.Error()] = err
}

// New creates an error using [New] and sets in to the underlying map.
// It will call err.Error() to make a key for the map.
func (e *Set) New(msg string, fields ...any) {
	err := New(msg, fields...)
	e.errs[err.Error()] = err
}

// Errorf creates an error using [Errorf] and sets in to the underlying map.
// It will call err.Error() to make a key for the map.
func (e *Set) Errorf(format string, args ...any) {
	err := Errorf(format, args...)
	e.errs[err.Error()] = err
}

// Wrap creates an error using [Wrap] and sets in to the underlying map.
// It will call err.Error() to make a key for the map.
func (e *Set) Wrap(err error, format string, fields ...any) {
	err = Wrap(err, format, fields...)
	e.errs[err.Error()] = err
}

// Wrapf creates an error using [Wrapf] and sets in to the underlying map.
// It will call err.Error() to make a key for the map.
func (e *Set) Wrapf(err error, format string, args ...any) {
	err = Wrapf(err, format, args...)
	e.errs[err.Error()] = err
}

// Has returns true if the [Set] contains the given error.
func (e *Set) Has(err error, errs ...error) bool {
	for _, e := range e.errs {
		if Is(e, err) {
			return true
		}
		for _, err2 := range errs {
			if Is(e, err2) {
				return true
			}
		}
	}
	return false
}

// Err returns current [Set] instance as error interface or nil if it is empty.
func (e *Set) Err() error {
	if len(e.errs) == 0 {
		return nil
	}
	return setError{e}
}

// Empty return true if the [Set] collector is empty.
func (e *Set) Empty() bool {
	return len(e.errs) == 0
}

// Clear removes an underlying map of errors.
func (e *Set) Clear() {
	e.errs = make(map[string]error)
}

// Len returns the number of errors in [Set].
func (e *Set) Len() int {
	return len(e.errs)
}

// SafeSet object is useful for collecting multiple unique errors from different goroutines into a single error,
// in which error messages are separated by a ";". It is safe for concurrent/parallel usage.
type SafeSet struct {
	set *Set
	mu  sync.Mutex
}

// NewSafeSet returns a new [SafeSet] instance with an empty underlying slice.
// Working with [SafeSet] will cause allocations, use [NewSafeSetWithCapacity] if you know the number of unique errors.
func NewSafeSet() *SafeSet {
	return &SafeSet{
		set: NewSet(),
	}
}

// NewSafeSetWithCapacity returns a new [SafeSet] instance with an initialized underlying slice.
// It may be useful if you know the number of errors and you want to optimize code.
func NewSafeSetWithCapacity(capacity int) *SafeSet {
	return &SafeSet{
		set: NewSetWithCapacity(capacity),
	}
}

// Add sets an error to the underlying map. It is safe for concurrent/parallel usage.
func (e *SafeSet) Add(err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.set.Add(err)
}

// New creates an error using [New] and sets it to the underlying map.
// It is safe for concurrent/parallel usage.
func (e *SafeSet) New(err string, fields ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.set.New(err, fields...)
}

// Errorf creates an error using [Errorf] and sets in to the underlying map.
// It is safe for concurrent/parallel usage.
func (e *SafeSet) Errorf(format string, args ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.set.Errorf(format, args...)
}

// Wrap creates an error using [Wrap] and sets in to the underlying map.
// It is safe for concurrent/parallel usage.
func (e *SafeSet) Wrap(err error, format string, fields ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.set.Wrap(err, format, fields...)
}

// Wrapf creates an error using [Wrapf] and sets in to the underlying map.
// It is safe for concurrent/parallel usage.
func (e *SafeSet) Wrapf(err error, format string, args ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.set.Wrapf(err, format, args...)
}

// Has returns true if the [SafeSet] contains the given error. It is safe for concurrent/parallel usage.
func (e *SafeSet) Has(err error) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.set.Has(err)
}

// Empty return true if the [SafeSet] collector is empty. It is safe for concurrent/parallel usage.
func (e *SafeSet) Empty() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.set.Empty()
}

// Err returns current [SafeSet] instance as error interface or nil if it is empty.
// It is safe for concurrent/parallel usage.
func (e *SafeSet) Err() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.set.Err()
}

// Clear removes underlying map of errors. It is safe for concurrent/parallel usage.
func (e *SafeSet) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.set.Clear()
}

// Len returns the number of errors in [SafeSet]. It is safe for concurrent/parallel usage.
func (e *SafeSet) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.set.Len()
}

type setError struct{ *Set }

func (e setError) Error() string {
	if len(e.errs) == 0 {
		return ""
	}
	errs := make([]error, 0, len(e.errs))
	for _, k := range e.errs {
		errs = append(errs, k)
	}
	return JoinErrors(errs...).Error()
}
