package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sayedulkrm/golang-jwt-project/config"
	"github.com/sayedulkrm/golang-jwt-project/helpers"
	"github.com/sayedulkrm/golang-jwt-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = config.OpenCollection(config.Client, "users")

var validate = validator.New()

func HashPassword(password string) {}

func ComparePassword(password string, hash string) {}

func Register() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}

		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count documents"})
			return
		}

		phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count documents phone"})
			return
		}

		if emailCount > 0 || phoneCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}

		// start from here

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, err := helpers.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.User_type, user.User_id)

		if err != nil {
			log.Panic(err) // Or handle the error in a way suitable for your application
		}

		user.Token = token
		user.Refresh_token = refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)

		if insertErr != nil {
			msg := fmt.Sprintf("Error: User not created %v", insertErr)

			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "User created successfully", "user_Id": resultInsertionNumber.InsertedID})

		// c.JSON(http.StatusOK, resultInsertionNumber)

	}

}

func Login() {}

func GetAllUsers() {

}

func GetUserById() {

}

func AdminGetAllUsers() gin.HandlerFunc {

	return func(c *gin.Context) {

		userId := c.Param("user_id")

		fmt.Println("Admin get all User ID", userId)

		if err := helpers.MatchUserType(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var user models.User

		//  idPrimitive, err := primitive.ObjectIDFromHex(userId)

		//  if err != nil {
		// 	 // Handle error (e.g., invalid ID format)
		// 	 c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		// 	 return
		//  }

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)

	}

}
