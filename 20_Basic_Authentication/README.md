# Lab Exercise 8.1 Exploring Basic Authentication Mechanisms


# Lab Exercise 8.1: 

Exploring Basic Authentication Mechanisms
In this exercise, we will explore how to directly reference an environment variable for authentication.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.

## Lab Environment Setup

We will install RabbitMQ and Hashicorp Vault in Kubernetes using their respective operators.
1. Install RabbitMQ operator:
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
Create a file name cluster.yaml with the following contents and apply it using the command below.
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
type: Opaque
---
apiVersion: rabbitmq.com/v1beta1

kind: RabbitmqCluster
metadata:
name: rabbitmq-cluster
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
5. Create authentication creds:
Create a file name secret.yaml with the following contents and apply it using the command below.
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
6. Create RabbitMQ producer:
Create a file name producer.yaml with the following contents and apply it using the command below.
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
7. Create RabbitMQ consumer:
Create a file name consumer.yaml with the following contents and apply it using the command below.
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
8. Install Hashicorp Vault:
```bash
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update hashicorp
helm upgrade -i vault hashicorp/vault --set "server.dev.enabled=true" -n vault --create-namespace
```
9. Verify Hashicorp Vault:
Wait until the pod is in running state, as shown below:
```bash
kubectl get pods -n vault
```
```text
NAME READY STATUS RESTARTS AGE
vault-0 1/1 Running 0 71s
vault-agent-injector-55748c487f-fchcn 1/1 Running 0 71s
```
10. Access Vault UI:
Following the execution of the command, navigate to the URL http://localhost:8200 in your web browser. Upon
prompt for a Token, please input root as the value.
```bash
kubectl port-forward -n vault vault-0 8200:8200
```
```text
Forwarding from 127.0.0.1:8200 -> 8200
Forwarding from [::1]:8200 -> 8200
```
Congratulations, you have successfully set up the required lab environment.

## Lab Exercise

1. Create ScaledObject:
In the below specification, we are using RabbitMQ trigger in our ScaledObject. For connecting to RabbitMQ,
we are required to pass a host address field. But instead of directly hardcoding the address, we are using
hostFromEnv field to refer to the environment variable from the pod.
Create a file name scaled-object-direct-secret.yaml with the following contents and apply it using
the command below.
Note: hostFromEnv field is specific to RabbitMQ scaler, check the scaler-specific documentation to learn how
to pass credentials securely.
```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
name: keda-rabbitmq
spec:
scaleTargetRef:
name: consumer-program
envSourceContainerName: consumer-program # Specifies the container with the
```
environment variable
triggers:
- type: rabbitmq
```yaml
metadata:
protocol: amqp
queueName: testqueue
hostFromEnv: RABBITMQ_URL # References the environment variable for
```
secure authentication
queueLength: "5"
```bash
kubectl apply -f scaled-object-direct-secret.yaml
```
2. Verify ScaledObject:
When a ScaledObject is created, KEDA conducts an initial check to confirm the availability of essential
variables, with host being a crucial variable for the RabbitMQ cluster. In the context of utilizing a RabbitMQ
scaler, KEDA uses the provided details to attempt a connection with the RabbitMQ cluster. If KEDA
successfully references the host variable from the pod's environment variables, the preliminary verifications are
deemed successful. Consequently, KEDA updates its READY status to true, indicating readiness, as evident
from the below output
The following output shows the status of the keda-rabbitmq ScaledObject, reflecting that it's ready to
function and actively monitor RabbitMQ for scaling actions.
```bash
kubectl get scaledobjects.keda.sh keda-rabbitmq
```
NAME SCALETARGETKIND SCALETARGETNAME MIN MAX TRIGGERS
AUTHENTICATION READY ACTIVE FALLBACK PAUSED AGE
keda-rabbitmq apps/v1.Deployment consumer-program rabbitmq
True False Unknown Unknown 7s
3. Clean up:
```bash
kubectl delete -f scaled-object-direct-secret.yaml
```

## Summary

In this exercise, we learned how to use environment variables for secure authentication in KEDA
ScaledObjects, specifically utilizing the hostFromEnv field to dynamically fetch the RabbitMQ host address
from the pod's environment. This approach ensures secure and flexible configuration for the RabbitMQ scaler
in KEDA, as demonstrated through the creation and validation of a ScaledObject, which successfully
connected to the RabbitMQ cluster and indicated readiness for scaling operations.