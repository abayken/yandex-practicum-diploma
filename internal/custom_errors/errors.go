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

type InvalidOrderNumber struct {
}

func (error *InvalidOrderNumber) Error() string {
	return "Invalid order number"
}

type OrderAlreadyAddedError struct {
	UserID int
}

func (error *OrderAlreadyAddedError) Error() string {
	return "Order added already"
}

type InsufficientFundsError struct {
}

func (error *InsufficientFundsError) Error() string {
	return "Insufficient funds"
}
