package types

type SKU struct {
	SKUID        string
	Prices       []Price
	UnitQuantity int
	Name         string
	ImageURL     string
}

type Price struct {
	OrderQuantity int
	Discount      float64
	Price         float64
}
