package dto

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	OsmGeojsonUrl      string
	OsmRefreshInterval time.Duration
	MatchRadiusMeters  float64
}

func NewConfig() Config {
	return Config{
		OsmGeojsonUrl:      getEnv("OSM_GEOJSON_URL", "https://openaedmap.org/api/v1/countries/PL.geojson"),
		OsmRefreshInterval: getEnvDuration("OSM_REFRESH_INTERVAL", 3*time.Hour),
		MatchRadiusMeters:  getEnvFloat("MATCH_RADIUS_METERS", 100),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvFloat(key string, fallback float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}
	return parsed
}
