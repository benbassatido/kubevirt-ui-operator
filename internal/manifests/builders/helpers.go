package builders

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func intstrFromInt32(i int32) intstr.IntOrString {
	return intstr.FromInt32(i)
}

func resourceQuantity(s string) resource.Quantity {
	return resource.MustParse(s)
}

func tcpProtocol() *corev1.Protocol {
	return new(corev1.ProtocolTCP)
}
