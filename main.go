package main

import "start-limit-scheduler/limiter"
import (
    "os"

    "k8s.io/kubernetes/cmd/kube-scheduler/app"
)


func main() {
    command := app.NewSchedulerCommand(
        app.WithPlugin(limiter.Name, limiter.New),
    )

    if err := command.Execute(); err != nil {
        os.Exit(1)
    }
}
