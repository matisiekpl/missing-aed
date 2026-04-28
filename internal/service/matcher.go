package service

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"

	"github.com/mwozniak/missing-aed/internal/dto"
)

type MatcherService interface {
	Missing() ([]dto.MissingAED, error)
}

type matcherService struct {
	osmService         OsmService
	competitionService CompetitionService
	config             dto.Config
}

func NewMatcherService(osmService OsmService, competitionService CompetitionService, config dto.Config) MatcherService {
	return &matcherService{
		osmService:         osmService,
		competitionService: competitionService,
		config:             config,
	}
}

func (m matcherService) Missing() ([]dto.MissingAED, error) {
	snapshot, ready := m.osmService.Snapshot()
	if !ready {
		return nil, dto.OsmNotReady
	}
	radiusMeters := m.config.MatchRadiusMeters
	radiusAngle := metersToAngle(radiusMeters)
	coverer := &s2.RegionCoverer{
		MinLevel: snapshot.Level,
		MaxLevel: snapshot.Level,
		MaxCells: 16,
	}
	missing := make([]dto.MissingAED, 0)
	for _, candidate := range m.competitionService.All() {
		latLng := s2.LatLngFromDegrees(candidate.Latitude, candidate.Longitude)
		center := s2.PointFromLatLng(latLng)
		cap := s2.CapFromCenterAngle(center, radiusAngle)
		covering := coverer.Covering(cap)
		if !hasNeighbor(snapshot, covering, center, radiusAngle) {
			missing = append(missing, candidate)
		}
	}
	return missing, nil
}

func hasNeighbor(snapshot *OsmSnapshot, covering s2.CellUnion, center s2.Point, radius s1.Angle) bool {
	for _, cellId := range covering {
		for _, index := range snapshot.Buckets[cellId] {
			if center.Distance(snapshot.Points[index]) <= radius {
				return true
			}
		}
	}
	return false
}
