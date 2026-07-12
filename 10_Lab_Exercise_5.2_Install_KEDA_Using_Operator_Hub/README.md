# Lab Exercise 5.2: Install KEDA Using Operator Hub

In this exercise, we install KEDA using the Operator Lifecycle Manager (OLM) and the OperatorHub marketplace catalog, configuring KEDA parameters via a custom `KedaController` resource.

### 🌐 Operator Lifecycle Manager (OLM) KEDA Deployment Flow

```mermaid
graph TD
    subgraph OLMControl["OLM System Namespace"]
        Catalog["OperatorHub Catalog Source<br/>(operatorhubio-catalog)"]
    end

    subgraph KedaNS["keda Namespace"]
        OG["OperatorGroup Resource<br/>(Defines target namespaces)"]
        Sub["Subscription Resource<br/>(my-keda points to catalog)"]
        CSV["ClusterServiceVersion (CSV)<br/>(keda.v2.12.1)"]
    end

    subgraph CustomController["KEDA Controller Resource"]
        KController["KedaController CR<br/>(Specifies QPS, Burst, Cert rotation)"]
        KedaDeploy["KEDA Daemon Deployments<br/>(operator, webhooks, metrics-apiserver)"]
    end

    Sub -->|1. Pulls operator bundle from| Catalog
    Sub -->|2. Bounds to scope of| OG
    Catalog -->|3. Resolves and installs| CSV
    CSV -->|4. Registers Custom Resource Definitions (CRDs)| KController
    KController -->|5. Activates & Configures| KedaDeploy
```

### 🛠️ Key Concepts & Design Decisions
1. **Operator Lifecycle Manager (OLM)**:
   - OLM helps manage the lifecycle (install, update, and manage permissions) of Kubernetes operators. It catalogues operators and installs dependencies (CRDs, ServiceAccounts, cluster roles) automatically.
2. **Subscription & OperatorGroup**:
   - A **Subscription** tells OLM which operator to pull, from which channel (e.g. `stable`), and which catalog source.
   - An **OperatorGroup** defines which namespaces the operator will watch. Setting `targetNamespaces: [keda]` scopes the operator permissions to the keda namespace.
3. **KedaController CR (Custom Resource)**:
   - Instead of configuring KEDA pods via environment variables or Helm values, OLM-based installation uses the `KedaController` custom resource. This CR defines parameters like API request rate-limiting (`kube-api-qps`, `kube-api-burst`) and TLS certificate rotation.

## Prerequisites

Kubernetes cluster with Metric Server installed as per Lab 1.

## Lab Exercise

1. Install Operator Lifecycle Manager (OLM):
The OLM is a tool to help manage the Operators running on your cluster.
```bash
curl -sL https://github.com/operator-framework/operator-lifecycle-manager/releases/downloa
```
d/v0.26.0/install.sh | bash -s v0.26.0
2. Create KEDA Namespace:
```bash
kubectl create namespace keda
```
3. Create OperatorGroup:
An OperatorGroup is required to deploy the KEDA Operator in a specific namespace, and to specify the target
namespaces where the KEDA Controller will be deployed. The OperatorGroup is used to define the scope of
the operator and to restrict its access to specific namespaces. The KEDA Operator is responsible for deploying
the KEDA Controller, which is triggered by the creation of a KedaController custom resource.
Create a file named keda-operator-group.yaml with the contents provided below and apply it using the
following command.
```yaml
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
name: keda-operator-group
namespace: keda
spec:
targetNamespaces:
- keda
```
```bash
kubectl apply -f keda-operator-group.yaml
```
4. Create KEDA Operator:
Create a file name keda-operator.yaml with the contents below and apply it using the provided command.
```yaml
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
name: my-keda
namespace: keda
spec:
channel: stable
name: keda
source: operatorhubio-catalog
sourceNamespace: olm
```
```bash
kubectl apply -f keda-operator.yaml
```
5. Verify operator is installed and running.
```bash
kubectl get csv -n keda
```
NAME DISPLAY VERSION REPLACES PHASE
keda.v2.12.1 KEDA 2.12.1 keda.v2.12.0 Succeeded
6. Create KEDA Controller:
Create a file name keda-controller.yaml with the contents below and apply it using the following
command. Note the highlighted arguments passed to KedaController, kube-api-qps (Sets the QPS rate for
throttling requests sent to the apiserver) and kube-api-burst ( Sets the burst for throttling requests sent to
the apiserver). We can pass many other additional configurations through controller args.
```yaml
apiVersion: keda.sh/v1alpha1
kind: KedaController

metadata:
name: keda
namespace: keda
spec:
admissionWebhooks:
logEncoder: console
logLevel: info
metricsServer:
logLevel: '0'
operator:
logEncoder: console
logLevel: info
args:
- "kube-api-qps=40"
- "kube-api-burst=50"
serviceAccount: null
watchNamespace: ''
```
```bash
kubectl apply -f keda-controller.yaml
```
7. Verify KEDA pods are running in the cluster using the command below:
```bash
kubectl get deployment -n keda
```
NAME READY UP-TO-DATE AVAILABLE AGE
keda-admission-webhooks 1/1 1 1 2d16h
keda-operator 1/1 1 1 2d16h
keda-olm-operator 1/1 1 1 2d16h
keda-operator-metrics-apiserver 1/1 1 1 2d16h
8. (Optional) Register your own CA in KEDA Operator Trusted Store
For configuring, custom CA & TLS certificates refer to Lab 5.1, Lab Exercise instruction
number 3.
Once the certificates and secrets are created, apply the modified controller resource below. As highlighted
below, we have disabled auto certificate generation --enable-cert-rotation=false
```yaml
apiVersion: keda.sh/v1alpha1
kind: KedaController
metadata:
name: keda
namespace: keda
spec:
admissionWebhooks:
logEncoder: console

logLevel: info
metricsServer:
logLevel: '0'
operator:
logEncoder: console
logLevel: info
args:
- "--enable-cert-rotation=false"
- "kube-api-qps=40"
- "kube-api-burst=50"
serviceAccount: null
watchNamespace: ''
```
9. (Optional) Uninstall KEDA
If you wish to try different installation methods, uninstall KEDA installed via Operator Hub first.
```bash
kubectl delete -f keda-controller.yaml
kubectl delete -f keda-operator.yaml
kubectl delete -f keda-operator-group.yaml
kubectl delete clusterserviceversion –all -n keda
```

## Summary

In this exercise, we completed the following:
- Installed the Operator Lifecycle Manager (OLM) to manage the Operators running on your cluster.
- Created a namespace for KEDA and an OperatorGroup to deploy the KEDA Operator in the
namespace, and specify the target namespaces where the KEDA Controller will be deployed.
- Created a Subscription resource to install the KEDA Operator from Operator Hub and verify that the
Operator is installed and running.
- Create a KedaController custom resource to deploy the KEDA Controller and verify that the KEDA pods
are running in the cluster.
- Modified the KedaController resource to disable auto certificate generation by adding the
--enable-cert-rotation=false argument to the KEDA Operator.