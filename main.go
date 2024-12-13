package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Template Renderer
type TemplateRenderer struct {
	templates *template.Template
}

// Render implements echo.Renderer interface
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Product struct định nghĩa cấu trúc sản phẩm
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// ProductStore quản lý danh sách sản phẩm
type ProductStore struct {
	mu       sync.RWMutex
	products map[int]Product
	nextID   int
}

// NewProductStore tạo store mới cho sản phẩm
func NewProductStore() *ProductStore {
	return &ProductStore{
		products: make(map[int]Product),
		nextID:   1,
	}
}

// Add thêm sản phẩm mới
func (ps *ProductStore) Add(product Product) int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	product.ID = ps.nextID
	ps.products[ps.nextID] = product
	ps.nextID++
	return product.ID
}

// Update cập nhật sản phẩm
func (ps *ProductStore) Update(id int, product Product) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, exists := ps.products[id]; !exists {
		return false
	}
	product.ID = id
	ps.products[id] = product
	return true
}

// Delete xóa sản phẩm
func (ps *ProductStore) Delete(id int) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	_, exists := ps.products[id]
	if !exists {
		return false
	}
	delete(ps.products, id)
	return true
}

// GetAll lấy danh sách sản phẩm
func (ps *ProductStore) GetAll() []Product {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	products := make([]Product, 0, len(ps.products))
	for _, product := range ps.products {
		products = append(products, product)
	}
	return products
}

// GetByID lấy sản phẩm theo ID
func (ps *ProductStore) GetByID(id int) (Product, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	product, exists := ps.products[id]
	return product, exists
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Khởi tạo templates
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer

	// Khởi tạo store sản phẩm
	productStore := NewProductStore()

	// Thêm một số sản phẩm mặc định
	productStore.Add(Product{Name: "Laptop", Price: 1000.0})
	productStore.Add(Product{Name: "Smartphone", Price: 500.0})

	// Trang chủ
	e.GET("/", func(c echo.Context) error {
		products := productStore.GetAll()
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"Products": products,
		})
	})

	e.GET("/products/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		product, exists := productStore.GetByID(id)
		if !exists {
			return c.HTML(http.StatusNotFound, "")
		}
		return c.Render(http.StatusOK, "product-row.html", product)
	})

	// Route thêm sản phẩm
	e.POST("/products", func(c echo.Context) error {
		name := c.FormValue("name")
		price, _ := strconv.ParseFloat(c.FormValue("price"), 64)

		newProduct := Product{
			Name:  name,
			Price: price,
		}
		id := productStore.Add(newProduct)
		newProduct.ID = id

		return c.Render(http.StatusCreated, "product-row.html", newProduct)
	})

	// Route xóa sản phẩm
	e.DELETE("/products/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		if productStore.Delete(id) {
			return c.HTML(http.StatusOK, "")
		}
		return c.HTML(http.StatusNotFound, "")
	})

	// Route chỉnh sửa sản phẩm
	e.PUT("/products/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		name := c.FormValue("name")
		price, _ := strconv.ParseFloat(c.FormValue("price"), 64)

		product, exists := productStore.GetByID(id)
		if !exists {
			return c.HTML(http.StatusNotFound, "")
		}

		product.Name = name
		product.Price = price
		productStore.Update(id, product)

		return c.Render(http.StatusOK, "product-row.html", product)
	})

	// Route hiển thị form chỉnh sửa
	e.GET("/products/:id/edit", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		product, exists := productStore.GetByID(id)
		if !exists {
			return c.HTML(http.StatusNotFound, "")
		}
		return c.Render(http.StatusOK, "product-edit.html", product)
	})

	// Phục vụ tệp tĩnh
	e.Static("/static", "static")

	// Khởi động server
	e.Start(":8080")
}
