package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/maargenton/go-cli"
)

type dummyLoadCmd struct {
	WorkerCount int           `yaml:"workerCount"  opts:"-j,--job, default: 1"                            desc:"number of concurrent tasks to run"`
	WorkPeriod  time.Duration `yaml:"workPeriod"   opts:"-w,--work-period, default: 5m, name: duration"  desc:"duration of work cycle"`
	SleepPeriod time.Duration `yaml:"sleepPeriod"  opts:"-s,--sleep-period, default: 5m, name: duration" desc:"duration of sleep between work cycles"`

	ServicePort int `yaml:"servicePort"  opts:"--service-port, default: 8080, env: SERVICE_PORT, name: port"  desc:"port number the main service endpoint"`
	MetricsPort int `yaml:"metricsPort"  opts:"--metrics-port, default: 8081, env: METRICS_PORT, name: port"  desc:"port number the service metrics and monitoring endpoint"`
	// Actions     []string `yaml:"actions" opts:"--actions, delim:\\,, default:foo\\,bar\\,foobar"`
}

func (options *dummyLoadCmd) Run() error {
	d, err := json.Marshal(options)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", string(d))
	return nil
}

func main() {
	cli.Run(&cli.Command{
		Handler:     &dummyLoadCmd{},
		Description: "Perform CPU intensive work with tunable duty cycle",
	})
}
