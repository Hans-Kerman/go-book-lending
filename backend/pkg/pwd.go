package pkg

import "golang.org/x/crypto/bcrypt"

const EncryCost = 12

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		EncryCost,
	)

	return string(hashedBytes), err
}

func CheckPassword(givenPwd, storedPwd string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(storedPwd),
		[]byte(givenPwd),
	)
}
