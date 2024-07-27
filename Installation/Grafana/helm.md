# Install using Helm

## Add helm repo

```javascript
helm repo add grafana https://grafana.github.io/helm-charts
```

## Update helm repo

```javascript
helm repo update`
```

## Install helm 

```javascript
helm install grafana grafana/grafana
```

## Expose Grafana Service

```javascript
kubectl expose service grafana — type=NodePort — target-port=3000 — name=grafana-ext
```
