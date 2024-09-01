# Install using Helm

## Add helm repo

```shell
helm repo add grafana https://grafana.github.io/helm-charts
```

## Update helm repo

```shell
helm repo update
```

## Install helm

```shell
helm install grafana grafana/grafana
```

## Expose Grafana Service

```shell
kubectl expose service grafana --type=NodePort --target-port=3000 --name=grafana-ext
```
