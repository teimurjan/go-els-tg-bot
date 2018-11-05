package errs

const NoSuchTrackingErrCode = 1

type Err struct {
	code int64
	msg  string
}

func NewErr(code int64, msg string) *Err {
	return &Err{
		code,
		msg,
	}
}

func (e *Err) Error() string {
	return e.msg
}
