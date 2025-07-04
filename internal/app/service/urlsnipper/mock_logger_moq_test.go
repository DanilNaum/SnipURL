// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package urlsnipper

import (
	"sync"
)

// Ensure, that loggerMock does implement logger.
// If this is not the case, regenerate this file with moq.
var _ logger = &loggerMock{}

// loggerMock is a mock implementation of logger.
//
//	func TestSomethingThatUseslogger(t *testing.T) {
//
//		// make and configure a mocked logger
//		mockedlogger := &loggerMock{
//			ErrorfFunc: func(s string, ifaceVals ...interface{})  {
//				panic("mock out the Errorf method")
//			},
//		}
//
//		// use mockedlogger in code that requires logger
//		// and then make assertions.
//
//	}
type loggerMock struct {
	// ErrorfFunc mocks the Errorf method.
	ErrorfFunc func(s string, ifaceVals ...interface{})

	// calls tracks calls to the methods.
	calls struct {
		// Errorf holds details about calls to the Errorf method.
		Errorf []struct {
			// S is the s argument value.
			S string
			// IfaceVals is the ifaceVals argument value.
			IfaceVals []interface{}
		}
	}
	lockErrorf sync.RWMutex
}

// Errorf calls ErrorfFunc.
func (mock *loggerMock) Errorf(s string, ifaceVals ...interface{}) {
	if mock.ErrorfFunc == nil {
		panic("loggerMock.ErrorfFunc: method is nil but logger.Errorf was just called")
	}
	callInfo := struct {
		S         string
		IfaceVals []interface{}
	}{
		S:         s,
		IfaceVals: ifaceVals,
	}
	mock.lockErrorf.Lock()
	mock.calls.Errorf = append(mock.calls.Errorf, callInfo)
	mock.lockErrorf.Unlock()
	mock.ErrorfFunc(s, ifaceVals...)
}

// ErrorfCalls gets all the calls that were made to Errorf.
// Check the length with:
//
//	len(mockedlogger.ErrorfCalls())
func (mock *loggerMock) ErrorfCalls() []struct {
	S         string
	IfaceVals []interface{}
} {
	var calls []struct {
		S         string
		IfaceVals []interface{}
	}
	mock.lockErrorf.RLock()
	calls = mock.calls.Errorf
	mock.lockErrorf.RUnlock()
	return calls
}
