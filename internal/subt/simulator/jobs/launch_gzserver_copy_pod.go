package jobs

import (
	"context"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
	subtapp "gitlab.com/ignitionrobotics/web/subt/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/subt/internal/subt/simulator/state"
)

// LaunchGazeboServerCopyPod launches a gazebo server copy pod.
var LaunchGazeboServerCopyPod = jobs.LaunchPods.Extend(actions.Job{
	Name:            "launch-gzserver-copy-pod",
	PreHooks:        []actions.JobFunc{setStartState, prepareGazeboCreateCopyPodInput},
	PostHooks:       []actions.JobFunc{checkLaunchPodsError, returnState},
	RollbackHandler: rollbackLaunchGazeboServerCopyPod,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
})

func rollbackLaunchGazeboServerCopyPod(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, err error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	name := subtapp.GetPodNameGazeboServerCopy(s.GroupID)
	ns := s.Platform().Store().Orchestrator().Namespace()

	_, _ = s.Platform().Orchestrator().Pods().Delete(resource.NewResource(name, ns, nil))

	return nil, nil
}

// prepareGazeboCreateCopyPodInput prepares the input to launch a copy pod for a gzserver.
func prepareGazeboCreateCopyPodInput(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	if !s.Platform().Store().Ignition().LogsCopyEnabled() {
		return jobs.LaunchPodsInput{}, nil
	}

	// Set up namespace
	namespace := s.Platform().Store().Orchestrator().Namespace()

	// Set up nameservers
	nameservers := s.Platform().Store().Orchestrator().Nameservers()

	// Set up secrets
	secretsName := s.Platform().Store().Ignition().SecretsName()
	secret, err := s.Platform().Secrets().Get(context.TODO(), secretsName, namespace)
	if err != nil {
		return nil, err
	}

	accessKey := string(secret.Data[s.Platform().Store().Ignition().AccessKeyLabel()])
	secretAccessKey := string(secret.Data[s.Platform().Store().Ignition().SecretAccessKeyLabel()])

	volumes := []pods.Volume{
		{
			Name:         "logs",
			HostPath:     "/tmp",
			HostPathType: pods.HostPathDirectoryOrCreate,
			MountPath:    s.Platform().Store().Ignition().SidecarContainerLogsPath(),
			SubPath:      "logs",
		},
	}

	return jobs.LaunchPodsInput{
		{
			Name:                          subtapp.GetPodNameGazeboServerCopy(s.GroupID),
			Namespace:                     namespace,
			Labels:                        subtapp.GetPodLabelsGazeboServerCopy(s.GroupID, s.ParentGroupID).Map(),
			RestartPolicy:                 pods.RestartPolicyAlways,
			TerminationGracePeriodSeconds: s.Platform().Store().Orchestrator().TerminationGracePeriod(),
			NodeSelector:                  subtapp.GetNodeLabelsGazeboServer(s.GroupID),
			Containers: []pods.Container{
				{
					Name:    subtapp.GetContainerNameGazeboServerCopy(),
					Image:   "infrastructureascode/aws-cli:latest",
					Command: []string{"tail", "-f", "/dev/null"},
					Volumes: volumes,
					EnvVars: subtapp.GetEnvVarsGazeboServerCopy(
						s.Platform().Store().Ignition().Region(),
						accessKey,
						secretAccessKey,
					),
				},
			},
			Volumes:     volumes,
			Nameservers: nameservers,
		},
	}, nil
}
