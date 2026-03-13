# GitOps & Observability Playground with ArgoCD, Kargo.io, Istio and OpenTelemetry
A playground project that demonstrates the collaboration between ArgoCD and Kargo.io.

Kargo.io monitors this GitHub repository, including the GitHub Container Registry, for changes,
renders new manifests, and pushes them back to the GitHub repository.

Furthermore, Kargo.io can create entire promotion pipelines so that changes are only rolled out to
an environment if they have been successfully deployed in a previous environment.

ArgoCD monitors the manifests for changes and deploys them accordingly in the cluster.

This project consists of a simple frontend and backend. The backend returns the current date, and the frontend renders it.

KinD with Tilt is used as the local development environment. The local development environment also uses Istio
to be as close as possible to a production system. To ensure that Istio also has a load balancer, `kind cloud provider` is used.

### Required Tooling for Local Development
* [Docker](https://www.docker.com/)
* [KinD](https://kind.sigs.k8s.io/)
* [KinD Cloud Provider](https://github.com/kubernetes-sigs/cloud-provider-kind)
* [Tilt.dev](https://tilt.dev/)
* [Helm](https://helm.sh/)


### Spinning Up the Local Development Environment
All of the above tools must be installed and available in the $PATH Environment Variable.

```bash
# spin up kind cluster with registry
local-dev/kind-with-registry.sh
# spin up local development environment
tilt up
```

Getting Load Balancer External-IP Address
```bash
kubectl -n istio-gateway get service istio-gateway
NAME            TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)                                      AGE
istio-gateway   LoadBalancer   10.96.133.29   172.18.0.3    15021:31704/TCP,80:30396/TCP,443:30233/TCP   5h46m
```

Sending a Request trough the Load Balancer / Istio
```bash
$GATEWAY_IP=$(kubectl -n istio-gateway get svc istio-gateway -o jsonpath='{.status.loadBalancer.ingress[].ip}')
curl --fail http://dev.local/api/time --resolve "dev.local:80:${GATEWAY_IP}"
```

The browser can also be used to send requests. To do this, the Gateway IP must be provided via DNS. Example `/etc/hosts`
````bash
<the-ip-of-the-gateway> dev.local
````

### ArgoCD ApplicationSet Settings
``` 
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: gitops-playground
  namespace: argo-cd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  goTemplate: true
  goTemplateOptions: [ "missingkey=error" ]
  generators:
    - git:
        repoURL: https://github.com/procinger/gitops-playground.git
        revision: HEAD
        directories:
          - path: manifests/stages/*
          - path: manifests/stages/ci*
            exclude: true
          - path: manifests/stages/local-dev
            exclude: true
  template:
    metadata:
      name: gitops-playground-{{`{{.path.basename}}`}}
      annotations:
        kargo.akuity.io/authorized-stage: "gitops-playground:{{`{{.path.basename}}`}}"
    spec:
      project: gitops-playground
      source:
        repoURL: https://github.com/procinger/gitops-playground.git
        path: ./manifests/stages/{{`{{.path.basename}}`}}
        targetRevision: stage/{{`{{.path.basename}}`}}
      destination:
        server: https://kubernetes.default.svc
        namespace: gitops-playground-{{`{{.path.basename}}`}}
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - CreateNamespace=true
        managedNamespaceMetadata:
          labels:
            pod-security.kubernetes.io/enforce: privileged
            istio-injection: enabled
```

### Kiali - Istio Service Mesh Console
To view and analyze the service mesh in a graph, Kiali should be installed. Kiali requires Prometheus and Jaeger as dependencies.

Kiali login tokens have a short lifetime, so a new token must be requested each time.
```
kubectl --namespace observability create token kiali  
eyJhbGciOiJSUzI1NiIsImtpZCI6IndrWVlPaElib0ZIa1ZpWVdKbUZQdlZYT3MyU1ROV2tGQV9BR2RwOExEMVUifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzY4NzczOTEwLCJpYXQiOjE3Njg3NzAzMTAsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiZTYxMTU1ZTgtNzkwNy00NWFiLWFhMjctMTRjNjBlYWFiYTU2Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJvYnNlcnZhYmlsaXR5Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImtpYWxpIiwidWlkIjoiODY1M2NiMDAtMzc1My00ZDZiLWJmOTAtZjc4ODYzNDA5ZGQzIn19LCJuYmYiOjE3Njg3NzAzMTAsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpvYnNlcnZhYmlsaXR5OmtpYWxpIn0.0VZZr6-rD-XZJmp2IyMp4iKAezMefw36JcpGNVZ2qCs-vjzl-yDez4yKKmGeL4_cPLYqOkQ92C_OP47Zm_07G9n0uCKFH5phPsAuHPllSyUJFZ_GT2ezjVWg3auFZp_SWKoiSNT6cWIOlVSNTO-Y1i5xZwL96_PUS9fWvv3n5eMnoBhHZ4KMrb9pFUlcx1LWr0r1EFfoalwMuS2TfmxI0C6cXYMHtpfZUQtP_6LeVqzbuFFs_zdRniDJ6sIiNzJ9VLDMkCUqutaKUaHyQKd2YmGkI4V_eVJJlYwzHu9jxGnvbdHf4z6NHeshRV4pNcOMFmTOtkJ1Q-G8gtrCnoZWFw 
```

<img width="1917" height="1005" alt="kiali" src="https://github.com/user-attachments/assets/e826d25b-8b94-4e72-a151-a06ac9cc0a1e" />

