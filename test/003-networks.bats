load helpers

@test "podman_network_info" {
    podman network create $TEST_NETWORK
    network_inspect=$(podman network inspect $TEST_NETWORK -f "{{ .Id }} {{ .Driver }} {{ .NetworkInterface }}")
    network_id=$(echo $network_inspect | awk '{print $1}')
    network_driver=$(echo $network_inspect | awk '{print $2}')
    network_interface=$(echo $network_inspect | awk '{print $3}')

    curl -s $METRICS_URL | grep "podman_network_info{driver=\"$network_driver\",id=\"${network_id:0:12}\",interface=\"$network_interface\",labels=\"\",name=\"$TEST_NETWORK\"} 1"
    if [ $? -ne 0 ] ; then
        false
        exit
    fi

    true
}
