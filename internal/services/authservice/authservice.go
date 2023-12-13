package authservice

import (
	"crypto/sha256"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lammer90/gofermart/internal/repository/userstorage"
)

type authenticationServiceImpl struct {
	userRepository userstorage.UserRepository
	privateKey     string
}

func New(userRepository userstorage.UserRepository, privateKey string) AuthenticationService {
	return &authenticationServiceImpl{userRepository: userRepository, privateKey: privateKey}
}

func (a *authenticationServiceImpl) CheckAuthentication(tokenString string) (err error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(a.privateKey), nil
		})
	if err != nil {
		return NotAuthorized
	}
	if !token.Valid || claims.Login == "" {
		return NotAuthorized
	}
	return nil
}

func (a *authenticationServiceImpl) ToRegisterUser(login, password string) (token string, err error) {
	existHash, err := a.userRepository.Find(login)
	if err != nil {
		return "", err
	}
	if existHash != "" {
		return "", UserAlreadyExist
	}

	newHash := buildHash(login, password)
	err = a.userRepository.Save(login, newHash)
	if err != nil {
		return "", err
	}

	token, err = buildJWTString(login, a.privateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *authenticationServiceImpl) ToLoginUser(login, password string) (token string, err error) {
	existHash, err := a.userRepository.Find(login)
	if err != nil {
		return "", err
	}
	if existHash == "" {
		return "", UserDidntFind
	}

	sentHash := buildHash(login, password)
	if existHash != sentHash {
		return "", NotAuthorized
	}

	token, err = buildJWTString(login, a.privateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func buildHash(login string, password string) string {
	src := []byte(login + ":" + password)
	newHashByte := sha256.Sum256(src)
	return string(newHashByte[:])
}

func buildJWTString(login, privateKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		Login:            login,
	})

	tokenString, err := token.SignedString([]byte(privateKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
