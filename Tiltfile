load('ext://restart_process', 'docker_build_with_restart')

local_resource('Install Cluster API',
               'make init-cluster-api'
)

local_resource('Wait Cluster API resources',
               'make wait-cluster-api-resources',
               resource_deps=[
                 'Install Cluster API'
               ]
)

local_resource('Populate Cluster',
               'make create-cluster',
               deps=['Makefile'],
               resource_deps=[
                 'Wait Cluster API resources'
               ]
)

local_resource('Build binary',
               'make all',
               deps=['Makefile'],
)

docker_build_with_restart('manager:test',
             '.',
             dockerfile='./Dockerfile.dev',
             entrypoint='/app/manager',
             live_update=[
               sync('./build/manager', '/app/manager')
             ]
)

k8s_yaml('.kubernetes/prod/manifest.yaml')

k8s_resource('manager', port_forwards=8080)
