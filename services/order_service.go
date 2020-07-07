package services

import (
	"fmt"
	"order/application"
	"order/clients/restaurant"
	"order/models"
	"order/repositories"

	"github.com/jinzhu/gorm"
)

type OrderService interface {
	FindAllOrderByUserID(userID int) (models.Orders, error)
}

// NewOrderService new order service
func NewOrderService() *orderService {
	return &orderService{}
}

type orderService struct {
	db               gorm.DB
	restaurantClient restaurant.Client
	orderRepository  repositories.OrderRepository
}

// SetOrderRepository set order Repository
func (s *orderService) SetOrderRepository(r repositories.OrderRepository) {
	s.orderRepository = r
}

func (s *orderService) getOrderRepository() repositories.OrderRepository {
	if s.orderRepository != nil {
		return s.orderRepository
	}

	s.orderRepository = repositories.NewOrderRepository(application.ResolveDB())
	return s.orderRepository
}

// SetRestaurantClient set restaurant client
func (s *orderService) SetRestaurantClient(c restaurant.Client) {
	s.restaurantClient = c
}

func (s *orderService) getRestaurantClient() restaurant.Client {
	if s.restaurantClient != nil {
		return s.restaurantClient
	}

	s.restaurantClient = restaurant.NewClient()
	return s.restaurantClient
}

// FindAllOrderByUserID find order by user id
func (s orderService) FindAllOrderByUserID(userID int) (models.Orders, error) {
	orders, err := s.getOrderRepository().FindAllOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return orders, nil
	}

	restaurantIDs := make([]int, 0, len(orders))
	for _, order := range orders {
		restaurantIDs = append(restaurantIDs, order.RestaurantID)
	}

	restaurants, err := s.getRestaurantClient().GetRestaurantsByIDs(restaurantIDs)
	if err != nil {
		return nil, err
	}

	restaurantsByID := make(map[int]*models.Restaurant)
	for _, restaurantItem := range restaurants {
		restaurantsByID[restaurantItem.ID] = restaurantItem
	}

	for _, order := range orders {
		restaurantItem, ok := restaurantsByID[order.RestaurantID]
		if !ok {
			return nil, fmt.Errorf("restaurantItem with ID %d not found", order.RestaurantID)
		}
		order.Restaurant = restaurantItem
	}
	return orders, nil
}
