package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	subtapp "gitlab.com/ignitionrobotics/web/subt/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/subt/internal/subt/simulator/state"
)

// WaitForCommsBridgePodIPs waits for the simulation comms bridge pods to have an IP.
var WaitForCommsBridgePodIPs = jobs.Wait.Extend(actions.Job{
	Name:       "wait-comms-bridge-pods-ips",
	PreHooks:   []actions.JobFunc{createWaitRequestForCommsBridgePod},
	PostHooks:  []actions.JobFunc{checkWaitError, returnState},
	InputType:  actions.GetJobDataType(&state.StartSimulation{}),
	OutputType: actions.GetJobDataType(&state.StartSimulation{}),
})

// createWaitRequestForCommsBridgePod is the pre hook in charge of passing the needed input to the Wait job.
func createWaitRequestForCommsBridgePod(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)

	store.SetState(s)

	res := resource.NewResource("", "", subtapp.GetPodLabelsBase(s.GroupID, nil))

	// Create wait for condition request
	// Since only the gazebo server pod has been created and already has an IP, we only need to wait until
	// comms bridge pods have an ip.
	req := s.Platform().Orchestrator().Pods().WaitForCondition(res, resource.HasIPStatusCondition)

	// Get timeout and poll frequency from store
	timeout := s.Platform().Store().Orchestrator().Timeout()
	pollFreq := s.Platform().Store().Orchestrator().PollFrequency()

	// Return new wait input
	return jobs.WaitInput{
		Request:       req,
		PollFrequency: pollFreq,
		Timeout:       timeout,
	}, nil
}
