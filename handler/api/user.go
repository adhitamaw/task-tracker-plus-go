package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Endpoint untuk melihat seluruh user (debug)
func (u *userAPI) ListUsers(c *gin.Context) {
	users, err := u.userService.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Kosongkan password sebelum dikirim ke client
	for i := range users {
		users[i].Password = ""
	}
	c.JSON(200, users)
}

type UserAPI interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserTaskCategory(c *gin.Context)
	ListUsers(c *gin.Context) // debug: lihat semua user
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Register(c *gin.Context) {
	var user model.UserRegister

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid decode json"})
		return
	}

	if user.Email == "" || user.Password == "" || user.Fullname == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "register data is empty"})
		return
	}

	var recordUser = model.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}

	_, err := u.userService.Register(&recordUser)
	if err != nil {
		if err.Error() == "email already exists" {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "error internal server: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{Message: "register success"})
}

func (u *userAPI) Login(c *gin.Context) {
	var user model.UserLogin

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid decode json"})
		return
	}

	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "email or password is empty"})
		return
	}

	var recordUser = model.User{
		Email:    user.Email,
		Password: user.Password,
	}

	token, err := u.userService.Login(&recordUser)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "user not found"})
			return
		}
		if err.Error() == "wrong email or password" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "wrong email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "error internal server: " + err.Error()})
		return
	}

	c.SetCookie("session_token", *token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
	})
}

func (u *userAPI) GetUserTaskCategory(c *gin.Context) {
	userTaskCategories, err := u.userService.GetUserTaskCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "error internal server"})
		return
	}

	c.JSON(http.StatusOK, userTaskCategories)
}
