package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sayedulkrm/golang-jwt-project/config"
	"github.com/sayedulkrm/golang-jwt-project/helpers"
	"github.com/sayedulkrm/golang-jwt-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = config.OpenCollection(config.Client, "users")

var validate = validator.New()

func HashPassword(password string) string {
	// Generate the hash from the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		// If an error occurs, log it and handle as needed
		log.Panic(err)
	}

	// Return the hashed password as a string
	return string(hashedPassword)
}

func ComparePassword(userTyingToLoginPassword string, databaseUserPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(databaseUserPassword), []byte(userTyingToLoginPassword))

	fmt.Println("Password eror", err)

	check := true

	msg := ""

	if err != nil {
		check = false
		msg = "Password is incorrect"

	}

	return check, msg

}

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

		password := HashPassword(user.Password)

		user.Password = password

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

		// Set expiration times
		accessTokenExpiry := 5 * time.Minute
		refreshTokenExpiry := 7 * 24 * time.Hour // 7 days

		// Set access token cookie with expiration time (5 minutes)
		c.SetCookie("access_token", token, int(accessTokenExpiry.Seconds()), "/", "localhost", false, true)

		// Set refresh token cookie with expiration time (7 days)
		c.SetCookie("refresh_token", refreshToken, int(refreshTokenExpiry.Seconds()), "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "User created successfully", "user_Id": resultInsertionNumber.InsertedID})

		// c.JSON(http.StatusOK, resultInsertionNumber)

	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var user models.User

		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password is incorrect"})
			return
		}

		isPasswordCorrect, msg := ComparePassword(user.Password, foundUser.Password)

		defer cancel()

		if !isPasswordCorrect {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		if foundUser.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Founduser.email Email or Password is incorrect"})
			return

		}

		token, refreshToken, err := helpers.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.User_id, foundUser.User_type)

		if err != nil {
			log.Panic(err) // Or handle the error in a way suitable for your application
		}

		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found err= usercollect"})
			return

		}

		// Set expiration times
		accessTokenExpiry := 5 * time.Minute
		refreshTokenExpiry := 7 * 24 * time.Hour // 7 days

		// Set access token cookie with expiration time (5 minutes)
		c.SetCookie("access_token", token, int(accessTokenExpiry.Seconds()), "/", "localhost", false, true)

		// Set refresh token cookie with expiration time (7 days)
		c.SetCookie("refresh_token", refreshToken, int(refreshTokenExpiry.Seconds()), "/", "localhost", false, true)

		c.JSON(http.StatusOK, foundUser)

	}

}

// Admin only

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := helpers.CheckUserType(c, "admin"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{"$match", bson.D{}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}}},
		}}
		projectStage := bson.D{{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		}}}

		pipeline := mongo.Pipeline{matchStage, groupStage, projectStage}
		result, err := userCollection.Aggregate(ctx, pipeline)
		if err != nil {
			// Log the error for debugging
			log.Println("Error fetching users:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "222 Failed to fetch users"})
			return
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			// Log the error for debugging
			log.Println("Error decoding users:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "555 Failed to fetch users"})
			return
		}

		if len(allUsers) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
			return
		}

		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUserById() gin.HandlerFunc {

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
