package limiter

import (
    "context"
    "fmt"

    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/kubernetes/pkg/scheduler/framework"
    clientset "k8s.io/client-go/kubernetes"
)

const (
    Name = "PodStartupLimiter"
)

type PodStartupLimiter struct {
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

func New(ctx context.Context, obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
    return &PodStartupLimiter{
        client: handle.ClientSet(),
        maxStartingPods: 3,
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

func isPodReady(pod *v1.Pod) bool {
    for _, cond := range pod.Status.Conditions {
        if cond.Type == v1.PodReady && cond.Status == v1.ConditionTrue {
            return true
        }
    }
    return false
}
