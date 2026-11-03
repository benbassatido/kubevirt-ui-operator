package manifests

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func boolPtr(b bool) *bool {
	return &b
}

func int32Ptr(i int32) *int32 {
	return &i
}

func intstrFromInt32(i int32) intstr.IntOrString {
	return intstr.FromInt32(i)
}

func resourceQuantity(s string) resource.Quantity {
	return resource.MustParse(s)
}

func tcpProtocol() *corev1.Protocol {
	p := corev1.ProtocolTCP
	return &p
}
