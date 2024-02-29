load helpers

@test "podman_pod_created_seconds" {
    podman pod create $TEST_POD
    pod_id=$(podman pod inspect $TEST_POD -f "{{ .Id }}")

    curl -s $METRICS_URL | grep "podman_pod_created_seconds{id=\"${pod_id:0:12}\"}"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}

@test "podman_pod_info" {
    pod_inspect=$(podman pod inspect $TEST_POD -f "{{ .Id }} {{ .Containers }}")
    pod_id=$(echo $pod_inspect | awk '{print $1}')
    pod_infra_id=$(echo $pod_inspect | awk '{print $2}' | awk -F{ '{print $2}')

    curl -s $METRICS_URL | grep "podman_pod_info{id=\"${pod_id:0:12}\",infra_id=\"${pod_infra_id:0:12}\",name=\"$TEST_POD\"}"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}

@test "podman_pod_state" {
    pod_id=$(podman pod inspect $TEST_POD -f "{{ .Id }}")

    curl -s $METRICS_URL | grep "podman_pod_state{id=\"${pod_id:0:12}\"} 0"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}

@test "podman_pod_containers" {
    pod_id=$(podman pod inspect $TEST_POD -f "{{ .Id }}")

    curl -s $METRICS_URL | grep "podman_pod_containers{id=\"${pod_id:0:12}\"} 1"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    podman container create --pod $TEST_POD $TEST_BUSYBOX_IMAGE
    sleep 4

    curl -s $METRICS_URL | grep "podman_pod_containers{id=\"${pod_id:0:12}\"} 2"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}
