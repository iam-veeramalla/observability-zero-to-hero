## 📊 What is OpenTelemetry?
- OpenTelemetry is an open-source observability framework for generating, collecting, and exporting telemetry data (traces, metrics, logs) to help monitor applications.

## 🛠️ How is it Different from Other Libraries?
- OpenTelemetry offers a unified standard for observability across multiple tools and vendors, unlike other libraries that may focus only on a specific aspect like tracing or metrics.

## ⏳ What Existed Before OpenTelemetry?
- Before OpenTelemetry, observability was typically managed using a combination of specialized tools for different aspects like
    - `Tracing`: Tools like Jaeger and Zipkin were used to track requests
    - `Metrics`: Solutions like Prometheus and StatsD were popular for collecting metrics
    - `Logging`: Tools like ELK Stack (Elasticsearch, Logstash, Kibana) or Fluentd were used to aggregate and analyze logs.
-  OpenTelemetry unified these by standardizing how telemetry data is collected and exported.
- Prior to OpenTelemetry, there were OpenTracing and OpenCensus, which OpenTelemetry merged to provide a more comprehensive and standardized observability solution.

## 🌐 Supported Programming Languages

OpenTelemetry supports several languages, including:

- **Go**
- **Java**
- **JavaScript**
- **Python**
- **C#**
- **C++**
- **Ruby**
- **PHP**
- **Swift**
- ...and others.

## Architecture

### 🖥️ Step 1: Create EKS Cluster

```bash
eksctl create cluster --name=observability \
                      --region=us-east-1 \
                      --zones=us-east-1a,us-east-1b \
                      --without-nodegroup
```
```bash
eksctl utils associate-iam-oidc-provider \
    --region us-east-1 \
    --cluster observability \
    --approve
```
```bash
eksctl create nodegroup --cluster=observability \
                        --region=us-east-1 \
                        --name=observability-ng-private \
                        --node-type=t3.medium \
                        --nodes-min=2 \
                        --nodes-max=3 \
                        --node-volume-size=20 \
                        --managed \
                        --asg-access \
                        --external-dns-access \
                        --full-ecr-access \
                        --appmesh-access \
                        --alb-ingress-access \
                        --node-private-networking

# Update ./kube/config file
aws eks update-kubeconfig --name observability
```

### 🔐 Step 2: Create IAM Role for Service Account
```bash
eksctl create iamserviceaccount \
    --name ebs-csi-controller-sa \
    --namespace kube-system \
    --cluster observability \
    --role-name AmazonEKS_EBS_CSI_DriverRole \
    --role-only \
    --attach-policy-arn arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy \
    --approve
```
- This command creates an IAM role for the EBS CSI controller.
- IAM role allows EBS CSI controller to interact with AWS resources, specifically for managing EBS volumes in the Kubernetes cluster.
- We will attach the Role with service account

### 📝 Step 3: Retrieve IAM Role ARN
```bash
ARN=$(aws iam get-role --role-name AmazonEKS_EBS_CSI_DriverRole --query 'Role.Arn' --output text)
```
- Command retrieves the ARN of the IAM role created for the EBS CSI controller service account.

### 📦 Step 4: Deploy EBS CSI Driver
```bash
eksctl create addon --cluster observability --name aws-ebs-csi-driver --version latest \
    --service-account-role-arn $ARN --force
```
- Above command deploys the AWS EBS CSI driver as an addon to your Kubernetes cluster.
- It uses the previously created IAM service account role to allow the driver to manage EBS volumes securely.


### 🧩 Step 5: Understand the Application
- We have two very simple microservice A (`microservice-a`) & B (`microservice-a`), Built with Golang using the Gin web framework for handling HTTP requests.
- **Microservice A** API Endpoints:
    - `GET /hello-a` – Returns a greeting message
    - `GET /call-b` – Calls another service (Service B) and returns its response
    - `GET /getme-coffee` – Fetches and returns data from an external coffee API
- **Microservice B** API Endpoints:
    - `GET /hello-b` – Returns a greeting message
    - `GET /call-a` – Calls another service (Service A) and returns its response
    - `GET /getme-coffee` – Fetches and returns data from an external coffee API
- Observability:
    - OpenTelemetry SDK integrated for tracing and metrics.
    - Metrics and traces are exported to the OpenTelemetry Collector via OTLP over HTTP.
- Instrumentation:
    - Uses OpenTelemetry middleware (otelgin) for automatic request tracing.
    - Instruments HTTP clients with otelhttp for distributed tracing of outbound requests.


### 🐳 Step 6: Dockerize & push it to the registry
```bash
# Dockerize microservice - a
docker build -t <<NAME_OF_YOUR_REPO>>:<<TAG>> microservice-a/

# Dockerize microservice - b
docker build -t <<NAME_OF_YOUR_REPO>>:<<TAG>> microservice-b/

# push both images
docker push  <<NAME_OF_YOUR_REPO>>:<<TAG>>
docker push  <<NAME_OF_YOUR_REPO>>:<<TAG>>
```


