load('ext://restart_process', 'docker_build_with_restart')

local_resource('Install CAPI dependencies',
               'make apply-capi-dependencies'
)

local_resource('Wait CAPI dependencies resources',
               'make wait-capi-dependencies-resources',
               resource_deps=[
                 'Install CAPI dependencies'
               ]
)

local_resource('Install Cluster API',
               'make apply-capi'
)

local_resource('Wait CAPI resources',
               'make wait-capi-resources',
               resource_deps=[
                 'Install Cluster API'
               ]
)

local_resource('Populate test clusters into CAPI',
               'make apply-test-clusters',
               deps=['Makefile'],
               resource_deps=[
                 'Wait CAPI resources'
               ]
)

local_resource('Build KaaS-manager binary',
               'make build',
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

k8s_yaml('.kubernetes/dev/manifest.yaml')

k8s_resource('manager', port_forwards=8080)
