package repositories

import (
	"database/sql/driver"
	"github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	"order/models"
	"testing"
)

func Test_orderRepository_FindAllOrdersByUserID(t *testing.T) {
	Convey("test FindAllOrdersByUserI", t, func() {
		var (
			orderRepo OrderRepository
			orders    models.Orders
			err       error

			userID = 5
		)
		tx, err := gorm.Open("testdb", "")
		So(err, ShouldBeNil)
		orderRepo = NewOrderRepository(tx)

		Convey("with no records in the database", func() {
			testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
				columns := []string{"total", "currency_code", "user_id", "restaurant_id", "placed_at"}
				result := ""
				return testdb.RowsFromCSVString(columns, result), nil
			})
			orders, err = orderRepo.FindAllOrdersByUserID(userID)
			So(err, ShouldBeNil)
			So(len(orders), ShouldEqual, 0)
		})

		Convey("returns only the records belonging to the user", func() {
			testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
				columns := []string{"total", "currency_code", "user_id", "restaurant_id"}
				result := `
		1000,GBP,5,9
		2500,GBP,5,8
		`
				return testdb.RowsFromCSVString(columns, result), nil
			})
			orders, err = orderRepo.FindAllOrdersByUserID(userID)
			So(err, ShouldBeNil)
			So(len(orders), ShouldEqual, 2)
			So(orders[0].RestaurantID, ShouldEqual, 9)
			So(orders[1].RestaurantID, ShouldEqual, 8)
		})
	})
}
