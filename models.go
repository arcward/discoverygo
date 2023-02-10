package discoverygo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Link is a link to another resource (see API spec)
type Link struct {
	Href      string `json:"href,omitempty"`
	Templated bool   `json:"templated,omitempty"`
}

// Links is a collection of links to other resources, for pagination
type Links struct {
	Self Link `json:"self,omitempty"`
	Next Link `json:"next,omitempty"`
	Prev Link `json:"prev,omitempty"`
}

// Page indicates the current page of a paginated response
type Page struct {
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
	Number        int `json:"number"`
}

// EmbeddedResponse is a collection of embedded resources from
// the "_embedded" field
type EmbeddedResponse struct {
	Events          []map[string]any `json:"events,omitempty"`
	Venues          []map[string]any `json:"venues,omitempty"`
	Attractions     []map[string]any `json:"attractions,omitempty"`
	Classifications []map[string]any `json:"classifications,omitempty"`
}

// PagedResponse is a response from the Discovery API - it can be paginated
// with the `NextPage` and `PrevPage` methods
type PagedResponse struct {
	Links    Links            `json:"_links,omitempty"`
	Page     Page             `json:"page"`
	Embedded EmbeddedResponse `json:"_embedded"`
}

// NextPage returns the next page of results from the Discovery API, for
// the given paged response
func (p *PagedResponse) NextPage(
	client *DiscoveryClient,
) (*PagedResponse, error) {
	if p.Page.Size*p.Page.Number >= 1000 {
		return nil, fmt.Errorf(
			"Max page depth reached (%d)",
			p.Page.Size*p.Page.Number,
		)
	}
	baseUrl := client.ApiUrl
	if p.Links.Next.Href == "" {
		return nil, nil
	}

	rel, _ := baseUrl.Parse(p.Links.Next.Href)
	q := rel.Query()
	q.Set("apikey", client.ApiKey)
	rel.RawQuery = q.Encode()

	resp, err := http.Get(rel.String())
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

	var rs PagedResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&rs)
	if decodeErr != nil {
		log.Println(decodeErr)
		return nil, decodeErr
	}
	return &rs, nil
}

// PreviousPage returns the previous page of results from the Discovery API, for
// the given paged response
func (p *PagedResponse) PreviousPage(
	client *DiscoveryClient,
) (*PagedResponse, error) {
	if p.Page.Size*p.Page.Number >= 1000 {
		return nil, fmt.Errorf(
			"Max page depth reached (%d)",
			p.Page.Size*p.Page.Number,
		)
	}
	baseUrl := client.ApiUrl
	if p.Links.Prev.Href == "" {
		return nil, nil
	}

	rel, _ := baseUrl.Parse(p.Links.Prev.Href)
	q := rel.Query()
	q.Set("apikey", client.ApiKey)
	rel.RawQuery = q.Encode()

	resp, err := http.Get(rel.String())
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

	var rs PagedResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&rs)
	if decodeErr != nil {
		log.Println(decodeErr)
		return nil, decodeErr
	}
	return &rs, nil
}
