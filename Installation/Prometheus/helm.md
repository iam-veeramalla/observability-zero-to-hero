# Install using Helm

## Add helm repo

```shell
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
```

## Update helm repo

```shell
helm repo update
```

## Install helm

```shell
helm install prometheus prometheus-community/prometheus
```

## Expose Prometheus Service

This is required to access prometheus-server using your browser.

```shell
kubectl expose service prometheus-server --type=NodePort --target-port=9090 --name=prometheus-server-ext
```
