package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/mwozniak/missing-aed/internal/dto"
)

type OpenAedMapClient interface {
	Fetch(ctx context.Context) ([]dto.AED, error)
}

type openAedMapClient struct {
	url        string
	httpClient *http.Client
}

func NewOpenAedMapClient(url string) OpenAedMapClient {
	return &openAedMapClient{
		url:        url,
		httpClient: &http.Client{Timeout: 5 * time.Minute},
	}
}

type geojsonFeatureCollection struct {
	Features []geojsonFeature `json:"features"`
}

type geojsonFeature struct {
	Geometry   geojsonGeometry        `json:"geometry"`
	Properties map[string]any `json:"properties"`
}

type geojsonGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func (o openAedMapClient) Fetch(ctx context.Context) ([]dto.AED, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, o.url, nil)
	if err != nil {
		return nil, err
	}
	response, err := o.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from openaedmap", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var collection geojsonFeatureCollection
	if err := json.Unmarshal(body, &collection); err != nil {
		return nil, err
	}
	aeds := make([]dto.AED, 0, len(collection.Features))
	for _, feature := range collection.Features {
		if feature.Geometry.Type != "Point" || len(feature.Geometry.Coordinates) < 2 {
			continue
		}
		longitude := feature.Geometry.Coordinates[0]
		latitude := feature.Geometry.Coordinates[1]
		identifier := extractOsmId(feature.Properties)
		aeds = append(aeds, dto.AED{
			ID:        identifier,
			Latitude:  latitude,
			Longitude: longitude,
			Source:    "osm",
		})
	}
	return aeds, nil
}

func extractOsmId(properties map[string]any) string {
	value, ok := properties["@osm_id"]
	if !ok {
		return ""
	}
	switch typed := value.(type) {
	case float64:
		return strconv.FormatInt(int64(typed), 10)
	case string:
		return typed
	}
	return ""
}
