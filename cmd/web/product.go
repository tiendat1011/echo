package web

import (
	"sort"
	"sync"
)

const (
	ASC  = "asc"
	DESC = "desc"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductStore struct {
	mu        sync.RWMutex
	products  map[int]Product
	nextID    int
	sortField string
	sortOrder string
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

func (ps *ProductStore) GetAllSorted(sortParam string) []Product {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	products := make([]Product, 0, len(ps.products))
	for _, product := range ps.products {
		products = append(products, product)
	}

	// Toggle sort order if same field is clicked
	if sortParam == ps.sortField {
		if ps.sortOrder == ASC {
			ps.sortOrder = DESC
		} else {
			ps.sortOrder = ASC
		}
	} else {
		ps.sortField = sortParam
		ps.sortOrder = ASC
	}

	// Sort based on field
	sort.Slice(products, func(i, j int) bool {
		var comparison bool
		switch ps.sortField {
		case "id":
			comparison = products[i].ID < products[j].ID
		case "name":
			comparison = products[i].Name < products[j].Name
		case "price":
			comparison = products[i].Price < products[j].Price
		default:
			comparison = products[i].ID < products[j].ID
		}

		if ps.sortOrder == DESC {
			return !comparison
		}
		return comparison
	})

	return products
}

func (ps *ProductStore) GetSortField() string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.sortField
}

func (ps *ProductStore) GetSortOrder() string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.sortOrder
}
