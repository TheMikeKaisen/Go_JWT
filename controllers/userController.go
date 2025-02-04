package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/TheMikeKaisen/go_JWT/database"
	"github.com/TheMikeKaisen/go_JWT/helpers"
	"github.com/TheMikeKaisen/go_JWT/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection("user")

var validate = validator.New()

// to hash passwords
func hashPassword(password string) (string, error) {
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		fmt.Println("Error hashing password")
		return "", hashErr
	}
	return string(hashedPassword), nil
}

func verifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. bind the incoming json body
		var user models.User
		bindError := c.ShouldBindJSON(&user)
		if bindError != nil {
			fmt.Println("Error while signing up!")
			c.JSON(500, gin.H{"error": bindError.Error()})
			return
		}

		// 2. hash password
		hashedPassword, hashErr := hashPassword(*user.Password)
		if hashErr != nil {
			c.JSON(500, gin.H{"error": hashErr.Error()})
		}
		// make user password as hashed password
		user.Password = &hashedPassword

		user.ID = primitive.NewObjectID()
		user.Created_At = time.Now()
		user.Updated_at = time.Now()
		user.User_id = user.ID.Hex()

		// validate the body
		validateErr := validate.Struct(user)
		if validateErr != nil {
			fmt.Println("Validate error")
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}

		// 3. create context
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 4. check if email or phone number already exists.
		countEmail, countErr := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if countErr != nil {
			fmt.Println("Error while counting documents!")
			c.JSON(500, gin.H{"error": countErr.Error()})
			return
		}
		// check if email or phone number already exists.
		countPhone, countErr := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if countErr != nil {
			fmt.Println("Error while counting documents!")
			c.JSON(500, gin.H{"error": countErr.Error()})
			return
		}

		if countEmail > 0 || countPhone > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or phone number already exists"})
			return
		}

		accessToken, refreshToken, tokenErr := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)

		if tokenErr != nil {
			return
		}

		user.Refresh_token = &refreshToken
		user.Token = &accessToken

		// 5. create user
		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			fmt.Println("Insert error")
			c.JSON(500, gin.H{"error": insertErr.Error()})
			return
		}

		// 6. return user
		c.JSON(http.StatusOK, user)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var user models.User

		bindErr := c.ShouldBindJSON(&user)
		if bindErr != nil {
			fmt.Println("Error binding body into user.")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
		defer cancel()

		var foundUser models.User

		userDecodeErr := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if userDecodeErr != nil {
			fmt.Println("User Decode Error", userDecodeErr)
			return
		}

		// check if password is correct
		verifyErr := verifyPassword(*foundUser.Password, *user.Password)
		if verifyErr != nil {
			fmt.Println("Email or password is wrong!")
			return
		}

		accessToken, refreshToken, tokenErr := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)

		if tokenErr != nil {
			return
		}

		updateErr := helpers.UpdateTokens(accessToken, refreshToken, foundUser.User_id)
		if updateErr != nil {
			return
		}
		findErr := userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if findErr != nil {
			fmt.Println("Find Error")
			return
		}

		c.JSON(200, foundUser)

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		// Role based access
		authorizeErr := helpers.AuthorizeRole(c, userId)
		if authorizeErr != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": authorizeErr.Error()})
			return
		}

		// create context
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User

		decodeErr := userCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
		if decodeErr != nil {
			fmt.Println("Error while finding user")
			c.JSON(500, gin.H{"error": decodeErr.Error()})
			return
		}
		c.JSON(http.StatusOK, user)

	}
}
