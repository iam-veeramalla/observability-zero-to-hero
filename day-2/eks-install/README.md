# **Project Overview**
This project automates the deployment of an EKS cluster on AWS and configures monitoring with the kube-prometheus-stack using Pulumi.

## Prerequisites

- **Pulumi**: Install the Pulumi CLI.
- **AWS CLI**: Install and configure by running `aws configure`.
- **Python**: Install Python 3.8+ and set up a virtual environment.
- **Dependencies**: Install dependencies by running:
  ```bash
  pip install -r requirements.txt
* Clone the Repository:
    ```bash
  git clone https://github.com/your-repo-name.git
  cd your-repo-name

## Pulumi Stack Setup: Initialize a Pulumi stack for the environment:
* Initialize a Pulumi stack for the environment
    ```bash
  pulumi stack init dev

* Configure environment settings (e.g., region, vpcCidr) in environments/<environment>.yaml.

## Running the Code
* Preview Changes
    ```bash
    pulumi preview

* Deploy
    ```bash
    pulumi up

## Access Monitoring Tools

Use port-forwarding commands to access Prometheus and Grafana UIs. (Check in other readme of day2)

* Verify deployment using:
    ```bash
    kubectl get all -n monitoring

## Cleanup 
pulumi destroy


