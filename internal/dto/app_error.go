package dto

import "fmt"

type AppError error

var (
	OsmNotReady = AppError(fmt.Errorf("osm data not loaded yet"))
)
