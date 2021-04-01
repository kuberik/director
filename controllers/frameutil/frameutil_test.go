package frameutil

import (
	"reflect"
	"testing"
	"time"

	directorv1 "github.com/kuberik/director/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSetFrameStatus(t *testing.T) {
	oneHourBefore := time.Now().Add(-1 * time.Hour)
	oneHourAfter := time.Now().Add(1 * time.Hour)

	tests := []struct {
		name     string
		statuses []directorv1.FrameStatus
		toAdd    directorv1.FrameStatus
		expected []directorv1.FrameStatus
	}{
		{
			name:     "should-add",
			statuses: []directorv1.FrameStatus{},
			toAdd:    directorv1.FrameStatus{Name: "second", State: directorv1.FrameState{Running: &directorv1.FrameStateRunning{StartedAt: metav1.Time{Time: oneHourBefore}}}},
			expected: []directorv1.FrameStatus{
				{Name: "second", State: directorv1.FrameState{Running: &directorv1.FrameStateRunning{StartedAt: metav1.Time{Time: oneHourBefore}}}},
			},
		},
		{
			name: "update-fields",
			statuses: []directorv1.FrameStatus{
				{Name: "first"},
				{Name: "second", State: directorv1.FrameState{Running: &directorv1.FrameStateRunning{StartedAt: metav1.Time{Time: oneHourBefore}}}},
				{Name: "third"},
			},
			toAdd: directorv1.FrameStatus{Name: "second", State: directorv1.FrameState{Finished: &directorv1.FrameStateFinished{Success: true, FinishedAt: metav1.Time{Time: oneHourAfter}}}},
			expected: []directorv1.FrameStatus{
				{Name: "first"},
				{Name: "second", State: directorv1.FrameState{Finished: &directorv1.FrameStateFinished{Success: true, FinishedAt: metav1.Time{Time: oneHourAfter}}}},
				{Name: "third"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			SetFrameStatus(&test.statuses, test.toAdd)
			if !reflect.DeepEqual(test.statuses, test.expected) {
				t.Error(test.statuses)
			}
		})
	}
}

// func TestRemoveFrameStatus(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		statuses    []directorv1.FrameStatus
// 		conditionType string
// 		expected      []directorv1.FrameStatus
// 	}{
// 		{
// 			name: "present",
// 			statuses: []directorv1.FrameStatus{
// 				{Type: "first"},
// 				{Type: "second"},
// 				{Type: "third"},
// 			},
// 			conditionType: "second",
// 			expected: []directorv1.FrameStatus{
// 				{Type: "first"},
// 				{Type: "third"},
// 			},
// 		},
// 		{
// 			name: "not-present",
// 			statuses: []directorv1.FrameStatus{
// 				{Type: "first"},
// 				{Type: "second"},
// 				{Type: "third"},
// 			},
// 			conditionType: "fourth",
// 			expected: []directorv1.FrameStatus{
// 				{Type: "first"},
// 				{Type: "second"},
// 				{Type: "third"},
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			RemoveFrameStatus(&test.statuses, test.conditionType)
// 			if !reflect.DeepEqual(test.statuses, test.expected) {
// 				t.Error(test.statuses)
// 			}
// 		})
// 	}
// }

func TestFindFrameStatus(t *testing.T) {
	tests := []struct {
		name      string
		statuses  []directorv1.FrameStatus
		frameName string
		expected  *directorv1.FrameStatus
	}{
		{
			name: "not-present",
			statuses: []directorv1.FrameStatus{
				{Name: "first"},
			},
			frameName: "second",
			expected:  nil,
		},
		{
			name: "present",
			statuses: []directorv1.FrameStatus{
				{Name: "first"},
				{Name: "second"},
			},
			frameName: "second",
			expected:  &directorv1.FrameStatus{Name: "second"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := FindFrameStatus(test.statuses, test.frameName)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Error(actual)
			}
		})
	}
}

// func TestIsFrameStatusPresentAndEqual(t *testing.T) {
// 	tests := []struct {
// 		name            string
// 		statuses      []directorv1.FrameStatus
// 		conditionType   string
// 		conditionStatus directorv1.FrameStatusStatus
// 		expected        bool
// 	}{
// 		{
// 			name: "doesnt-match-true",
// 			statuses: []directorv1.FrameStatus{
// 				{Type: "first", Status: directorv1.FrameStatusUnknown},
// 			},
// 			conditionType:   "first",
// 			conditionStatus: directorv1.FrameStatusTrue,
// 			expected:        false,
// 		},
// 		{
// 			name: "does-match-true",
// 			statuses: []directorv1.FrameStatus{
// 				{Type: "first", Status: directorv1.FrameStatusTrue},
// 			},
// 			conditionType:   "first",
// 			conditionStatus: directorv1.FrameStatusTrue,
// 			expected:        true,
// 		},
// 		{
// 			name: "does-match-false",
// 			statuses: []directorv1.FrameStatus{
// 				{Type: "first", Status: directorv1.FrameStatusFalse},
// 			},
// 			conditionType:   "first",
// 			conditionStatus: directorv1.FrameStatusFalse,
// 			expected:        true,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			actual := IsFrameStatusPresentAndEqual(test.statuses, test.conditionType, test.conditionStatus)
// 			if actual != test.expected {
// 				t.Error(actual)
// 			}

// 		})
// 	}
// }
