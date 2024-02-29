load helpers

@test "podman_system_api_version" {
    api_version=$(podman system info -f json | jq -r '.version.Version' | awk -F . '{print $1,$2}' | tr " " ".")
    output=$(curl -s $METRICS_URL | grep ^podman_system_api_version.*)

    echo $output | grep "podman_system_api_version{version=\"$api_version"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}

@test "podman_system_buildah_version" {
    buildah_version=$(podman system info -f json | jq -r '.host.buildahVersion' | awk -F . '{print $1,$2}' | tr " " ".")
    output=$(curl -s $METRICS_URL | grep ^podman_system_buildah_version.*)

    echo $output | grep "podman_system_buildah_version{version=\"$buildah_version"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}

@test "podman_system_conmon_version" {
    conmon_version=$(podman system info -f json | jq -r '.host.conmon.version' | awk -F, '{print $1}' | awk '{print $3}')
    output=$(curl -s $METRICS_URL | grep ^podman_system_conmon_version.*)

    echo $output | grep "podman_system_conmon_version{version=\"$conmon_version\"} 1"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}

@test "podman_system_runtime_version" {
    runtime_version=$(podman system info -f json | jq -r '.host.ociRuntime.version' | head -1)
    output=$(curl -s $METRICS_URL | grep ^podman_system_runtime_version.*)

    echo $output | grep "podman_system_runtime_version{version=\"$runtime_version\"} 1"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}
