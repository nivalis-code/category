package option

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSome(t *testing.T) {
	type foo struct {
		bar int
		baz string
	}
	testStruct := foo{6, "test"}

	testSomeString := Some("test")
	testSomeInt := Some(5)
	testSomeFloat := Some(65.2)
	testSomeStruct := Some(testStruct)

	assert.True(t, testSomeString.IsSome())
	assert.True(t, testSomeInt.IsSome())
	assert.True(t, testSomeFloat.IsSome())
	assert.True(t, testSomeStruct.IsSome())
}

func TestUnwrap(t *testing.T) {
	testSome := Some("test")
	testNothing := Nothing[string]()

	unwrappedSome, unwrappedSomeErr := testSome.Unwrap()
	unwrappedNothing, unwrappedNothingErr := testNothing.Unwrap()

	assert.Equal(t, unwrappedSome, "test")
	assert.Nil(t, unwrappedSomeErr)
	assert.Equal(t, unwrappedNothing, "")
	assert.NotNil(t, unwrappedNothingErr)
}

func TestIsNothing(t *testing.T) {
	testSome := Some("test")
	testNothing := Nothing[string]()

	assert.False(t, testSome.IsNothing())
	assert.True(t, testNothing.IsNothing())
}

func TestIsSome(t *testing.T) {
	testSome := Some("test")
	testNothing := Nothing[string]()

	assert.True(t, testSome.IsSome())
	assert.False(t, testNothing.IsSome())
}

func TestLift(t *testing.T) {
	testFunction := func(x string) int {
		return len(x)
	}

	liftedFunction := Lift(testFunction)

	testSome := Some("test")
	testNothing := Nothing[string]()

	assert.Equal(t, Some(4), liftedFunction(testSome))
	assert.Equal(t, Nothing[int](), liftedFunction(testNothing))
}

func TestLift2(t *testing.T) {
	testFunction := func(x string, y int) float64 {
		return float64(len(x)-y) * 1.5
	}

	liftedFunction := Lift2(testFunction)

	testSomeString := Some("test")
	testNothingString := Nothing[string]()
	testSomeInt := Some(2)
	testNothingInt := Nothing[int]()

	assert.Equal(t, Some(3.0), liftedFunction(testSomeString, testSomeInt))
	assert.Equal(t, Nothing[float64](), liftedFunction(testNothingString, testSomeInt))
	assert.Equal(t, Nothing[float64](), liftedFunction(testSomeString, testNothingInt))
	assert.Equal(t, Nothing[float64](), liftedFunction(testNothingString, testNothingInt))
}

func TestBind(t *testing.T) {
	testFunction := func(x string) Option[float64] {
		n := len(x)
		if n == 0 {
			return Nothing[float64]()
		}
		return Some(1 / float64(n))
	}

	testSome := Some("test")
	testSomeEmpty := Some("")
	testNothing := Nothing[string]()

	assert.Equal(t, Some(1/4.0), Bind(testFunction)(testSome))
	assert.Equal(t, Nothing[float64](), Bind(testFunction)(testSomeEmpty))
	assert.Equal(t, Nothing[float64](), Bind(testFunction)(testNothing))
}

func TestWrap(t *testing.T) {
	testFunction := func(x int) (float64, error) {
		if x == 0 {
			return 0, errors.New("divide by zero")
		}

		return 1 / float64(x), nil
	}

	assert.Equal(t, Some(1/4.0), Wrap(testFunction(4)))
	assert.Equal(t, Nothing[float64](), Wrap(testFunction(0)))
}
