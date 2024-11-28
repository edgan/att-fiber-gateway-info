# Helm
`att-fiber-gatway-info` can be run in Kubernetes via the helm chart.

## Values
* flags.password
  * This must be set
  * Use single quotes to avoid issues with special characters in the password
* flags.statsdipport
  * 127.0.0.1 in the container and 127.0.0.1 on the host aren't the same
  * "bind_host 0.0.0.0" may need to be set to get statsd to listen on all
addresses
* flags.url
  * Set if your gateway's ip address isn't the defaul of 192.168.1.254

## Usage
```
cd helm
helm package .
helm install att-fiber-gateway-info att-fiber-gateway-info-0.1.0.tgz -f values.yaml
```

## Checking on it
```
kubectl get pods -l app=att-fiber-gateway-info
kubectl get pods -l app=att-fiber-gateway-info -o yaml
kubectl logs -l app=att-fiber-gateway-info
```
