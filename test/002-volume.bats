load helpers

@test "podman_volume_created_seconds" {
    podman volume create $TEST_VOLUME

    curl -s $METRICS_URL | grep "podman_volume_created_seconds{name=\"$TEST_VOLUME\"}"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}

@test "podman_volume_info" {
    vol_inspect=$(podman volume inspect $TEST_VOLUME -f "{{ .Driver }} {{ .Mountpoint }}")
    vol_driver=$(echo $vol_inspect | awk '{print $1}')
    vol_mountpoint=$(echo $vol_inspect | awk '{print $2}')

    curl -s $METRICS_URL | grep "podman_volume_info{driver=\"$vol_driver\",mount_point=\"$vol_mountpoint\",name=\"$TEST_VOLUME\"} 1"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}
