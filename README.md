# GitOps Demo Project with ArgoCD and Kargo.io
A demo project that demonstrates the collaboration between ArgoCD and Kargo.io.

Kargo.io monitors this GitHub repository, including the GitHub Container Registry, for changes,
renders new manifests, and pushes them back to the GitHub repository.

Furthermore, Kargo.io can create entire promotion pipelines so that changes are only rolled out to
an environment if they have been successfully deployed in a previous environment.

ArgoCD monitors the manifests for changes and deploys them accordingly in the cluster.

This project consists of a simple frontend and backend. The backend returns the current date, and the frontend renders it.

KinD with Tilt is used as the local development environment. The local development environment also uses Istio
to be as close as possible to a production system. To ensure that Istio also has a load balancer, `kind cloud provider` is used.

### Required Tooling 
* [Docker](https://www.docker.com/)
* [KinD](https://kind.sigs.k8s.io/)
* [KinD Cloud Provider](https://github.com/kubernetes-sigs/cloud-provider-kind)
* [Tilt.dev](https://tilt.dev/)


### Spinning Up the Local Development Environment
All of the above tools must be installed and available in the $PATH Environment Variable.

```bash
tilt up
```

Getting Load Balancer IP Address
```bash
 kubectl -n istio-gateway get service istio-gateway
```

Sending a Request trough the Load Balancer / Istio
```bash
$GATEWAY_IP=$(kubectl -n istio-gateway get svc istio-gateway -o jsonpath='{.status.loadBalancer.ingress[].ip}')
curl --fail http://dev.local/api/time --resolve "dev.local:80:${GATEWAY_IP}"
```

The browser can also be used to send requests. To do this, the Gateway IP must be provided via DNS. Example /etc/hosts
````bash
<the-ip-of-the-gateway> dev.local
````
