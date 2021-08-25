package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/subt/internal/subt/summaries"
	"gitlab.com/ignitionrobotics/web/subt/internal/subt/tracks"
	"gitlab.com/ignitionrobotics/web/subt/pkg/runsim"
)

// Services extends a generic application services interface to add SubT services.
type Services interface {
	application.Services
	Tracks() tracks.Service
	Summaries() summaries.Service
	RunningSimulations() runsim.Manager
}

// services is a Services implementation.
type services struct {
	application.Services
	tracks             tracks.Service
	summaries          summaries.Service
	runningSimulations runsim.Manager
}

func (s *services) RunningSimulations() runsim.Manager {
	return s.runningSimulations
}

func (s *services) Summaries() summaries.Service {
	return s.summaries
}

// Tracks returns a Track service.
func (s *services) Tracks() tracks.Service {
	return s.tracks
}

// NewServices initializes a new Services implementation using a base generic service.
func NewServices(base application.Services, tracks tracks.Service, summaries summaries.Service) Services {
	return &services{
		Services:           base,
		tracks:             tracks,
		summaries:          summaries,
		runningSimulations: runsim.NewManager(),
	}
}
