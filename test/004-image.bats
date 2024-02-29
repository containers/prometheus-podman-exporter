load helpers

@test "podman_image_created_seconds" {
    podman image pull $TEST_BUSYBOX_IMAGE
    image_inspect=$(podman image inspect $TEST_BUSYBOX_IMAGE -f "{{ .Id }} {{ .Digest }}")
    image_id=$(echo $image_inspect | awk '{print $1}')

    sleep 4
    curl -s $METRICS_URL | grep "podman_image_created_seconds{id=\"${image_id:0:12}\",repository=\"$TEST_BUSYBOX_IMAGE\",tag=\"latest\"}"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}

@test "podman_image_info" {
    image_inspect=$(podman image inspect $TEST_BUSYBOX_IMAGE -f "{{ .Id }} {{ .Digest }}")
    image_id=$(echo $image_inspect | awk '{print $1}')
    image_digest=$(echo $image_inspect | awk '{print $2}')

    curl -s $METRICS_URL | grep "podman_image_info{digest=\"$image_digest\",id=\"${image_id:0:12}\",parent_id=\"\",repository=\"$TEST_BUSYBOX_IMAGE\",tag=\"latest\"}"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}

@test "podman_image_size" {
    image_inspect=$(podman image inspect $TEST_BUSYBOX_IMAGE -f "{{ .Id }} {{ .Size }}")
    image_id=$(echo $image_inspect | awk '{print $1}')
    image_size=$(echo $image_inspect | awk '{print $2}')

    curl -s $METRICS_URL | grep "podman_image_size{id=\"${image_id:0:12}\",repository=\"$TEST_BUSYBOX_IMAGE\",tag=\"latest\"} ${image_size:0:1}"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}
