# Lab Exercise 7.3 Testing Scaling Strategies


# Lab Exercise 7.3: Testing Scaling Strategies

KEDA provides different scaling strategies such as default, custom and accurate. In our previous two exercises
we had used the default strategy. In this exercise we will observe how KEDA handles custom and accurate
strategy.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Familiarity with RabbitMQ.
3. Access to a Kubernetes environment with KEDA Metric Server installed as per Lab 5.
4. Completion of Lab Exercises 7.1 and 7.2.

## Lab Exercise

1. Custom Strategy:
As discussed in chapter 7, the custom strategy allows for a more tailored approach to scaling. You can specify
parameters like customScalingQueueLengthDeduction and customScalingRunningJobPercentage
to fine-tune the scaling behavior.
Create a file name scaled-job-custom-scaling.yaml with the following contents and apply it using the
command below.
```yaml
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
pollingInterval: 10
successfulJobsHistoryLimit: 100
failedJobsHistoryLimit: 100
maxReplicaCount: 100
scalingStrategy:
strategy: "custom"
customScalingQueueLengthDeduction: 1
customScalingRunningJobPercentage: "0.5"
triggers:
- type: rabbitmq
metadata:
protocol: amqp
queueName: testqueue
mode: QueueLength
value: "1"
authenticationRef:
name: keda-trigger-auth-rabbitmq-conn
```
```bash
kubectl apply -f scaled-job-custom-scaling.yaml
```
2. Generate messages in RabbitMQ:
sed 's/value: "15"/value: "20"/' rabbitmq-producer.yaml | kubectl create -f -
3. Watch the creation of Kubernetes pods:
Execute the following commands to observe the number of pods created by ScaledJob.
watch "kubectl get pods --no-headers -l
scaledjob.keda.sh/name=rabbitmq-scaledjob"
```text
NAME READY STATUS RESTARTS AGE
rabbitmq-scaledjob-5rfxq-l4bsk 1/1 Running 0 52m
rabbitmq-scaledjob-5xkq8-th42z 1/1 Running 0 52m
rabbitmq-scaledjob-62vmg-pt9gd 1/1 Running 0 52m
rabbitmq-scaledjob-6msw7-md4m9 1/1 Running 0 52m
rabbitmq-scaledjob-7vsmx-ln2hc 1/1 Running 0 52m
rabbitmq-scaledjob-b2rcv-g9hg5 1/1 Running 0 52m
rabbitmq-scaledjob-cr6gd-tczfz 1/1 Running 0 52m
rabbitmq-scaledjob-csjxm-4vhbr 1/1 Running 0 52m
rabbitmq-scaledjob-dgpq9-fk6xt 1/1 Running 0 52m
rabbitmq-scaledjob-dkmrl-66xgk 1/1 Running 0 52m
rabbitmq-scaledjob-dq9qk-gztbp 1/1 Running 0 52m
rabbitmq-scaledjob-dwc4h-lrg8t 1/1 Running 0 52m
rabbitmq-scaledjob-grh77-x8nbs 1/1 Running 0 52m
rabbitmq-scaledjob-gvqlr-m9l76 1/1 Running 0 52m
rabbitmq-scaledjob-nxf6s-7kbwm 1/1 Running 0 52m
rabbitmq-scaledjob-pdsn5-fdf4m 1/1 Running 0 52m
rabbitmq-scaledjob-q82wr-lznwx 1/1 Running 0 52m
rabbitmq-scaledjob-qj8wd-vnrrh 1/1 Running 0 52m
rabbitmq-scaledjob-qm9hh-ccnz7 1/1 Running 0 52m
rabbitmq-scaledjob-rrfln-tmp7j 1/1 Running 0 52m
rabbitmq-scaledjob-ssbhv-2rz6t 1/1 Running 0 52m
rabbitmq-scaledjob-wq5c7-djthc 1/1 Running 0 52m
rabbitmq-scaledjob-zqfz4-l77tp 1/1 Running 0 52m
As discussed in the chapter, below are the formulas used by KEDA to calculate the number of pods to be
created at every polling interval.
MaxScale:
maxScale = min(scaledJob.MaxReplicaCount(), divideWithCeil(queueLength,
targetAverageValue))

Custom Strategy:
min(maxScale-int64(*s.CustomScalingQueueLengthDeduction)-int64(float64(runningJob
Count)*(*s.CustomScalingRunningJobPercentage)), maxReplicaCount)
Parameters:
CustomScalingQueueLengthDeduction: 1
CustomScalingRunningJobPercentage: "0.5"
targetAverageValue: 1
It is difficult to understand the Custom scaling behavior just by looking at the creation of pods. The content
below is taken from the KEDA operator logs, to help you understand the scaling decision taken by KEDA.
Polling Interval 1: With 19 items in the queue and no running or pending jobs, the MaxScale formula
yields 19 based on the minimum of MaxReplicaCount (configured to 100) or queueLength divided by
targetAverageValue. Applying the Custom Strategy Formula, the effective number of pods to create is
calculated as 18. KEDA then proceeds to scale up by creating 18 pods. This calculation is repeated at
subsequent polling intervals to dynamically adjust pod numbers in response to workload changes.
-------------------------
Polling Interval 1
"s0-rabbitmq-testqueue": 20
"maxValue": 20
"Number of running Jobs": 0
"Number of pending Jobs ": 0
"Effective number of max jobs": 19
"Number of jobs": 19
-------------------------
Polling Interval 2
"s0-rabbitmq-testqueue": 15
"maxValue": 15
"Number of running Jobs": 18
"Number of pending Jobs ": 14
"Effective number of max jobs": 5
"Number of jobs": 5
-------------------------
Polling Interval 3
"s0-rabbitmq-testqueue": 6,
"maxValue": 6
"Number of running Jobs": 23
"Number of pending Jobs ": 10

"Effective number of max jobs": 0
"Number of jobs": 0
-------------------------
Polling Interval 4
"s0-rabbitmq-testqueue": 0
"maxValue": 0
"Number of running Jobs": 23
"Number of pending Jobs ": 0
-------------------------
Polling Interval 5
"s0-rabbitmq-testqueue": 0
"maxValue": 0
"Number of running Jobs": 23
"Number of pending Jobs ": 0
```

