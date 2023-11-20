package xerrors

// X is a pair of errors, one for a user
// and one for a developer.
type X struct {
	// User is the error that should be shown to a user.
	User error
	// System is the error that should be logged.
	System error
}

// Error implements the error interface.
func (x X) Error() string {
	return x.User.Error()
}
