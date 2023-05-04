package result

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

	testOkString := Ok("test")
	testOkInt := Ok(5)
	testOkFloat := Ok(65.2)
	testOkStruct := Ok(testStruct)

	assert.True(t, testOkString.IsOk())
	assert.True(t, testOkInt.IsOk())
	assert.True(t, testOkFloat.IsOk())
	assert.True(t, testOkStruct.IsOk())
}

func TestUnwrap(t *testing.T) {
	testOk := Ok("test")
	testError := Error[string](errors.New("test error"))

	unwrappedOk, unwrappedOkErr := testOk.Unwrap()
	unwrappedError, unwrappedErrorErr := testError.Unwrap()

	assert.Equal(t, unwrappedOk, "test")
	assert.Nil(t, unwrappedOkErr)
	assert.Equal(t, unwrappedError, "")
	assert.NotNil(t, unwrappedErrorErr)
}

func TestIsError(t *testing.T) {
	testOk := Ok("test")
	testError := Error[string](errors.New("test error"))

	assert.False(t, testOk.IsError())
	assert.True(t, testError.IsError())
}

func TestIsOk(t *testing.T) {
	testOk := Ok("test")
	testError := Error[string](errors.New("test error"))

	assert.True(t, testOk.IsOk())
	assert.False(t, testError.IsOk())
}

func TestLift(t *testing.T) {
	testFunction := func(x string) int {
		return len(x)
	}

	liftedFunction := Lift(testFunction)

	testOk := Ok("test")
	testError := Error[string](errors.New("test error"))

	assert.Equal(t, Ok(4), liftedFunction(testOk))
	assert.Equal(t, Error[int](errors.New("test error")), liftedFunction(testError))
}

func TestLift2(t *testing.T) {
	testFunction := func(x string, y int) float64 {
		return float64(len(x)-y) * 1.5
	}

	liftedFunction := Lift2(testFunction)

	testOkString := Ok("test")
	testErrorString := Error[string](errors.New("test error 1"))
	testOkInt := Ok(2)
	testErrorInt := Error[int](errors.New("test error 2"))

	assert.Equal(t, Ok(3.0), liftedFunction(testOkString, testOkInt))
	assert.Equal(t, Error[float64](errors.New("test error 1")), liftedFunction(testErrorString, testOkInt))
	assert.Equal(t, Error[float64](errors.New("test error 2")), liftedFunction(testOkString, testErrorInt))
	assert.Equal(t, Error[float64](errors.New("test error 1")), liftedFunction(testErrorString, testErrorInt))
}

func TestBind(t *testing.T) {
	testFunction := func(x string) Result[float64] {
		n := len(x)
		if n == 0 {
			return Error[float64](errors.New("divide by zero"))
		}
		return Ok(1 / float64(n))
	}

	testSome := Ok("test")
	testSomeEmpty := Ok("")
	testNothing := Error[string](errors.New("test error"))

	assert.Equal(t, Ok(1/4.0), Bind(testFunction)(testSome))
	assert.Equal(t, Error[float64](errors.New("divide by zero")), Bind(testFunction)(testSomeEmpty))
	assert.Equal(t, Error[float64](errors.New("test error")), Bind(testFunction)(testNothing))
}

func TestWrap(t *testing.T) {
	testFunction := func(x int) (float64, error) {
		if x == 0 {
			return 0, errors.New("divide by zero")
		}

		return 1 / float64(x), nil
	}

	assert.Equal(t, Ok(1/4.0), Wrap(testFunction(4)))
	assert.Equal(t, Error[float64](errors.New("divide by zero")), Wrap(testFunction(0)))
}

func TestKliesli(t *testing.T) {
	testFunction1 := func(s string) Result[int] {
		if s == "fnord" {
			return Error[int](errors.New("forbidden string"))
		}
		
		return Ok(len(s))
	}
	testFunction2 := func(x int) Result[float64] {
		if x == 0 {
			return Error[float64](errors.New("divide by zero"))
		}

		return Ok(1/float64(x))
	}

	composition := Kliesli(testFunction1, testFunction2)

	assert.Equal(t, Ok(1/4.0), composition("test"))
	assert.True(t, composition("fnord").IsError())
	assert.True(t, composition("").IsError())
}
