load helpers

@test "podman_container_created_seconds" {
    podman container create --pod $TEST_POD --name $TEST_CNT $TEST_HTTP_IMAGE

    sleep 4

    pod_id=$(podman pod inspect $TEST_POD -f "{{ .Id }}")
    cnt_id=$(podman container inspect $TEST_CNT -f "{{ .Id }}")

    curl -s $METRICS_URL | grep "podman_container_created_seconds{id=\"${cnt_id:0:12}\",pod_id=\"${pod_id:0:12}\",pod_name=\"$TEST_POD\"}"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}

@test "podman_container_state" {
    pod_id=$(podman pod inspect $TEST_POD -f "{{ .Id }}")
    cnt_id=$(podman container inspect $TEST_CNT -f "{{ .Id }}")

    curl -s $METRICS_URL | grep "podman_container_state{id=\"${cnt_id:0:12}\",pod_id=\"${pod_id:0:12}\",pod_name=\"$TEST_POD\"} 0"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    podman container start $TEST_CNT
    curl -s $METRICS_URL | grep "podman_container_state{id=\"${cnt_id:0:12}\",pod_id=\"${pod_id:0:12}\",pod_name=\"$TEST_POD\"} 2"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}

@test "podman_container_exited_seconds" {
    pod_id=$(podman pod inspect $TEST_POD -f "{{ .Id }}")
    cnt_id=$(podman container inspect $TEST_CNT -f "{{ .Id }}")

    curl -s $METRICS_URL | grep "podman_container_exited_seconds{id=\"${cnt_id:0:12}\",pod_id=\"${pod_id:0:12}\",pod_name=\"$TEST_POD\"}"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}

@test "podman_container_exit_code" {
    pod_id=$(podman pod inspect $TEST_POD -f "{{ .Id }}")
    cnt_id=$(podman container inspect $TEST_CNT -f "{{ .Id }}")

    curl -s $METRICS_URL | grep "podman_container_exit_code{id=\"${cnt_id:0:12}\",pod_id=\"${pod_id:0:12}\",pod_name=\"$TEST_POD\"}"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}


@test "podman_container_info" {
    pod_id=$(podman pod inspect $TEST_POD -f "{{ .Id }}")
    cnt_id=$(podman container inspect $TEST_CNT -f "{{ .Id }}")

    curl -s $METRICS_URL | grep "podman_container_info{id=\"${cnt_id:0:12}\",image=\"$TEST_HTTP_IMAGE:latest\",name=\"$TEST_CNT\",pod_id=\"${pod_id:0:12}\",pod_name=\"$TEST_POD\",ports=\"\"} 1"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    podman container create --name $TEST_CNT_LABEL --label label01=value01 $TEST_BUSYBOX_IMAGE
    cnt_id=$(podman container inspect $TEST_CNT_LABEL -f "{{ .Id }}")

    curl -s $METRICS_URL | grep "podman_container_info{id=\"${cnt_id:0:12}\",image=\"$TEST_BUSYBOX_IMAGE:latest\",label01=\"value01\",name=\"$TEST_CNT_LABEL\",pod_id=\"\",pod_name=\"\",ports=\"\"} 1"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}
