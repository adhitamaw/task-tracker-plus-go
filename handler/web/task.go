package web

import (
	"a21hc3NpZ25tZW50/client"
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"embed"
	"net/http"
	"path"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
)

type TaskWeb interface {
	TaskPage(c *gin.Context)
	TaskAddProcess(c *gin.Context)
	TaskDeleteProcess(c *gin.Context)
}

type taskWeb struct {
	taskClient     client.TaskClient
	sessionService service.SessionService
	userService    service.UserService
	embed          embed.FS
}

func NewTaskWeb(taskClient client.TaskClient, sessionService service.SessionService, userService service.UserService, embed embed.FS) *taskWeb {
	return &taskWeb{taskClient, sessionService, userService, embed}
}

func (t *taskWeb) TaskPage(c *gin.Context) {
	var email string
	if temp, ok := c.Get("email"); ok {
		if contextData, ok := temp.(string); ok {
			email = contextData
		}
	}

	session, err := t.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	tasks, err := t.taskClient.TaskList(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	// Get categories for dropdown
	categoryClient := client.NewCategoryClient()
	categories, err := categoryClient.CategoryList(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	var dataTemplate = map[string]interface{}{
		"email":      email,
		"tasks":      tasks,
		"categories": categories,
	}

	var funcMap = template.FuncMap{
		"exampleFunc": func() int {
			return 0
		},
	}

	var header = path.Join("views", "general", "header.html")
	var filepath = path.Join("views", "main", "task.html")

	temp, err := template.New("task.html").Funcs(funcMap).ParseFS(t.embed, filepath, header)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	err = temp.Execute(c.Writer, dataTemplate)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
	}
}

func (t *taskWeb) TaskAddProcess(c *gin.Context) {
	var email string
	if temp, ok := c.Get("email"); ok {
		if contextData, ok := temp.(string); ok {
			email = contextData
		}
	}

	session, err := t.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	// Tangkap semua data form
	title := c.Request.FormValue("title")
	deadline := c.Request.FormValue("deadline")
	priorityStr := c.Request.FormValue("priority")
	status := c.Request.FormValue("status")
	categoryIDStr := c.Request.FormValue("category-id")

	// Validasi data form
	if title == "" {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Title cannot be empty")
		return
	}

	if deadline == "" {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Deadline cannot be empty")
		return
	}

	if status == "" {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Status cannot be empty")
		return
	}

	// Konversi string ke integer
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Invalid priority value")
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil || categoryID <= 0 {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Invalid category ID")
		return
	}

	user, err := t.userService.GetUserByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	task := model.Task{
		Title:      title,
		Deadline:   deadline,
		Priority:   priority,
		Status:     status,
		CategoryID: categoryID,
		UserID:     user.ID,
	}

	statusCode, err := t.taskClient.AddTask(session.Token, task)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	if statusCode == 200 || statusCode == 201 {
		c.Redirect(http.StatusSeeOther, "/client/task")
	} else {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Add Task Failed with status: "+strconv.Itoa(statusCode))
	}
}

func (t *taskWeb) TaskDeleteProcess(c *gin.Context) {
	var email string
	if temp, ok := c.Get("email"); ok {
		if contextData, ok := temp.(string); ok {
			email = contextData
		}
	}

	session, err := t.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get task ID from URL parameter
	taskIDStr := c.Param("id")
	if taskIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	statusCode, err := t.taskClient.DeleteTask(session.Token, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if statusCode == 200 {
		c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
	}
}
