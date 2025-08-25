package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

func (r *ApplicationPool) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ApplicationPool struct {
	Name                  string              `json:"name"`
	ID                    string              `json:"id"`
	Status                string              `json:"status"`
	AutoStart             bool                `json:"auto_start"`
	PipelineMode          string              `json:"pipeline_mode"`
	ManagedRuntimeVersion string              `json:"managed_runtime_version"`
	Enable32BitWin64      bool                `json:"enable_32bit_win64"`
	QueueLength           int64               `json:"queue_length"`
	CPU                   CPU                 `json:"cpu"`
	ProcessModel          ProcessModel        `json:"process_model"`
	Identity              Identity            `json:"identity"`
	Recycling             Recycling           `json:"recycling"`
	RapidFailProtection   RapidFailProtection `json:"rapid_fail_protection"`
	ProcessOrphaning      ProcessOrphaning    `json:"process_orphaning"`
}

type CPU struct {
	Limit int64 `json:"limit"`
	//LimitInterval            int64  `json:"limit_interval"`
	Action                   string `json:"action"`
	ProcessorAffinityEnabled bool   `json:"processor_affinity_enabled"`
	ProcessorAffinityMask32  string `json:"processor_affinity_mask32"`
	ProcessorAffinityMask64  string `json:"processor_affinity_mask64"`
}

type Identity struct {
	IdentityType    string `json:"identity_type"`
	Username        string `json:"username"`
	LoadUserProfile bool   `json:"load_user_profile"`
}

type ProcessModel struct {
	//IdleTimeout       int64  `json:"idle_timeout"`
	MaxProcesses   int64 `json:"max_processes"`
	PingingEnabled bool  `json:"pinging_enabled"`
	//PingInterval      int64  `json:"ping_interval"`
	//PingResponseTime  int64  `json:"ping_response_time"`
	//ShutdownTimeLimit int64  `json:"shutdown_time_limit"`
	//StartupTimeLimit  int64  `json:"startup_time_limit"`
	IdleTimeoutAction string `json:"idle_timeout_action"`
}

type ProcessOrphaning struct {
	Enabled            bool   `json:"enabled"`
	OrphanActionExe    string `json:"orphan_action_exe"`
	OrphanActionParams string `json:"orphan_action_params"`
}

type RapidFailProtection struct {
	Enabled                  bool   `json:"enabled"`
	LoadBalancerCapabilities string `json:"load_balancer_capabilities"`
	//Interval                 int64  `json:"interval"`
	MaxCrashes         int64  `json:"max_crashes"`
	AutoShutdownExe    string `json:"auto_shutdown_exe"`
	AutoShutdownParams string `json:"auto_shutdown_params"`
}

type Recycling struct {
	DisableOverlappedRecycle     bool            `json:"disable_overlapped_recycle"`
	DisableRecycleOnConfigChange bool            `json:"disable_recycle_on_config_change"`
	LogEvents                    LogEvents       `json:"log_events"`
	PeriodicRestart              PeriodicRestart `json:"periodic_restart"`
}

type LogEvents struct {
	Time           bool `json:"time"`
	Requests       bool `json:"requests"`
	Schedule       bool `json:"schedule"`
	Memory         bool `json:"memory"`
	IsapiUnhealthy bool `json:"isapi_unhealthy"`
	OnDemand       bool `json:"on_demand"`
	ConfigChange   bool `json:"config_change"`
	PrivateMemory  bool `json:"private_memory"`
}

type PeriodicRestart struct {
	//TimeInterval  int64         `json:"time_interval"`
	PrivateMemory int64         `json:"private_memory"`
	RequestLimit  int64         `json:"request_limit"`
	VirtualMemory int64         `json:"virtual_memory"`
	Schedule      []interface{} `json:"schedule"`
}

func (client Client) ReadAppPool(ctx context.Context, id string) (*ApplicationPool, error) {
	url := fmt.Sprintf("/api/webserver/application-pools/%s", id)
	var appPool ApplicationPool
	if err := getJson(ctx, client, url, &appPool); err != nil {
		return nil, err
	}
	return &appPool, nil
}

func (client Client) DeleteAppPool(ctx context.Context, id string) error {
	url := fmt.Sprintf("/api/webserver/application-pools/%s", id)
	return httpDelete(ctx, client, url)
}
