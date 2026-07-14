# Lab Exercise 7.4 Pausing and Resuming Scaling Operations


# Lab Exercise 7.4: 

Pausing and Resuming Scaling Operations
As discussed in Chapter 7, It can be useful to instruct KEDA to pause autoscaling of objects, if you want to do
cluster maintenance or you want to avoid resource starvation by removing non-mission-critical workloads. You
can enable this by adding the following annotation to your ScaledObject definition:
```yaml
metadata:
annotations:
autoscaling.keda.sh/paused-replicas: "0"
autoscaling.keda.sh/paused: "true"
```
The presence of these annotations will pause autoscaling, no matter what number of replicas is provided.
The annotation autoscaling.keda.sh/paused will pause scaling immediately and use the current instance
count while the annotation autoscaling.keda.sh/paused-replicas: "<number>" will scale your
current workload to a specified amount of replicas and pause autoscaling. You can set the value of replicas for
an object to be paused to any arbitrary number.
In this exercise, we will observe how KEDA reacts when you pause or resume autoscaling. We will be reusing
the ScaledJob and Producer application created in our previous exercise (Lab Exercise 7.3).

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Familiarity with RabbitMQ.
3. Access to a Kubernetes environment with KEDA Metric Server installed as per Lab 5.
4. Completion of Lab Exercises 7.1, 7.2 and 7.3.

## Lab Exercise

1. Pausing autoscaling:
As you observed in exercise 7.1, KEDA created 15 pods corresponding to 15 messages in KEDA. Now let’s
add a pause annotation using the following command and observe how KEDA reacts.
```bash
kubectl annotate scaledjob rabbitmq-scaledjob autoscaling.keda.sh/paused="true"
```
--overwrite
2. Generate messages in RabbitMQ:
```bash
kubectl create -f rabbitmq-producer.yaml
```
3. Observe autoscaling behavior:
Because we have paused autoscaling, if you try to list the pods created by ScaledJob you will get an empty
response.
```bash
kubectl get pods -l=scaledjob.keda.sh/name=rabbitmq-scaledjob
```
No resources found in default namespace.
Also, you can observe from the following output that KEDA ScaledJob has been paused.
```bash
kubectl get scaledjobs.keda.sh rabbitmq-scaledjob
```
NAME MIN MAX TRIGGERS AUTHENTICATION READY ACTIVE
PAUSED AGE
rabbitmq-scaledjob 100 rabbitmq keda-trigger-... True False
True 10h
4. Resuming autoscaling:
For resuming autoscaling, just set the value of pause annotation to false using this command:
```bash
kubectl annotate scaledjob rabbitmq-scaledjob autoscaling.keda.sh/paused="false"
```
--overwrite
You can verify if autoscaling operation has been resumed if the Paused column is false as shown below.
```bash
kubectl get scaledjobs.keda.sh rabbitmq-scaledjob
```
NAME MIN MAX TRIGGERS AUTHENTICATION READY ACTIVE
PAUSED AGE
rabbitmq-scaledjob 100 rabbitmq keda-trigger-... True False
False 10h

## Summary

In this exercise we learned how to pause and resume KEDA's auto scaling operations, which is useful for
cluster maintenance or managing resource allocation. This functionality provides fine-grained control over
scaling behavior in dynamic environments.

## Clean Up

1. kubectl delete jobs --all --wait
2. kubectl delete -f rabbitmq_cluster.yaml
3. kubectl delete -f rabbitmq-creds-secret.yaml
4. kubectl delete -f rabbitmq-consumer-script.yaml
5. kubectl delete -f rabbitmq-producer.yaml
6. kubectl delete -f scaled-job.yaml