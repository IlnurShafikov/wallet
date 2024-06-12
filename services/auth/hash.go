package auth

import "golang.org/x/crypto/bcrypt"

type BcryptHashing struct {
	secret string
}

func NewBcryptHashing(secret string) *BcryptHashing {
	return &BcryptHashing{secret: secret}

}

func (p *BcryptHashing) HashPassword(password string) ([]byte, error) {
	passwordWithSecret := password + p.secret
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithSecret), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashPassword, nil
}

func (p *BcryptHashing) Verify(password string, hashPassword []byte) error {
	passwordWithSecret := password + p.secret
	err := bcrypt.CompareHashAndPassword(hashPassword, []byte(passwordWithSecret))
	if err != nil {
		return ErrWrongRePassword
	}

	return nil
}
