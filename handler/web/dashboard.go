package web

import (
	"a21hc3NpZ25tZW50/client"
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"embed"
	"net/http"
	"path"
	"text/template"

	"github.com/gin-gonic/gin"
)

type DashboardWeb interface {
	Dashboard(c *gin.Context)
}

type dashboardWeb struct {
	sessionService service.SessionService
	taskClient     client.TaskClient
	userService    service.UserService
	embed          embed.FS
}

func NewDashboardWeb(sessionService service.SessionService, taskClient client.TaskClient, userService service.UserService, embed embed.FS) *dashboardWeb {
	return &dashboardWeb{sessionService, taskClient, userService, embed}
}

func (d *dashboardWeb) Dashboard(c *gin.Context) {
	var email string
	if temp, ok := c.Get("email"); ok {
		if contextData, ok := temp.(string); ok {
			email = contextData
		}
	}

	session, err := d.sessionService.GetSessionByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	// Get tasks for the logged-in user (filtered by user ID)
	tasks, err := d.taskClient.TaskList(session.Token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "Error getting user tasks: " + err.Error(),
		})
		return
	}

	user, err := d.userService.GetUserByEmail(email)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	categoryClient := client.NewCategoryClient()
	categories, err := categoryClient.CategoryList(session.Token)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	categoryByID := make(map[int]string)
	for _, category := range categories {
		categoryByID[category.ID] = category.Name
	}

	// Convert tasks to UserTaskCategory format for template compatibility
	var userTaskCategories []model.UserTaskCategory
	if tasks != nil {
		for _, task := range tasks {
			categoryName := categoryByID[task.CategoryID]
			if categoryName == "" {
				categoryName = "Unknown"
			}

			userTaskCategory := model.UserTaskCategory{
				ID:       task.ID,
				Fullname: user.Fullname,
				Email:    email,
				Task:     task.Title,
				Deadline: task.Deadline,
				Priority: task.Priority,
				Status:   task.Status,
				Category: categoryName,
			}
			userTaskCategories = append(userTaskCategories, userTaskCategory)
		}
	}

	dataLength := len(userTaskCategories)

	var dataTemplate = map[string]interface{}{
		"email":                email,
		"user_task_categories": userTaskCategories,
		"data_count":           dataLength,
		"has_sample_data":      false,
	}

	var funcMap = template.FuncMap{
		"exampleFunc": func() int {
			return 0
		},
	}

	var header = path.Join("views", "general", "header.html")
	var filepath = path.Join("views", "main", "dashboard.html")

	t, err := template.New("dashboard.html").Funcs(funcMap).ParseFS(d.embed, filepath, header)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	err = t.Execute(c.Writer, dataTemplate)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
	}
}
