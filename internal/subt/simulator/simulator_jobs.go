package simulator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/subt/internal/subt/simulator/jobs"
)

// JobsStartSimulation groups the jobs needed to start a simulation.
var JobsStartSimulation = actions.Jobs{
	jobs.CheckSimulationPendingStatus,
	jobs.CheckStartSimulationIsNotParent,
	jobs.CheckSimulationNoErrors,
	jobs.SetSimulationStatusToLaunchInstances,
	jobs.LaunchInstances,
	jobs.SetRate,
	jobs.SetSimulationStatusToWaitInstances,
	jobs.WaitForInstances,
	jobs.SetSimulationStatusToWaitNodes,
	jobs.WaitForNodes,
	jobs.SetSimulationStatusToLaunchPods,
	jobs.CreateNetworkPolicyGazeboServer,
	jobs.LaunchGazeboServerPod,
	jobs.WaitForGazeboServerPod,
	jobs.GetGazeboIP,
	jobs.LaunchMoleBridgePod,
	jobs.LaunchMappingServerPod,
	jobs.LaunchGazeboServerCopyPod,
	jobs.LaunchMappingServerCopyPod,
	jobs.WaitForMappingServerPod,
	jobs.GetMappingServerIP,
	jobs.LaunchWebsocketService,
	jobs.WaitUpstream,
	jobs.ConfigureIngressGloo,
	jobs.CreateNetworkPolicyCommsBridges,
	jobs.CreateNetworkPolicyFieldComputers,
	jobs.CreateNetworkPolicyMappingServer,
	jobs.LaunchCommsBridgePods,
	jobs.LaunchCommsBridgeCopyPods,
	jobs.WaitForCommsBridgePodIPs,
	jobs.GetCommsBridgePodIP,
	jobs.WaitForCommsBridgePodsReady,
	jobs.LaunchFieldComputerPods,
	jobs.SetSimulationStatusToWaitPods,
	jobs.WaitSimulationPods,
	jobs.SetWebsocketConnection,
	jobs.AddRunningSimulation,
	jobs.SetSimulationStatusToRunning,
}

// JobsStopSimulation groups the jobs needed to stop a simulation.
var JobsStopSimulation = actions.Jobs{
	jobs.CheckSimulationTerminateRequestedStatus,
	jobs.CheckStopSimulationIsNotParent,
	jobs.SetStoppedAt,
	jobs.ChargeUser,
	jobs.DisconnectWebsocket,
	jobs.RemoveRunningSimulation,
	jobs.SetSimulationStatusToProcessingResults,
	jobs.UploadLogs,
	jobs.ReadScore,
	jobs.ReadStats,
	jobs.ReadRunData,
	jobs.SaveScore,
	jobs.SaveSummary,
	jobs.SendSummaryEmail,
	jobs.SetSimulationStatusToDeletingPods,
	jobs.RemoveIngressRulesGloo,
	jobs.RemoveWebsocketService,
	jobs.RemoveNetworkPolicies,
	jobs.RemovePods,
	jobs.SetSimulationStatusToDeletingNodes,
	jobs.RemoveInstances,
	jobs.SetSimulationStatusToTerminated,
}
