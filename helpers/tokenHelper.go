package helpers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/TheMikeKaisen/go_JWT/database"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection("user")


type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.RegisteredClaims
}

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (string, string, error) {

	// 1. create claim
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_type:  userType,
		Uid:        uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	// 2. create refresh token claim if necessary
	refreshtokenClaim := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
	}

	// 3. create the token
	accessToken, accessError := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("SECRET")))

	// 4. create refresh token
	refreshToken, refreshError := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshtokenClaim).SignedString([]byte(os.Getenv("SECRET")))

	if accessError != nil {
		fmt.Println("Cannot Create access token!")
		return "", "", accessError
	}
	if refreshError != nil {
		fmt.Println("Cannot Create refresh token!")
		return "", "", refreshError
	}

	// 5. return tokens
	return accessToken, refreshToken, nil

}

func UpdateTokens(accessToken string, refreshToken string, userId string) error {

	updateObj := primitive.D{
		bson.E{Key: "token", Value: accessToken},           // Correctly using keyed fields
		bson.E{Key: "refresh_token", Value: refreshToken},
		bson.E{Key: "updated_at", Value: time.Now()},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()

	_, updateErr := userCollection.UpdateOne(ctx, bson.M{"user_id":userId}, bson.D{{Key: "$set", Value: updateObj}})

	if updateErr != nil {
		fmt.Println("Error updating tokens")
		return updateErr
	}

	return nil
	
}

func ValidateToken(signedToken string) (*SignedDetails, string) {
    token, err := jwt.ParseWithClaims(
        signedToken,
        &SignedDetails{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("SECRET")), nil
        },
    )

    if err != nil {
        return nil, err.Error() // Return error immediately
    }

    claims, ok := token.Claims.(*SignedDetails)
    if !ok {
        return nil, "the token is invalid"
    }

	// compare both expiry and current date
    if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, "token is expired"
	}

    return claims, "" // No errors, return claims
}

