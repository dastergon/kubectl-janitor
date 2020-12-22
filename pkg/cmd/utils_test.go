package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsPodHealthy(t *testing.T) {
	ts := time.Now()

	tests := []struct {
		name string
		pod  corev1.Pod
		want bool
	}{
		{
			name: "Pod waiting for containers is expected to marked as unhealthy",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodPending,
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.PodScheduled,
							Status:             corev1.ConditionFalse,
							LastTransitionTime: metav1.NewTime(ts.Add(-time.Minute * 2)),
						},
					},
					ContainerStatuses: []corev1.ContainerStatus{{
						State: corev1.ContainerState{
							Waiting: &corev1.ContainerStateWaiting{
								Reason: "CreateContainerConfigError",
							},
						},
					}},
				},
			},
			want: false,
		},
		{
			name: "Pod with a Failed status is expected to be unhealthy ",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodFailed,
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.PodReady,
							Status:             corev1.ConditionFalse,
							LastTransitionTime: metav1.NewTime(ts.Add(-time.Minute * 2)),
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Pod with a Succeeded status and a successful exit code should be considered healthy",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodSucceeded,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name: "test-env",
							State: corev1.ContainerState{
								Terminated: &corev1.ContainerStateTerminated{
									Reason:   "aReason",
									ExitCode: 0,
								},
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "Pod with a Succeeded status but exit code 1 should be considered as unhealthy",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodSucceeded,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name: "test-env",
							State: corev1.ContainerState{
								Terminated: &corev1.ContainerStateTerminated{
									ExitCode: 1,
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Pod with an Uknown status is expected to be unhealthy",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodUnknown,
				},
			},
			want: false,
		},
		{
			name: "Pod in a Running status but waiting for containers is considered unhealthy",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name: "test-env",
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason: "ImagePullBackOff",
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Pod in with Running status but a False condition (i.e., Not Ready) is considered unhealthy",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.PodReady,
							Status:             corev1.ConditionFalse,
							LastTransitionTime: metav1.NewTime(ts.Add(-time.Minute * 2)),
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := isPodHealthy(tc.pod)
			assert.Equal(t, tc.want, got)
		})
	}
}
func TestGetPodStatus(t *testing.T) {
	ts := time.Now()

	tests := []struct {
		name string
		pod  corev1.Pod
		want string
	}{
		{
			name: "Pod should have a Pending status",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodPending,
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.PodScheduled,
							Status:             corev1.ConditionFalse,
							LastTransitionTime: metav1.NewTime(ts.Add(-time.Minute * 2)),
						},
					},
				},
			},
			want: "Pending",
		},
		{
			name: "Pod should have a CreateContainerConfigError status",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodPending,
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.PodScheduled,
							Status:             corev1.ConditionFalse,
							LastTransitionTime: metav1.NewTime(ts.Add(-time.Minute * 2)),
						},
					},
					ContainerStatuses: []corev1.ContainerStatus{{
						State: corev1.ContainerState{
							Waiting: &corev1.ContainerStateWaiting{
								Reason: "CreateContainerConfigError",
							},
						},
					}},
				},
			},
			want: "CreateContainerConfigError",
		},
		{
			name: "Pod should have Failed and reported aReason",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodFailed,
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.PodScheduled,
							Status:             corev1.ConditionFalse,
							LastTransitionTime: metav1.NewTime(ts.Add(-time.Minute * 2)),
						},
					},
					ContainerStatuses: []corev1.ContainerStatus{{
						State: corev1.ContainerState{
							Terminated: &corev1.ContainerStateTerminated{
								Reason:   "aReason",
								ExitCode: 1,
							},
						},
					}},
				},
			},
			want: "aReason",
		},
		{
			name: "Pod should have a Succeeded status",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodSucceeded,
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.PodReady,
							Status:             corev1.ConditionFalse,
							LastTransitionTime: metav1.NewTime(ts.Add(-time.Minute * 2)),
						},
					},
				},
			},
			want: "Succeeded",
		},
		{
			name: "Pod should have a Failed status",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodFailed,
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.PodReady,
							Status:             corev1.ConditionFalse,
							LastTransitionTime: metav1.NewTime(ts.Add(-time.Minute * 2)),
						},
					},
				},
			},
			want: "Failed",
		},
		{
			name: "Pod should have a Running status",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.PodReady,
							Status:             corev1.ConditionTrue,
							LastTransitionTime: metav1.NewTime(ts.Add(-time.Minute * 2)),
						},
					},
				},
			},
			want: "Running",
		},
		{
			name: "Pod should be having a Terminating status",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
					DeletionTimestamp: &metav1.Time{
						Time: ts,
					},
				},
			},
			want: "Terminating",
		},
		{
			name: "Pod should have an Init:Error status",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodFailed,
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.PodInitialized,
							Status:             corev1.ConditionFalse,
							LastTransitionTime: metav1.NewTime(ts.Add(-time.Minute * 2)),
						},
					},
				},
			},
			want: "Init:Error",
		},
		{
			name: "Pod should have terminated successfully a reason",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodSucceeded,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name: "test-env",
							State: corev1.ContainerState{
								Terminated: &corev1.ContainerStateTerminated{
									Reason:   "aReason",
									ExitCode: 0,
								},
							},
						},
					},
				},
			},
			want: "aReason",
		},
		{
			name: "Pod should have a ImagePullBackOff status",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name: "test-env",
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason: "ImagePullBackOff",
								},
							},
						},
					},
				},
			},
			want: "ImagePullBackOff",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := getPodStatus(tc.pod)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIsPodWaitingContainers(t *testing.T) {
	tests := []struct {
		name string
		pod  corev1.Pod
		want bool
	}{
		{
			name: "Pod should be waiting containers",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name: "test-env",
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason: "ImagePullBackOff",
								},
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "Pod is expected to run without waiting for containers",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodSucceeded,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name: "test-env",
							State: corev1.ContainerState{
								Terminated: &corev1.ContainerStateTerminated{
									ExitCode: 0,
								},
							},
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := isPodWaitingContainers(tc.pod)
			assert.Equal(t, got, tc.want)
		})
	}
}

