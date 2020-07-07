package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FindOrdersForUser(c *gin.Context) {
	p := &Provider{}
	p.FindOrdersForUser(c)
}

// FindOrdersForUser is the provider method that gets the orders for a user from
// the user's ID.
func (p *Provider) FindOrdersForUser(c Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	orders, err := p.getOrderService().FindAllOrderByUserID(userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, orders)
}
