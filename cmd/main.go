package main

import (
	"github.com/Seven11Eleven/time-tracker-test-task/internal/controllers"
	db "github.com/Seven11Eleven/time-tracker-test-task/internal/database"

	"github.com/gin-gonic/gin"
)

func main(){
	dbpool := db.ConnectDatabase()
	defer dbpool.Close()



	db.Pool = dbpool
	
	userRepo := db.NewUserRepository(dbpool)
	taskRepo := db.NewTaskRepository(dbpool)

	userController := controllers.NewUserController(userRepo)
	taskController := controllers.NewTaskController(taskRepo)

	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/users", userController.GetUsers)
		api.POST("/users", userController.AddUser)
        api.PUT("/users/:userID", userController.UpdateUser)
        api.DELETE("/users/:userID", userController.DeleteUser)

        api.GET("/users/:userID/tasks", taskController.GetUserTasksByPeriod)
        api.POST("/tasks/start", taskController.StartTask)
        api.POST("/tasks/end/:taskID", taskController.EndTask)
	}

	router.Run(":8080")

}