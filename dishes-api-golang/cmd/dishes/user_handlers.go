package main

import (
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	// Getting data form request body
	var reqBody struct {
		Username  string
		Firstname string
		Lastname  string
		Email     string
		Password  string
	}

	err := c.Bind(&reqBody)
	if err != nil {
		return
	}

	// Create a user
	user := models.User{
		Username:  reqBody.Username,
		Firstname: reqBody.Firstname,
		Lastname:  reqBody.Lastname,
		Email:     reqBody.Email,
		Password:  reqBody.Password,
	}

	result := initialize.DB.Create(&user)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(201, gin.H{
		"user": user,
	})
}

func GetAllUsers(c *gin.Context) {
	var users []models.User
	initialize.DB.Find(&users)

	c.JSON(200, gin.H{
		"users": users,
	})
}

func GetUserById(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	initialize.DB.Find(&user, id)

	c.JSON(200, gin.H{
		"user": user,
	})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var reqBody struct {
		Username  string
		Firstname string
		Lastname  string
		Email     string
		Password  string
	}

	err := c.Bind(&reqBody)
	if err != nil {
		return
	}
	var user models.User
	initialize.DB.Find(&user, id)

	initialize.DB.Model(&user).Updates(models.User{
		Username:  reqBody.Username,
		Firstname: reqBody.Firstname,
		Lastname:  reqBody.Lastname,
		Email:     reqBody.Email,
		Password:  reqBody.Password})
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")

	initialize.DB.Delete(&models.User{}, id)
	c.Status(200)
}
