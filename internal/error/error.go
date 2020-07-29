package error

type OutError string

func (e OutError) Error() string {
	return string(e)
}

var (
	ErrINetAlreadySet = OutError("this network is already set")
	ErrInternal       = OutError("internal server error")
	ErrInvalidParams  = OutError("invalid params")
	ErrEmptyLogin     = OutError("empty login param")
	ErrEmptyPassword  = OutError("empty password param")
	ErrEmptyIP        = OutError("empty IP param")
)