### 🗂️ Step 7: Create Namespace for observability components
```bash
kubectl create namespace olly
```

### 📚 Step 8: Install Elasticsearch on K8s
helm repo add elastic https://helm.elastic.co

helm install elasticsearch \
 --set replicas=1 \
 --set volumeClaimTemplate.storageClassName=gp2 \
 --set persistence.labels.enabled=true elastic/elasticsearch -n olly


### 📜 Step 9: Export Elasticsearch CA Certificate
- This command retrieves the CA certificate from the Elasticsearch master certificate secret and decodes it, saving it to a ca-cert.pem file.
```bash
kubectl get secret elasticsearch-master-certs -n olly -o jsonpath='{.data.ca\.crt}' | base64 --decode > ca-cert.pem
```

###  🔑 Step 10: Create ConfigMap for Jaeger's TLS Certificate
- Creates a ConfigMap in the olly namespace, containing the CA certificate to be used by Jaeger for TLS.
```bash
kubectl create configmap jaeger-tls --from-file=ca-cert.pem -n olly
```

### 🛡️ Step 11: Create Secret for Elasticsearch TLS
- Creates a Kubernetes Secret in the tracing namespace, containing the CA certificate for Elasticsearch TLS communication.
```bash
kubectl create secret generic es-tls-secret --from-file=ca-cert.pem -n olly
```

### 🔍 Step 12: Retrieve Elasticsearch Username & Password
```bash
# for username
kubectl get secrets --namespace=olly elasticsearch-master-credentials -ojsonpath='{.data.username}' | base64 -d
# for password
kubectl get secrets --namespace=olly elasticsearch-master-credentials -ojsonpath='{.data.password}' | base64 -d
```
- Retrieves the password for the Elasticsearch cluster's master credentials from the Kubernetes secret.
- 👉 **Note**: Please write down the password for future reference


### 🕵️‍♂️ Step 13: Install Jaeger with Custom Values
- 👉 **Note**: Please update the `password` field and other related field in the `jaeger-values.yaml` file with the password retrieved previous step at step 12: (i.e NJyO47UqeYBsoaEU)"
-  Command installs Jaeger into the olly namespace using a custom jaeger-values.yaml configuration file. Ensure the password is updated in the file before installation.
```bash
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm repo update

helm install jaeger jaegertracing/jaeger -n olly --values jaeger-values.yaml
```

### 🌐 Step 14: Access UI - Port Forward Jaeger Query Service
kubectl port-forward svc/jaeger-query 8080:80 -n olly



### 📈 Step 15: Install Opentelemetry-collector
helm install otel-collector open-telemetry/opentelemetry-collector -n olly --values otel-collector-values.yaml


### 📊 Step 16: Install prometheus
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm install  prometheus prometheus-community/prometheus -n olly --values prometheus-values.yaml
```

### 🚀 Step 17: Deploy the applicaiton
- ***Note:*** - Review the Kubernetes manifest files located in `./k8s-manifest`. and you should change image name & tag with your own image
```bash
kubectl apply -k k8s-manifests/
```
- 👉 ***Note***: wait for 5 minutes till you load balancer comes in running state

## 🔄 Step 18: Generate Load
- Script: `test.sh` takes two load balancer DNS addresses as input arguments and alternates requests between them using curl.
- `test.sh` Continuously sends random HTTP requests every second to predefined routes on two provided load balancer DNSs
- ***Note:*** Keep the script running in another terminal to quickly gather metrics & traces.

```bash
./test.sh http://Microservice_A_LOAD_BALANCER_DNS http://Microservice_B_LOAD_BALANCER_DNS
```

### 📊 Step 19: Access the UI of Prometheus
```bash
kubectl port-forward svc/prometheus-server 9090:80 -n olly
```
- Look for your application's metrics like `request_count`, `request_duration_ms`, `active_requests` and other to monitor request rates & performance.


### 🕵️‍♂️ Step 20: Access the UI of Jaeger
```bash
kubectl port-forward svc/jaeger-query 8080:80 -n olly
```
-  Look for traces from the service name microservice-a, microservice-b and operations such as `[/hello-a, /call-b, and /getme-coffee]` or `[/hello-b, /call-a, and /getme-coffee]` to monitor request flows and dependencies.

## ✅ Conclusion
- By following the above steps, you have successfully set up an observability stack using OpenTelemetry on an EKS cluster. This setup allows you to monitor your microservices effectively through integrated tracing, metrics, and logging.

## 🧼 Clean Up
```bash
helm uninstall prometheus -n olly
helm uninstall otel-collector -n olly
helm uninstall jaeger -n olly
helm uninstall elasticsearch -n olly

<!-- Delete all the pvc & pv -->

kubectl delete -k k8s-manifests/


kubectl delete ns olly

eksctl delete cluster --name observability
```