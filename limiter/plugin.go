package limiter

import (
    "context"
    "fmt"

    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/runtime"
    framework "k8s.io/kubernetes/pkg/scheduler/framework"
    clientset "k8s.io/client-go/kubernetes"
)

const (
    Name = "PodStartupLimiter"
)

type PodStartupLimiter struct {
    handle          framework.Handle
    client          clientset.Interface
    maxStartingPods int
}

type Args struct {
    MaxStartingPods int `json:"maxStartingPods"`
}

var _ framework.FilterPlugin = &PodStartupLimiter{}

func (pl *PodStartupLimiter) Name() string {
    return Name
}

func New(obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
    args := obj.(*Args)

    return &PodStartupLimiter{
        handle:          handle,
        client:          handle.ClientSet(),
        maxStartingPods: args.MaxStartingPods,
    }, nil
}

func (pl *PodStartupLimiter) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
    pods := nodeInfo.Pods
    starting := 0

    for _, pi := range pods {
        p := pi.Pod
        if p.Status.Phase == v1.PodPending || (p.Status.Phase == v1.PodRunning && !isPodReady(p)) {
            starting++
        }
    }

    if starting >= pl.maxStartingPods {
        return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("node has %d starting pods, max allowed is %d", starting, pl.maxStartingPods))
    }

    return framework.NewStatus(framework.Success, "")
}
