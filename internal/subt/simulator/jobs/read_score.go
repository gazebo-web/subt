package jobs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	subtapp "gitlab.com/ignitionrobotics/web/subt/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/subt/internal/subt/simulator/state"
	"strconv"
	"strings"
)

// ReadScore is a job in charge of reading the score from a gzserver copy pod for the simulation that is being processed.
var ReadScore = &actions.Job{
	Name:       "read-simulation-score",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    readScore,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// readScore is the main execute function for the ReadScore job.
func readScore(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	path := fmt.Sprintf("%s/gazebo/score.yml", s.Platform().Store().Ignition().SidecarContainerLogsPath())

	body, err := readFileContentFromPod(
		s.Platform().Orchestrator().Pods(),
		subtapp.GetPodNameGazeboServerCopy(s.GroupID),
		s.Platform().Store().Orchestrator().Namespace(),
		path,
	)
	if err != nil {
		return nil, err
	}

	// Parse the score
	score, err := strconv.ParseFloat(strings.TrimSpace(string(body)), 64)
	if err != nil {
		return nil, err
	}

	// Set the score
	s.Score = &score
	store.SetState(s)

	// Return state
	return s, nil
}
