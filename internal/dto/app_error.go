package dto

import "fmt"

type AppError error

var (
	OsmNotReady           = AppError(fmt.Errorf("osm data not loaded yet"))
	CommentEmpty          = AppError(fmt.Errorf("komentarz nie może być pusty"))
	CommentTooLong        = AppError(fmt.Errorf("komentarz może mieć maksymalnie 200 znaków"))
	CommentNodeIDRequired = AppError(fmt.Errorf("nodeID jest wymagane"))
)
