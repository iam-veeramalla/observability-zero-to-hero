This is setting up a **complete Kubernetes monitoring stack** on **Amazon EKS** using **Prometheus, Grafana, and Alertmanager**.

Think of it like this:

```
AWS Cloud
│
├── EKS Cluster (Your Kubernetes Cluster)
│      │
│      ├── Worker Nodes (EC2 Instances)
│      │
│      ├── Your Applications
│      │
│      └── Monitoring Stack
│             ├── Prometheus (Collects metrics)
│             ├── Grafana (Visualizes metrics)
│             └── Alertmanager (Sends alerts)
```

Let's go through every line.

---

# Step 1 : Create an EKS Cluster

Before anything, we need a Kubernetes cluster.

Imagine you're building a city.

Before building hospitals, schools, and houses, you first need land.

The EKS Cluster is that land.

---

## Prerequisites

### AWS CLI

```
aws
```

AWS CLI is simply a command line tool that talks to AWS.

Instead of opening AWS Console every time...

You type

```
aws ec2 describe-instances
```

and AWS replies.

Without AWS CLI, your computer cannot communicate with AWS services.

---

### Configure AWS CLI

```
aws configure
```

This asks you for

```
AWS Access Key
AWS Secret Key
Region
Output Format
```

Example

```
AWS Access Key ID: XXXXX
AWS Secret Access Key: XXXXX
Default region: us-east-1
Default output: json
```

These credentials are stored locally.

Now whenever you type

```
aws ....
```

AWS knows who you are.

---

## Install eksctl

`eksctl` is a tool created specifically for creating EKS clusters.

Instead of manually creating

* VPC
* IAM Roles
* Security Groups
* Cluster
* Node Groups

it does everything automatically.

Without eksctl you would need dozens of AWS Console steps.

---

## Install kubectl

This is Kubernetes' command-line tool.

Think of it as:

```
kubectl
      ↓
talks to Kubernetes API Server
```

Whenever you type

```
kubectl get pods
```

it asks Kubernetes

> "Show me all the running Pods."

---

# Create Cluster

```
eksctl create cluster \
--name=observability \
--region=us-east-1 \
--zones=us-east-1a,us-east-1b \
--without-nodegroup
```

Let's decode every option.

---

## eksctl create cluster

Means

> Create a brand new EKS cluster.

---

## --name=observability

Cluster name.

AWS will create

```
observability
```

as your Kubernetes cluster.

---

## --region=us-east-1

Create it in

```
Virginia
```

AWS Region.

---

## --zones=us-east-1a,us-east-1b

AWS regions contain multiple Availability Zones.

```
us-east-1

├── us-east-1a
├── us-east-1b
├── us-east-1c
├── us-east-1d
```

Instead of using only one data center,

Kubernetes spreads resources across

```
1a
1b
```

This gives High Availability.

If one zone fails,

the cluster still works.

---

## --without-nodegroup

This is interesting.

Normally a cluster has worker machines.

```
Cluster

API Server

Worker Node
Worker Node
Worker Node
```

Here we're saying

> Create only the control plane.

Do NOT create worker nodes yet.

Why?

Because we want more control over the worker node configuration later.

---

So after this command

You only have

```
Control Plane

No worker nodes.
```

---

# Associate IAM OIDC Provider

```
eksctl utils associate-iam-oidc-provider \
--region us-east-1 \
--cluster observability \
--approve
```

This sounds complicated but is actually simple.

---

Normally

Pods inside Kubernetes

cannot directly access AWS services.

For example

```
Pod

↓

S3 Bucket
```

Permission denied.

---

OIDC enables

```
IAM Role

↓

Pod

↓

AWS Services
```

Now Pods can securely access AWS.

This is required for many Kubernetes add-ons.

---

## --approve

Without this

eksctl asks

```
Are you sure?

yes/no
```

This automatically answers

```
yes
```

---

# Create Node Group

Now we're creating actual worker machines.

```
eksctl create nodegroup
```

Remember

Earlier we created only the control plane.

Now we're creating the workers.

---

```
--cluster=observability
```

Attach these nodes to the cluster.

---

```
--region=us-east-1
```

Same region.

---

```
--name=observability-ng-private
```

Worker node group name.

---

```
--node-type=t3.medium
```

Every worker node will be

```
EC2 t3.medium
```

Specifications

```
2 vCPU

4 GB RAM
```

---

```
--nodes-min=2
```

Minimum

```
2 EC2 instances
```

---

```
--nodes-max=3
```

Auto Scaling can increase to

```
3 nodes
```

when load increases.

---

```
--node-volume-size=20
```

Each EC2 gets

```
20 GB
```

disk.

---

```
--managed
```

AWS manages upgrades,

health,

replacement,

etc.

---

```
--asg-access
```

Allows Kubernetes Cluster Autoscaler to communicate with AWS Auto Scaling Groups.

---

```
--external-dns-access
```

Allows ExternalDNS to automatically manage DNS records in Route 53.

---

```
--full-ecr-access
```

Allows worker nodes to pull container images from Amazon ECR.

---

```
--appmesh-access
```

Allows integration with AWS App Mesh.

---

```
--alb-ingress-access
```

Allows the AWS Load Balancer Controller to create Application Load Balancers (ALBs) for Kubernetes Ingress resources.

---

```
--node-private-networking
```

Worker nodes receive only private IP addresses.

```
Internet

↓

Load Balancer

↓

Private Worker Nodes
```

Safer than exposing the nodes directly to the internet.

---

Now the cluster looks like

