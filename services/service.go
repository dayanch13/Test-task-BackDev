package services

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
	"time"
)

const secret = "go-jwt-auth"

type UserClaims struct {
	ID        int
	UserEmail string
	UserIP    string
	jwt.RegisteredClaims
}

func NewUserClaims(id int, email string, ip string, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error genarating tokenID %w", err)
	}

	return &UserClaims{
		ID:        id,
		UserEmail: email,
		UserIP:    ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}

func CreateToken(id int, email string, ip string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserClaims(id, email, ip, duration)
	if err != nil {
		return "", nil, fmt.Errorf("error get claims &w", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil, err
	}

	return tokenStr, claims, nil
}

func VerifyToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing &w", err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "dayko130201@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetHeader("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, "dayko130201@gmail.com", "1234567890.")
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
