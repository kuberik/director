package frameutil

import (
	directorv1 "github.com/kuberik/director/api/v1"
)

// FindFrameStatus finds a frame status in statuses.
func FindFrameStatus(statuses []directorv1.FrameStatus, frameName string) *directorv1.FrameStatus {
	for i := range statuses {
		if statuses[i].Name == frameName {
			return &statuses[i]
		}
	}

	return nil
}

// SetFrameStatus sets the corresponding status in statuses to newStatus.
// statuses must be non-nil.
// 1. if the status of the specified type already exists (all fields of the existing status are updated to
//    newStatus, LastTransitionTime is set to now if the new status differs from the old status)
// 2. if a status of the specified type does not exist (LastTransitionTime is set to now() if unset, and newStatus is appended)
func SetFrameStatus(statuses *[]directorv1.FrameStatus, newStatus directorv1.FrameStatus) {
	if statuses == nil {
		return
	}
	existingStatus := FindFrameStatus(*statuses, newStatus.Name)
	if existingStatus == nil {
		*statuses = append(*statuses, newStatus)
		return
	}

	existingStatus.State = newStatus.State
}
