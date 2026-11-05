package hcoversion

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	Group   = "hco.kubevirt.io"
	Version = "v1"
	Kind    = "HyperConverged"

	OperatorVersionName = "operator"
)

var (
	GVK = schema.GroupVersionKind{Group: Group, Version: Version, Kind: Kind}
	GV  = schema.GroupVersion{Group: Group, Version: Version}
)

// HyperConverged is a minimal view of hco.kubevirt.io/HyperConverged for version discovery.
type HyperConverged struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Status            HyperConvergedStatus `json:"status"`
}

type HyperConvergedStatus struct {
	Versions []ComponentVersion `json:"versions,omitempty"`
}

type ComponentVersion struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// HyperConvergedList is a list of HyperConverged resources.
type HyperConvergedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []HyperConverged `json:"items"`
}

func (in *HyperConverged) DeepCopyInto(out *HyperConverged) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Status.Versions != nil {
		out.Status.Versions = make([]ComponentVersion, len(in.Status.Versions))
		copy(out.Status.Versions, in.Status.Versions)
	}
}

func (in *HyperConverged) DeepCopy() *HyperConverged {
	if in == nil {
		return nil
	}
	out := new(HyperConverged)
	in.DeepCopyInto(out)
	return out
}

func (in *HyperConverged) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *HyperConvergedList) DeepCopyInto(out *HyperConvergedList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		out.Items = make([]HyperConverged, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}

func (in *HyperConvergedList) DeepCopy() *HyperConvergedList {
	if in == nil {
		return nil
	}
	out := new(HyperConvergedList)
	in.DeepCopyInto(out)
	return out
}

func (in *HyperConvergedList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
