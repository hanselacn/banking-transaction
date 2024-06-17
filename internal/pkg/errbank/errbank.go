package errbank

type ErrRequestISE string

func NewErrRequestISE(message string) ErrRequestISE {
	return ErrRequestISE(message)
}

func (e ErrRequestISE) Error() string {
	return string(e)
}

type ErrUnprocessableEntity string

func NewErrUnprocessableEntity(message string) ErrUnprocessableEntity {
	return ErrUnprocessableEntity(message)
}

func (e ErrUnprocessableEntity) Error() string {
	return string(e)
}

type ErrNotFound string

func NewErrNotFound(message string) ErrNotFound {
	return ErrNotFound(message)
}

func (e ErrNotFound) Error() string {
	return string(e)
}

type ErrForbidden string

func NewErrForbidden(message string) ErrForbidden {
	return ErrForbidden(message)
}

func (e ErrForbidden) Error() string {
	return string(e)
}

type ErrConflict string

func NewErrConflict(message string) ErrConflict {
	return ErrConflict(message)
}

func (e ErrConflict) Error() string {
	return string(e)
}

type ErrUnauthorized string

func NewErrUnauthorized(message string) ErrUnauthorized {
	return ErrUnauthorized(message)
}

func (e ErrUnauthorized) Error() string {
	return string(e)
}

type ErrTooManyRequest string

func NewErrTooManyRequest(message string) ErrTooManyRequest {
	return ErrTooManyRequest(message)
}

func (e ErrTooManyRequest) Error() string {
	return string(e)
}

type ErrBadRequest string

func NewErrBadRequest(message string) ErrBadRequest {
	return ErrBadRequest(message)
}

func (e ErrBadRequest) Error() string {
	return string(e)
}

type ErrServiceUnavailable string

func NewErrServiceUnavailable(message string) ErrServiceUnavailable {
	return ErrServiceUnavailable(message)
}

func (e ErrServiceUnavailable) Error() string {
	return string(e)
}
