# Install using APT repository in Ubuntu or Debian

## Install packages 

```sudo apt-get install -y apt-transport-https software-properties-common wget```

## Import GPG key

```sudo mkdir -p /etc/apt/keyrings/```
```wget -q -O - https://apt.grafana.com/gpg.key | gpg --dearmor | sudo tee /etc/apt/keyrings/grafana.gpg > /dev/null```

## Add repository for stable release

```echo "deb [signed-by=/etc/apt/keyrings/grafana.gpg] https://apt.grafana.com stable main" | sudo tee -a /etc/apt/sources.list.d/grafana.list```

## Update the list of available packages

```sudo apt-get update```

## Install Grafana OSS 

```sudo apt-get install grafana```