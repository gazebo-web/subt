package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/subt/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/subt/internal/subt/simulator/state"
)

// GetCommsBridgePodIP is a job in charge of getting the IP from the simulation's comms bridge pods.
// WaitForCommsBridgePodIPs should be run before running this job.
var GetCommsBridgePodIP = &actions.Job{
	Name:       "get-comms-bridge-pod-ip",
	PreHooks:   []actions.JobFunc{setStartState},
	Execute:    getCommsBridgeIP,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
}

// getCommsBridgeIP gets all coms bridge server pod IPs and assigns them to the start simulation state.
func getCommsBridgeIP(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	ips := make([]string, len(robots))
	for i := range robots {
		robotID := application.GetRobotID(i)

		name := application.GetPodNameCommsBridge(s.GroupID, robotID)
		ns := s.Platform().Store().Orchestrator().Namespace()

		ip, err := s.Platform().Orchestrator().Pods().GetIP(name, ns)
		if err != nil {
			return nil, err
		}

		ips[i] = ip
	}

	s.CommsBridgeIPs = ips
	store.SetState(s)

	return s, nil
}
