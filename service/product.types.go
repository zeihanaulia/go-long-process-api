package service

type UpdaterRequest struct {
	Data []Product
}

type Product struct {
	ID      int
	StoreID string
	SKU     string
	Name    string
}
