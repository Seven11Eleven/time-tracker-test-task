package controllers

import (
	"net/http"
	"strconv"

	db "github.com/Seven11Eleven/time-tracker-test-task/internal/database"
	"github.com/Seven11Eleven/time-tracker-test-task/internal/models"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userRepo *db.UserRepository
}

func NewUserController(userRepo *db.UserRepository) *UserController {
	return &UserController{userRepo: userRepo}
}

func (uc *UserController) GetUsers(c *gin.Context) {
	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
			return
		}
		limit = parsedLimit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset value"})
			return
		}
		offset = parsedOffset
	}

	filter := make(map[string]interface{})
	allowedFilters := []string{"passport_number", "surname", "name", "patronymic", "address"}

	for _, filterReq := range allowedFilters {
		if filterVal := c.Query(filterReq); filterVal != "" {
			filter[filterReq] = filterVal
		}
	}

	users, err := uc.userRepo.GetUsers(c, filter, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}


func (uc *UserController) AddUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	if err := uc.userRepo.CreateUser(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Пользователь добавлен!"})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	if err := uc.userRepo.UpdateUser(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Информация пользователя была успешно изменена!"})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверная айдишка пользователя"})
		return
	}
	if err := uc.userRepo.DeleteUser(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    c.JSON(http.StatusOK, gin.H{"msg":"Пользователь был успешно удален"})
}
