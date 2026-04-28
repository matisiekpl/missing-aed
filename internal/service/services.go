package service

import (
	"github.com/mwozniak/missing-aed/internal/client"
	"github.com/mwozniak/missing-aed/internal/dto"
	"github.com/mwozniak/missing-aed/internal/repository"
)

type Services interface {
	Osm() OsmService
	Competition() CompetitionService
	Matcher() MatcherService
	Comment() CommentService
}

type services struct {
	osmService         OsmService
	competitionService CompetitionService
	matcherService     MatcherService
	commentService     CommentService
}

func NewServices(clients client.Clients, repositories repository.Repositories, config dto.Config) (Services, error) {
	competitionService, err := NewCompetitionService()
	if err != nil {
		return nil, err
	}
	osmService := NewOsmService(clients.OpenAedMap(), config)
	matcherService := NewMatcherService(osmService, competitionService, config)
	commentService := NewCommentService(repositories.Comment())
	return &services{
		osmService:         osmService,
		competitionService: competitionService,
		matcherService:     matcherService,
		commentService:     commentService,
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

func (s services) Comment() CommentService {
	return s.commentService
}
