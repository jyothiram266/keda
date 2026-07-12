# Mastering Kubernetes Event-Driven Autoscaling with KEDA (LFS257)

Welcome to the hands-on lab repository for **LFS257: Mastering Kubernetes Event-Driven Autoscaling with KEDA** by The Linux Foundation.

This repository organizes all the course labs and exercises into structured directories. Each directory contains a dedicated `README.md` with cleaned instructions, code snippets, manifests, and commands.

---

## 📚 Course Labs & Exercises

| Lab Index | Exercise Title | Link |
| :---: | :--- | :---: |
| **01** | Set Up the Course Lab Environment | [View Lab Guide](./01_Lab_1_Set_Up_the_Course_Lab_Environment/README.md) |
| **02** | Containerize and Deploy Sample Application in Kubernetes | [View Lab Guide](./02_Lab_Exercise_2.1_Containerize_and_Deploy_Sample_Application_in_Kubernetes/README.md) |
| **03** | Configure and Test HPA | [View Lab Guide](./03_Lab_Exercise_2.2_Configure_and_Test_HPA/README.md) |
| **04** | Configure and Test VPA | [View Lab Guide](./04_Lab_Exercise_2.3_Configure_and_Test_VPA/README.md) |
| **05** | Containerize and Deploy Instrumented Sample Application in Kubernetes | [View Lab Guide](./05_Lab_Exercise_3.1_Containerize_and_Deploy_Instrumented_Sample_Application_in_Kubernetes/README.md) |
| **06** | Setting Up Prometheus in Kubernetes | [View Lab Guide](./06_Lab_Exercise_3.2_Setting_Up_Prometheus_in_Kubernetes/README.md) |
| **07** | Installing and Configuring the Prometheus Adapter | [View Lab Guide](./07_Lab_Exercise_3.3_Installing_and_Configuring_the_Prometheus_Adapter/README.md) |
| **08** | Testing Autoscaling with Custom Metrics | [View Lab Guide](./08_Lab_Exercise_3.4_Testing_Autoscaling_with_Custom_Metrics/README.md) |
| **09** | Install KEDA Using Helm Chart | [View Lab Guide](./09_Lab_Exercise_5.1_Install_KEDA_Using_Helm_Chart/README.md) |
| **10** | Install KEDA Using Operator Hub | [View Lab Guide](./10_Lab_Exercise_5.2_Install_KEDA_Using_Operator_Hub/README.md) |
| **11** | Install KEDA Using YAML Files | [View Lab Guide](./11_Lab_Exercise_5.3_Install_KEDA_Using_YAML_Files/README.md) |
| **12** | Post Install Verification | [View Lab Guide](./12_Lab_Exercise_5.4_Post_Install_Verification/README.md) |
| **13** | Kafka Cluster Setup | [View Lab Guide](./13_Lab_Exercise_6.1_Kafka_Cluster_Setup/README.md) |
| **14** | Autoscale Based on Kafka Consumer Lag | [View Lab Guide](./14_Lab_Exercise_6.2_Autoscale_Based_on_Kafka_Consumer_Lag/README.md) |
| **15** | Advance Configuration | [View Lab Guide](./15_Lab_Exercise_6.3_Advance_Configuration/README.md) |
| **16** | Configuring ScaledJobs for RabbitMQ | [View Lab Guide](./16_Lab_Exercise_7.1_Configuring_ScaledJobs_for_RabbitMQ/README.md) |
| **17** | Implementing Different Rollout Strategies | [View Lab Guide](./17_Lab_Exercise_7.2_Implementing_Different_Rollout_Strategies/README.md) |
| **18** | Testing Scaling Strategies | [View Lab Guide](./18_Lab_Exercise_7.3_Testing_Scaling_Strategies/README.md) |
| **19** | Pausing and Resuming Scaling Operations | [View Lab Guide](./19_Lab_Exercise_7.4_Pausing_and_Resuming_Scaling_Operations/README.md) |
| **20** | Exploring Basic Authentication Mechanisms | [View Lab Guide](./20_Lab_Exercise_8.1_Exploring_Basic_Authentication_Mechanisms/README.md) |
| **21** | Implementing TriggerAuthentication | [View Lab Guide](./21_Lab_Exercise_8.2_Implementing_TriggerAuthentication/README.md) |
| **22** | Implementing TriggerAuthentication Referencing a Secret | [View Lab Guide](./22_Lab_Exercise_8.3_Implementing_TriggerAuthentication_Referencing_a_Secret/README.md) |
| **23** | Integrating External Secret Management Solutions | [View Lab Guide](./23_Lab_Exercise_8.4_Integrating_External_Secret_Management_Solutions/README.md) |
| **24** | Understanding and Configuring ClusterTriggerAuthentication | [View Lab Guide](./24_Lab_Exercise_8.5_Understanding_and_Configuring_ClusterTriggerAuthentication/README.md) |
| **25** | Visualize KEDA AutoScaling Metrics in Grafana | [View Lab Guide](./25_Lab_Exercise_9.1_Visualize_KEDA_AutoScaling_Metrics_in_Grafana/README.md) |
| **26** | Pause and Unpause the ScaledObject | [View Lab Guide](./26_Lab_Exercise_9.2_Pause_and_Unpause_the_ScaledObject/README.md) |
| **27** | Monitoring AutoScaling Errors and Alerts | [View Lab Guide](./27_Lab_Exercise_9.3_Monitoring_AutoScaling_Errors_and_Alerts/README.md) |
| **28** | Setting Up and Testing Fallback Mechanisms | [View Lab Guide](./28_Lab_Exercise_10.1_Setting_Up_and_Testing_Fallback_Mechanisms/README.md) |
| **29** | Implementing Basic Scaling Modifiers | [View Lab Guide](./29_Lab_Exercise_10.2_Implementing_Basic_Scaling_Modifiers/README.md) |
| **30** | Querying Metrics from KEDA Metrics Server | [View Lab Guide](./30_Lab_Exercise_10.3_Querying_Metrics_from_KEDA_Metrics_Server/README.md) |
| **31** | Developing the External Scaler in Go | [View Lab Guide](./31_Lab_Exercise_11.1_Developing_the_External_Scaler_in_Go/README.md) |
| **32** | Containerizing and Deploying the External Scaler in Kubernetes | [View Lab Guide](./32_Lab_Exercise_11.2_Containerizing_and_Deploying_the_External_Scaler_in_Kubernetes/README.md) |
| **33** | Testing and Observing Custom Scaler | [View Lab Guide](./33_Lab_Exercise_11.3_Testing_and_Observing_Custom_Scaler/README.md) |

---

## ⚙️ Core Stack & Prerequisites
- **Kubernetes**: kind v1.27
- **Autoscaler**: KEDA (Kubernetes Event-Driven Autoscaling)
- **Monitoring**: Prometheus, Prometheus Adapter, Grafana
- **Development**: Go 1.21, Docker
