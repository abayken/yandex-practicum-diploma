package helpers

import "github.com/theplant/luhn"

type LuhnChecker struct {
}

func (checker LuhnChecker) IsValid(number int) bool {
	return luhn.Valid(number)
}
