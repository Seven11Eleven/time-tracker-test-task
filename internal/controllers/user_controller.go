package controllers

import (
	"net/http"
	"strconv"

	db "github.com/Seven11Eleven/time-tracker-test-task/internal/database"
	"github.com/Seven11Eleven/time-tracker-test-task/internal/logger"
	"github.com/Seven11Eleven/time-tracker-test-task/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	userRepo *db.UserRepository
}

func NewUserController(userRepo *db.UserRepository) *UserController {
	return &UserController{userRepo: userRepo}
}

// GetUsers godoc
// @Summary Получить всех пользователей
// @Description Возвращает список всех пользователей с возможностью фильтрации, пагинации и сортировки
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Количество записей для возврата" default(1)
// @Param offset query int false "Смещение записей для возврата" default(0)
// @Param passport_number query string false "Номер паспорта"
// @Param surname query string false "Фамилия"
// @Param name query string false "Имя"
// @Param patronymic query string false "Отчество"
// @Param address query string false "Адрес"
// @Success 200 {array} models.User
// @Failure 400 {object} gin.H{"error": "Invalid limit value"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /users [get]

func (uc *UserController) GetUsers(c *gin.Context) {
	limit := 1
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"limit": limitStr,
				"error": err,
			}).Error("Неверное значение указанного лимита")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
			return
		}
		limit = parsedLimit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"offset": offsetStr,
				"error":  err,
			}).Error("Указано неверное значение отступа ")
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
		logger.Logger.WithFields(logrus.Fields{
			"filter": filter,
			"limit": limit,
			"offset": offset,
			"error": err,
		}).Error("Не удалось получить информацию о пользователях")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// AddUser godoc
// @Summary Добавить нового пользователя
// @Description Добавляет нового пользователя в систему
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Данные пользователя"
// @Success 201 {object} gin.H{"msg": "Пользователь добавлен!"}
// @Failure 400 {object} gin.H{"error": "Неверный запрос"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /users [post]

func (uc *UserController) AddUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Не удалось забиндить модель с данными")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	if err := uc.userRepo.CreateUser(c, &user); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"user": user,
			"error": err,
		}).Error("Произошла ошибка при попытке создать пользователя")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "упс, не получилось создать пользователя"})
		return
	}

	logger.Logger.WithFields(logrus.Fields{
		"user": user,
	}).Info("Ура, пользователь был создан и добавлен!")

	c.JSON(http.StatusCreated, gin.H{"msg": "Пользователь добавлен!"})
}

// @Summary     Update user
// @Description Update user information
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       user body     models.User true "User to update"
// @Success     200  {object} models.User
// @Failure     400  {object} gin.H
// @Failure     500  {object} gin.H
// @Router      /users [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Не удалось забиндить модель с данными")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	if err := uc.userRepo.UpdateUser(c, &user); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"user":  user,
			"error": err,
		}).Error("Произошла ошибка при попытке обновить информацию об пользователе")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.WithFields(logrus.Fields{
		"user": user,
	}).Info("Информация об пользователе была успешно обновлена")
	c.JSON(http.StatusOK, gin.H{"msg": "Информация пользователя была успешно изменена!"})
}

// @Summary     Delete user
// @Description Delete a user by ID
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       userID path     int true "User ID"
// @Success     200    {object} gin.H
// @Failure     400    {object} gin.H
// @Failure     500    {object} gin.H
// @Router      /users/{userID} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"userID": c.Param("userID"),
			"error": err,
		}).Error("Неверный юзер-айди")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверная айдишка пользователя"})
		return
	}
	if err := uc.userRepo.DeleteUser(c, userID); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"userID": c.Param("userID"),
			"error": err,
		}).Error("Не получилось удалить юзера")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Logger.WithFields(logrus.Fields{
		"userID": c.Param("userID"),
		"error": err,
	}).Info("Юзер был удален")
	c.JSON(http.StatusOK, gin.H{"msg": "Пользователь был успешно удален"})
}
