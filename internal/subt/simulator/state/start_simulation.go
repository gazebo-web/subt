package state

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport/ign"
	subtapp "gitlab.com/ignitionrobotics/web/subt/internal/subt/application"
)

// StartSimulation is the state of the action that starts a simulation.
type StartSimulation struct {
	state.PlatformGetter
	state.ServicesGetter
	platform             platform.Platform
	services             subtapp.Services
	GroupID              simulations.GroupID
	GazeboServerPod      resource.Resource
	CreateMachinesInput  []machines.CreateMachinesInput
	CreateMachinesOutput []machines.CreateMachinesOutput
	ParentGroupID        *simulations.GroupID
	UpstreamName         string
	GazeboServerIP       string
	WebsocketConnection  ignws.PubSubWebsocketTransporter
	CommsBridgeIPs       []string
	MappingServerIP      string
}

// Platform returns the underlying platform.
func (s *StartSimulation) Platform() platform.Platform {
	return s.platform
}

// Services returns the underlying application services.
func (s *StartSimulation) Services() application.Services {
	return s.services
}

// SubTServices returns the subt specific application services.
func (s *StartSimulation) SubTServices() subtapp.Services {
	return s.services
}

// NewStartSimulation initializes a new state for starting simulations.
func NewStartSimulation(platform platform.Platform, services subtapp.Services, groupID simulations.GroupID) *StartSimulation {
	return &StartSimulation{
		platform: platform,
		services: services,
		GroupID:  groupID,
	}
}
