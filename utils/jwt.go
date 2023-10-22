package utils

import (
	"crypto/rsa"
	"go-restapi/app"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (u *utils) GetPrivateKey() (*rsa.PrivateKey, error) {
	pem, err := os.ReadFile("./config/private.pem")
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(pem)
}

func (u *utils) GetAccessToken(key *rsa.PrivateKey, data app.TokenData, expireHour int) (string, int64, error) {
	var (
		t           *jwt.Token
		expDuration time.Duration
		exp         time.Time
		s           string
		err         error
	)

	now := time.Now()
	expDuration = time.Hour * time.Duration(expireHour)
	exp = time.Now().Add(expDuration)
	t = jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"iss":       "go-restapi",
			"sub":       data.UserID,
			"userID":    data.UserID,
			"username":  data.Username,
			"firstname": data.FirstName,
			"lastname":  data.LastName,
			"role":      data.Role,
			"iat":       now.Unix(),
			"exp":       exp.Unix(),
		})
	s, err = t.SignedString(key)
	if err != nil {
		return "", 0, err
	}
	return s, exp.Unix(), nil
}
