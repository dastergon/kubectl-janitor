package cmd

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
)

// getAge returns the age of an object.
func getAge(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}
	return duration.HumanDuration(time.Since(timestamp.Time))
}

// getPodStatus returns the current status of a Pod.
func getPodStatus(pod corev1.Pod) string {
	switch pod.Status.Phase {
	case corev1.PodSucceeded:
		for _, status := range pod.Status.ContainerStatuses {
			if status.State.Terminated != nil {
				return string(status.State.Terminated.Reason)
			}
		}
	case corev1.PodFailed:
		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.PodInitialized && condition.Status == corev1.ConditionFalse {
				return "Init:Error"
			}
			for _, status := range pod.Status.ContainerStatuses {
				if status.State.Terminated != nil {
					return string(status.State.Terminated.Reason)
				}

			}
		}
	case corev1.PodRunning:
		for _, status := range pod.Status.ContainerStatuses {
			if status.State.Waiting != nil {
				return string(status.State.Waiting.Reason)
			}
		}
	case corev1.PodPending:
		for _, status := range pod.Status.ContainerStatuses {
			if status.State.Waiting != nil {
				return string(status.State.Waiting.Reason)
			}
		}
	default:
		if pod.DeletionTimestamp != nil && !pod.DeletionTimestamp.IsZero() {
			return "Terminating"
		}
	}

	return string(pod.Status.Phase)
}

// isPodWaitingContainers checks whether one of the containers
//  in the Pod are waiting for an operation.
func isPodWaitingContainers(pod corev1.Pod) bool {
	for _, st := range pod.Status.ContainerStatuses {
		if st.State.Waiting != nil {
			return true
		}
	}
	return false
}

// isPodHealthy checks whether the Pod is in a healthy state.
func isPodHealthy(pod corev1.Pod) bool {
	switch pod.Status.Phase {
	case corev1.PodSucceeded:
		for _, status := range pod.Status.ContainerStatuses {
			if status.State.Terminated != nil && status.State.Terminated.ExitCode != 0 {
				return false
			}
		}
	case corev1.PodPending:
		if isPodWaitingContainers(pod) {
			return false
		}
	case corev1.PodRunning:
		for _, condition := range pod.Status.Conditions {
			if condition.Status == corev1.ConditionFalse {
				return false
			}
		}

		if isPodWaitingContainers(pod) {
			return false
		}

	default:
		return false
	}

	return true
}

// isPodReady checks if the Pod is ready.
func isPodReady(pod corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// writeResults consolidates the final output in the out io.Writer.
func writeResults(out io.Writer, headers []string, matrix [][]string, namespace string, noHeader bool) {
	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)
	defer w.Flush()

	if len(matrix) == 0 {
		if namespace == "" {
			fmt.Fprintln(w, "No resources found")
		} else {
			fmt.Fprintf(w, "No resources found in %s namespace\n", namespace)
		}
	} else {
		if !noHeader {
			if namespace == "" {
				headers = append([]string{"NAMESPACE"}, headers...)
			}
			fmt.Fprintln(w, strings.Join(headers, "\t"))
		}

		for _, row := range matrix {
			fmt.Fprintln(w, strings.Join(row, "\t"))
		}
	}
}
