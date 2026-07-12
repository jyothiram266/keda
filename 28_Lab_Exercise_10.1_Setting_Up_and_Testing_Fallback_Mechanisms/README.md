# Lab Exercise 10.1 Setting Up and Testing Fallback Mechanisms


# Lab Exercise 10.1: Setting Up and Testing Fallback Mechanisms

In this exercise we will learn how to implement and validate a fallback strategy for a Kubernetes deployment
using KEDA. The exercise will guide you through the steps of monitoring a deployment’s ability to scale based
on RabbitMQ metrics and enforce a predefined replica count when these metrics are unavailable.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.

## Lab Environment Setup

We will install RabbitMQ in Kubernetes using the RabbitMQ Cluster Kubernetes Operator.
1. Install the RabbitMQ operator:
We will use RabbitMQ operator to provision and manage our cluster. Use the following command to install
operator:
```bash
kubectl apply -f
```
"https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-op
erator.yml"
2. Verify RabbitMQ operator:
Wait until the pod is in running state, as shown below:
```bash
kubectl get pods -n rabbitmq-system
```
```text
NAME READY STATUS RESTARTS AGE
rabbitmq-cluster-operator-ccf488f4c-nqrwn 1/1 Running 0 18s
```
3. Create RabbitMQ cluster:
Create a file named cluster.yaml with the following contents and apply it using the command below.
```yaml
apiVersion: v1
kind: Namespace
metadata:
name: rabbitmq
---
apiVersion: v1
data:
default_user.conf:
```
ZGVmYXVsdF91c2VyID0gZGVmYXVsdF91c2VyX2htR1pGaGRld3E2NVA0ZElkeDcKZGVmYXVsdF9wYXNzI
D0gcWM5OG40aUdEN01ZWE1CVkZjSU8ybXRCNXZvRHVWX24K
password: cWM5OG40aUdEN01ZWE1CVkZjSU8ybXRCNXZvRHVWX24=
username: ZGVmYXVsdF91c2VyX2htR1pGaGRld3E2NVA0ZElkeDc=
```yaml
kind: Secret
metadata:
name: my-secret
namespace: rabbitmq
type: Opaque
---
apiVersion: rabbitmq.com/v1beta1
kind: RabbitmqCluster
metadata:
name: rabbitmq-cluster
namespace: rabbitmq
spec:
secretBackend:
externalSecret:
name: "my-secret"
```
```bash
kubectl apply -f cluster.yaml
```
4. Verify RabbitMQ cluster:
```bash
kubectl get pods -n rabbitmq
```
```text
NAME READY STATUS RESTARTS AGE
rabbitmq-cluster-server-0 1/1 Running 0 63s
```
5. Create Authentication creds:
```yaml
apiVersion: v1
kind: Secret
metadata:
name: keda-rabbitmq-secret
type: Opaque
data:
host:
```
"YW1xcDovL2RlZmF1bHRfdXNlcl9obUdaRmhkZXdxNjVQNGRJZHg3OnFjOThuNGlHRDdNWVhNQlZGY0lP
Mm10QjV2b0R1Vl9uQHJhYmJpdG1xLWNsdXN0ZXIucmFiYml0bXEuc3ZjLmNsdXN0ZXIubG9jYWw6NTY3M
g=="
```bash
kubectl apply -f secret.yaml
```
6. Create RabbitMQ Producer:
Create a file named producer.yaml with the following contents and apply it using the command below.
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
value: "1"
- name: RABBITMQ_URL
valueFrom:
secretKeyRef:
name: keda-rabbitmq-secret

key: host
```
```bash
kubectl create -f producer.yaml
```
7. Create RabbitMQ Consumer:
Create a file named consumer.yaml with the contents below and apply it using the following command.
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
name: consumer-script-config
data:
consumer-script.sh: |
```
#! /bin/bash
currentMessage=""
handle_sigterm() {
if [ -n "$currentMessage" ]; then
```bash
echo "SIGTERM signal received while processing a message."
curl -X POST http://result-analyzer-service:8080/kill/count -s
echo "Kill count HTTP request sent."
```
else
```bash
echo "SIGTERM signal received, but no message was being processed."
```
fi
exit 0
}
trap 'handle_sigterm' SIGTERM
while true; do
```bash
echo "Waiting for message...\n"
```
if ! currentMessage=$(amqp-consume --url="$RABBITMQ_URL" -q "testqueue" -c
1 cat); then
```bash
echo "Error occurred during message consumption. Exiting...\n"
```
continue
fi
```bash
echo "Message received, processing: $currentMessage \n"
```
i=1
while [ $i -le 360 ]; do
```bash
echo "Encoding video $i"
sleep 1
```
i=$((i+1))
done
currentMessage=""
```bash
curl -X POST http://result-analyzer-service:8080/create/count -s
echo "Waiting for next message...\n"
```
done
```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
name: consumer-program
spec:
replicas: 1
selector:
matchLabels:
app: consumer-program
template:
metadata:
labels:
app: consumer-program
spec:
containers:
- name: consumer-program
image: ghcr.io/kedify/blog05-cli-consumer-program:latest
command: ["/bin/bash"]
args: ["/scripts/consumer-script.sh"]
volumeMounts:
- name: script-volume
mountPath: /scripts
env:
- name: RABBITMQ_URL
valueFrom:
secretKeyRef:
name: keda-rabbitmq-secret
key: host

volumes:
- name: script-volume
configMap:
name: consumer-script-config
```
```bash
kubectl apply -f consumer.yaml
```
Congratulations, you have successfully set up the required lab environment.