func TestIsPodReady(t *testing.T) {
	ts := time.Now()

	tests := []struct {
		name string
		pod  corev1.Pod
		want bool
	}{
		{
			name: "Pod is expected to be ready",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.PodReady,
							Status:             corev1.ConditionTrue,
							LastTransitionTime: metav1.NewTime(ts.Add(-time.Minute * 2)),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "Pod is not expected to be ready",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Pod",
						},
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.PodReady,
							Status:             corev1.ConditionFalse,
							LastTransitionTime: metav1.NewTime(ts.Add(-time.Minute * 2)),
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := isPodReady(tc.pod)
			assert.Equal(t, got, tc.want)
		})
	}
}

func TestPrintResults(t *testing.T) {
	type args struct {
		headers   []string
		matrix    [][]string
		namespace string
		noHeader  bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "expect no headers",
			args: args{
				headers:   []string{"NAME", "STATUS", "AGE"},
				matrix:    [][]string{{"tester", "Running", "8m"}},
				namespace: "default",
				noHeader:  true,
			},
			want: "tester   Running   8m\n",
		},
		{
			name: "expect headers",
			args: args{
				headers:   []string{"NAME", "STATUS", "AGE"},
				matrix:    [][]string{{"tester", "Running", "8m"}},
				namespace: "default",
				noHeader:  false,
			},
			want: "NAME     STATUS    AGE\ntester   Running   8m\n",
		},
		{
			name: "expect headers with Namespace column",
			args: args{
				headers:   []string{"NAME", "STATUS", "AGE"},
				matrix:    [][]string{{"production", "tester", "Running", "8m"}},
				namespace: "",
				noHeader:  false,
			},
			want: "NAMESPACE    NAME     STATUS    AGE\nproduction   tester   Running   8m\n",
		},
		{
			name: "expect no resources found from one namespace",
			args: args{
				headers:   []string{"NAME", "STATUS", "AGE"},
				matrix:    [][]string{},
				namespace: "default",
				noHeader:  false,
			},
			want: "No resources found in default namespace\n",
		},
		{
			name: "expect no resources found from all namespaces",
			args: args{
				headers:   []string{"NAME", "STATUS", "AGE"},
				matrix:    [][]string{},
				namespace: "",
				noHeader:  false,
			},
			want: "No resources found\n",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			out := &bytes.Buffer{}
			writeResults(out, tc.args.headers, tc.args.matrix, tc.args.namespace, tc.args.noHeader)
			assert.Equal(t, tc.want, out.String())
		})
	}
}

func TestGetAge(t *testing.T) {
	//ts := time.Now()

	tests := []struct {
		name      string
		timestamp metav1.Time
		want      string
	}{
		{
			name:      "zero timestamp passed",
			timestamp: metav1.Time{},
			want:      "<unknown>",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := getAge(tc.timestamp)
			assert.Equal(t, tc.want, got)
		})
	}
}
