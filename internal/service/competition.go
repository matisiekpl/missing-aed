package service

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"strconv"

	backend "github.com/mwozniak/missing-aed"
	"github.com/mwozniak/missing-aed/internal/dto"
)

type CompetitionService interface {
	All() []dto.MissingAED
}

type competitionService struct {
	aeds []dto.MissingAED
}

type competitionRecord struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	City        string  `json:"city"`
	Description string  `json:"description"`
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lon"`
}

func NewCompetitionService() (CompetitionService, error) {
	aeds, err := loadCompetition()
	if err != nil {
		return nil, err
	}
	return &competitionService{aeds: aeds}, nil
}

func (c competitionService) All() []dto.MissingAED {
	return c.aeds
}

func loadCompetition() ([]dto.MissingAED, error) {
	entries, err := fs.ReadDir(backend.CompetitionFiles, "competition")
	if err != nil {
		return nil, err
	}
	aeds := make([]dto.MissingAED, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := backend.CompetitionFiles.ReadFile("competition/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", entry.Name(), err)
		}
		var record competitionRecord
		if err := json.Unmarshal(data, &record); err != nil {
			return nil, fmt.Errorf("parse %s: %w", entry.Name(), err)
		}
		if record.Latitude == 0 && record.Longitude == 0 {
			continue
		}
		aeds = append(aeds, dto.MissingAED{
			ID:          strconv.Itoa(record.ID),
			Latitude:    record.Latitude,
			Longitude:   record.Longitude,
			Name:        record.Name,
			Address:     record.Address,
			City:        record.City,
			Description: record.Description,
		})
	}
	return aeds, nil
}