```
Control Plane

↓

Node 1
Node 2
```

Pods can finally run.

---

# Update kubeconfig

```
aws eks update-kubeconfig --name observability
```

This is one of the most important commands.

`kubectl` doesn't automatically know which cluster to manage. It reads a configuration file called **kubeconfig**, usually located at:

```
~/.kube/config
```

Before running this command:

```
kubectl
      │
      └── "I don't know which cluster to talk to."
```

After running it, AWS fetches the cluster's endpoint and authentication details and writes them into your kubeconfig.

Now:

```
kubectl
      │
      ▼
~/.kube/config
      │
      ▼
EKS API Server
```

From this point onward, commands like:

```bash
kubectl get nodes
kubectl get pods
kubectl apply -f deployment.yaml
```

will operate on your `observability` cluster.

---

# Step 2 : Install kube-prometheus-stack

First, Helm is used to install applications on Kubernetes. You can think of Helm as **apt** on Ubuntu or **npm** for Node.js, but for Kubernetes.

```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
```

This tells Helm:

> "Add a repository named `prometheus-community` that contains many ready-made Kubernetes applications."

It doesn't install anything yet; it just registers the repository.

Then:

```bash
helm repo update
```

downloads the latest list of available charts (packages) from all configured repositories, similar to `apt update`.

---

# Step 3 : Deploy the monitoring stack

Create a dedicated namespace:

```bash
kubectl create ns monitoring
```

Namespaces logically separate applications within the same cluster.

Instead of everything being mixed together:

```
Cluster
│
├── default
├── kube-system
└── monitoring
```

all monitoring components will live inside the `monitoring` namespace.

Then:

```bash
cd day-2
```

moves into the directory containing your configuration file.

Now install the chart:

```bash
helm install monitoring prometheus-community/kube-prometheus-stack \
-n monitoring \
-f ./custom_kube_prometheus_stack.yml
```

Breaking it down:

* `helm install` → install a Helm chart.
* `monitoring` → the name of this Helm release.
* `prometheus-community/kube-prometheus-stack` → the chart to install.
* `-n monitoring` → install into the `monitoring` namespace.
* `-f custom_kube_prometheus_stack.yml` → override default settings using your custom values file (for example, Grafana password, storage size, resource limits, or enabled components).

This deploys many Kubernetes objects, including Deployments, StatefulSets, Services, ConfigMaps, ServiceAccounts, CRDs, Prometheus, Grafana, Alertmanager, and exporters.

---

# Step 4 : Verify the installation

```bash
kubectl get all -n monitoring
```

asks Kubernetes to show all major resources in the `monitoring` namespace:

* Pods
* Services
* Deployments
* ReplicaSets
* StatefulSets

You should see components like:

```
grafana
prometheus
alertmanager
node-exporter
kube-state-metrics
```

If the pods are in the `Running` state, the installation was successful.

---

# Access the UIs

These services run **inside the cluster**, so your laptop cannot reach them directly.

`kubectl port-forward` temporarily creates a tunnel from your computer to a Kubernetes Service.

### Prometheus

```bash
kubectl port-forward service/prometheus-operated -n monitoring 9090:9090
```

```
Browser localhost:9090
        │
        ▼
Your Laptop
        │
        ▼
kubectl tunnel
        │
        ▼
Prometheus Service
```

If you're connected to a remote EC2 instance, adding:

```bash
--address 0.0.0.0
```

allows other machines to access the forwarded port.

### Grafana

```bash
kubectl port-forward service/monitoring-grafana -n monitoring 8080:80
```

Visit:

```
http://localhost:8080
```

Login:

```
Username: admin
Password: prom-operator
```

Grafana displays dashboards using data collected by Prometheus.

### Alertmanager

```bash
kubectl port-forward service/alertmanager-operated -n monitoring 9093:9093
```

Visit:

```
http://localhost:9093
```

Alertmanager receives alerts from Prometheus, groups them, suppresses duplicates, and can send notifications through email, Slack, PagerDuty, and many other integrations.

---

# Step 5 : Clean up

If you're done, remove everything to avoid AWS charges.

Remove the Helm release:

```bash
helm uninstall monitoring --namespace monitoring
```

This deletes all monitoring components but leaves the namespace.

Delete the namespace:

```bash
kubectl delete ns monitoring
```

This removes any remaining resources inside it.

Finally, delete the EKS cluster:

```bash
eksctl delete cluster --name observability
```

This tears down the control plane, node groups, networking resources, and associated infrastructure created by `eksctl`.

---

# Complete workflow

```text
1. Install tools
        │
        ▼
AWS CLI + kubectl + eksctl + Helm
        │
        ▼
2. Create EKS Control Plane
        │
        ▼
3. Enable IAM OIDC
        │
        ▼
4. Create Worker Nodes
        │
        ▼
5. Configure kubectl (kubeconfig)
        │
        ▼
6. Add Helm Repository
        │
        ▼
7. Create monitoring Namespace
        │
        ▼
8. Install kube-prometheus-stack
        │
        ▼
9. Verify Pods and Services
        │
        ▼
10. Port-forward to access:
      ├── Prometheus (Metrics Storage)
      ├── Grafana (Dashboards)
      └── Alertmanager (Alerts)
        │
        ▼
11. When finished, uninstall the chart and delete the cluster.
```

Once you understand this setup, you'll have a solid foundation for how production Kubernetes monitoring is typically deployed: **Prometheus collects metrics**, **Grafana visualizes them**, and **Alertmanager notifies you when something goes wrong**, all running inside an **Amazon EKS** cluster.
