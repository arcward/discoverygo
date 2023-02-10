package discoverygo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// DiscoveryApiUrl is the base URL to the Ticketmaster Discovery API
const DiscoveryApiUrl = "https://app.ticketmaster.com/discovery/v2"

// DiscoveryClient is a client for the Ticketmaster Discovery API
type DiscoveryClient struct {
	// Base URL to the Discovery API
	ApiUrl url.URL
	// API key (consumer key)
	ApiKey string
}

// EventsUrl returns the URL to the events endpoint, with
// the API key added as a query parameter
func (d *DiscoveryClient) EventsUrl() url.URL {
	eventsUrl := d.ApiUrl.JoinPath("events")
	if d.ApiKey == "" {
		return *eventsUrl
	}
	q := eventsUrl.Query()
	q.Set("apikey", d.ApiKey)
	eventsUrl.RawQuery = q.Encode()
	return *eventsUrl
}

// VenuesUrl returns the URL to the venues endpoint, with
// the API key added as a query parameter
func (d *DiscoveryClient) VenuesUrl() url.URL {
	venuesUrl := d.ApiUrl.JoinPath("venues")
	if d.ApiKey == "" {
		return *venuesUrl
	}
	q := venuesUrl.Query()
	q.Set("apikey", d.ApiKey)
	venuesUrl.RawQuery = q.Encode()
	return *venuesUrl
}

// GetEvent returns an event by its ID
// See: https://developer.ticketmaster.com/products-and-docs/apis/discovery-api/v2/#event-details-v2
func (d *DiscoveryClient) GetEvent(id string) (*map[string]any, error) {
	baseEventUrl := d.EventsUrl()
	eventUrl := baseEventUrl.JoinPath(id)
	log.Printf("Querying: %s", eventUrl)
	resp, err := http.Get(eventUrl.String())
	defer resp.Body.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("Status code: %v", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"Status code: %d: %s",
			resp.StatusCode,
			body,
		)
	}
	var rs map[string]any
	decodeErr := json.NewDecoder(resp.Body).Decode(&rs)
	if decodeErr != nil {
		log.Println(decodeErr)
		return nil, decodeErr
	}
	return &rs, nil
}

// SearchEvents returns a list of events matching the given query parameters
func (d *DiscoveryClient) SearchEvents(
	queryParams QueryParams,
) (*PagedResponse, error) {
	eventsUrl, err := queryParams.UpdateURL(d.EventsUrl(), d.ApiKey)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(eventsUrl.String())
	defer resp.Body.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Printf("Status code: %v", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Status code: %d: %s", resp.StatusCode, body)
	}
	//var rs map[string]any
	var rs PagedResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&rs)
	if decodeErr != nil {
		log.Println(decodeErr)
		return nil, decodeErr
	}
	return &rs, nil
}

// QueryParams is a struct that holds the query parameters for the Discovery API
type QueryParams struct {
	Id                 string `json:"id,omitempty"`
	Sort               string `json:"sort,omitempty"`
	Page               string `json:"page,omitempty"`
	Size               string `json:"size,omitempty"`
	Locale             string `json:"locale,omitempty"`
	Keyword            string `json:"keyword,omitempty"`
	IncludeTest        string `json:"includeTest,omitempty"`
	IncludeTBA         string `json:"includeTBA,omitempty"`
	IncludeTBD         string `json:"includeTBD,omitempty"`
	VenueID            string `json:"venueId,omitempty"`
	StartDateTime      string `json:"startDateTime,omitempty"`
	EndDateTime        string `json:"endDateTime,omitempty"`
	CountryCode        string `json:"countryCode,omitempty"`
	StateCode          string `json:"stateCode,omitempty"`
	AttractionID       string `json:"attractionId,omitempty"`
	SegmentID          string `json:"segmentId,omitempty"`
	SegmentName        string `json:"segmentName,omitempty"`
	ClassificationID   string `json:"classificationId,omitempty"`
	ClassificationName string `json:"classificationName,omitempty"`
	MarketID           string `json:"marketId,omitempty"`
	PromoterID         string `json:"promoterId,omitempty"`
	DmaID              string `json:"dmaId,omitempty"`
	LatLong            string `json:"latlong,omitempty"`
	Radius             string `json:"radius,omitempty"`
	Unit               string `json:"unit,omitempty"`
}

// UpdateURL updates the given URL with the query parameters, and includes
// the API key as a query parameter
func (q QueryParams) UpdateURL(u url.URL, apikey string) (*url.URL, error) {
	var qp map[string]string
	inrec, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}
	unmarshalError := json.Unmarshal(inrec, &qp)
	if unmarshalError != nil {
		return nil, unmarshalError
	}

	query := u.Query()
	query.Set("apikey", apikey)
	for field, val := range qp {
		if val != "" {
			query.Add(field, val)
		}
	}
	u.RawQuery = query.Encode()
	return &u, nil
}

// redactUrl replaces the API key in the given URL with the string "REDACTED"
func redactUrl(u url.URL) string {
	query := u.Query()
	_, exists := query["apikey"]
	if exists {
		query.Set("apikey", "REDACTED")
	}
	u.RawQuery = query.Encode()
	return u.String()
}
