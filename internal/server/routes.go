package server

import (
	"io"
	"net/http"
	"strconv"
	"test/cmd/web"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("cmd/web/*.html")),
	}
	e.Renderer = renderer

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/", s.handleIndex)
	e.POST("/products", s.handleAddProduct)
	e.GET("/products/:id", s.handleGetProduct)
	e.PUT("/products/:id", s.handleUpdateProduct)
	e.DELETE("/products/:id", s.handleDeleteProduct)
	e.GET("/products/:id/edit", s.handleEditProduct)

	return e
}

func (s *Server) handleIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Products": s.productStore.GetAll(),
	})
}

func (s *Server) handleAddProduct(c echo.Context) error {
	name := c.FormValue("name")
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)

	product := web.Product{
		Name:  name,
		Price: price,
	}
	id := s.productStore.Add(product)
	product.ID = id

	return c.Render(http.StatusOK, "product-row.html", product)
}

func (s *Server) handleGetProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	product, exists := s.productStore.GetByID(id)
	if !exists {
		return c.NoContent(http.StatusNotFound)
	}
	return c.Render(http.StatusOK, "product-row.html", product)
}

func (s *Server) handleUpdateProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	name := c.FormValue("name")
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)

	product := web.Product{
		ID:    id,
		Name:  name,
		Price: price,
	}

	if !s.productStore.Update(id, product) {
		return c.NoContent(http.StatusNotFound)
	}

	return c.Render(http.StatusOK, "product-row.html", product)
}

func (s *Server) handleDeleteProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if !s.productStore.Delete(id) {
		return c.NoContent(http.StatusNotFound)
	}
	return c.NoContent(http.StatusOK)
}

func (s *Server) handleEditProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	product, exists := s.productStore.GetByID(id)
	if !exists {
		return c.NoContent(http.StatusNotFound)
	}
	return c.Render(http.StatusOK, "product-edit.html", product)
}
