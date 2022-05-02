package creds

import "github.com/brianvoe/sjwt"

type Creds struct {
}

type UserModel struct {
	Login string
	Id    int
}

const jwtKey = "diploma"

func (creds Creds) BuildJWT(model UserModel) string {
	claims := sjwt.New()
	claims.Set("login", model.Login)
	claims.Set("id", model.Id)
	jwt := claims.Generate([]byte(jwtKey))

	return jwt
}

func (creds Creds) Id(token string) (int, error) {
	claims, err := sjwt.Parse(token)

	if err != nil {
		return 0, err
	}

	id, err := claims.GetInt("id")

	return id, err
}
