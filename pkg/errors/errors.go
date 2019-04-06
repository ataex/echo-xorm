package errors

const (
	defaultStatus = 503 // http.StatusUnavailable
	zeroStatus    = 0
)

// errorCodeMessage is an error type with integer and text
type errorCodeMessage struct {
	code int
	msg  string
}

//-----------------------------------------------------------------------------
// standard error interface

// New is a constructor, fits standard error interface
// uses 503 as default code
func New(m string) error {
	return &errorCodeMessage{code: defaultStatus, msg: m}
}

// Error() fits standard error interface
func (e *errorCodeMessage) Error() string {
	return e.msg
}

//-----------------------------------------------------------------------------
// custom additions to errors

// Code() returns error code
func (e *errorCodeMessage) Code() int {
	return e.code
}

// NewWithCode is a constructor for non-default custom error code
func NewWithCode(c int, m string) error {
	return &errorCodeMessage{code: c, msg: m}
}

// NewWithPrefix is a constructor-wrapper from error (first parameter)
// second parameter is a prefix message
func NewWithPrefix(e error, m string) error {
	if e == nil {
		return nil
	}
	code := defaultStatus
	ecm, ok := e.(*errorCodeMessage)
	if ok {
		code = ecm.code
	}
	return &errorCodeMessage{code: code, msg: m + "; " + e.Error()}
}

// Decompose returns error code and message
// if error contains 'code' then function extracts it
// else function returns defaultStatus as a code
// In case of nil error the function will return nils
func Decompose(e error) (int, string) {
	if e == nil {
		return zeroStatus, "not an error"
	}
	ecm, ok := e.(*errorCodeMessage)
	if ok {
		return ecm.code, ecm.msg
	}
	// else
	return defaultStatus, e.Error()
}
