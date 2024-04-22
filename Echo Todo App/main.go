package main

import (
  "html/template"
  "io"
  "net/http"
  "strconv"

  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
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
  e := echo.New()
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }
  db.AutoMigrate(&Todo{})

  renderer := &TemplateRenderer{
    templates: template.Must(template.ParseGlob("views/*.html")),
  }
  e.Renderer = renderer

  e.GET("/", home)
  e.POST("/add", add)
  e.GET("/update/:id", update)
  e.GET("/delete/:id", delete)

  e.Logger.Fatal(e.Start(":8080"))
}

type TemplateRenderer struct {
  templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  return t.templates.ExecuteTemplate(w, name, data)
}

func home(c echo.Context) error {
  var todoList []Todo
  db.Find(&todoList)
  return c.Render(http.StatusOK, "base.html", map[string]interface{}{
    "todoList": todoList,
  })
}

func add(c echo.Context) error {
  title := c.FormValue("title")
  newTodo := Todo{Title: title, Complete: false}
  db.Create(&newTodo)
  return c.Redirect(http.StatusFound, "/")
}

func update(c echo.Context) error {
  idStr := c.Param("id")
  id, err := strconv.Atoi(idStr)
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "invalid ID")
  }

  var todo Todo
  db.First(&todo, id)
  todo.Complete = !todo.Complete
  db.Save(&todo)
  return c.Redirect(http.StatusFound, "/")
}

func delete(c echo.Context) error {
  idStr := c.Param("id")
  id, err := strconv.Atoi(idStr)
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "invalid ID")
  }

  var todo Todo
  db.Delete(&todo, id)
  return c.Redirect(http.StatusFound, "/")
}