# Lab Exercise 9.2 Pause and Unpause the ScaledObject


# Lab Exercise 9.2: Pause and Unpause the ScaledObject

This exercise demonstrates the effect of pausing/unpausing autoscaling activity on KEDA metrics. These
metrics can be used for writing customs auto scaling alerts.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Familiarity with Prometheus and Grafana.
3. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.
4. Completion of Lab Exercise 9.1.

## Lab Exercise

1. Pausing Autoscaling:
As discussed in a previous chapter, to pause auto scaling in KEDA you have to add the annotations below to
the ScaledObject.
annotations:
autoscaling.keda.sh/paused: "true"
We will utilize the ScaledObject created in the previous exercise (Lab Exercise 9.1) and execute the following
commands to annotate the ScaledObject with the above-mentioned annotations.
```bash
kubectl annotate scaledobject keda-rabbitmq autoscaling.keda.sh/paused="true"
```
--overwrite
2. (Optional) Verify ScaledObject:
Execute the following command and ensure the PAUSED columns shows True.
```bash
kubectl get scaledobjects.keda.sh keda-rabbitmq
```
NAME SCALETARGETKIND SCALETARGETNAME MIN MAX TRIGGERS
AUTHENTICATION READY ACTIVE FALLBACK PAUSED AGE
keda-rabbitmq apps/v1.Deployment consumer-program 1 10 rabbitmq
keda-trigger-auth-rabbitmq-conn True False False True 10h
3. Verify ScaledObject from Prometheus:
Go to the Prometheus metrics explored UI and search for this metric: keda_scaled_object_paused. You will
observe a result with the value set to one and labels containing ScaledObject: keda-rabbitmq. The value
one indicates the ScaledObject has paused autoscaling.
Prometheus: ScaledObject has Paused Autoscaling
4. Unpausing ScaledObject:
Execute the following command to unpause KEDA autoscaling.
```bash
kubectl annotate scaledobject keda-rabbitmq autoscaling.keda.sh/paused="false"
```
--overwrite
5. Verify ScaledObject:
Execute the same command on Prometheus UI as done previously. You will see a result as shown in the
image below. The value of the metric will be set to zero, indicating auto scaling is active.
Prometheus: Auto Scaling is Active

## Summary

In this exercise we learned how to pause or unpause autoscaling activity in KEDA and its effect on the metric
keda_scaled_object_paused. This metric can be further used for writing custom auto scaling alerts.