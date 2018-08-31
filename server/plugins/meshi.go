package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/otiai10/amesh/server/middlewares"
)

// Meshi ...
type Meshi struct {
	YelpAPIKey string
}

// Match ...
func (meshi Meshi) Match(ctx context.Context, texts []string) bool {
	if len(texts) == 0 {
		return false
	}
	return texts[0] == "meshi"
}

// Method ...
func (meshi Meshi) Method() string {
	return "meshi"
}

// TaskValues ...
func (meshi Meshi) TaskValues(ctx context.Context, texts []string) url.Values {
	values := url.Values{}
	values.Set("location", strings.Join(texts[1:], "+"))
	return values
}

// Exec ...
func (meshi Meshi) Exec(ctx context.Context, r *http.Request) (string, error) {
	query := url.Values{}
	location := r.FormValue("location")
	query.Add("location", location)
	query.Add("locale", "ja_JP")
	query.Add("open_now", "true")
	req, err := http.NewRequest("GET", "https://api.yelp.com/v3/businesses/search?"+query.Encode(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", meshi.YelpAPIKey))

	res, err := middlewares.HTTPClient(ctx).Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	response := new(YelpAPIResponse)
	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return "", err
	}

	if len(response.Businesses) == 0 {
		return "", fmt.Errorf("not found for location query: %v", location)
	}

	rand.Seed(time.Now().Unix())
	business := response.Businesses[rand.Intn(len(response.Businesses))]

	text := fmt.Sprintf(
		"%s\nreview: %d, rating: %.1f\n%s\n%s",
		business.Name, business.ReviewCount, business.Rating, business.URL, business.ImageURL,
	)
	return text, nil
}

// TODO: どっかに持ってく
type (
	// YelpAPIResponse ...
	YelpAPIResponse struct {
		Businesses []YelpBusiness `json:"businesses"`
		Region     YelpRegion     `json:"region"`
	}
	// YelpBusiness ...
	YelpBusiness struct {
		ID          string  `json:"id"`
		Alias       string  `json:"alias"`
		Name        string  `json:"name"`
		URL         string  `json:"url"`
		ImageURL    string  `json:"image_url"`
		ReviewCount int     `json:"review_count"`
		Rating      float32 `json:"rating"`
	}
	// YelpRegion ...
	YelpRegion struct {
		Center struct {
			Longitude float32 `json:"longitude"`
			Latitude  float32 `json:"latitude"`
		} `json:"center"`
	}
)
