package errs

// NoSuchTrackingErrCode is used for identifying invalid tracking error
const NoSuchTrackingErrCode = 1

// Err is a custom error object
type Err struct {
	code int64
	msg  string
}

// NewErr creates as new Err instance
func NewErr(code int64, msg string) *Err {
	return &Err{
		code,
		msg,
	}
}

func (e *Err) Error() string {
	return e.msg
}
