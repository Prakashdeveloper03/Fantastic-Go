package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

type Todo struct {
	ID       uint `gorm:"primaryKey"`
	Title    string
	Complete bool
}

func main() {
	r := gin.Default()

	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Todo{})

	renderer := TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	r.SetHTMLTemplate(renderer.templates)

	r.GET("/", home)
	r.POST("/add", add)
	r.GET("/update/:id", update)
	r.GET("/delete/:id", delete)

	r.Run(":8080")
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w http.ResponseWriter, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func home(c *gin.Context) {
	var todoList []Todo
	db.Find(&todoList)
	c.HTML(http.StatusOK, "base.html", gin.H{
		"todoList": todoList,
	})
}

func add(c *gin.Context) {
	title := c.PostForm("title")
	newTodo := Todo{Title: title, Complete: false}
	db.Create(&newTodo)
	c.Redirect(http.StatusFound, "/")
}

func update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid ID")
		return
	}

	var todo Todo
	db.First(&todo, id)
	todo.Complete = !todo.Complete
	db.Save(&todo)
	c.Redirect(http.StatusFound, "/")
}

func delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid ID")
		return
	}

	var todo Todo
	db.Delete(&todo, id)
	c.Redirect(http.StatusFound, "/")
}