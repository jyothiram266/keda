# Lab Exercise 7.1 Configuring ScaledJobs for RabbitMQ


# Lab Exercise 7.1: 

Configuring ScaledJobs for RabbitMQ
This exercise introduces the configuration and deployment of ScaledJobs using KEDA to dynamically manage
workloads in response to messages in a RabbitMQ queue. The core of the exercise involves creating a
ScaledJob resource that instructs KEDA to automatically scale job instances based on the volume of
messages, showcasing a practical implementation of event-driven autoscaling in a Kubernetes environment.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Familiarity with RabbitMQ.
3. Access to a Kubernetes environment with KEDA Metric Server installed as per Lab 5.

## Lab Environment Setup

We will install RabbitMQ cluster in Kubernetes using RabbitMQ operator
1. Install Operator:
We will use RabbitMQ operator to provision and manage our cluster. Use the following command to install
operator.
```bash
kubectl apply -f
```
"https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-op
erator.yml"
2. Verify Operator:
Wait until the pod is in running state, as shown below:
```bash
kubectl get pods -n rabbitmq-system
```
```text
NAME READY STATUS RESTARTS AGE
rabbitmq-cluster-operator-ccf488f4c-nqrwn 1/1 Running 0 18s
```
3. Create a cluster:
Create a file name rabbitmq_cluster.yaml with the contents below and apply it using the following
command.
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
kubectl apply -f rabbitmq_cluster.yaml
```
4. Verify the cluster:
Wait until the pod is in running state, as shown below:
```bash
kubectl get pods -n rabbitmq
```
```text
NAME READY STATUS RESTARTS AGE
rabbitmq-cluster-server-0 1/1 Running 0 63s
Congratulations, you have successfully set up the required lab environment! Now we can move on to the
exercise for this lab.
```

## Lab Exercise

1. Create a Kubernetes secret for storing RabbitMQ credentials.
This secret holds the RabbitMQ address in the form of URL string, containing both username, password and
RabbitMQ address.
Create a file name rabbitmq-creds-secret.yaml with the following contents and apply it using the
command below.
```yaml
apiVersion: v1
kind: Secret
metadata:
name: keda-rabbitmq-secret
data:
host:
```
"YW1xcDovL2RlZmF1bHRfdXNlcl9obUdaRmhkZXdxNjVQNGRJZHg3OnFjOThuNGlHRDdNWVhNQlZGY0lP
Mm10QjV2b0R1Vl9uQHJhYmJpdG1xLWNsdXN0ZXIucmFiYml0bXEuc3ZjLmNsdXN0ZXIubG9jYWw6NTY3M
g=="
```bash
kubectl apply -f rabbitmq-creds-secret.yaml
```
2. Create a Kubernetes Configmap for storing RabbitMQ Consumer bash script.
To consume messages from RabbitMQ, we will be using the below bash script. We are using amqp-consume
CLI to read messages from RabbitMQ. This command blocks execution until a message is available to read
from RabbitMQ. Once the message is available it starts processing (note: to simulate some processing action,
our script sleeps for a 5 minute duration). Once processing is completed, the script posts the result to the
result-analyzer and exits.
Create a file name rabbitmq-consumer-script.yaml with the following contents and apply it using the
command below.
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
```bash
echo "Waiting for message...\n"
```
if ! currentMessage=$(amqp-consume --url="$RABBITMQ_URL" -q "testqueue" -c 1
cat); then
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
```bash
kubectl apply -f rabbitmq-consumer-script.yaml
```
3. Create a Producer application.
The job described below ensures the creation of a RabbitMQ queue named testqueue and also establishes
the corresponding exchange channel. Additionally, it binds a route to facilitate message routing from this
channel to the designated queue.
The job also populates the queue with messages, the quantity of which is determined by the value of the
MESSAGE_COUNT environment variable. With current value of 15, it will generate 15 messages in testqueue.
Create a file name rabbitmq-producer.yaml with the following contents and apply it using the command
below.
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
value: "15"
- name: RABBITMQ_URL
valueFrom:
secretKeyRef:
name: keda-rabbitmq-secret
key: host
```
```bash
kubectl create -f rabbitmq-producer.yaml
```
4. Create a ScaledJob.
The following configuration creates a ScaledJob. Here is a detailed breakdown:
- Trigger Authentication for RabbitMQ: The TriggerAuthentication kind is used to securely pass the
RabbitMQ connection details to KEDA. It references a secret (keda-rabbitmq-secret) to obtain the
RabbitMQ host address. This ensures that sensitive information like the host address is not exposed
directly in the YAML file, but is instead stored securely in a Kubernetes secret.
- Consumer Script Execution from ConfigMap: The ScaledJob spec defines a job that runs a
consumer script. The script itself is stored in a ConfigMap (consumer-script-config), which is
mounted as a volume inside the container at /scripts. The consumer program image
(ghcr.io/kedify/blog05-cli-consumer-program:latest) is instructed to execute the script
(consumer-script.sh) upon starting, ensuring that the job's main task is to process messages from
RabbitMQ.
- Trigger Section for RabbitMQ in ScaledJob: The triggers field in the ScaledJob spec specifies the
conditions under which KEDA should scale the job. It's configured to respond to the length of the
RabbitMQ queue named “testqueue”. The value field indicates that a new job should be created for
every message in the queue. KEDA checks the queue length every 10 seconds (pollingInterval),
and the job scaling is subject to a maximum of 100 replicas (maxReplicaCount). The authentication
details for RabbitMQ are linked through the authenticationRef, pointing to the previously defined
TriggerAuthentication.
Create a file name scaled-job.yaml with the following contents and apply it using the command below.
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
kind: ScaledJob
metadata:
name: rabbitmq-scaledjob
namespace: default

