// Package result provides functions for dealing with computations that may fail

package result

// A Result is either an Ok value, or it's an error.
type Result[T any] interface {
	// Unwrap returns the contained value of a Ok-valued Result.
	Unwrap() (T, error)
	// IsOk returns true if a result is Ok-valued.
	IsOk() bool
	// IsError returns true if a result is Error-valued.
	IsError() bool
}

type result[T any] struct {
	maybeValue T
	maybeErr   error
}

func (res result[T]) Unwrap() (T, error) {
	if res.maybeErr != nil {
		var zero T
		return zero, res.maybeErr
	}
	return res.maybeValue, nil
}

func (res result[T]) IsOk() bool {
	return res.maybeErr == nil
}

func (res result[T]) IsError() bool {
	return res.maybeErr != nil
}

// Ok creates an Ok-valued result.
func Ok[T any](x T) Result[T] {
	return result[T]{x, nil}
}

// Error creates an Error-valued result.
func Error[T any](err error) Result[T] {
	var zero T
	return result[T]{zero, err}
}

// Lift lifts a function to the Result applicative functor. The lifted function returns an Error for an Error argument, and an
// Ok-valued Result containing the result of applying the base function to the inner value of an Ok-valued argument.
func Lift[T any, R any](f func(T) R) func(Result[T]) Result[R] {
	return func(x Result[T]) Result[R] {
		xValue, err := x.Unwrap()
		if err != nil {
			return Error[R](err)
		}

		return Ok(f(xValue))
	}
}

// Lift2 lifts a binary function to the Result applicative functor. The lifted function returns an Error if either argument is
// Error-valued, and an Ok-valued Result containing the result of applying the base function to the inner values of
// two Ok-valued arguments.
func Lift2[A any, B any, C any](f func(A, B) C) func(Result[A], Result[B]) Result[C] {
	return func(x Result[A], y Result[B]) Result[C] {
		xValue, err := x.Unwrap()
		if err != nil {
			return Error[C](err)
		}

		yValue, err := y.Unwrap()
		if err != nil {
			return Error[C](err)
		}

		return Ok(f(xValue, yValue))
	}
}

// Bind creates a function that lifts a function returning a Result and joins the result. The resulting function
// returns nothing for an Error-valued argument, or the result of applying the base function to the inner value of an
// Ok-valued input.
func Bind[T any, R any](f func(T) Result[R]) func(Result[T]) Result[R] {
	return func(x Result[T]) Result[R] {
		xValue, err := x.Unwrap()
		if err != nil {
			return Error[R](err)
		}

		return f(xValue)
	}
}

// Wrap wraps a Go-idiomatic value-and-error pair in a Result. Error information is preserved.
func Wrap[T any](value T, err error) Result[T] {
	if err != nil {
		return Error[T](err)
	}

	return Ok(value)
}

// Kliesli is Kliesli composition of two functions for the Result monad
func Kliesli[A any, B any, C any](f func(A) Result[B], g func(B) Result[C]) func(A) Result[C] {
	return func(x A) Result[C] {
		return Bind(g)(f(x))
	}
}
