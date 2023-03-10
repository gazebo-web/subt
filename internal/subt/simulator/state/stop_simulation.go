package state

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	subtapp "gitlab.com/ignitionrobotics/web/subt/internal/subt/application"
)

// StopSimulation is the state of the action that stops a simulation.
type StopSimulation struct {
	platform     platform.Platform
	services     subtapp.Services
	GroupID      simulations.GroupID
	Score        *float64
	Stats        simulations.Statistics
	Summary      simulations.Summary
	RunData      string
	PodList      []resource.Resource
	UpstreamName string
}

// Platform returns the underlying platform.
func (s *StopSimulation) Platform() platform.Platform {
	return s.platform
}

// Services returns the underlying application services.
func (s *StopSimulation) Services() application.Services {
	return s.services
}

// SubTServices returns the subt specific application services.
func (s *StopSimulation) SubTServices() subtapp.Services {
	return s.services
}

// NewStopSimulation initializes a new state for stopping simulations.
func NewStopSimulation(platform platform.Platform, services subtapp.Services, groupID simulations.GroupID) *StopSimulation {
	return &StopSimulation{
		platform: platform,
		services: services,
		GroupID:  groupID,
	}
}