spec:
jobTargetRef:
template:
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
restartPolicy: Never
pollingInterval: 10 # How often KEDA will check the RabbitMQ queue
successfulJobsHistoryLimit: 100 # Number of successful jobs to keep
failedJobsHistoryLimit: 100 # Number of failed jobs to keep
maxReplicaCount: 100 # Maximum number of jobs that KEDA can create
scalingStrategy:
strategy: "default" # Scaling strategy (default, custom, or accurate)
triggers:
- type: rabbitmq
metadata:
protocol: amqp
queueName: testqueue
mode: QueueLength
value: "1" # Number of messages per job
authenticationRef:
name: keda-trigger-auth-rabbitmq-conn
```
```bash
kubectl apply -f scaled-job.yaml
```
5. Monitoring the scaling behavior:
Note: A ScaledJob does not create a HPA resource like ScaledObject, the pods created are totally handled by
ScaledJob resource.
As we created 15 messages in RabbitMQ testqueue and configured ScaledJob to create one pod per
message, you can see below it has created 15 pods to process 15 messages.
```bash
kubectl get pods -l=scaledjob.keda.sh/name=rabbitmq-scaledjob
```
```text
NAME READY STATUS RESTARTS AGE
rabbitmq-scaledjob-7nw9k-dhcs9 1/1 Running 0 25s
rabbitmq-scaledjob-dt8tr-52mkv 1/1 Running 0 26s
rabbitmq-scaledjob-gnxb9-ggz6n 1/1 Running 0 27s
rabbitmq-scaledjob-h5mzq-q6tsh 1/1 Running 0 27s
rabbitmq-scaledjob-hjbfq-6jppm 1/1 Running 0 25s
rabbitmq-scaledjob-hjqk9-g78lx 1/1 Running 0 26s
rabbitmq-scaledjob-k8fcc-n5zpk 1/1 Running 0 25s
rabbitmq-scaledjob-kd9t2-g6skd 1/1 Running 0 27s
rabbitmq-scaledjob-kh5pd-qkcjt 1/1 Running 0 25s
rabbitmq-scaledjob-kkbg4-k2vrx 1/1 Running 0 24s
rabbitmq-scaledjob-km7b8-w85n2 1/1 Running 0 27s
rabbitmq-scaledjob-n6ksh-cmkbt 1/1 Running 0 25s
rabbitmq-scaledjob-n7wgf-smppx 1/1 Running 0 27s
rabbitmq-scaledjob-rvtg8-mq4xv 1/1 Running 0 26s
rabbitmq-scaledjob-zfn96-fwvcc 1/1 Running 0 25s
```
6. Clean up:
```bash
kubectl delete jobs --all --wait
```

## Summary

In this exercise, we set up a ScaledJob in KEDA to process messages from a RabbitMQ queue, which involved
the creation of Kubernetes secrets for RabbitMQ credentials, a ConfigMap for a consumer script, and a
producer application to populate the queue. We applied the scaled job configuration, which triggers KEDA to
dynamically create pods based on the number of messages in the queue, demonstrating an efficient,
event-driven scaling mechanism for workload processing.