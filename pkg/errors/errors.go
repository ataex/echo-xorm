package errors

const (
	defaultStatus = 503 // http.StatusUnavailable
	zeroStatus    = 0
)

// errorCodeMsg is an error type with integer and text
type errorCodeMsg struct {
	code int
	msg  string
}

//-----------------------------------------------------------------------------
// standard error interface

// New is a constructor, fits standard error interface
// uses 503 as default code
func New(m string) error {
	return &errorCodeMsg{code: defaultStatus, msg: m}
}

// Error() fits standard error interface
func (e *errorCodeMsg) Error() string {
	return e.msg
}

//-----------------------------------------------------------------------------
// custom additions to errors

// NewWithCode is a constructor for non-default custom error code
func NewWithCode(c int, m string) error {
	return &errorCodeMsg{code: c, msg: m}
}

// Code() returns error code
func (e *errorCodeMsg) Code() int {
	return e.code
}

// AddPrefix adds prefix to error message for tracing the context
func (e *errorCodeMsg) AddPrefix(m string) {
	if e == nil {
		return
	}
	if len(m) != 0 {
		e.msg = m + ";" + e.msg
	}
}

// Decompose returns error code and message
func (e *errorCodeMsg) Decompose() (int, string) {
	if e == nil {
		return zeroStatus, "not an error"
	}
	return e.code, e.msg
}
