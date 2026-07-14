# Mastering Kubernetes Event-Driven Autoscaling with KEDA (LFS257)

Welcome to the hands-on lab repository for **LFS257: Mastering Kubernetes Event-Driven Autoscaling with KEDA** by The Linux Foundation.

This repository organizes all the course labs and exercises into structured directories. Each directory contains a dedicated `README.md` with cleaned instructions, code snippets, manifests, and commands.

---

## 📚 Course Labs & Exercises

| Lab Index | Exercise Title | Link |
| :---: | :--- | :---: |
| **01** | Set Up the Course Lab Environment | [View Lab Guide](./01_Setup_Lab_Env/README.md) |
| **02** | Containerize and Deploy Sample Application in Kubernetes | [View Lab Guide](./02_Deploy_App_K8s/README.md) |
| **03** | Configure and Test HPA | [View Lab Guide](./03_Configure_HPA/README.md) |
| **04** | Configure and Test VPA | [View Lab Guide](./04_Configure_VPA/README.md) |
| **05** | Containerize and Deploy Instrumented Sample Application in Kubernetes | [View Lab Guide](./05_Deploy_Instrumented_App/README.md) |
| **06** | Setting Up Prometheus in Kubernetes | [View Lab Guide](./06_Setup_Prometheus/README.md) |
| **07** | Installing and Configuring the Prometheus Adapter | [View Lab Guide](./07_Prometheus_Adapter/README.md) |
| **08** | Testing Autoscaling with Custom Metrics | [View Lab Guide](./08_Autoscaling_Custom_Metrics/README.md) |
| **09** | Install KEDA Using Helm Chart | [View Lab Guide](./09_Install_KEDA_Helm/README.md) |
| **10** | Install KEDA Using Operator Hub | [View Lab Guide](./10_Install_KEDA_Operator/README.md) |
| **11** | Install KEDA Using YAML Files | [View Lab Guide](./11_Install_KEDA_YAML/README.md) |
| **12** | Post Install Verification | [View Lab Guide](./12_Post_Install_Verify/README.md) |
| **13** | Kafka Cluster Setup | [View Lab Guide](./13_Kafka_Cluster_Setup/README.md) |
| **14** | Autoscale Based on Kafka Consumer Lag | [View Lab Guide](./14_Autoscale_Kafka_Lag/README.md) |
| **15** | Advance Configuration | [View Lab Guide](./15_Advance_Configuration/README.md) |
| **16** | Configuring ScaledJobs for RabbitMQ | [View Lab Guide](./16_ScaledJobs_RabbitMQ/README.md) |
| **17** | Implementing Different Rollout Strategies | [View Lab Guide](./17_Rollout_Strategies/README.md) |
| **18** | Testing Scaling Strategies | [View Lab Guide](./18_Testing_Scaling_Strategies/README.md) |
| **19** | Pausing and Resuming Scaling Operations | [View Lab Guide](./19_Pausing_Resuming_Scaling/README.md) |
| **20** | Exploring Basic Authentication Mechanisms | [View Lab Guide](./20_Basic_Authentication/README.md) |
| **21** | Implementing TriggerAuthentication | [View Lab Guide](./21_TriggerAuthentication/README.md) |
| **22** | Implementing TriggerAuthentication Referencing a Secret | [View Lab Guide](./22_TriggerAuthentication_Secret/README.md) |
| **23** | Integrating External Secret Management Solutions | [View Lab Guide](./23_External_Secret_Management/README.md) |
| **24** | Understanding and Configuring ClusterTriggerAuthentication | [View Lab Guide](./24_ClusterTriggerAuthentication/README.md) |
| **25** | Visualize KEDA AutoScaling Metrics in Grafana | [View Lab Guide](./25_Visualize_Metrics_Grafana/README.md) |
| **26** | Pause and Unpause the ScaledObject | [View Lab Guide](./26_Pause_Unpause_ScaledObject/README.md) |
| **27** | Monitoring AutoScaling Errors and Alerts | [View Lab Guide](./27_Monitoring_Errors_Alerts/README.md) |
| **28** | Setting Up and Testing Fallback Mechanisms | [View Lab Guide](./28_Fallback_Mechanisms/README.md) |
| **29** | Implementing Basic Scaling Modifiers | [View Lab Guide](./29_Scaling_Modifiers/README.md) |
| **30** | Querying Metrics from KEDA Metrics Server | [View Lab Guide](./30_Querying_Metrics_Server/README.md) |
| **31** | Developing the External Scaler in Go | [View Lab Guide](./31_External_Scaler_Go/README.md) |
| **32** | Containerizing and Deploying the External Scaler in Kubernetes | [View Lab Guide](./32_Deploying_External_Scaler/README.md) |
| **33** | Testing and Observing Custom Scaler | [View Lab Guide](./33_Testing_Custom_Scaler/README.md) |

---

## ⚙️ Core Stack & Prerequisites
- **Kubernetes**: kind v1.27
- **Autoscaler**: KEDA (Kubernetes Event-Driven Autoscaling)
- **Monitoring**: Prometheus, Prometheus Adapter, Grafana
- **Development**: Go 1.21, Docker
