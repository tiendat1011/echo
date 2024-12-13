package web

import (
	"sync"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductStore struct {
	mu       sync.RWMutex
	products map[int]Product
	nextID   int
}

func NewProductStore() *ProductStore {
	store := &ProductStore{
		products: make(map[int]Product),
		nextID:   1,
	}
	// Add sample products
	store.Add(Product{Name: "Laptop", Price: 1000.0})
	store.Add(Product{Name: "Smartphone", Price: 500.0})
	return store
}

func (ps *ProductStore) Add(product Product) int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	product.ID = ps.nextID
	ps.products[ps.nextID] = product
	ps.nextID++
	return product.ID
}

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

func (ps *ProductStore) Delete(id int) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, exists := ps.products[id]; !exists {
		return false
	}
	delete(ps.products, id)
	return true
}

func (ps *ProductStore) GetAll() []Product {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	products := make([]Product, 0, len(ps.products))
	for _, product := range ps.products {
		products = append(products, product)
	}
	return products
}

func (ps *ProductStore) GetByID(id int) (Product, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	product, exists := ps.products[id]
	return product, exists
}
