package service

import (
	"github.com/mwozniak/missing-aed/internal/client"
	"github.com/mwozniak/missing-aed/internal/dto"
)

type Services interface {
	Osm() OsmService
	Competition() CompetitionService
	Matcher() MatcherService
}

type services struct {
	osmService         OsmService
	competitionService CompetitionService
	matcherService     MatcherService
}

func NewServices(clients client.Clients, config dto.Config) (Services, error) {
	competitionService, err := NewCompetitionService()
	if err != nil {
		return nil, err
	}
	osmService := NewOsmService(clients.OpenAedMap(), config)
	matcherService := NewMatcherService(osmService, competitionService, config)
	return &services{
		osmService:         osmService,
		competitionService: competitionService,
		matcherService:     matcherService,
	}, nil
}

func (s services) Osm() OsmService {
	return s.osmService
}

func (s services) Competition() CompetitionService {
	return s.competitionService
}

func (s services) Matcher() MatcherService {
	return s.matcherService
}
