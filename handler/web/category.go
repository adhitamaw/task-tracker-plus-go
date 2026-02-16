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

type CategoryWeb interface {
	Category(c *gin.Context)
	AddCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}

type categoryWeb struct {
	categoryClient client.CategoryClient
	sessionService service.SessionService
	embed          embed.FS
}

func NewCategoryWeb(categoryClient client.CategoryClient, sessionService service.SessionService, embed embed.FS) *categoryWeb {
	return &categoryWeb{categoryClient, sessionService, embed}
}

func (c *categoryWeb) Category(ctx *gin.Context) {
	var email string
	if temp, ok := ctx.Get("email"); ok {
		if contextData, ok := temp.(string); ok {
			email = contextData
		}
	}

	session, err := c.sessionService.GetSessionByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	categories, err := c.categoryClient.CategoryList(session.Token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var dataTemplate = map[string]interface{}{
		"email":      email,
		"categories": categories,
	}

	var funcMap = template.FuncMap{
		"exampleFunc": func() int {
			return 0
		},
	}

	var header = path.Join("views", "general", "header.html")
	var filepath = path.Join("views", "main", "category.html")

	t, err := template.New("category.html").Funcs(funcMap).ParseFS(c.embed, filepath, header)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	err = t.Execute(ctx.Writer, dataTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
	}
}

func (c *categoryWeb) AddCategory(ctx *gin.Context) {
	var email string
	if temp, ok := ctx.Get("email"); ok {
		if contextData, ok := temp.(string); ok {
			email = contextData
		}
	}

	session, err := c.sessionService.GetSessionByEmail(email)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	name := ctx.Request.FormValue("name")
	if name == "" {
		ctx.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Category name cannot be empty")
		return
	}

	status, err := c.categoryClient.AddCategory(session.Token, name)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	if status == 200 || status == 201 {
		ctx.Redirect(http.StatusSeeOther, "/client/category")
	} else {
		ctx.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Add Category Failed!")
	}
}

func (c *categoryWeb) DeleteCategory(ctx *gin.Context) {
	var email string
	if temp, ok := ctx.Get("email"); ok {
		if contextData, ok := temp.(string); ok {
			email = contextData
		}
	}

	session, err := c.sessionService.GetSessionByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get category ID from URL parameter
	categoryIDStr := ctx.Param("id")
	if categoryIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Category ID is required"})
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	statusCode, err := c.categoryClient.DeleteCategory(session.Token, strconv.Itoa(categoryID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if statusCode == 200 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
	}
}
