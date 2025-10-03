package security

import "golang.org/x/crypto/bcrypt"

var pepper = []byte("my-secret")

func Hash(password string) (string, error) {
	passwordWithPepper := append([]byte(password), pepper...)
	hashed, err := bcrypt.GenerateFromPassword(passwordWithPepper, bcrypt.DefaultCost)
	return string(hashed), err
}

func CheckPassword(password, hashedPassword string) bool {
	passwordWithPepper := append([]byte(password), pepper...)
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), passwordWithPepper) == nil
}
