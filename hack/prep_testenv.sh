#!/usr/bin/env bash

INSTALL_RPMS="podman fuse-overlayfs openssh-clients cpp git-core sqlite make gcc glib2-devel glibc-devel glibc-static device-mapper-devel"

echo -e "\n\n# Added during image build" >> /etc/dnf/dnf.conf
echo -e "minrate=100\ntimeout=60\n" >> /etc/dnf/dnf.conf

dnf -y makecache
dnf -y update
rpm --setcaps shadow-utils 2>/dev/null
dnf -y install fedora-repos-rawhide
dnf -y install $INSTALL_RPMS --exclude container-selinux --enablerepo=rawhide
dnf clean all

rm -fv /etc/machine-id /var/lib/systemd/random-seed /var/lib/dnf/repos/*/countme
rm -fv /usr/lib/systemd/profile.d/*
rm -rf /var/cache /var/log/dnf* /var/log/hawkey.log /var/log/yum.*

useradd podman -d /home/podman/
mkdir -p /home/podman/.local/share/containers
chown podman:podman -R /home/podman

echo -e "root:1:65535\npodman:1:999\npodman:1001:64535" > /etc/subuid
echo -e "root:1:65535\npodman:1:999\npodman:1001:64535"

cat << EOF > /etc/containers/containers.conf
[containers]
netns="host"
userns="host"
ipcns="host"
utsns="host"
cgroupns="host"
cgroups="disabled"
log_driver = "k8s-file"
[engine]
cgroup_manager = "cgroupfs"
events_logger="file"
runtime="crun"
EOF

cat <<EOF > /home/podman/.config/containers/containers.conf
[containers]
volumes = [
	"/proc:/proc",
]
default_sysctls = []
EOF

chmod 644 /etc/containers/containers.conf

sed -e 's|^#mount_program|mount_program|g' \
    -e 's|^mountopt[[:space:]]*=.*$|mountopt = "nodev,fsync=0"|g' \
    /usr/share/containers/storage.conf \
    > /etc/containers/storage.conf

printf '/run/secrets/etc-pki-entitlement:/run/secrets/etc-pki-entitlement\n/run/secrets/rhsm:/run/secrets/rhsm\n' > /etc/containers/mounts.conf

mkdir -p /var/lib/shared/overlay-images
mkdir -p /var/lib/shared/overlay-layers
mkdir -p /var/lib/shared/vfs-images
mkdir -p /var/lib/shared/vfs-layers
touch /var/lib/shared/overlay-images/images.lock
touch /var/lib/shared/overlay-layers/layers.lock
touch /var/lib/shared/vfs-images/images.lock
touch /var/lib/shared/vfs-layers/layers.lock
