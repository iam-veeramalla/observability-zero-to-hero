# Install using Helm

## Add helm repo

`helm repo add prometheus-community https://prometheus-community.github.io/helm-charts`

## Update helm repo

`helm repo update`

## Install helm 

`helm install prometheus prometheus-community/prometheus`

## Expose Prometheus Service

This is required to access prometheus-server using your browser.

`kubectl expose service prometheus-server --type=NodePort --target-port=9090 --name=prometheus-server-ext`

