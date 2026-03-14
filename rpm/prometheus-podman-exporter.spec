%bcond_without check
%bcond_without bundled
%if 0%{?rhel}
%bcond_without bundled
%endif

%if %{defined rhel} && 0%{?rhel} < 10
%define gobuild(o:) go build -buildmode pie -compiler gc -tags="rpm_crashtraceback ${BUILDTAGS:-}" -ldflags "-linkmode=external -compressdwarf=false ${LDFLAGS:-} -B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \\n') -extldflags '%__global_ldflags'" -a -v -x %{?**};
%endif

%if %{with bundled}
%global gomodulesmode   GO111MODULE=on
%endif

# https://github.com/containers/prometheus-podman-exporter
%global goipath         github.com/containers/prometheus-podman-exporter
Version: 0

%gometa -f

%global goname prometheus-podman-exporter

%global common_description %{expand:
Prometheus exporter for podman environments exposing containers, pods, images,
volumes and networks information.}

%global golicenses      LICENSE
%global godocs          CODE_OF_CONDUCT.md CONTRIBUTING.md MAINTAINERS.md\\\
                        README.md SECURITY.md

Name:           %{goname}
Release:        %{?autorelease}%{!?autorelease:1%{?dist}}
Summary:        Prometheus exporter for podman environment

License:        Apache-2.0 AND MPL-2.0 AND BSD-3-Clause AND BSD-2-Clause AND MIT AND Unlicense AND CC-BY-SA-4.0 AND ISC
URL:             %{gourl}
Source0:         %{gosource}
Source1:         vendor-%{version}.tar.gz

%if 0%{?fedora} && ! 0%{?rhel}
BuildRequires: pkgconfig(libbtrfsutil)
%endif
BuildRequires: gcc
BuildRequires: glibc-devel
BuildRequires: glibc-static
BuildRequires: git-core
%if 0%{?rhel} >= 9
BuildRequires: go-rpm-macros
%endif
BuildRequires: golang
BuildRequires: make
BuildRequires: pkgconfig(devmapper)
BuildRequires: pkgconfig(glib-2.0)
BuildRequires: pkgconfig(gpgme)
BuildRequires: pkgconfig(libassuan)
%if 0%{?fedora} >= 37
BuildRequires: shadow-utils-subid-devel
%endif

%description %{common_description}

%prep
%goprep %{?with_bundledc:-k}
%if %{with bundled}
%setup -q -T -D -a 1 -n %{name}-%{version}
%endif

%if %{without bundled}
%generate_buildrequires
%go_generate_buildrequires
%endif

%build
%if %{with bundled}
export GOFLAGS="-mod=vendor"
%endif

%if 0%{?rhel} >= 9
export BUILDTAGS="exclude_graphdriver_btrfs btrfs_noversion systemd libtrust_openssl"
%endif

%if 0%{?fedora}
export BUILDTAGS="systemd"
%endif

export LDFLAGS="-X %{goipath}/cmd.buildVersion=%{version} -X %{goipath}/cmd.buildRevision=%{release} -X %{goipath}/cmd.buildBranch=main"

%gobuild -o %{gobuilddir}/bin/prometheus-podman-exporter %{goipath}

%install
install -m 0755 -vd                     %{buildroot}%{_bindir}
install -m 0755 -vp %{gobuilddir}/bin/* %{buildroot}%{_bindir}/
install -m 0755 -vd                     %{buildroot}%{_unitdir}
install -m 0755 -vd                     %{buildroot}%{_userunitdir}
install -m 0755 -vd                     %{buildroot}%{_sysconfdir}/sysconfig/
install -m 0644 -vp ./contrib/systemd/system/%{name}.sysconfig %{buildroot}%{_sysconfdir}/sysconfig/%{name}
install -m 0644 -vp ./contrib/systemd/system/%{name}.service %{buildroot}%{_unitdir}/%{name}.service
install -m 0644 -vp ./contrib/systemd/user/%{name}.service %{buildroot}%{_userunitdir}/%{name}.service

%post
%systemd_user_post %{name}.service
%systemd_post %{name}.service

%preun
%systemd_user_preun %{name}.service
%systemd_preun %{name}.service

%if %{with check}
%check
%endif

%files
%license LICENSE
%doc CODE_OF_CONDUCT.md CONTRIBUTING.md MAINTAINERS.md README.md SECURITY.md
%{_bindir}/%{name}
%{_unitdir}/%{name}.service
%{_userunitdir}/%{name}.service
%config(noreplace) %{_sysconfdir}/sysconfig/%{name}

%changelog
%autochangelog
