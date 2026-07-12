# Lab Exercise 10.2 Implementing Basic Scaling Modifiers


# Lab Exercise 10.2: Implementing Basic Scaling Modifiers

In this exercise we will learn how to create composite metrics using scaling modifiers in KEDA. To create a
composite metric, we will configure KEDA to use two RabbitMQ scalers, each using different RabbitMQ
clusters to get their metrics.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.
3. Completion of Lab Exercise 10.1.

## Lab Exercise

1. Create Second RabbitMQ Cluster:
To create the second cluster, reproduce the lab environment steps number 3 and 4 from Lab Exercise 10.1 to
create the second rabbitmq cluster. Just replace the namespace name from rabbitmq to rabbitmq-2.
sed 's/rabbitmq$/rabbitmq-2/' cluster.yaml | kubectl apply -f -
2. Create authentication creds for second cluster:
Create a Kuberenetes secret which will store second RabbitMQ cluster authentication credentials, used by
KEDA and the producer program.
Create a file named secret-2.yaml with the following contents and apply it using the following command.
```yaml
apiVersion: v1
kind: Secret
metadata:
name: keda-rabbitmq-secret-2
data:

host:
```
"YW1xcDovL2RlZmF1bHRfdXNlcl9obUdaRmhkZXdxNjVQNGRJZHg3OnFjOThuNGlHRDdNWVhNQlZGY0lP
Mm10QjV2b0R1Vl9uQHJhYmJpdG1xLWNsdXN0ZXIucmFiYml0bXEtMi5zdmMuY2x1c3Rlci5sb2NhbDo1N
jcy"
```bash
kubectl apply -f secret-2.yaml
```
3. Create messages in RabbitMQ cluster 1:
To create messages, use the producer program created in the lab setup to generate messages for RabbitMQ
cluster 1.
sed 's/value: "1"/value: "10"/' producer.yaml | kubectl create -f -
4. Create messages in RabbitMQ cluster 2:
Create a file named producer-2.yaml with the contents below and apply it using the following command.
```yaml
apiVersion: batch/v1
kind: Job
metadata:
generateName: producer-program-
spec:
template:
metadata:
labels:
app: producer-program
spec:
restartPolicy: Never
containers:
- name: producer-program
image: ghcr.io/kedify/blog05-python-producer-program:latest
env:
- name: MESSAGE_COUNT
value: "20"
- name: RABBITMQ_URL
valueFrom:
secretKeyRef:
name: keda-rabbitmq-secret-2
key: host
```
```bash
kubectl create -f producer-2.yaml
```
5. Create a ScaledObject with scaling modifier:
As discussed in the chapter, to create a composite metric we will be configuring a scaling modifier. In the below
configuration a formula is defined to calculate the average of two metrics, metric1 and metric2, by
dividing their sum by 2. The target for this composite metric is set at 5, indicating that the HPA will aim to
maintain this average value across the metrics.
Create a file named scaled-object-scaling-modifiers.yaml with the following contents and apply it
using the command below.
```yaml
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
name: keda-trigger-auth-rabbitmq-conn
namespace: default
spec:
secretTargetRef:
- parameter: host
name: keda-rabbitmq-secret
key: host
---
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
name: keda-trigger-auth-rabbitmq-conn-2
namespace: default
spec:
secretTargetRef:
- parameter: host
name: keda-rabbitmq-secret-2
key: host
---
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
name: keda-rabbitmq
namespace: default

spec:
scaleTargetRef:
apiVersion: apps/v1
kind: Deployment
name: consumer-program
minReplicaCount: 1
maxReplicaCount: 10
advanced:
scalingModifiers:
formula: "(metric1 + metric2) / 2"
target: "5"
metricType: "AverageValue"
triggers:
- type: rabbitmq
name: metric1
metadata:
protocol: amqp
queueName: testqueue
queueLength: "5"
authenticationRef:
name: keda-trigger-auth-rabbitmq-conn
- type: rabbitmq
name: metric2
metadata:
protocol: amqp
queueName: testqueue
queueLength: "10"
authenticationRef:
name: keda-trigger-auth-rabbitmq-conn-2
```
```bash
kubectl apply -f scaled-object-scaling-modifiers.yaml
```
6. Verify ScaledObject:
Execute the following command and ensure the READY column shows True.
```bash
kubectl get scaledobjects.keda.sh keda-rabbitmq
```
NAME SCALETARGETKIND SCALETARGETNAME MIN MAX TRIGGERS
AUTHENTICATION READY ACTIVE FALLBACK PAUSED AGE
keda-rabbitmq apps/v1.Deployment consumer-program 1 10 rabbitmq
keda-trigger-auth-rabbitmq-conn True True False Unknown 15m
7. Observe HPA and scaling behavior:
Execute the following command and observe how HPA reacts. As soon as the target threshold goes above
five, HPA starts scaling up and tries to maintain the target value at five.
```bash
kubectl get hpa --watch
```
NAME REFERENCE TARGETS
MINPODS MAXPODS REPLICAS AGE
keda-hpa-keda-rabbitmq Deployment/consumer-program <unknown>/5 (avg)
1 10 0 0s
keda-hpa-keda-rabbitmq Deployment/consumer-program 20/5 (avg)
1 10 1 0s
keda-hpa-keda-rabbitmq Deployment/consumer-program 5/5 (avg)
1 10 4 60s

## Summary

In this exercise we learned about creating a composite metric in the ScaledObject resource using scaling
modifiers.