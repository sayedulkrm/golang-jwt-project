package helpers

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sayedulkrm/golang-jwt-project/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func ValidateToken(accessToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JWT_SECRET_KEY), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, errors.New("the token is invalid")
	}

	// Convert claims.ExpiresAt to Unix timestamp (int64)
	expiresAtUnix := claims.ExpiresAt.Time.Unix()

	// Get the current time in Unix timestamp (seconds)
	currentTimeUnix := time.Now().Unix()

	// Compare the expiry time of the token with the current time
	if expiresAtUnix < currentTimeUnix {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}

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

func UpdateAllTokens(accessToken string, refreshToken string, userId string) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: accessToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: refreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true

	filter := bson.M{"user_id": userId}

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
	}

}
