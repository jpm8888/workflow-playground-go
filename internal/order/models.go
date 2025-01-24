package order

type Order struct {
	gorm.Model
	UserID    uint
	Status    string
	Amount    float64
	ProductID uint
	Quantity  int
}
