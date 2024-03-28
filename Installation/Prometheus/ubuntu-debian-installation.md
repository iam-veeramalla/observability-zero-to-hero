# Linux Installation for Prometheus

## Using pre-compiled binaries

```wget https://github.com/prometheus/prometheus/releases/download/v2.51.0/prometheus-2.51.0.linux-amd64.tar.gz```

## Extracting files

```tar -xvf prometheus-2.51.0-rc.0.linux-amd64.tar.gz```

## Go inside prometheus

```cd prometheus-2.51.0-rc.0.linux-amd64/```

## Run prometheus in background

After you're on the main directory and found out `prometheus`, execute the file

```./prometheus &```