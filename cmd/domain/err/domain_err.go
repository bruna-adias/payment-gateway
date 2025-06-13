package exceptions

type DomainError struct {
	error
	reason string
}

func NewDomainError(reason string) *DomainError {
	return &DomainError{
		reason: reason,
	}
}

func (e *DomainError) Error() string {
	return e.reason
}
