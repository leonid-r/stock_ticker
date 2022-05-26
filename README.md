
# Stock Ticker

Provides stock infromataion for requested days and avarage closing price

## Getting Started

### Dependencies

* golang 1.18
* Docker engine v20.10.14
* Kubernetes cluster v1.24.0

### Building/Installing

* Pull the code from repo
* Run from root folder `docker build .`
* Push image to the docker registry (You will need link to the image to update deployment)

### Deployment prerequisites 

* Create AKS cluster (v1.24.0)
* Add label to the nodes where the application will be deployed `kubectl label nodes <node-name> app=stock-ticker`
* Verify label `kubectl get nodes --show-labels`
* Install nginx controller https://kubernetes.github.io/ingress-nginx/deploy/#quick-start
* Update deployment with new image location if required (https://github.com/leonid-r/stock_ticker/blob/dbef67a11a8e8a9573c4a05aa7ed8486eb0309ae/manifests/deployment.yaml#L20)
* Update config map if required with new parameters (https://github.com/leonid-r/stock_ticker/blob/dbef67a11a8e8a9573c4a05aa7ed8486eb0309ae/manifests/configmap.yaml#L7-L8)
* Add API key to secret file (https://github.com/leonid-r/stock_ticker/blob/dbef67a11a8e8a9573c4a05aa7ed8486eb0309ae/manifests/deployment.yaml#L20)
* Get EXTERNAL-IP field from command: `kubectl get service ingress-nginx-controller --namespace=ingress-nginx`
* Add entry to the host file `<NGINX-EXTERNAL-IP> stock-ticker.io`

### Deployment

* Deploy manifests
```
cd manifests
kubectl apply -f namespace.yaml
kubectl apply -f .
```

## Validation

`curl http://stock-ticker.io/data`




