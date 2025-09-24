package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var words = []string{
	"river","silver","captain","forest","rocket","amber","pioneer","sunset",
	"harbor","island","novel","canvas","granite","thunder","meadow","pluto",
}

// GeneratePass returns human-readable passphrase with two digits at end.
func GeneratePass() string {
	w := func() string {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
		return words[n.Int64()]
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(90))
	return fmt.Sprintf("%s-%s-%s%02d", w(), w(), w(), 10+n.Int64())
}

// HashPassword hashes using bcrypt with default cost.
func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword verifies bcrypt hash.
func CheckPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

// Slugify is used for username base.
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, ".", "")
	return s
}
