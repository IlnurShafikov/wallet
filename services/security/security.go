package security

type Password interface {
	PasswordVerify
	PasswordHasher
}

type PasswordVerify interface {
	Verify(password string, hashPassword []byte) error
}

type PasswordHasher interface {
	HashPassword(password string) ([]byte, error)
}