## Clean Up

```bash
kubectl delete jobs --all --wait
```
4. Implement Accurate Strategy:
To configure accurate strategy in our ScaledJob, just copy the scaled-job-custom-scaling.yaml file
and create a new file named scaled-job-accurate-scaling.yaml with the scalingStrategy
section replaced with following contents and apply it using the command below.
scalingStrategy:
strategy: "accurate"
```bash
kubectl apply -f scaled-job-accurate-scaling.yaml
```
5. Generate messages In RabbitMQ
sed 's/value: "15"/value: "20"/' rabbitmq-producer.yaml | kubectl create -f -
6. Watch creation of Kubernetes Pods:
Execute the following commands to observe the number of pods created by ScaledJob.
watch "kubectl get pods --no-headers -l
scaledjob.keda.sh/name=rabbitmq-scaledjob"
```text
NAME READY STATUS RESTARTS AGE
rabbitmq-scaledjob-59bcl-6x4ss 1/1 Running 0 70s
rabbitmq-scaledjob-5k97p-89bf7 1/1 Running 0 70s
rabbitmq-scaledjob-6hfdw-ddf6d 1/1 Running 0 71s
rabbitmq-scaledjob-6zwtl-jbcsf 1/1 Running 0 70s
rabbitmq-scaledjob-8jk5z-cnc84 1/1 Running 0 70s
rabbitmq-scaledjob-bbmk5-wbg6b 1/1 Running 0 71s
rabbitmq-scaledjob-hrkqw-6gl6b 1/1 Running 0 70s
rabbitmq-scaledjob-lx5mx-glq9r 1/1 Running 0 70s
rabbitmq-scaledjob-mqb8v-d8cxp 1/1 Running 0 70s
rabbitmq-scaledjob-n5942-7fhm7 1/1 Running 0 70s
rabbitmq-scaledjob-nv9dr-xrpqq 1/1 Running 0 70s
rabbitmq-scaledjob-pjrfw-n87q2 1/1 Running 0 71s
rabbitmq-scaledjob-rgx9v-r6vvv 1/1 Running 0 71s
rabbitmq-scaledjob-w7jc8-tr47j 1/1 Running 0 70s
rabbitmq-scaledjob-wmg74-zt6pw 1/1 Running 0 71s
rabbitmq-scaledjob-xg9hn-5wzfb 1/1 Running 0 70s
rabbitmq-scaledjob-xsmh5-gzsmv 1/1 Running 0 71s
rabbitmq-scaledjob-xx7rm-bsx2m 1/1 Running 0 70s
rabbitmq-scaledjob-zn5r2-pvc5m 1/1 Running 0 70s
As discussed in the chapter, below are the formulas used by KEDA to calculate the number of pods to be
created at every polling interval for accurate strategy.
if (maxScale + runningJobCount) > maxReplicaCount {
return maxReplicaCount - runningJobCount
}
return maxScale - pendingJobCount
The content below is taken from the KEDA operator logs, to help you understand the scaling decision taken by
KEDA.
Polling Interval 1: With 19 items in the queue and no running or pending jobs, the MaxScale formula
yields 19 based on the minimum of MaxReplicaCount (configured to 100) or queueLength divided by
targetAverageValue. Applying the Accurate Strategy Formula, the effective number of pods to create is
calculated as 19. KEDA then proceeds to scale up by creating 19 pods. This calculation is repeated at
subsequent polling intervals to dynamically adjust pod numbers in response to workload changes.
-------------------------
Polling Interval 1
"s0-rabbitmq-testqueue": 20
"maxValue": 20
"Number of running Jobs": 0

"Number of pending Jobs ": 0
"Effective number of max jobs": 20
"Number of jobs": 20
-------------------------
Polling Interval 2
"s0-rabbitmq-testqueue": 12
"maxValue": 12
"Number of running Jobs": 19
"Number of pending Jobs ": 13
"Effective number of max jobs": 0
"Number of jobs": 0
-------------------------
Polling Interval 3
"s0-rabbitmq-testqueue": 2,
"maxValue": 2
"Number of running Jobs": 19
"Number of pending Jobs ": 3
"Effective number of max jobs": 0
"Number of jobs": 0
-------------------------
Polling Interval 4
"s0-rabbitmq-testqueue": 0
"maxValue": 0
"Number of running Jobs": 19
"Number of pending Jobs ": 0
-------------------------
Polling Interval 5
"s0-rabbitmq-testqueue": 0
"maxValue": 0
"Number of running Jobs": 29
"Number of pending Jobs ": 0
```

## Summary

In this exercise we explored KEDA's custom and accurate scaling strategies for managing workload in
Kubernetes. By configuring a ScaledJob with custom parameters, we observed how KEDA adjusts pod
creation in response to changes in the metric value, demonstrating control over autoscaling of ScaledJobs.