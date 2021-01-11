package autoheal

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

//AutoHealer spins to find unhealty containers to restart
type AutoHealer struct {
	containerLabel     string
	startPeriod        int
	interval           int
	defaultStopTimeOut int
	dockerSock         string

	cli *client.Client
	ctx context.Context
}

//NewAutoHealer creates new AutoHealer
func NewAutoHealer() *AutoHealer {
	ah := &AutoHealer{}
	ah.readEnvs()
	ah.ctx = context.Background()
	var err error
	ah.cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return ah
}

func (ah *AutoHealer) readEnvs() {
	var err error
	ah.containerLabel = os.Getenv(autoHealContainerLabelEnv)
	ah.startPeriod, err = strconv.Atoi(getEnv(autoHealStartPeriodEnv, "0"))
	checkErrorForIntegerEnv(err, autoHealStartPeriodEnv)

	ah.interval, err = strconv.Atoi(getEnv(autoHealIntervalEnv, "5"))
	checkErrorForIntegerEnv(err, autoHealIntervalEnv)

	ah.defaultStopTimeOut, err = strconv.Atoi(getEnv(autoHealdefaultStopTimeOutEnv, "0"))
	checkErrorForIntegerEnv(err, autoHealIntervalEnv)

	ah.dockerSock = getEnv(autoHealdockerSockEnv, "/var/run/docker.sock")
}

func (ah *AutoHealer) resetUnhealthy(container *types.Container) {
	duration := time.Second * time.Duration(ah.defaultStopTimeOut)
	ah.cli.ContainerRestart(ah.ctx, container.ID, &duration)
}

func (ah *AutoHealer) spinOnce() {
	args := filters.NewArgs(
		filters.Arg("health", "unhealthy"),
	)
	containerListOptions := types.ContainerListOptions{
		Filters: args,
	}

	containers, err := ah.cli.ContainerList(ah.ctx, containerListOptions)

	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		if ah.containerLabel == "" || valueInList(ah.containerLabel, container.Labels) {
			ah.resetUnhealthy(&container)
		}

	}
}

//Spin Endless loop to list unhealty containers to restart
func (ah *AutoHealer) Spin() {
	for {
		ah.spinOnce()
		time.Sleep(time.Second * time.Duration(ah.interval))
	}
}
