import os
import pulumi
import pulumi_eks as eks
import pulumi_aws as aws
import pulumi_kubernetes.helm.v3 as helm
import yaml

# Load configuration
config = pulumi.Config()
environment = config.require("environment")

# Load YAML configuration based on environment
with open(f"./environments/{environment}.yaml", 'r') as file:
    env_config = yaml.safe_load(file)

# Set variables from YAML config
vpc_cidr = env_config["vpcCidr"]
subnet_cidr1 = env_config["subnetCidr1"]
subnet_cidr2 = env_config["subnetCidr2"]
desired_capacity = env_config["desiredCapacity"]
region = env_config["region"]

# Create VPC
vpc = aws.ec2.Vpc(
    f"eks-vpc-{environment}",
    cidr_block=vpc_cidr,
    enable_dns_support=True,
    enable_dns_hostnames=True
)

# Create subnets
subnet1 = aws.ec2.Subnet(
    f"subnet-1-{environment}",
    vpc_id=vpc.id,
    cidr_block=subnet_cidr1,
    availability_zone=f"{region}a"
)
subnet2 = aws.ec2.Subnet(
    f"subnet-2-{environment}",
    vpc_id=vpc.id,
    cidr_block=subnet_cidr2,
    availability_zone=f"{region}b"
)

# Create EKS cluster
cluster = eks.Cluster(
    f"my-cluster-{environment}",
    vpc_id=vpc.id,
    subnet_ids=[subnet1.id, subnet2.id],
    instance_type="t3.medium",
    desired_capacity=desired_capacity,
    min_size=1,
    max_size=5,
    enabled_cluster_log_types=["api", "audit", "authenticator"],
    version="1.29"
)

# Install kube-prometheus-stack using Helm
prometheus_chart = helm.Chart(
    "prometheus",
    helm.ChartOpts(
        chart="kube-prometheus-stack",
        fetch_opts=helm.FetchOpts(
            repo="https://prometheus-community.github.io/helm-charts"
        ),
        namespace="monitoring",
        values={
            "global": {
                "rbac": {"create": True},
            }
        },
    ),
    opts=pulumi.ResourceOptions(depends_on=[cluster])
)

# Generate kubeconfig with error handling for missing keys
kubeconfig = cluster.kubeconfig.apply(lambda config: f"""
apiVersion: v1
clusters:
- cluster:
    server: {config.get('server', '<default-server-url>')}  # Use a default if 'server' is missing
    certificate-authority-data: {config.get('certificateAuthorityData', '<default-ca-data>')}
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: aws
  name: aws
current-context: aws
kind: Config
users:
- name: aws
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1  # Updated here
      command: aws
      args:
        - "eks"
        - "get-token"
        - "--region"
        - "{region}"
        - "--cluster-name"
        - "{config.get('clusterName', 'default-cluster-name')}"
""")

# Export outputs for further inspection and use
pulumi.export("kubeconfig", kubeconfig)
pulumi.export("cluster_endpoint", cluster.core.cluster.endpoint)
pulumi.export("cluster_name", cluster.core.cluster.name)
pulumi.export("prometheus_chart", prometheus_chart)
