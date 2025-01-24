package order

import "gorm.io/gorm"

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) CreateOrder(userID uint, productID uint, quantity int) (*Order, error) {
	order := &Order{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		Status:    "pending",
	}
	return order, s.db.Create(order).Error
}
