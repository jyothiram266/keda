# Lab Exercise 8.5 Understanding and Configuring ClusterTriggerAuthentication


# Lab Exercise 8.5: Understanding and Configuring

ClusterTriggerAuthentication
In this exercise, we will explore how to use Kubernetes secrets with ClusterTriggerAuthentication CRD for
secure authentication.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.
3. Completion of Lab Exercises 8.1, 8.2, 8.3 and 8.4.

## Lab Exercise

1. Create ScaledObject and ClusterTriggerAuthentication
Create a file name scaled-object-cluster-trigger.yaml with the following contents and apply it using
the command below.
Please note the below ScaledObject is being created in the keda namespace.
```yaml
apiVersion: keda.sh/v1alpha1
kind: ClusterTriggerAuthentication
metadata:
name: rabbitmq-auth
spec:
secretTargetRef:
- parameter: host
key: host
name: keda-rabbitmq-secret
---
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:

name: keda-rabbitmq
namespace: keda
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
kind: ClusterTriggerAuthentication
name: rabbitmq-auth
```
```bash
kubectl apply -f scaled-object-cluster-trigger.yaml
```
2. Verify ScaledObject:
Execute the following command to view the result.
```bash
kubectl get scaledobjects.keda.sh keda-rabbitmq
```
NAME SCALETARGETKIND SCALETARGETNAME MIN MAX TRIGGERS
AUTHENTICATION READY ACTIVE FALLBACK PAUSED AGE
keda-rabbitmq apps/v1.Deployment consumer-program rabbitmq
False False Unknown Unknown 7s
From the above output, the READY column is marked as false because the ClusterTriggerAuthentication is
searching for the keda-rabbitmq-secret in the keda namespace where KEDA is installed, but the secret
exists in the default namespace (created during lab setup). To resolve this, configure the KEDA operator to
recognize the default namespace by setting KEDA_CLUSTER_OBJECT_NAMESPACE env variable to default.
3. Delete ScaledObject:
```bash
kubectl delete -f scaled-object-cluster-trigger.yaml
```
4. Editing KEDA Operator:
Execute the following command to set the environment variable value to default.
```bash
kubectl set env -n keda deployment/keda-operator -c keda-operator
```
KEDA_CLUSTER_OBJECT_NAMESPACE=default
5. Re-Apply ScaledObject and observe result:
```bash
kubectl apply -f scaled-object-cluster-trigger.yaml
kubectl get scaledobjects.keda.sh keda-rabbitmq
```
NAME SCALETARGETKIND SCALETARGETNAME MIN MAX TRIGGERS
AUTHENTICATION READY ACTIVE FALLBACK PAUSED AGE
keda-rabbitmq apps/v1.Deployment consumer-program rabbitmq
True False Unknown Unknown 7s
6. Clean up:
```bash
kubectl delete -f scaled-object-cluster-trigger.yaml
```

## Summary

In this exercise, we configured ClusterTriggerAuthentication with Kubernetes secrets for secure, cluster-wide
authentication in KEDA. The exercise guided you through creating a ClusterTriggerAuthentication and a
ScaledObject for a RabbitMQ scaler. Upon initial deployment, the READY status was false due to namespace
mismatch. By adjusting the KEDA operator's environment variable (KEDA_CLUSTER_OBJECT_NAMESPACE) to
the correct namespace and reapplying the configuration, the READY status successfully transitioned to true,
validating the correct setup and functionality of cluster-wide authentication in KEDA.