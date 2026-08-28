package hash

import "golang.org/x/crypto/bcrypt"

// dummyHash is compared against when a user is not found, so login timing does
// not reveal whether an email exists (mitigates user-enumeration timing attacks).
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("timing-attack-mitigation"), bcrypt.DefaultCost)

func Password(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Compare(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

func DummyCompare(plain string) {
	_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(plain))
}
