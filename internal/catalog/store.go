package catalog

// Product is shared by REST and gRPC handlers (FR-01, FR-02).
type Product struct {
	ID          string
	Name        string
	Price       float64
	Description string
}

// Store is an in-memory catalog — replace with Aurora in later lessons.
type Store struct {
	products []Product
}

func NewStore() *Store {
	return &Store{
		products: []Product{
			{ID: "1", Name: "Cheeky Mug", Price: 12.99, Description: "Ceramic mug"},
			{ID: "2", Name: "Cheeky Tote", Price: 8.50, Description: "Canvas tote bag"},
		},
	}
}

func (s *Store) List() []Product {
	out := make([]Product, len(s.products))
	copy(out, s.products)
	return out
}

func (s *Store) Get(id string) (Product, bool) {
	for _, p := range s.products {
		if p.ID == id {
			return p, true
		}
	}
	return Product{}, false
}
