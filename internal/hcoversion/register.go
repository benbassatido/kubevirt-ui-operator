package hcoversion

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// AddToScheme registers HyperConverged types for client/cache access.
func AddToScheme(s *runtime.Scheme) error {
	s.AddKnownTypes(GV, &HyperConverged{}, &HyperConvergedList{})
	metav1.AddToGroupVersion(s, GV)
	return nil
}
