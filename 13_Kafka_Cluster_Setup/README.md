# Lab Exercise 6.1: Kafka Cluster Setup

In this exercise, we will install a Kafka Cluster in Kubernetes and verify it is working by deploying a producer and a consumer application.

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
We will be using the Strimzi kubernetes operator to create and manage our Kafka cluster. Use the command below to install the operator.
```bash
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
```

3. Verify operator:
Execute the following command and ensure that the operator is in running state, the READY column should show 1/1.
```bash
kubectl get deployments strimzi-cluster-operator -n kafka
```
```text
NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
strimzi-cluster-operator   1/1     1            1           12m
```

4. Create a Kafka cluster.
Once the operator is in running state. The configuration below utilizes the Kafka CRD of the Strimzi operator to create a cluster called `my-cluster` along with a `KafkaNodePool` using KRaft mode.
Create a file named `kafka-cluster.yaml` with the following contents and apply it using the command below.
```yaml
apiVersion: kafka.strimzi.io/v1
kind: KafkaNodePool
metadata:
  name: pool-a
  namespace: kafka
  labels:
    strimzi.io/cluster: my-cluster
spec:
  replicas: 3
  roles:
    - controller
    - broker
  storage:
    type: jbod
    volumes:
      - id: 0
        type: ephemeral
---
apiVersion: kafka.strimzi.io/v1
kind: Kafka
metadata:
  name: my-cluster
  namespace: kafka
  annotations:
    strimzi.io/kraft: enabled
    strimzi.io/node-pools: enabled
spec:
  kafka:
    version: 4.2.0
    config:
      offsets.topic.replication.factor: 3
      transaction.state.log.replication.factor: 3
      transaction.state.log.min.isr: 2
      default.replication.factor: 3
      min.insync.replicas: 2
    listeners:
      - name: plain
        port: 9092
        type: internal
        tls: false
      - name: tls
        port: 9093
        type: internal
        tls: true
  entityOperator:
    topicOperator: {}
    userOperator: {}
```
```bash
kubectl apply -f kafka-cluster.yaml
```

5. Verify that the Kafka cluster is up and running.
Execute the command below and ensure that the READY in the output column is True.
```bash
kubectl get kafka.kafka.strimzi.io/my-cluster -n kafka
```
```text
NAME         DESIRED KAFKA REPLICAS   DESIRED ZK REPLICAS   READY
my-cluster   3                        3                     True
```

6. Create Kafka topics.
The configuration below creates a Kafka topic called `my-topic` by utilizing the `KafkaTopic` CRD provided by the Strimzi operator.
Create a file named `topic.yaml` with the following contents and apply it using the command provided.
```yaml
apiVersion: kafka.strimzi.io/v1
kind: KafkaTopic
metadata:
  name: my-topic
  namespace: kafka
  labels:
    strimzi.io/cluster: my-cluster
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

7. Deploy Python Scripts ConfigMap.
To run the producer and consumer applications using Python, create a ConfigMap named `kafka-app-scripts` that contains the source code for producing and consuming messages.
Create a file named `kafka-app-configmap.yaml` with the following contents and apply it using the command below.
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: kafka-app-scripts
  namespace: kafka
data:
  producer.py: |
    import os
    import time
    import sys
    from kafka import KafkaProducer

    bootstrap_servers = os.environ.get('BOOTSTRAP_SERVERS', 'my-cluster-kafka-bootstrap.kafka.svc:9092')
    topic = os.environ.get('TOPIC', 'my-topic')
    message_count = int(os.environ.get('MESSAGE_COUNT', '20'))
    delay_ms = int(os.environ.get('DELAY_MS', '100'))

    print(f"Connecting to Kafka server: {bootstrap_servers}, topic: {topic}")
    sys.stdout.flush()

    for i in range(15):
        try:
            producer = KafkaProducer(bootstrap_servers=bootstrap_servers)
            break
        except Exception as e:
            print(f"Failed to connect to Kafka (attempt {i+1}/15): {e}")
            sys.stdout.flush()
            time.sleep(5)
    else:
        print("Could not connect to Kafka after multiple attempts.")
        sys.exit(1)

    print("Connected to Kafka successfully.")
    sys.stdout.flush()

    for i in range(message_count):
        msg = f"Hello from Python Kafka Producer-{i}"
        producer.send(topic, value=msg.encode('utf-8'))
        print(f"Message sent: {msg}")
        sys.stdout.flush()
        if delay_ms > 0:
            time.sleep(delay_ms / 1000.0)

    producer.flush()
    producer.close()
    print("Finished producing messages.")
    sys.stdout.flush()

  consumer.py: |
    import os
    import sys
    import time
    from kafka import KafkaConsumer

    bootstrap_servers = os.environ.get('BOOTSTRAP_SERVERS', 'my-cluster-kafka-bootstrap.kafka.svc:9092')
    topic = os.environ.get('TOPIC', 'my-topic')
    group = os.environ.get('GROUP', 'my-group')

    print(f"Python consumer starting, connecting to Kafka Server: bootstrapServer={bootstrap_servers}, topic={topic}, group={group}")
    sys.stdout.flush()

    for i in range(15):
        try:
            consumer = KafkaConsumer(
                topic,
                bootstrap_servers=bootstrap_servers,
                group_id=group,
                auto_offset_reset='latest'
            )
            break
        except Exception as e:
            print(f"Failed to connect to Kafka (attempt {i+1}/15): {e}")
            sys.stdout.flush()
            time.sleep(5)
    else:
        print("Could not connect to Kafka after multiple attempts.")
        sys.exit(1)

    print("Sarama consumer up and running!...")
    sys.stdout.flush()

    try:
        for message in consumer:
            val = message.value.decode('utf-8')
            print(f"Message received: value={val}, topic={message.topic}, partition={message.partition}, offset={message.offset}")
            sys.stdout.flush()
    except KeyboardInterrupt:
        pass
    finally:
        consumer.close()
```
```bash
kubectl apply -f kafka-app-configmap.yaml
```

