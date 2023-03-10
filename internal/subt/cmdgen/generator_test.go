package cmdgen

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
	"time"
)

func TestGenerateGazebo(t *testing.T) {
	token := "test-token"
	maxConn := 500
	seed := 5678

	world := "cloudsim_sim.ign;worldName:=tunnel_circuit_practice_01"

	fakeRobotA := fake.NewRobot("testA", "X1")
	fakeRobotB := fake.NewRobot("testB", "X2")

	duration := 1500 * time.Second

	cmd := Gazebo(GazeboConfig{
		World:                   world,
		WorldMaxSimSeconds:      duration,
		Seed:                    &seed,
		AuthorizationToken:      &token,
		MaxWebsocketConnections: maxConn,
		Robots:                  []simulations.Robot{fakeRobotA, fakeRobotB},
		Marsupials: []simulations.Marsupial{
			simulations.NewMarsupial(fakeRobotA, fakeRobotB),
		},
		RosEnabled: true,
	})

	assert.Equal(t, "cloudsim_sim.ign", cmd[0])
	assert.Equal(t, "worldName:=tunnel_circuit_practice_01", cmd[1])
	assert.Equal(t, fmt.Sprintf("durationSec:=%d", int(duration.Seconds())), cmd[2])
	assert.Equal(t, "headless:=true", cmd[3])
	assert.Equal(t, fmt.Sprintf("seed:=%d", seed), cmd[4])
	assert.Equal(t, fmt.Sprintf("websocketAuthKey:=%s", token), cmd[5])
	assert.Equal(t, fmt.Sprintf("websocketAdminAuthKey:=%s", token), cmd[6])
	assert.Equal(t, fmt.Sprintf("websocketMaxConnections:=%d", maxConn), cmd[7])

	assert.Equal(t, fmt.Sprintf("robotName1:=%s", fakeRobotA.GetName()), cmd[8])
	assert.Equal(t, fmt.Sprintf("robotConfig1:=%s", fakeRobotA.GetKind()), cmd[9])

	assert.Equal(t, fmt.Sprintf("robotName2:=%s", fakeRobotB.GetName()), cmd[10])
	assert.Equal(t, fmt.Sprintf("robotConfig2:=%s", fakeRobotB.GetKind()), cmd[11])

	assert.Equal(t, fmt.Sprintf("marsupial1:=%s:%s", fakeRobotA.GetName(), fakeRobotB.GetName()), cmd[12])

	assert.Equal(t, fmt.Sprintf("ros:=%t", true), cmd[13])
}

func TestGenerateCommsBridge(t *testing.T) {
	const (
		firstWorld  = "cloudsim_sim.ign;worldName:=tunnel_circuit_01;circuit:=tunnel"
		secondWorld = "cloudsim_sim.ign;worldName:=tunnel_circuit_02;circuit:=tunnel"
		thirdWorld  = "cloudsim_sim.ign;worldName:=tunnel_circuit_03;circuit:=tunnel"
	)

	cmd, err := CommsBridge(CommsBridgeConfig{
		World:          firstWorld,
		RobotNumber:    0,
		Robot:          fake.NewRobot("X1", "X1_CONFIG_A"),
		ChildMarsupial: true,
	})
	assert.IsType(t, []string{}, cmd)
	assert.NotNil(t, cmd)
	assert.Nil(t, err)
	assert.NotEmpty(t, cmd[0])
	assert.Equal(t, "worldName:=tunnel_circuit_01", cmd[0])
	assert.Equal(t, "robotName1:=X1", cmd[1])
	assert.Equal(t, "robotConfig1:=X1_CONFIG_A", cmd[2])
	assert.Equal(t, "headless:=true", cmd[3])
	assert.Equal(t, "marsupial:=true", cmd[4])

	cmd, err = CommsBridge(CommsBridgeConfig{
		World: secondWorld,
		Robot: fake.NewRobot("X1", "X1_CONFIG_A"),
	})
	assert.Equal(t, "worldName:=tunnel_circuit_02", cmd[0])

	cmd, err = CommsBridge(CommsBridgeConfig{
		World: thirdWorld,
		Robot: fake.NewRobot("X1", "X1_CONFIG_A"),
	})
	assert.Equal(t, "worldName:=tunnel_circuit_03", cmd[0])

	cmd, err = CommsBridge(CommsBridgeConfig{
		World: "",
	})
	assert.Equal(t, ErrEmptyWorld, err)

}

func TestGenerateMapAnalysis(t *testing.T) {
	const (
		world = "cloudsim_sim.ign;worldName:=tunnel_circuit_01;circuit:=tunnel"
	)

	// It should return error if the world is empty
	_, err := MapAnalysis(MapAnalysisConfig{
		World: "",
	})
	assert.Error(t, err)

	// It should return error if the robot slice is empty.
	_, err = MapAnalysis(MapAnalysisConfig{
		World:  world,
		Robots: nil,
	})
	assert.Error(t, err)

	cmd, err := MapAnalysis(MapAnalysisConfig{
		World: world,
		Robots: []simulations.Robot{
			fake.NewRobot("X1", "X1_CONFIG_A"),
			fake.NewRobot("X2", "X2_CONFIG_A"),
		},
	})
	require.NoError(t, err)

	assert.Len(t, cmd, 4)
	assert.Equal(t, "pcd:=tunnel_circuit_01.pcd", cmd[0])
	assert.Equal(t, "gt:=tunnel_circuit_01.csv", cmd[1])
	assert.Equal(t, "robot:=X1", cmd[2])
	assert.Equal(t, "robot:=X2", cmd[3])
}
