package client

import "github.com/mwozniak/missing-aed/internal/dto"

type Clients interface {
	OpenAedMap() OpenAedMapClient
}

type clients struct {
	openAedMapClient OpenAedMapClient
}

func NewClients(config dto.Config) Clients {
	return &clients{
		openAedMapClient: NewOpenAedMapClient(config.OsmGeojsonUrl),
	}
}

func (c clients) OpenAedMap() OpenAedMapClient {
	return c.openAedMapClient
}
