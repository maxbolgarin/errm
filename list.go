package errm

import (
	"sync"
)

// List object is useful for collecting multiple errors into a single error,
// in which error messages are separated by a ";". This object is not safe for concurrent/parallel usage.
type List struct {
	errs []error
}

// NewList returns a new [List] instance with an empty underlying slice.
// Working with [List] will cause allocations, use [NewListWithCapacity] if you know the number of errors.
func NewList() *List {
	return &List{}
}

// NewListWithCapacity returns a new [List] instance with an initialized underlying slice.
// It may be useful if you know the number of errors and you want to optimize code.
func NewListWithCapacity(capacity int) *List {
	return &List{errs: make([]error, 0, capacity)}
}

// Add appends an error to the underlying slice. It is noop if you provide an empty error.
func (e *List) Add(err error) {
	if err == nil {
		return
	}
	e.errs = append(e.errs, err)
}

// New creates an error using [New] and appends in to the underlying slice.
func (e *List) New(err string, fields ...any) {
	e.errs = append(e.errs, New(err, fields...))
}

// Errorf creates an error using [Errorf] and appends in to the underlying slice.
func (e *List) Errorf(format string, args ...any) {
	e.errs = append(e.errs, Errorf(format, args...))
}

// Wrap creates an error using [Wrap] and appends in to the underlying slice.
func (e *List) Wrap(err error, format string, fields ...any) {
	e.errs = append(e.errs, Wrap(err, format, fields...))
}

// Wrapf creates an error using [Wrapf] and appends in to the underlying slice.
func (e *List) Wrapf(err error, format string, args ...any) {
	e.errs = append(e.errs, Wrapf(err, format, args...))
}

// Has returns true if the [List] contains the given error.
func (e *List) Has(err error, errs ...error) bool {
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

// Err returns current [List] instance as error interface or nil if it is empty.
func (e *List) Err() error {
	if len(e.errs) == 0 {
		return nil
	}
	return listError{e}
}

// Empty returns true if the [List] collector is empty.
func (e *List) Empty() bool {
	return len(e.errs) == 0
}

// NotEmpty returns true if the [List] collector has errors.
func (e *List) NotEmpty() bool {
	return len(e.errs) != 0
}

// Clear removes an underlying slice of errors.
func (e *List) Clear() {
	e.errs = nil
}

// Len returns the number of errors in [List].
func (e *List) Len() int {
	return len(e.errs)
}

// SafeList object is useful for collecting multiple errors from different goroutines into a single error,
// in which error messages are separated by a ";". It is safe for concurrent/parallel usage.
type SafeList struct {
	List *List
	mu   sync.Mutex
}

// NewSafeList returns a new [SafeList] instance with an empty underlying slice.
// Working with [SafeList] will cause allocations, use [NewSafeListWithCapacity] if you know the number of errors.
func NewSafeList() *SafeList {
	return &SafeList{
		List: NewList(),
	}
}

// NewSafeListWithCapacity returns a new [SafeList] instance with an initialized underlying slice.
// It may be useful if you know the number of errors and you want to optimize code.
func NewSafeListWithCapacity(capacity int) *SafeList {
	return &SafeList{
		List: NewListWithCapacity(capacity),
	}
}

// Add appends an error to the underlying slice. It is noop if you provide an empty error.
// It is safe for concurrent/parallel usage.
func (e *SafeList) Add(err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.List.Add(err)
}

// New creates an error using [New] and appends in to the underlying slice.
// It is safe for concurrent/parallel usage.
func (e *SafeList) New(err string, fields ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.List.New(err, fields...)
}

// Errorf creates an error using [Errorf] and appends in to the underlying slice.
// It is safe for concurrent/parallel usage.
func (e *SafeList) Errorf(format string, args ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.List.Errorf(format, args...)
}

// Wrap creates an error using [Wrap] and appends in to the underlying slice.
// It is safe for concurrent/parallel usage.
func (e *SafeList) Wrap(err error, format string, fields ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.List.Wrap(err, format, fields...)
}

// Wrapf creates an error using [Wrapf] and appends in to the underlying slice.
// It is safe for concurrent/parallel usage.
func (e *SafeList) Wrapf(err error, format string, args ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.List.Wrapf(err, format, args...)
}

// Has returns true if the [SafeList] contains the given error. It is safe for concurrent/parallel usage.
func (e *SafeList) Has(err error) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.List.Has(err)
}

// Empty return true if the [SafeList] collector is empty. It is safe for concurrent/parallel usage.
func (e *SafeList) Empty() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.List.Empty()
}

// NotEmpty return true if the [SafeList] collector has errors. It is safe for concurrent/parallel usage.
func (e *SafeList) NotEmpty() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.List.NotEmpty()
}

// Err returns current [SafeList] instance as error interface or nil if it is empty.
// It is safe for concurrent/parallel usage.
func (e *SafeList) Err() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.List.Err()
}

// Clear removes underlying slice of errors. It is safe for concurrent/parallel usage.
func (e *SafeList) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.List.Clear()
}

// Len returns the number of errors in [SafeList]. It is safe for concurrent/parallel usage.
func (e *SafeList) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.List.Len()
}

type listError struct{ *List }

func (e listError) Error() string {
	if len(e.errs) == 0 {
		return ""
	}
	return JoinErrors(e.errs...).Error()
}
