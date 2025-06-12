package main

import (
	"net/http"					// allows HTTP requests 
	"github.com/gin-gonic/gin"	// GIN API FRAMEWORK 
	"log"						// Logging steps 
	"sync"						// Used for mutex - stops multiple threads from accessing the shared data at the same time (like a semaphore)
)


// structs are used for grouping data 
// as no methods i dont wanna use a obj
type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}


var (
	tasks []Task 			  // holds all tasks in memory 
	idCounter = 1			  // keep track of the task ID's
	mutex     = &sync.Mutex{} // for safe concurrent access
)

func main() {
	router := gin.Default()	  // set up gin

	// GET /tasks
	router.GET("/tasks", func(c *gin.Context) {
		log.Println("GET /tasks called")
		c.JSON(http.StatusOK, tasks) // sends back 200 ok and tasks to client
	})

	// POST /tasks
	router.POST("/tasks", func(c *gin.Context) { // pointer to the gin context (context is basically eveything about the http request)
		var newTask Task 						 // placeholder for new incoming request

		if err := c.ShouldBindJSON(&newTask); err != nil {  // try to map the incoming JSON to Task struct
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// locks this piece of code so theres a consistent ID is kept 
		mutex.Lock()
		newTask.ID = idCounter // update new inputted task 
		idCounter++ // update struct 
		tasks = append(tasks, newTask) //update global task array
		mutex.Unlock() // unlock when finshed 

		log.Printf("Task added: %+v\n", newTask)
		c.JSON(http.StatusCreated, newTask)
	})

	// Start server
	log.Println("Starting server on :8080")
	router.Run(":8080")
}
