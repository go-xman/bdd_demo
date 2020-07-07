package restaurant

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"order/models"
	"os"
	"strconv"
	"strings"
)

type Client interface {
	GetRestaurantsByIDs(ids []int) (models.Restaurants, error)
}

func NewClient() Client {
	return &client{}
}

type client struct {
	baseURL string
}

func (c *client) GetRestaurantsByIDs(ids []int) (models.Restaurants, error) {
	if len(ids) == 0 {
		return []*models.Restaurant{}, nil
	}

	idStrings := make([]string, 0, len(ids))
	for _, id := range ids {
		idStrings = append(idStrings, strconv.Itoa(id))
	}

	url := fmt.Sprintf(
		"%s/v1/restaurants?id=%s",
		c.getBaseURL(),
		strings.Join(idStrings, ","),
	)

	ret, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if ret.StatusCode != http.StatusOK {
		return nil, errors.New("err retrieving restaurants from RestaurantService")
	}

	body, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		return nil, err
	}

	parsedBody := models.Restaurants{}
	if err = json.Unmarshal(body, &parsedBody); err != nil {
		return nil, err
	}

	return parsedBody, nil
}

func (c *client) SetBaseUrl(url string) {
	c.baseURL = url
}

func (c *client) getBaseURL() string {
	if c.baseURL != "" {
		return c.baseURL
	}
	c.baseURL = os.Getenv("RESTAURANT_SERVICE_BASE_URL")
	return c.baseURL
}
