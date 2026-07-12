# Lab Exercise 8.2 Implementing TriggerAuthentication


# Lab Exercise 8.2: Implementing TriggerAuthentication

In this exercise, we will explore how to use environment variables with TriggerAuthentication CRD for secure
authentication.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.
3. Completion of Lab Exercise 8.1.

## Lab Exercise

1. Create ScaledObject:
As discussed in the chapter, In the configuration detailed below, we employ a TriggerAuthentication resource,
utilizing environment variables as the authentication source. This resource is linked to the ScaledObject via the
authenticationRef field. The TriggerAuthentication CRD specifies the required host parameter for the
RabbitMQ scaler, referencing the environment variable RABBITMQ_URL and the consumer-program container
within the pod.
Create a file name scaled-object-trigger-auth-env.yaml with the following contents and apply it
using the command below.
```yaml
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
name: rabbitmq-trigger-auth
spec:
env:
- parameter: host
name: RABBITMQ_URL
containerName: consumer-program
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
kubectl apply -f scaled-object-trigger-auth-env.yaml
```
2. Verify ScaledObject:
```bash
kubectl get scaledobjects.keda.sh keda-rabbitmq
```
NAME SCALETARGETKIND SCALETARGETNAME MIN MAX TRIGGERS
AUTHENTICATION READY ACTIVE FALLBACK PAUSED AGE
keda-rabbitmq apps/v1.Deployment consumer-program rabbitmq
True False Unknown Unknown 7s
3. Clean up:
```bash
kubectl delete -f scaled-object-trigger-auth-env.yaml
```

## Summary

In this exercise we learned how to use TriggerAuthentication Custom Resource Definition (CRD) with
ScaledObject to securely handle authentication, utilizing environment variables for secure and dynamic
authentication. Upon successful deployment, KEDA confirms the readiness of the ScaledObject for scaling
operations based on RabbitMQ triggers, as reflected in the READY status in the kubectl output.