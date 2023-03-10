package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	subtapp "gitlab.com/ignitionrobotics/web/subt/internal/subt/application"
	subt "gitlab.com/ignitionrobotics/web/subt/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/subt/internal/subt/simulator/state"
)

// LaunchFieldComputerPods launches the list of field computer pods.
var LaunchFieldComputerPods = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-field-computer-pods",
	PreHooks:        []actions.JobFunc{setStartState, prepareFieldComputerPodInput},
	PostHooks:       []actions.JobFunc{checkLaunchPodsError, returnState},
	RollbackHandler: rollbackLaunchFieldComputerPods,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

func rollbackLaunchFieldComputerPods(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	robots, err := s.Services().Simulations().GetRobots(s.GroupID)
	if err != nil {
		return nil, err
	}

	for i := range robots {
		name := subtapp.GetPodNameFieldComputer(s.GroupID, subtapp.GetRobotID(i))
		ns := s.Platform().Store().Orchestrator().Namespace()

		_, _ = s.Platform().Orchestrator().Pods().Delete(resource.NewResource(name, ns, nil))
	}

	return nil, nil
}

// prepareFieldComputerPodInput prepares the input for the generic LaunchPods job to launch field computer pods.
func prepareFieldComputerPodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	sim, err := s.Services().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	subtSim := sim.(subt.Simulation)

	podInputs := make([]pods.CreatePodInput, len(subtSim.GetRobots()))

	for i, r := range subtSim.GetRobots() {
		robotID := subtapp.GetRobotID(i)
		// Create field computer input
		allowPrivilegesEscalation := false
		podInputs[i] = pods.CreatePodInput{
			Name:                          subtapp.GetPodNameFieldComputer(s.GroupID, robotID),
			Namespace:                     s.Platform().Store().Orchestrator().Namespace(),
			Labels:                        subtapp.GetPodLabelsFieldComputer(s.GroupID, s.ParentGroupID).Map(),
			RestartPolicy:                 pods.RestartPolicyNever,
			TerminationGracePeriodSeconds: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
			NodeSelector:                  subtapp.GetNodeLabelsFieldComputer(s.GroupID, r),
			Containers: []pods.Container{
				{
					Name:                     subtapp.GetContainerNameFieldComputer(),
					Image:                    r.GetImage(),
					AllowPrivilegeEscalation: &allowPrivilegesEscalation,
					EnvVars:                  subtapp.GetEnvVarsFieldComputer(r.GetName(), s.CommsBridgeIPs[i]),
					EnvVarsFrom:              subtapp.GetEnvVarsFromSourceFieldComputer(),
					ResourceLimits: map[pods.ResourceName]string{
						pods.ResourceMemory: "114Gi",
					},
				},
			},
		}
	}

	return jobs.LaunchPodsInput(podInputs), nil
}
