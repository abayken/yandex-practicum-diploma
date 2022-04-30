package custom_errors

type AlreadyExistsUserError struct {
}

func (error *AlreadyExistsUserError) Error() string {
	return "Such a user exists already"
}

type InvalidCredentialsError struct {
}

func (error *InvalidCredentialsError) Error() string {
	return "Incorrect pair of credentials"
}
