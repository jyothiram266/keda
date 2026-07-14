# Lab Exercise 9.1 Visualize KEDA AutoScaling Metrics in Grafana


# Lab Exercise 9.1: Visualize KEDA

AutoScaling Metrics In Grafana
In this exercise, you will learn how to visualize and interact with KEDA metrics exported in Prometheus format,
offering a practical experience in monitoring and managing application scalability.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Familiarity with Prometheus and Grafana.
3. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.

## Lab Environment Setup

We will install RabbitMQ, Prometheus and Grafana in the Kubernetes cluster using their respective operator:
1. Install Operator:
The Prometheus Operator provides Kubernetes native deployment and management of Prometheus and
related monitoring components. The setup is simplified by using Helm, a package manager for Kubernetes.
# Add the Prometheus Helm chart repository
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
```
# Update the repository
```bash
helm repo update
```
# Install the Prometheus stack
```bash
helm upgrade -i prometheus-stack prometheus-community/kube-prometheus-stack -n monitoring --create-namespace
```
2. Verify Operator:
```bash
kubectl get pods -n monitoring
```
```text
NAME READY STATUS
RESTARTS
alertmanager-prometheus-stack-kube-prom-alertmanager-0 2/2 Running
0
prometheus-prometheus-stack-kube-prom-prometheus-0 2/2 Running
0
prometheus-stack-grafana-6c99cbfccb-8zc7r 3/3 Running
0
prometheus-stack-kube-prom-operator-7dfbbf8df-qtvjk 1/1 Running
0
prometheus-stack-kube-state-metrics-556d4c4c5d-wwjzk 1/1 Running
0
prometheus-stack-prometheus-node-exporter-xnzgk 1/1 Running
0
```
3. Expose KEDA metrics in Prometheus format:
The following command configures KEDA to export metrics in Prometheus format.
```bash
helm upgrade -i keda kedacore/keda --namespace keda --create-namespace --set
```
prometheus.metricServer.enabled=true --set prometheus.operator.enabled=true --set
prometheus.webhooks.enabled=true
4. Monitor KEDA services using Pod Monitor:
Create a Pod Monitor resource such that prometheus can start collecting metrics from KEDA pods.
Create a file named pod-monitors.yaml with the following contents and apply it using the command below.
```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
name: keda-operator
namespace: monitoring # Namespace where Prometheus Operator is installed
labels:
release: prometheus-stack # Adjust to match the label of your Prometheus
```
instance
```yaml
spec:
selector:
matchLabels:
app: keda-operator
name: keda-operator
namespaceSelector:
matchNames:
- keda

podMetricsEndpoints:
- port: metrics
interval: 15s
path: /metrics
relabelings:
- action: labelmap
regex: __meta_kubernetes_pod_label_(.+)
---
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
name: keda-admission-webhooks
namespace: monitoring # Namespace where Prometheus Operator is installed
labels:
release: prometheus-stack # Adjust to match the label of your Prometheus
```
instance
```yaml
spec:
selector:
matchLabels:
app: keda-admission-webhooks
name: keda-admission-webhooks
namespaceSelector:
matchNames:
- keda
podMetricsEndpoints:
- port: metrics
interval: 15s
path: /metrics
relabelings:
- action: labelmap
regex: __meta_kubernetes_pod_label_(.+)
---
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
name: keda-operator-metrics-apiserver
namespace: monitoring # Namespace where Prometheus Operator is installed
labels:
release: prometheus-stack # Adjust to match the label of your Prometheus
```
instance
```yaml
spec:
selector:
matchLabels:
app: keda-operator-metrics-apiserver
namespaceSelector:

