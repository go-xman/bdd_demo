package services_test

import (
	"order/clients/restaurant"
	"order/models"
	"order/repositories"
	"order/services"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Services Suite")
}

var _ = Describe("Services", func() {
	var (
		restaurantClient restaurant.Client
		orderRepo        repositories.OrderRepository
		orderService     services.OrderService
		orders           models.Orders
		ctrl             *gomock.Controller
		err              error

		userID = 5
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	JustBeforeEach(func() {
		orderServiceImpl := services.NewOrderService()
		orderServiceImpl.SetOrderRepository(orderRepo)
		orderServiceImpl.SetRestaurantClient(restaurantClient)
		orderService = orderServiceImpl
	})

	Describe("FindAllOrdersByUserID", func() {
		Describe("with no records in the database", func() {
			BeforeEach(func() {
				orderRepoMock := repositories.NewMockOrderRepository(ctrl)
				orderRepoMock.EXPECT().FindAllOrdersByUserID(gomock.Eq(userID))
				orderRepo = orderRepoMock
			})

			It("return an empty slice of orders", func() {
				orders, err = orderService.FindAllOrderByUserID(userID)
				Expect(err).To(BeNil())
				Expect(len(orders)).To(Equal(0))
			})
		})

		Describe("when a few records exist", func() {
			// 设置order的mock数据
			BeforeEach(func() {
				order1 := &models.Order{
					Total:        1000,
					CurrencyCode: "GBP",
					UserID:       userID,
					RestaurantID: 8,
					PlacedAt:     time.Now().Add(-72 * time.Hour),
				}
				order2 := &models.Order{
					Total:        2500,
					CurrencyCode: "GBP",
					UserID:       userID,
					RestaurantID: 9,
					PlacedAt:     time.Now().Add(-36 * time.Hour),
				}
				orderRepoMock := repositories.NewMockOrderRepository(ctrl)
				orderRepoMock.EXPECT().FindAllOrdersByUserID(gomock.Eq(userID)).
					Return(models.Orders{order1, order2}, error(nil))
				orderRepo = orderRepoMock
			})

			Describe("when not all Restaurants can be found", func() {
				BeforeEach(func() {
					restaurantClientMock := restaurant.NewMockClient(ctrl)
					restaurantClientMock.EXPECT().
						GetRestaurantsByIDs([]int{8, 9}).
						Return(models.Restaurants{}, error(nil))
					restaurantClient = restaurantClientMock
				})

				It("returns only the records belonging to the user", func() {
					orders, err = orderService.FindAllOrderByUserID(userID)
					Expect(err).To(MatchError("restaurantItem with ID 8 not found"))
				})
			})

			Describe("when all Restaurants can be found", func() {
				BeforeEach(func() {
					restaurant1 := &models.Restaurant{
						ID:   9,
						Name: "Nando's",
					}

					restaurant2 := &models.Restaurant{
						ID:   8,
						Name: "KFC",
					}
					restaurantClientMock := restaurant.NewMockClient(ctrl)
					restaurantClientMock.EXPECT().
						GetRestaurantsByIDs([]int{8, 9}).
						Return(models.Restaurants{restaurant1, restaurant2}, nil)
					restaurantClient = restaurantClientMock
				})

				It("returns only the records belonging to the user", func() {
					orders, err = orderService.FindAllOrderByUserID(userID)
					Expect(err).To(BeNil())
					Expect(len(orders)).To(Equal(2))

					// 名称和排序
					Expect(orders[0].Restaurant.Name).To(Equal("KFC"))
					Expect(orders[0].Total).To(Equal(1000))
					Expect(orders[1].Restaurant.Name).To(Equal("Nando's"))
					Expect(orders[1].Total).To(Equal(2500))
				})
			})
		})
	})

	AfterEach(func() {
		ctrl.Finish()
	})
})