## Lab Exercise

1. Check the current replicas of Consumer program:
```bash
kubectl get pods -l=app=consumer-program
```
```text
NAME READY STATUS RESTARTS AGE
consumer-program-6c4fc8b66b-pgdzl 1/1 Running 0 22s
```
2. Create ScaledObject with fallback configured:
As discussed in the chapter, in the configuration below we have configured the deployment to fall back to 6
replicas after 3 consecutive failures in retrieving metrics from RabbitMQ.
Create a file named scaled-object-fallback.yaml with the following contents and apply it using the
command below.
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
fallback:
failureThreshold: 3
replicas: 5
triggers:
- type: rabbitmq
metadata:
protocol: amqp
queueName: testqueue
queueLength: "5"
authenticationRef:
name: keda-trigger-auth-rabbitmq-conn
```
```bash
kubectl apply -f scaled-object-fallback.yaml
```
3. Check status of ScaledObject:
In the output below you can observe that columns READY and ACTIVE are set to True while column
FALLBACK is set to False — indicating ScaledObject operation is working as intended.
```bash
kubectl get scaledobjects.keda.sh keda-rabbitmq
```
NAME SCALETARGETKIND SCALETARGETNAME MIN MAX TRIGGERS
AUTHENTICATION READY ACTIVE FALLBACK PAUSED AGE
keda-rabbitmq apps/v1.Deployment consumer-program 1 10 rabbitmq
True True False Unknown 60s
4. Simulate RabbitMQ cluster DownTime:
To force ScaledObject to trigger the fallback mechanism, we will delete the RabbitMQ cluster created during
the lab setup. With this action, ScaledObject will no longer be able to communicate with RabbitMQ.
Execute the following command to delete the cluster.
```bash
kubectl delete rabbitmqcluster/rabbitmq-cluster -n rabbitmq
```
5. Check Status of ScaledObject and Observe Pods:
Based on the settings in the ScaledObject, if there are three successive instances where metric retrieval fails
(which would take approximately two minutes, depending on the polling frequency, with the default set to every
30 seconds), KEDA will activate the fallback process and adjust the number of consumer replicas to six.
In the below output you can observe that columns READY and ACTIVE are set to True and False while
column FALLBACK is set to True — indicating ScaledObject has triggered the fallback mechanism.
```bash
kubectl get scaledobjects.keda.sh keda-rabbitmq
```
NAME SCALETARGETKIND SCALETARGETNAME MIN MAX TRIGGERS
AUTHENTICATION READY ACTIVE FALLBACK PAUSED AGE
keda-rabbitmq apps/v1.Deployment consumer-program 1 10 rabbitmq
True False True Unknown 60s
```bash
kubectl get pods -l=app=consumer-program
```
```text
NAME READY STATUS RESTARTS AGE
consumer-program-6c4fc8b66b-7g2rf 1/1 Running 0 3m57s
consumer-program-6c4fc8b66b-mqjj4 1/1 Running 0 3m57s
consumer-program-6c4fc8b66b-pgdzl 1/1 Running 0 14m
consumer-program-6c4fc8b66b-rrjms 1/1 Running 0 3m57s
consumer-program-6c4fc8b66b-xkwhf 1/1 Running 0 3m57s
You can also observe the reason for failure by describing the ScaledObject resource using the following
command.
```
```bash
kubectl describe scaledobjects.keda.sh keda-rabbitmq | grep -A 10 -i "Events:"
```
Events:
Type Reason Age From Message
```yaml
---- ------ ---- ---- -------
```
Normal KEDAScalersStarted 34m keda-operator Started scalers
watch Normal ScaledObjectReady 34m
keda-operator ScaledObject is ready for scaling Warning KEDAScalerFailed
33m (x3 over 33m) keda-operator error establishing rabbitmq connection: dial
tcp 10.102.184.138:5672: connect: connection refused
Normal KEDAScalersStarted 32m (x2 over 34m) keda-operator Scaler rabbitmq
is built. Warning KEDAScalerFailed 4m56s (x80 over 31m)
keda-operator error establishing rabbitmq connection: dial tcp: lookup
rabbitmq-cluster.rabbitmq.svc.cluster.local on 10.96.0.10:53: no such host
6. Clean up:
Delete the resource created above and recreate the RabbitMQ cluster:
```bash
kubectl apply -f cluster.yaml
kubectl delete -f scaled-object-fallback.yaml
kubectl scale deployment consumer-program --replicas=1
```

## Summary

In this exercise we configured a fallback mechanism for a Kubernetes deployment using KEDA. Initially, the
current pod replicas were checked, and then a ScaledObject was created with the fallback mechanism to scale
the consumer program to 6 replicas if it fails to fetch metrics from RabbitMQ three times. The exercise
simulated RabbitMQ downtime and observed KEDA triggering the fallback mechanism.