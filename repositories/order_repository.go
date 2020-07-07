package repositories

import (
	"github.com/jinzhu/gorm"
	"order/models"
)

type OrderRepository interface {
	FindAllOrdersByUserID(userID int) (models.Orders, error)
}

// NewOrderRepository new order Repository
func NewOrderRepository(db *gorm.DB) orderRepository {
	r := orderRepository{}
	r.SetDB(db)
	return r
}

type orderRepository struct {
	db *gorm.DB
}

func (r *orderRepository) SetDB(db *gorm.DB) {
	r.db = db
}

func (r *orderRepository) GetDB() (db *gorm.DB) {
	return r.db
}

func (r orderRepository) FindAllOrdersByUserID(userID int) (models.Orders, error) {
	orders := models.Orders{}
	err := r.GetDB().Table("orders").Where("user_id=?", userID).
		Order("placed_at DESC").Find(&orders).Error
	return orders, err
}
