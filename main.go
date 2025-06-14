package main

import (
	"fmt"
	"log"      // Logging steps
	"net/http" // allows HTTP requests
	"sync"     // Used for mutex - stops multiple threads from accessing the shared data at the same time (like a semaphore)

	"github.com/gin-gonic/gin" // GIN API FRAMEWORK
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

	// DELETE /tasks/:id
	router.DELETE("/tasks/:id", func(c *gin.Context) {
		taskIDParam := c.Param("id") // Get the :id parameter from URL
		var deleted bool //gonna use this to control my mutex 

		mutex.Lock()
		defer mutex.Unlock()

		for i, t := range tasks {
			if fmt.Sprintf("%d", t.ID) == taskIDParam {
				tasks = append(tasks[:i], tasks[i+1:]...) // Remove the task
				deleted = true // for mutex 
				break
			}
		}

		if deleted {
			log.Printf("Deleted", taskIDParam)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		}
	})


	// Start server
	log.Println("Starting server on :8080")
	router.Run(":8080")
}
