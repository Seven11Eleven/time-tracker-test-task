package controllers

import (
	"net/http"
	"strconv"
	"time"

	db "github.com/Seven11Eleven/time-tracker-test-task/internal/database"
	"github.com/Seven11Eleven/time-tracker-test-task/internal/models"
	"github.com/gin-gonic/gin"
)

type TaskController struct{
	taskRepo *db.TaskRepository
}

func NewTaskController(taskRepo *db.TaskRepository) *TaskController{
	return &TaskController{taskRepo: taskRepo}
}

func (tc *TaskController) GetUserTasksByPeriod(c *gin.Context){
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"айди юзера неверное"})
		return
	}

	start, err := time.Parse(time.RFC3339, c.Query("start"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time"})
        return
    }

    end, err := time.Parse(time.RFC3339, c.Query("end"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time"})
        return
    }

	tasks, err := tc.taskRepo.GetUserTasksByPeriod(c, userID, start, end)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) StartTask(c *gin.Context){
	var req models.Request

	if err := c.BindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	if err := tc.taskRepo.StartTask(c, int(req.UserID), req.Description); err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Задача взята, тайм-треккинг начат"})
}

func(tc *TaskController) EndTask(c *gin.Context){
	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверная айдишка начатой задачи"})
		return
	}

	if err := tc.taskRepo.EndTask(c, taskID); err !=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Выполнение задачи окончено, тайм-треккинг остановлен"})
}