# Lab Exercise 6.1 Kafka Cluster Setup


# Lab Exercise 6.1: Kafka Cluster Setup

In this exercise, we will install Kafka Cluster in Kubernetes and verify it is working by deploying a producer and
a consumer application.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Familiarity with Kafka.
3. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.

## Lab Exercise

1. Create Kafka namespace:
```bash
kubectl create namespace kafka
```
2. Install Kafka operator.
We will be using the Strimzi kubernetes operator to create and manage our Kafka cluster. Use the command
below to install the operator.
```bash
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
```
3. Verify operator:
Execute the following command and ensure that the operator is in running state, the READY column should
show 1/1.
```bash
kubectl get deployments strimzi-cluster-operator -n kafka
```
NAME READY UP-TO-DATE AVAILABLE AGE
strimzi-cluster-operator 1/1 1 1 12m
4. Create a Kafka cluster.
Once the operator is in running state. The configuration below utilizes the Kafka CRD of the Strimzi operator
and creates a cluster called my-cluster.
Kafka Specification:
a. Number of partitions: 5
b. Number of brokers: 3
Create a file named kafka-cluster.yaml with the following contents and apply it using the command
below.
```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: Kafka
metadata:
name: my-cluster
namespace: kafka
spec:
kafka:
config:
offsets.topic.replication.factor: 3
transaction.state.log.replication.factor: 3
transaction.state.log.min.isr: 2
default.replication.factor: 3
min.insync.replicas: 2
inter.broker.protocol.version: "3.6"
storage:
type: ephemeral
listeners:
- name: plain
port: 9092
type: internal
tls: false
- name: tls
port: 9093
type: internal
tls: true
version: 3.6.0
replicas: 3
entityOperator:
topicOperator: {}
userOperator: {}
zookeeper:
storage:
type: ephemeral

replicas: 3
```
```bash
kubectl apply -f kafka-cluster.yaml
```
5. Verify that the Kafka cluster is up and running.
Execute the command below and ensure that the READY in the output column is True.
watch kubectl get kafka.kafka.strimzi.io/my-cluster -n kafka
NAME DESIRED KAFKA REPLICAS DESIRED ZK REPLICAS READY
my-cluster 3 3 True
6. Create Kafka topics.
The configuration below creates a Kafka topic called my-topic by utilizing KafkaTopic CRD provided by the
Strimzi operator.
Create a file named topic.yaml with the following contents and apply it using the command provided.
```yaml
kind: KafkaTopic
apiVersion: kafka.strimzi.io/v1beta2
metadata:
name: my-topic
labels:
strimzi.io/cluster: my-cluster
namespace: kafka
spec:
partitions: 5
replicas: 3
config:
retention.ms: 604800000
segment.bytes: 1073741824
```
```bash
kubectl apply -f topic.yaml
```
7. Create Kafka consumers.
The below configuration creates a Kubernetes deployment using a public container image
(quay.io/zroubalik/kafka-app:latest). This application is configured to read Kafka messages from a
topic called “my-topic” (created in our previous step).
Create a file named consumer.yaml with the contents below and apply it using the following command.
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
name: kafka-amqstreams-consumer
labels:
app: kafka-amqstreams-consumer
spec:
selector:
matchLabels:
app: kafka-amqstreams-consumer
template:
metadata:
labels:
app: kafka-amqstreams-consumer
spec:
containers:
- name: kafka-amqstreams-consumer
image: quay.io/zroubalik/kafka-app:latest
imagePullPolicy: IfNotPresent
securityContext:
allowPrivilegeEscalation: false
runAsNonRoot: true
capabilities:
drop:
- ALL
seccompProfile:
type: RuntimeDefault
env:
- name: BOOTSTRAP_SERVERS
value: my-cluster-kafka-bootstrap.kafka.svc:9092
resources:
requests:
cpu: 100m
memory: 100Mi
limits:
cpu: 500m
memory: 500Mi
command:
- /kafkaconsumerapp
```
```bash
kubectl apply -f consumer.yaml
```
8. Create Kafka Producer:
The below configuration creates a Kubernetes Job using a public container image (scal). This application is
configured to write messages into Kafka. The number of messages to be created, delay between messages
and topic name is configured via environment variables.
In the below configuration we have set the MESSAGE_COUNT environment variable to 20, which will produce
20 messages.
Create a file name producer.yaml with the following contents and apply it using the following command.
```yaml
apiVersion: batch/v1
kind: Job
metadata:
generateName: kafka-amqstreams-producer-
namespace: kafka
spec:
parallelism: 1
completions: 1
backoffLimit: 1
template:
metadata:
name: kafka-amqstreams-producer
labels:
app: kafka-amqstreams-producer
spec:
restartPolicy: Never
containers:
- name: kafka-amqstreams-producer
image: quay.io/zroubalik/kafka-app:latest
imagePullPolicy: IfNotPresent
securityContext:
allowPrivilegeEscalation: false
runAsNonRoot: true
capabilities:
drop:
- ALL
seccompProfile:
type: RuntimeDefault
command: ["/kafkaproducerapp"]
env:
- name: BOOTSTRAP_SERVERS
value: my-cluster-kafka-bootstrap.kafka.svc:9092
- name: TOPIC
value: my-topic
- name: MESSAGE_COUNT
value: "20"
- name: DELAY_MS
value: "100"
```
```bash
kubectl create -f producer.yaml
```
9. Verify logs on consumer:
Once messages are created by the producer you can execute the command below to check if the consumer is
reading the newly created messages. If you see similar logs, we have successfully verified that our Kafka
setup is working.
```bash
kubectl logs deployment/kafka-amqstreams-consumer
```
2024/01/30 10:09:43 Go consumer starting, connecting to Kafka Server:
bootstrapServer=my-cluster-kafka-bootstrap.kafka.svc:9092, topic=my-topic,
group=my-group, sasl=false
2024/01/30 10:09:53 Consumer group handler setup
2024/01/30 10:09:53 Sarama consumer up and running!...
2024/01/30 10:11:04 Message received: value=Hello from Go Kafka Sarama-0,
topic=my-topic, partition=1, offset=0
2024/01/30 10:11:04 Message received: value=Hello from Go Kafka Sarama-1,
topic=my-topic, partition=0, offset=3
2024/01/30 10:11:04 Message received: value=Hello from Go Kafka Sarama-2,
topic=my-topic, partition=4, offset=2
2024/01/30 10:11:04 Message received: value=Hello from Go Kafka Sarama-3,
topic=my-topic, partition=3, offset=3
2024/01/30 10:11:04 Message received: value=Hello from Go Kafka Sarama-4,
topic=my-topic, partition=0, offset=4
2024/01/30 10:11:04 Message received: value=Hello from Go Kafka Sarama-5,
topic=my-topic, partition=0, offset=5
2024/01/30 10:11:04 Message received: value=Hello from Go Kafka Sarama-6,
topic=my-topic, partition=1, offset=1
2024/01/30 10:11:05 Message received: value=Hello from Go Kafka Sarama-7,
topic=my-topic, partition=0, offset=6

## Summary

Now that you have completed this exercise, Kafka Cluster is set up in Kubernetes using the Strimzi operator,
which involved creating a Kafka namespace, installing the operator, and deploying a Kafka cluster. Then we
created a Kafka topic named 'my-topic' for message storage, followed by the deployment of consumer and
producer applications to validate the Kafka setup by reading and writing messages to the topic.