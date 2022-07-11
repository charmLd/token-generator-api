package error

var (
	DuplicateEntryStr          = "duplicate entry"
	InvalidUserStr             = "invalid user error"
	FailedUserInsertionStr     = "insertion of user detail is failed"
	FailedUserRoleInsertionStr = "insertion of user role is failed"
	FailedAppRoleInsertionStr  = "insertion of app role is failed"
)

type DuplicateEntryError struct{}

func (e DuplicateEntryError) Error() string {
	return DuplicateEntryStr
}

type InvalidUserError struct{}

func (e InvalidUserError) Error() string {
	return InvalidUserStr
}

type FailedUserInsertionError struct{}

func (e FailedUserInsertionError) Error() string {
	return InvalidSignatureStr
}

type FailedUserRoleInsertionError struct{}

func (e FailedUserRoleInsertionError) Error() string {
	return FailedUserRoleInsertionStr
}

type FailedAppRoleInsertionError struct{}

func (e FailedAppRoleInsertionError) Error() string {
	return FailedAppRoleInsertionStr
}
