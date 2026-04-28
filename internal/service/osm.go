package service

import (
	"context"
	"sync"
	"time"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"github.com/sirupsen/logrus"

	"github.com/mwozniak/missing-aed/internal/client"
	"github.com/mwozniak/missing-aed/internal/dto"
)

const (
	earthRadiusMeters = 6371010.0
	osmIndexLevel     = 16
)

type OsmService interface {
	Refresh(ctx context.Context) error
	Snapshot() (*OsmSnapshot, bool)
	Start(ctx context.Context)
}

type OsmSnapshot struct {
	Points  []s2.Point
	Buckets map[s2.CellID][]int
	Level   int
}

type osmService struct {
	client   client.OpenAedMapClient
	config   dto.Config
	mutex    sync.RWMutex
	snapshot *OsmSnapshot
}

func NewOsmService(openAedMapClient client.OpenAedMapClient, config dto.Config) OsmService {
	return &osmService{
		client: openAedMapClient,
		config: config,
	}
}

func (o *osmService) Refresh(ctx context.Context) error {
	logrus.Info("refreshing OSM AED data")
	aeds, err := o.client.Fetch(ctx)
	if err != nil {
		return err
	}
	snapshot := buildOsmSnapshot(aeds, osmIndexLevel)
	o.mutex.Lock()
	o.snapshot = snapshot
	o.mutex.Unlock()
	logrus.Infof("loaded %d OSM AEDs into S2 index", len(aeds))
	return nil
}

func (o *osmService) Snapshot() (*OsmSnapshot, bool) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	if o.snapshot == nil {
		return nil, false
	}
	return o.snapshot, true
}

func (o *osmService) Start(ctx context.Context) {
	if err := o.Refresh(ctx); err != nil {
		logrus.Errorf("initial OSM refresh failed: %v", err)
	}
	ticker := time.NewTicker(o.config.OsmRefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := o.Refresh(ctx); err != nil {
				logrus.Errorf("OSM refresh failed: %v", err)
			}
		}
	}
}

func buildOsmSnapshot(aeds []dto.AED, level int) *OsmSnapshot {
	points := make([]s2.Point, len(aeds))
	buckets := make(map[s2.CellID][]int, len(aeds))
	for index, aed := range aeds {
		latLng := s2.LatLngFromDegrees(aed.Latitude, aed.Longitude)
		point := s2.PointFromLatLng(latLng)
		points[index] = point
		cellId := s2.CellIDFromLatLng(latLng).Parent(level)
		buckets[cellId] = append(buckets[cellId], index)
	}
	return &OsmSnapshot{Points: points, Buckets: buckets, Level: level}
}

func metersToAngle(meters float64) s1.Angle {
	return s1.Angle(meters / earthRadiusMeters)
}
