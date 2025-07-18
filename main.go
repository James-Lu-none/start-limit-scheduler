package main

import (
    "os"

    "k8s.io/kubernetes/cmd/kube-scheduler/app"
    "github.com/yourrepo/scheduler/limiter"
    "k8s.io/component-base/configz"
    "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
)

func main() {
    command := app.NewSchedulerCommand(
        app.WithPlugin(limiter.Name, limiter.New),
    )

    if err := command.Execute(); err != nil {
        os.Exit(1)
    }
}
