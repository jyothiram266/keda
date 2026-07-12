# Lab Exercise 6.3 Advance Configuration


# Lab Exercise 6.3: Advance Configuration

In this exercise, we will explore advanced KEDA ScaledObject configuration. Participants will learn how to
customize the Horizontal Pod Autoscaler (HPA), manage scaling activities during various operational states,
and fine-tune performance parameters. Through hands-on tasks, we'll observe the effects of renaming HPA,
pausing scaling, adjusting cooldown periods, and managing idle replicas on the dynamic scaling behavior of a
Kafka consumer application.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Familiarity with Kafka.
3. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.
4. Completion of Lab Exercises 6.1 and 6.2.

## Lab Exercise

1. Renaming HPA:
By default when you create a ScaledObject, Keda creates and manages a corresponding HPA resource and
the default name for this HPA is keda-hpa-{scaled-object-name}.
KEDA gives you the ability to modify the HPA resource created by KEDA using the advanced configuration in
the ScaledObject specification (highlighted in the below YAML)
```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
name: kafka-amqstreams-consumer-scaledobject
spec:
minReplicaCount: 1
maxReplicaCount: 5
scaleTargetRef:
name: kafka-amqstreams-consumer

triggers:
- type: apache-kafka
metadata:
topic: my-topic
bootstrapServers: my-cluster-kafka-bootstrap.kafka.svc:9092
consumerGroup: my-group
lagThreshold: "1"
offsetResetPolicy: "latest"
advanced:
horizontalPodAutoscalerConfig:
name: kafka-amqstreams-consumer
```
Once you apply the above yaml using kubectl apply the HPA resource
keda-hpa-kafka-amqstreams-consumer-scaledobject will be renamed to
kafka-amqstreams-consumer.
```bash
kubectl get hpa kafka-amqstreams-consumer
```
To resume autoscaling, just remove the above annotations. To know more about HPA advance configuration
refer this documentation.
2. Pause scaling:
It can be useful to instruct KEDA to pause autoscaling of objects if you want to do cluster maintenance or you
want to avoid resource starvation by removing non-mission-critical workloads. You can enable this by adding
the below annotation to your ScaledObject definition:
```yaml
metadata:
annotations:
autoscaling.keda.sh/paused-replicas: "0"
autoscaling.keda.sh/paused: "true"
```
The presence of these annotations will pause autoscaling no matter what number of replicas is provided.
The annotation autoscaling.keda.sh/paused will pause scaling immediately and use the current instance
count, while the annotation autoscaling.keda.sh/paused-replicas: "<number>" will scale your
current workload to specified amount of replicas and pause autoscaling. You can set the value of replicas for
an object to be paused to any arbitrary number.
Below, you will find the updated ScaledObject. Apply it using:
```bash
kubectl apply -f scaledobject.yaml
```
```yaml
apiVersion: keda.sh/v1alpha1

kind: ScaledObject
metadata:
name: kafka-amqstreams-consumer-scaledobject
annotations:
autoscaling.keda.sh/paused: "true"
spec:
minReplicaCount: 1
maxReplicaCount: 5
cooldownPeriod: 5
scaleTargetRef:
name: kafka-amqstreams-consumer
triggers:
- type: apache-kafka
metadata:
topic: my-topic
bootstrapServers: my-cluster-kafka-bootstrap.kafka.svc:9092
consumerGroup: my-group
lagThreshold: "1"
offsetResetPolicy: "latest"
```
After creating the ScaledObject, you will observe that HPA resource
keda-hpa-kafka-amqstreams-consumer-scaledobject which is managed by KEDA, is deleted. Once
the HPA is removed, autoscaling is stopped.
You can verify the existence of HPA with the following command.
```bash
kubectl get hpa keda-hpa-kafka-amqstreams-consumer-scaledobject --watch
```
Error from server (NotFound): horizontalpodautoscalers.autoscaling
"keda-hpa-kafka-amqstreams-consumer-scaledobject" not found
3. Modify and observe cooldown periods.
In our previous Lab Exercise 6.2, we observed that HPA took 5 mins to scale to zero after reaching replica
count of 1. This happens because 5 mins is the default value of coolDownPeriod. We will modify this config
and observe what is the impact of it on scale to zero.
Create a file name scaled-object-cooldown-period.yaml with the following contents and apply it using
the command below.
```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
name: kafka-amqstreams-consumer-scaledobject
spec:

minReplicaCount: 1
maxReplicaCount: 5
cooldownPeriod: 60
scaleTargetRef:
name: kafka-amqstreams-consumer
triggers:
- type: apache-kafka
metadata:
topic: my-topic
bootstrapServers: my-cluster-kafka-bootstrap.kafka.svc:9092
consumerGroup: my-group
lagThreshold: "1"
offsetResetPolicy: "latest"
```
```bash
kubectl apply -f scaled-object-cooldown-period.yaml
```
Once the ScaledObject is created, generate Kafka messages using command
```bash
kubectl create producer.yaml.
```
Let’s monitor the scaling behavior as done previously, using the following command. Run these in two separate
terminal instances.
```bash
kubectl get hpa keda-hpa-kafka-amqstreams-consumer-scaledobject --watch
kubectl get events --watch --field-selector
```
involvedObject.kind=HorizontalPodAutoscaler,involvedObject.name=keda-hpa-kafka-am
qstreams-consumer-scaledobject
NAME REFERENCE TARGETS MINPODS MAXPODS
REPLICAS AGE
keda-hpa-kafka... kafka-amqstreams... <unknown>/1 (avg) 1 5
0 5s
keda-hpa-kafka... kafka-amqstreams... <unknown>/1 (avg) 1 5
0 47s
keda-hpa-kafka... kafka-amqstreams... 5/1 (avg) 1 5
1 5m47
0s Normal SuccessfulRescale
horizontalpodautoscaler/keda-hpa-kafka-amqstreams-consumer-scaledobject New
size: 4; reason: external metric s0-kafka-my-topic above target
keda-hpa-kafka... kafka-amqstreams... 1250m/1 (avg) 1 5
4 6m47
0s Normal SuccessfulRescale
horizontalpodautoscaler/keda-hpa-kafka-amqstreams-consumer-scaledobject New
size: 5; reason: external metric s0-kafka-my-topic above target
keda-hpa-kafka... kafka-amqstreams... 800m/1 (avg) 1 5
5 7m47
keda-hpa-kafka... kafka-amqstreams... 0/1 (avg) 1 5
5 8m47
keda-hpa-kafka... kafka-amqstreams... 0/1 (avg) 1 5
5 11m
0s Normal SuccessfulRescale
horizontalpodautoscaler/keda-hpa-kafka-amqstreams-consumer-scaledobject New
size: 4; reason: All metrics below target
keda-hpa-kafka... kafka-amqstreams... 0/1 (avg) 1 5
4 12m
0s Normal SuccessfulRescale
horizontalpodautoscaler/keda-hpa-kafka-amqstreams-consumer-scaledobject New
size: 1; reason: All metrics below target
keda-hpa-kafka... kafka-amqstreams... 0/1 (avg) 1 5
1 13m
keda-hpa-kafka... kafka-amqstreams... <unknown>/1 (avg) 1 5
0 24m
The output above is the combination of the preceding commands. The events are chronologically ordered as
they happened.
As you can see, the final scaling event, which adjusted the replica count to 1, occurred at the 13-minute mark.
Following that, adhering to the cooldown period configured for 10 minutes, the scaling was further reduced to
zero.
4. Managing idle replicas:
The idleReplicaCount in a KEDA ScaledObject specifies the minimum number of replicas that should remain
running even when there's no workload to process, effectively setting a baseline for idle conditions. This
parameter is particularly useful for ensuring that a certain number of pods are always ready to quickly respond
to sudden spikes in workload without the initial delay of scaling from zero, thereby improving the
responsiveness and resilience of the system.
Create a file name scaled-object-ideal-replica-count.yaml with the following contents and apply it
using the command below.
```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
name: kafka-amqstreams-consumer-scaledobject
spec:
minReplicaCount: 1
maxReplicaCount: 5
idleReplicaCount: 0
scaleTargetRef:
name: kafka-amqstreams-consumer
triggers:
- type: apache-kafka
metadata:
topic: my-topic
bootstrapServers: my-cluster-kafka-bootstrap.kafka.svc:9092
consumerGroup: my-group
lagThreshold: "1"
offsetResetPolicy: "latest"
```
```bash
kubectl apply -f scaled-object-ideal-replica-count.yaml
```
Once the ScaledObject is created, generate Kafka messages using the command
```bash
kubectl create producer.yaml
```
Let’s monitor the scaling behavior as done previously, using the following command. Run these in two separate
terminal instances.
```bash
kubectl get hpa keda-hpa-kafka-amqstreams-consumer-scaledobject --watch
kubectl get events --watch --field-selector
```
involvedObject.kind=HorizontalPodAutoscaler,involvedObject.name=keda-hpa-kafka-am
qstreams-consumer-scaledobject
NAME REFERENCE TARGETS MINPODS MAXPODS
REPLICAS AGE
keda-hpa-kafka... kafka-amqstreams... <unknown>/1 (avg) 1 5
0 14s
keda-hpa-kafka... kafka-amqstreams... <unknown>/1 (avg) 1 5
0 60s
0s Normal SuccessfulRescale
horizontalpodautoscaler/keda-hpa-kafka-amqstreams-consumer-scaledobject New
size: 1; reason: Current number of replicas below Spec.MinReplicas
keda-hpa-kafka... kafka-amqstreams... 1500m/1 (avg) 1 5
1 2m
0s Normal SuccessfulRescale
horizontalpodautoscaler/keda-hpa-kafka-amqstreams-consumer-scaledobject New
size: 3; reason: external metric s0-kafka-my-topic above target
keda-hpa-kafka... kafka-amqstreams... 1667m/1 (avg) 1 5
3 3m1s
1s Normal SuccessfulRescale
horizontalpodautoscaler/keda-hpa-kafka-amqstreams-consumer-scaledobject New
size: 5; reason: external metric s0-kafka-my-topic above target
keda-hpa-kafka... kafka-amqstreams... 0/1 (avg) 1 5
5 4m1s
0s Normal SuccessfulRescale
horizontalpodautoscaler/keda-hpa-kafka-amqstreams-consumer-scaledobject New
size: 2; reason: All metrics below target
keda-hpa-kafka... kafka-amqstreams... 0/1 (avg) 1 5
2 8m1s
keda-hpa-kafka... kafka-amqstreams... 0/1 (avg) 1 5
1 9m1s
keda-hpa-kafka... kafka-amqstreams... <unknown>/1 (avg) 1 5
0 13m
The above output is the combination of the above commands, the events are chronologically ordered as they
happened.
As you can see, it's evident that at the initial mark of 60 seconds when the target metrics were not known, the
replica count was maintained at 0, aligning with the configured ideal replica count. As the workload began to
increase, KEDA and HPA promptly scaled the deployment up to the minimum number of replicas by the
2-minute mark. Subsequently, as the demand subsided, the system adjusted by scaling down to the minimum
replica count at the 8-minute mark. Eventually, when the workload was completely alleviated, the replica count
reverted to the ideal replica count by the 13-minute mark.

## Summary

This exercise demonstrated various advanced configurations that KEDA offers for dealing with nuanced
autoscaling use cases, such as pausing autoscaling, configuring the underlying HPA from KEDA ScaledObject
CRD, working with idleReplicas and cooldown Period.

## Clean Up

```bash
kubectl delete scaledobjects.keda.sh kafka-amqstreams-consumer-scaledobject
kubectl delete deployments.apps kafka-amqstreams-consumer
kubectl delete kafka -n kafka my-cluster
kubectl delete namespace kafka
```