matchNames:
- keda
podMetricsEndpoints:
- port: metrics
interval: 15s
path: /metrics
relabelings:
- action: labelmap
regex: __meta_kubernetes_pod_label_(.+)
```
```bash
kubectl apply -f pod-monitors.yaml
```
5. Access Prometheus UI locally and verify KEDA targets:
To access the Prometheus UI locally, you can use kubectl port-forward to forward the Prometheus service port
to your local machine.
```bash
kubectl port-forward service/prometheus-stack-kube-prom-prometheus -n monitoring
```
9090:9090
Once the port forwarding is set up, you can access the Prometheus UI by navigating to http://localhost:9090 in
your web browser.
To verify KEDA targets, open the Prometheus targets page by accessing this URL
http://localhost:9090/targets?search=keda
Prometheus Targets
The three targets shown in the image above should be visible.
6. Access Grafana UI locally:
To access the Grafana UI locally, you can use kubectl port-forward to forward the Grafana service port to your
local machine.
```bash
kubectl port-forward service/prometheus-stack-grafana -n monitoring 3000:80
```
Once the port forwarding is set up, you can access the Prometheus UI by navigating to http://localhost:3000 in
your web browser.
Authentication Credentials
○ Username: admin
○ Password: prom-operator
7. Import PreMade KEDA Dashboard in Grafana:
Copy the dashboard JSON file available at Github.
Open the Grafana import page by accessing this URL: http://localhost:3000/dashboard/import
Paste the content in the UI model named Import via dashboard JSON model and click on the Import button
to import the KEDA dashboard.
8. Install RabbitMQ as defined in the environment setup section of Lab Exercise 8.1, from Step 1 - Step 7.
Congratulations, you have successfully set up the required lab environment.

## Lab Exercise

1. Create ScaledObject:
The following configuration creates a ScaledObject resource that targets the consumer-program deployment
and uses the RabbitMQ scaler to make scaling decisions. The ScaledObject will scale the number of replicas
between 1 to 10 based on the queueLength of the queue called “testqueue”.
Create a file named scaledobject.yaml with the following contents and apply it using the command below.
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
kind: ScaledObject
metadata:
name: keda-rabbitmq
namespace: default
spec:
scaleTargetRef:
apiVersion: apps/v1
kind: Deployment
name: consumer-program
minReplicaCount: 1
maxReplicaCount: 10
triggers:
- type: rabbitmq
metadata:
protocol: amqp
queueName: testqueue
queueLength: "5"
authenticationRef:
name: keda-trigger-auth-rabbitmq-conn
```
```bash
kubectl apply -f scaledobject.yaml
```
2. Verify ScaledObject:
Execute the following command and ensure the READY column shows True.
```bash
kubectl get scaledobjects.keda.sh keda-rabbitmq
```
NAME SCALETARGETKIND SCALETARGETNAME MIN MAX TRIGGERS
AUTHENTICATION READY ACTIVE FALLBACK PAUSED AGE
keda-rabbitmq apps/v1.Deployment consumer-program 1 10 rabbitmq
keda-trigger-auth-rabbitmq-conn True False False Unknown 9h
3. Create Messages in RabbitMQ:
Execute the command below to create 20 messages in RabbitMQ.
sed 's/value: "1"/value: "20"/' producer.yaml | kubectl create -f -
4. (Optional) Observe AutoScaling behavior:
You can optionally use the following command to watch HPA changes as done in the previous labs. However,
for this exercise we will observe these changes in the Grafana dashboard.
```bash
kubectl get hpa keda-hpa-keda-rabbitmq --watch
```
5. Observe AutoScaling behavior in the Grafana Dashboard:
Open the grafana KEDA dashboard that we imported during the lab setup. From the drop down menus select
the following:
- namespace as default,
- ScaledObject as keda-rabbitmq or all
- scaler as rabbitMQScaler
- metric as s0-rabbitmq-testqueue
- time range to Last 5 minutes.
In the following image, you can see as we generated messages in RabbitMQ cluster, the metric value
s0-rabbitmq-testqueue (inside Scale Target Grafana Panel) started increasing.
Auto Scaling Behavior In Grafana Dashboard
In the image below, you can see the response to the increased scale target value. KEDA begins to increase
pod replicas.
Auto Scaling Behavior In Grafana Dashboard
The image below provides another view of the same actions over a period of 15 minutes.
Auto Scaling Behavior In Grafana Dashboard

## Summary

This exercise demonstrates configuring a ScaledObject in KEDA for autoscaling a Kubernetes deployment
based on RabbitMQ queue length, and visualizing the auto-scaling metrics such as Scale Target and Current
Replicas in Grafana.