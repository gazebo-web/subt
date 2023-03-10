package jobs

import (
	"context"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport"
	subt "gitlab.com/ignitionrobotics/web/subt/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/subt/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/subt/pkg/runsim"
)

// AddRunningSimulation is job in charge of adding a running simulation to the list of running simulations.
var AddRunningSimulation = &actions.Job{
	Name:            "add-running-simulation",
	PreHooks:        []actions.JobFunc{setStartState},
	Execute:         addRunningSimulation,
	PostHooks:       []actions.JobFunc{returnState},
	RollbackHandler: revertAddingRunningSimulation,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
}

// revertAddingRunningSimulation reverts all the changes made while adding a running simulation.
func revertAddingRunningSimulation(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, _ error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	if s.SubTServices().RunningSimulations().GetTransporter(s.GroupID).IsConnected() {
		s.SubTServices().RunningSimulations().Free(s.GroupID)
	}

	err := s.SubTServices().RunningSimulations().Remove(s.GroupID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// addRunningSimulation is the main function of the AddRunningSimulation job.
func addRunningSimulation(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	// Get the simulation
	sim, err := s.SubTServices().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	// Parse it as a subt simulation
	subtSim := sim.(subt.Simulation)

	// Get the track for the given simulation
	t, err := s.SubTServices().Tracks().Get(subtSim.GetTrack(), subtSim.GetWorldIndex(), subtSim.GetRunIndex())
	if err != nil {
		return nil, err
	}

	// Initialize a new RunningSimulation.
	rs := runsim.NewRunningSimulation(sim)

	err = s.WebsocketConnection.Subscribe(t.WarmupTopic, func(message transport.Message) {
		_ = rs.ReadWarmup(context.Background(), message)
	})
	if err != nil {
		return nil, err
	}

	// Add the running simulation and websocket connection to the Running Simulation manager.
	err = s.SubTServices().RunningSimulations().Add(s.GroupID, rs, s.WebsocketConnection)
	if err != nil {
		return nil, err
	}

	return s, nil
}
