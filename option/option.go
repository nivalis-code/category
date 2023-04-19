// Package option provides functions for dealing with values that may or may not exist.

package option

import (
	"errors"
)

// An Option is either some value, or it's nothing.
type Option[T any] interface {
	// Unwrap returns the contained value of a Some-valued Option.
	Unwrap() (T, error)
	// IsNothing returns true if an Option is Nothing-valued
	IsNothing() bool
	// IsSome returns true if the Option is Some-valued
	IsSome() bool
}

type option[T any] struct {
	maybeValue T
	isSome     bool
}

func (opt option[T]) Unwrap() (T, error) {
	if !opt.isSome {
		var zero T
		return zero, errors.New("unwrapped nothing-value")
	}
	return opt.maybeValue, nil
}

func (opt option[T]) IsNothing() bool {
	return !opt.isSome
}

func (opt option[T]) IsSome() bool {
	return opt.isSome
}

// Some creates a some-valued Option.
func Some[T any](x T) Option[T] {
	return option[T]{x, true}
}

// Nothing creates a nothing-valued Option.
func Nothing[T any]() Option[T] {
	var zero T
	return option[T]{zero, false}
}

// Lift lifts a function to the Option category. The lifted function returns Nothing for a Nothing argument, and a
// Some-valued Option containing the result of applying the base function to the inner value of a Some-valued argument.
func Lift[T any, R any](f func(T) R) func(Option[T]) Option[R] {
	return func(x Option[T]) Option[R] {
		xValue, err := x.Unwrap()
		if err != nil {
			return Nothing[R]()
		}

		return Some(f(xValue))
	}
}

// Lift2 lifts a binary function to the Option category. The lifted function returns Nothing if either argument is
// Nothing-valued, and a Some-valued Option containing the result of applying the base function to the inner values of
// two Some-valued arguments.
func Lift2[A any, B any, C any](f func(A, B) C) func(Option[A], Option[B]) Option[C] {
	return func(x Option[A], y Option[B]) Option[C] {
		xValue, err := x.Unwrap()
		if err != nil {
			return Nothing[C]()
		}

		yValue, err := y.Unwrap()
		if err != nil {
			return Nothing[C]()
		}

		return Some(f(xValue, yValue))
	}
}

// Bind creates a function that lifts a function returning an Option and joins the result. The resulting function
// returns nothing for a Nothing argument, ir the result of applying the base function to the inner value of a
// Some-valued argument.
func Bind[T any, R any](f func(T) Option[R]) func(Option[T]) Option[R] {
	return func(x Option[T]) Option[R] {
		xValue, err := x.Unwrap()
		if err != nil {
			return Nothing[R]()
		}

		return f(xValue)
	}
}

// Wrap wraps a Go-idiomatic value-and-error pair in an Option. Any error information is discarded.
func Wrap[T any](value T, err error) Option[T] {
	if err != nil {
		return Nothing[T]()
	}

	return Some(value)
}
