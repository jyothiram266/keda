# Lab Exercise 8.3 Implementing TriggerAuthentication Referencing a Secret


# Lab Exercise 8.3: Implementing

TriggerAuthentication Referencing a Secret
In this exercise, we will explore how to use Kubernetes secrets with TriggerAuthentication CRD for secure
authentication.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.
3. Completion of Lab Exercises 8.1 and 8.2.

## Lab Exercise

1. Create ScaledObject:
As discussed in the chapter, in the configuration detailed below, we employ a TriggerAuthentication resource,
utilizing Kubernetes secrets (created during environment setup) as the authentication source. This resource is
linked to the ScaledObject via the authenticationRef field. The TriggerAuthentication CRD specifies the
required host parameter for the RabbitMQ scaler, referencing the Kubernetes secret name and key “host”.
Create a file name scaled-object-trigger-auth-secret.yaml with the following contents and apply it
using the command below.
```yaml
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
name: rabbitmq-trigger-auth
spec:
secretTargetRef:
- parameter: host
name: keda-rabbitmq-secret
key: host
---

apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
name: keda-rabbitmq
spec:
scaleTargetRef:
name: consumer-program
triggers:
- type: rabbitmq
metadata:
protocol: amqp
queueName: testqueue
queueLength: "5"
authenticationRef:
name: rabbitmq-trigger-auth
```
```bash
kubectl apply -f scaled-object-trigger-auth-secret.yaml
```
2. Verify Scaled Object:
```bash
kubectl get scaledobjects.keda.sh keda-rabbitmq
```
NAME SCALETARGETKIND SCALETARGETNAME MIN MAX TRIGGERS
AUTHENTICATION READY ACTIVE FALLBACK PAUSED AGE
keda-rabbitmq apps/v1.Deployment consumer-program rabbitmq
True False Unknown Unknown 7s
3. Clean up:
```bash
kubectl delete -f scaled-object-trigger-auth-secret.yaml
```

## Summary

In this exercise, we learned how to use Kubernetes secrets with the TriggerAuthentication CRD. The
TriggerAuthentication resource is linked to the ScaledObject, referencing a specific key in the Kubernetes
secret for the required host parameter. The successful deployment and verification process is evidenced by the
ScaledObject's READY status, indicating readiness for auto-scaling based on RabbitMQ workload.