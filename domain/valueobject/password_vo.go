package valueobject

import (
    "golang.org/x/crypto/bcrypt"
)

type PasswordVO struct {
    hash []byte
}

// NewPasswordVO takes a plain text password, validates length, and hashes
func NewPasswordVO(plain string) (PasswordVO, error) {
    if len(plain) < 6 {
        return PasswordVO{}, errors.New("password must be at least 6 characters")
    }
    hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
    if err != nil {
        return PasswordVO{}, err
    }
    return PasswordVO{hash: hashed}, nil
}

// initializes from an existing hash (e.g. from the DB)
func NewPasswordVOWithHash(hash string) PasswordVO {
    return PasswordVO{hash: []byte(hash)}
}

// Hash returns the hashed string (to store in DB)
func (p PasswordVO) Hash() string {
    return string(p.hash)
}

// checks a plaintext against the stored hash
func (p PasswordVO) Verify(plain string) bool {
    err := bcrypt.CompareHashAndPassword(p.hash, []byte(plain))
    return err == nil
}
