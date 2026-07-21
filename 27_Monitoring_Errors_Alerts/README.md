# Lab Exercise 9.3: Monitoring AutoScaling Errors and Alerts

In this exercise, we will learn how to identify and detect autoscaling related errors in KEDA using Prometheus and Grafana.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Familiarity with Prometheus and Grafana.
3. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.
4. Completion of Lab Exercises 9.1 and 9.2.

## Lab Exercise

### 1. Deleting the RabbitMQ Cluster

We will utilize the ScaledObject created in Exercise 9.1 of this lab. To introduce an error in the current RabbitMQ based scaler, we will delete the deployed RabbitMQ cluster.

KEDA periodically fetches the metrics (queue length) from RabbitMQ at every polling interval, so in the next polling call, KEDA will detect that there is a problem in the connectivity and KEDA will report an error in the scaler.

Execute the following command to delete the cluster.

```bash
kubectl delete -f ../25_Visualize_Metrics_Grafana/cluster.yaml
```

### 2. (Optional) Verify ScaledObject

Execute the following command. You will observe that the READY column becomes False, indicating some error in ScaledObject. You can get the detailed error message by describing the ScaledObject.

```bash
kubectl get scaledobjects.keda.sh keda-rabbitmq
```

*Expected Output:*
```text
NAME            SCALETARGETKIND      SCALETARGETNAME    MIN   MAX   READY   ACTIVE   FALLBACK   PAUSED   TRIGGERS   AUTHENTICATIONS                   AGE
keda-rabbitmq   apps/v1.Deployment   consumer-program   1     10    False   False    False      False    rabbitmq   keda-trigger-auth-rabbitmq-conn   10h
```

### 3. Observe ScaledObject errors in the Grafana Dashboard

Open the Grafana dashboard and observe the ScaledObject Errors panel (the value depends upon the `keda_scaled_object_errors` metric).

You can observe that the metric values have increased from zero, indicating an error in the ScaledObject named `keda-rabbitmq`.

*(Grafana Dashboard: Scaled Object Errors Panel)*

We can create custom alert rules using the `keda_scaled_object_errors` metric to send notifications to the SRE team whenever such autoscaling incidents happen.

## Summary

In this exercise we learned how to identify autoscaling-related errors using the KEDA metric `keda_scaled_object_errors`. This metric can be further used for writing custom auto scaling alerts.

## Clean Up

Run these commands from the `25_Visualize_Metrics_Grafana` folder to completely clean up:

```bash
kubectl delete job -l app=producer-program || true
kubectl delete -f consumer.yaml
kubectl delete -f secret.yaml
kubectl delete -f pod-monitors.yaml
kubectl delete -f scaledobject.yaml
```