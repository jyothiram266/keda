# Lab Exercise 7.2 Implementing Different Rollout Strategies


# Lab Exercise 7.2: 

Implementing Different Rollout Strategies
In this exercise, we will explore the nuances of managing and updating ScaledJob resources in Kubernetes
with KEDA, focusing on the impact of different rollout strategies on running jobs. The activity begins with
observing the default strategy's effect on pods when a ScaledJob configuration is modified, leading to the
immediate termination of existing pods. It then transitions to implementing a gradual rollout strategy,
demonstrating a methodical approach to applying configuration changes that minimizes disruption and allows
for seamless continuation of workload processing.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Familiarity with RabbitMQ.
3. Access to a Kubernetes environment with KEDA Metric Server installed as per Lab 5.
4. Completion of Lab Exercises 7.1.

## Lab Exercise

1. Observing the default strategy:
We will use the ScaledJob created in the previous exercise. Similar to the previous exercise (Lab Exercise
7.1), create messages in RabbitMQ using the following command. KEDA will respond accordingly by creating
15 pods.
```bash
kubectl create -f rabbitmq-producer.yaml
```
Wait until all the 15 pods are in Running state
Now let’s observe what happens when you modify the ScaledJob, which has 15 pods currently running. Run
the command below in a separate terminal for observing how pods behave when we apply changes to the
ScaledJob.
watch "kubectl get pods --no-headers -l
scaledjob.keda.sh/name=rabbitmq-scaledjob"
Make changes to our previous ScaledJob by changing maxReplicacount from 100 to 50 (or you can modify
the pod spec template as well) in the scaled-job.yaml file and apply it using the command below.
maxReplicaCount: 50
```bash
kubectl apply -f scaled-job.yaml
```
When you execute the command above, you will observe in the other terminal that all pods are getting
terminated as shown in the output below. This is due to the default rolling strategy.
```text
NAME READY STATUS RESTARTS
rabbitmq-scaledjob-6rtpd-x6vvv 0/1 Terminating 0
rabbitmq-scaledjob-75sqz-29zvf 0/1 Terminating 0
rabbitmq-scaledjob-cqzq5-q8d4n 0/1 Terminating 0
rabbitmq-scaledjob-czrwv-wsw7x 0/1 Terminating 0
rabbitmq-scaledjob-vftn8-h22q9 1/1 Terminating 0
rabbitmq-scaledjob-z8zt9-j6s6q 1/1 Terminating 0
```
2. Implementing a gradual strategy:
To configure a gradual rollout strategy in our ScaledJob, just copy the scaled-job.yaml file and create a
new file named scaled-job-gradual-rollout.yaml and add the contents below to the spec section of
the ScaledJob.
rollout:
strategy: gradual
Execute the following command to apply changes
```bash
kubectl apply -f scaled-job-gradual-rollout.yaml
```
3. Observing the gradual rollout strategy:
Generate messages in RabbitMQ using the following command.
```bash
kubectl create -f rabbitmq-producer.yaml
```
Wait until all the 15 pods are in Running state.
Now let’s observe what happens when you modify the ScaledJob, which has 15 pods currently running. Run
the command below in a separate terminal for observing how pods behave when we apply changes to the
ScaledJob.
watch "kubectl get pods --no-headers -l
scaledjob.keda.sh/name=rabbitmq-scaledjob"
Make changes to our previous ScaledJob by changing maxReplicacount to 70 in the
scaled-job-gradual-rollout.yaml file and apply it using the following command.
```bash
kubectl apply -f scaled-job-gradual-rollout.yaml
```
Compared to the previous strategy, you will observe that no pods were killed. The 15 pods will be terminated
only after complete execution. This result shows that gradual rollout strategy applies the changes only to the
newly created pods.

## Summary

This exercise demonstrated the implementation of different rollout strategies for ScaledJob resources managed
by KEDA in a Kubernetes environment. Initially, modifying a ScaledJob configuration with the default strategy
led to the immediate termination of all running pods, illustrating the default rolling update behavior.
Subsequently, by adopting a gradual rollout strategy, changes to the ScaledJob were applied in a less
disruptive manner, allowing existing pods to complete their execution before termination, which showcased a
more controlled approach to updating job configurations in response to workload changes.