update_settings(max_parallel_updates = 3, k8s_upsert_timeout_secs = 300, suppress_unused_image_warnings = None)
# install required cluster tooling
include('./local-dev/Tiltfile')

include('./frontend/Tiltfile')
include('./backend/Tiltfile')

k8s_yaml(kustomize('./manifests/stages/local-dev'))
k8s_resource(
  workload = 'frontend',
  labels = ['App'],
  port_forwards = '8080:8080',
  resource_deps = ['istio-gateway', 'otel-collector'],
)
k8s_resource(
  workload = 'backend',
  labels = ['App'],
  port_forwards = '8081:8080',
  resource_deps = ['istio-gateway', 'otel-collector'],
)

k8s_resource(
  objects=[
    'virtualservice:virtualservice',
    'strict-mtls:peerauthentication',
    'isolation:authorizationpolicy',
    'instrumentation:instrumentation',
    'otel:telemetry',
  ],
  new_name='istio-config',
  labels = ['App'],
  resource_deps = ['istio-gateway', 'wait-otel-operator-ready'],
)