8. Create Kafka consumers.
The below configuration creates a Kubernetes deployment using a Python container image to run the consumer script. This application is configured to read Kafka messages from the topic `my-topic`.
Create a file named `consumer.yaml` with the contents below and apply it using the following command.
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kafka-amqstreams-consumer
  namespace: kafka
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
        image: python:3.11-slim
        imagePullPolicy: IfNotPresent
        command:
        - /bin/sh
        - -c
        - |
          pip install --user --no-cache-dir kafka-python-ng
          python /app/consumer.py
        securityContext:
          allowPrivilegeEscalation: false
          runAsNonRoot: true
          runAsUser: 1000
          runAsGroup: 1000
          capabilities:
            drop:
            - ALL
          seccompProfile:
            type: RuntimeDefault
        env:
        - name: BOOTSTRAP_SERVERS
          value: my-cluster-kafka-bootstrap.kafka.svc:9092
        - name: TOPIC
          value: my-topic
        - name: GROUP
          value: my-group
        - name: HOME
          value: /tmp
        resources:
          requests:
            cpu: 100m
            memory: 100Mi
          limits:
            cpu: 500m
            memory: 500Mi
        volumeMounts:
        - name: script-vol
          mountPath: /app
      volumes:
      - name: script-vol
        configMap:
          name: kafka-app-scripts
```
```bash
kubectl apply -f consumer.yaml
```

9. Create Kafka Producer.
The configuration below creates a Kubernetes Job using a Python container image to produce 20 messages into Kafka.
Create a file named `producer.yaml` with the following contents and apply it using the command below.
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
        image: python:3.11-slim
        imagePullPolicy: IfNotPresent
        command:
        - /bin/sh
        - -c
        - |
          pip install --user --no-cache-dir kafka-python-ng
          python /app/producer.py
        securityContext:
          allowPrivilegeEscalation: false
          runAsNonRoot: true
          runAsUser: 1000
          runAsGroup: 1000
          capabilities:
            drop:
            - ALL
          seccompProfile:
            type: RuntimeDefault
        env:
        - name: BOOTSTRAP_SERVERS
          value: my-cluster-kafka-bootstrap.kafka.svc:9092
        - name: TOPIC
          value: my-topic
        - name: MESSAGE_COUNT
          value: "20"
        - name: DELAY_MS
          value: "100"
        - name: HOME
          value: /tmp
        volumeMounts:
        - name: script-vol
          mountPath: /app
      volumes:
      - name: script-vol
        configMap:
          name: kafka-app-scripts
```
```bash
kubectl create -f producer.yaml
```

10. Verify logs on consumer.
Once messages are created by the producer, you can execute the command below to check if the consumer is reading the newly created messages.
```bash
kubectl logs deployment/kafka-amqstreams-consumer -n kafka
```
```text
Python consumer starting, connecting to Kafka Server: bootstrapServer=my-cluster-kafka-bootstrap.kafka.svc:9092, topic=my-topic, group=my-group
Sarama consumer up and running!...
Message received: value=Hello from Python Kafka Producer-0, topic=my-topic, partition=4, offset=0
Message received: value=Hello from Python Kafka Producer-1, topic=my-topic, partition=3, offset=0
Message received: value=Hello from Python Kafka Producer-2, topic=my-topic, partition=1, offset=0
Message received: value=Hello from Python Kafka Producer-3, topic=my-topic, partition=0, offset=0
...
Message received: value=Hello from Python Kafka Producer-19, topic=my-topic, partition=3, offset=7
```

## Summary

Now that you have completed this exercise, Kafka Cluster is set up in Kubernetes using the Strimzi operator, which involved creating a Kafka namespace, installing the operator, and deploying a Kafka cluster. Then we created a Kafka topic named 'my-topic' for message storage, followed by the deployment of consumer and producer applications to validate the Kafka setup by reading and writing messages to the topic.