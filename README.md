# kubermeter

This repo maintains a JMeter + Kubernetes project that aims to help JMeter take advantage of Kubernetes horizontal scaling capabilities.


## Try it on a kind cluster
To illustrate how to run this locally:
```
kind create cluster --name kubermeter --config=<(cat <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
EOF
)
```

I recommend to install this storageclass
```
kubectl apply -f - <<EOF
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: local-storage
provisioner: kubernetes.io/no-provisioner
volumeBindingMode: WaitForFirstConsumer
EOF
```

NOTE: This project supports using EFS storageclasses for EKS. Helm chart can be customized further to include other storageclass types.

### Run the helm command:
```
helm upgrade --install kubermeter --set namespaceSuffix=123 --set destination.target=onprem .
```


## APPENDIX FOR KCD
I installed prometheus locally to use it as a test endpoint for KCD presentation:

```
helm repo add prometheus-community/prometheus
helm install prometheus prometheus-community/prometheus --namespace prometheus
kubectl port-forward svc/prometheus-server 9090:80 -n prometheus
```

A Beautiful report file will get generated after executing, to extract it:
```
kubectl cp jmeter-namespace-666/jmeter-0:/mnt/nfs/reports.tar.gz reports.tar.gz
tar -xzvf reports.tar.gz
```
