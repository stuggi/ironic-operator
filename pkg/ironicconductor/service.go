package ironicconductor

import (
	ironicv1 "github.com/openstack-k8s-operators/ironic-operator/api/v1beta1"
	ironic "github.com/openstack-k8s-operators/ironic-operator/pkg/ironic"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Service - Service for conductor pod services
func Service(
	serviceName string,
	instance *ironicv1.IronicConductor,
	serviceLabels map[string]string,
) *corev1.Service {

	var ports []corev1.ServicePort

	// RPC Transport is json-rpc so expose the service
	if instance.Spec.RPCTransport == "json-rpc" {
		jsonRPCPort := corev1.ServicePort{
			Name:     ironic.JSONRPCComponent,
			Port:     8089,
			Protocol: corev1.ProtocolTCP,
		}
		ports = append(ports, jsonRPCPort)
	}

	// There is no provision network so expose the deploy HTTP interface
	// as a service to enable virtual media boot
	if instance.Spec.ProvisionNetwork == "" {
		httpbootPort := corev1.ServicePort{
			Name:     ironic.HttpbootComponent,
			Port:     8088,
			Protocol: corev1.ProtocolTCP,
		}
		ports = append(ports, httpbootPort)
	}

	if len(ports) == 0 {
		return nil
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: instance.Namespace,
			Labels:    serviceLabels,
		},
		Spec: corev1.ServiceSpec{
			Selector: serviceLabels,
			Ports:    ports,
		},
	}
}

/*
// Route - Route for httpboot service when no provisioning network
func Route(
	serviceName string,
	instance *ironicv1.IronicConductor,
	routeLabels map[string]string,
) *routev1.Route {
	serviceRef := routev1.RouteTargetReference{
		Kind: "Service",
		Name: serviceName,
	}
	routePort := &routev1.RoutePort{
		TargetPort: intstr.FromString(ironic.HttpbootComponent),
	}
	return &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: instance.Namespace,
			Labels:    routeLabels,
		},
		Spec: routev1.RouteSpec{
			To:   serviceRef,
			Port: routePort,
		},
	}
}
*/
