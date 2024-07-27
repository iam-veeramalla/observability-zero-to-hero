# Install using Helm

## Add helm repo

```javascript
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
```

## Update helm repo

```javascript
helm repo update
```

## Install helm 

```javascript
helm install prometheus prometheus-community/prometheus
```
## Expose Prometheus Service

This is required to access prometheus-server using your browser.
```javascript 
kubectl expose service prometheus-server --type=NodePort --target-port=9090 --name=prometheus-server-ext
```
