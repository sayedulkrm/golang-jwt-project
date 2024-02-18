package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sayedulkrm/golang-jwt-project/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	FirstName string
	LastName  string
	Email     string
	UserId    string
	User_type string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection = config.OpenCollection(config.Client, "user")

var JWT_SECRET_KEY string = os.Getenv("JWT_SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, userId string, userType string) (singedToken string, signedRefreshToken string, err error) {
	clamis := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserId:    userId,
		User_type: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 1)),
		},
	}

	refreshClamis := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 168)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, clamis).SignedString([]byte(JWT_SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClamis).SignedString([]byte(JWT_SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}
