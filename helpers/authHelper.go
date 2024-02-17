package helpers

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// These 2 function if Admin want to see the data then admin can see any user data. IF owner (user) want to see the data then he can see only his data.

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")

	err = nil

	if userType != role {
		err = errors.New("Unauthorized to access this resource")
		return err

	}

	return err
}

func MatchUserType(c *gin.Context, userId string) (err error) {

	fmt.Println("Heyyyyy am from auth helper C", c)

	userType := c.GetString("user_type")
	fmt.Println("Heyyyyy am from auth helper userType", userType)

	uid := c.GetString("uid")

	err = nil

	if userType == "user" && uid != userId {
		err = errors.New("Unauthorized To Access This Resource")
		return err
	}

	err = CheckUserType(c, userType)

	return err

}
