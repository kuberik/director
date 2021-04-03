/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MovieSpec defines the desired state of Movie
type MovieSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +optional
	Props []Prop `json:"props,omitempty"`
	// +optional
	Screenplay `json:"screenplay,omitempty"`
}

// MovieStatus defines the observed state of Movie
type MovieStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// +optional
	FrameStatuses []FrameStatus `json:"frameStatuses,omitempty"`
}

type Prop struct {
	Name string `json:"name,omitempty"`
	// +optional
	Resources []runtime.RawExtension `json:"resources,omitempty"`
}

type Screenplay struct {
	Frames []Frame `json:"frames,omitempty"`
}

type Frame struct {
	Name   string          `json:"name,omitempty"`
	Action batchv1.JobSpec `json:"action,omitempty"`
	// +optional
	CausedBy []string `json:"causedBy,omitempty"`
	// +optional
	Props []string `json:"props,omitempty"`
}

type FrameStatus struct {
	Name  string     `json:"name,omitempty"`
	State FrameState `json:"state,omitempty"`
}

// FrameState holds a possible state of a frame.
// Only one of its members may be specified.
type FrameState struct {
	// Details about a running frame
	// +optional
	Running *FrameStateRunning `json:"running,omitempty"`
	// Details about a finished frame
	// +optional
	Finished *FrameStateFinished `json:"finished,omitempty"`
}

type FrameStateRunning struct {
	// Time at which execution of the frame started
	// +optional
	StartedAt metav1.Time `json:"startedAt,omitempty"`
}

type FrameStateFinished struct {
	// Indicated if the frame finished successfully
	// +optional
	Success bool `json:"success,omitempty"`
	// Time at which execution of the frame started
	// +optional
	StartedAt metav1.Time `json:"startedAt,omitempty"`
	// Time at which the frame finished
	// +optional
	FinishedAt metav1.Time `json:"finishedAt,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Movie is the Schema for the movies API
type Movie struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MovieSpec   `json:"spec,omitempty"`
	Status MovieStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MovieList contains a list of Movie
type MovieList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Movie `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Movie{}, &MovieList{})
}
