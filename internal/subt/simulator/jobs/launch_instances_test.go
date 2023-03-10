package jobs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	simfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	envfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	subtapp "gitlab.com/ignitionrobotics/web/subt/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/subt/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/subt/internal/subt/tracks"
	"testing"
	"time"
)

func TestLaunchInstances(t *testing.T) {
	// Initialize database
	db, err := gorm.GetTestDBFromEnvVars()
	defer db.Close()

	// If the database fails to connect, fail instantly.
	require.NoError(t, err)

	// Migrate db for actions
	actions.CleanAndMigrateDB(db)

	// Initialize simulation
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := simfake.NewSimulation(gid, simulations.StatusPending, simulations.SimSingle, nil, "test", 1*time.Minute, nil, nil)

	// Initialize fake simulation service
	svc := simfake.NewService()
	app := application.NewServices(svc, nil, nil)

	svc.On("Get", gid).Return(sim, error(nil)).Once()

	svc.On("GetRobots", gid).Return(
		[]simulations.Robot{
			simfake.NewRobot("TEST-X1", "X1"),
		},
		error(nil),
	).Once()

	// Configure machine config fake env store
	machineConfigStore := envfake.NewFakeMachines()

	configStore := envfake.NewFakeStore(machineConfigStore, nil, nil, nil)

	machineConfigStore.On("InstanceProfile").Return("arn::test::1234")
	machineConfigStore.On("KeyName").Return("testKey")
	machineConfigStore.On("Type").Return("g3.4xlarge")
	machineConfigStore.On("SidecarType").Return("c5.4xlarge")
	machineConfigStore.On("BaseImage").Return("osrf/test-image")
	machineConfigStore.On("BaseImageGPU").Return("osrf/test-image")
	machineConfigStore.On("FirewallRules").Return([]string{"sg-12345"})
	machineConfigStore.On("NamePrefix").Return("sim")
	machineConfigStore.On("ClusterName").Return("cloudsim-cluster")

	machineConfigStore.On("Tags", sim, "gzserver", "gzserver").Return([]machines.Tag{
		{
			Resource: "instance",
			Map: map[string]string{
				"app": "test",
			},
		},
	})

	machineConfigStore.On("Tags", sim, "field-computer", "fc-TEST-X1").Return([]machines.Tag{
		{
			Resource: "instance",
			Map: map[string]string{
				"app": "test",
			},
		},
	})

	machineConfigStore.On("SubnetAndZone").Return("subnet-test", "zone-test")
	machineConfigStore.On("InitScript").Return("bash")
	machineConfigStore.On("Limit").Return(-1)

	// Configure mocked machines interface
	machines := &instancesLauncher{}

	// Initialize platform
	p, _ := platform.NewPlatform("test", platform.Components{
		Machines: machines,
		Store:    configStore,
	})

	tracksService := tracks.NewService(nil, nil, nil)

	subt := subtapp.NewServices(app, tracksService, nil)

	// Create initial state
	initialState := state.NewStartSimulation(p, subt, gid)

	// Pass the initial state to the action store
	s := actions.NewStore(&initialState)

	// Run the job
	_, err = LaunchInstances.Run(s, db, &actions.Deployment{}, initialState)

	// Check there are no errors and that the machines component has been called once.
	assert.NoError(t, err)
	assert.Equal(t, 1, machines.TimesCalled)
}

type instancesLauncher struct {
	TimesCalled int
	machines.Machines
}

// Create mocks the create method of the machines.Machines interface.
func (i *instancesLauncher) Create(input []machines.CreateMachinesInput) ([]machines.CreateMachinesOutput, error) {
	i.TimesCalled++
	var output []machines.CreateMachinesOutput
	for _, in := range input {
		var out machines.CreateMachinesOutput
		for i := 0; i < int(in.MaxCount); i++ {
			out.Instances = append(out.Instances, fmt.Sprintf("%s-%d", in.Type, i))
		}
		output = append(output, out)
	}
	return output, nil
}

func (i *instancesLauncher) Count(input machines.CountMachinesInput) int {
	return 0
}
