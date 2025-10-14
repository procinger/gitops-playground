# install required cluster tooling
include('./local-dev/Tiltfile')

include('./frontend/Tiltfile')
include('./backend/Tiltfile')

k8s_yaml(kustomize('./manifests/stages/local-dev'))
k8s_resource(
  workload = 'frontend',
  labels = ['App'],
  port_forwards = '8080:8080',
  resource_deps = ['istio-gateway'],
)
k8s_resource(
  workload = 'backend',
  labels = ['App'],
  port_forwards = '8081:8080',
  resource_deps = ['istio-gateway'],
)

k8s_resource(
  objects=[
    'virtualservice:virtualservice',
    'strict-mtls:peerauthentication',
    'isolation:authorizationpolicy',
  ],
  new_name='istio-config',
  labels = ['App'],
  resource_deps = ['istio-gateway'],
)
