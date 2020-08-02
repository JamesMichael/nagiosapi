%undefine _missing_build_ids_terminate_build
%define debug_package %{nil}

Name:         nagios-api
Version:      0.1.0
Release:      2%{?dist}
Summary:      Nagios API Server
URL:          https://github.com/JamesMichael/nagiosapi
Source0:      https://github.com/JamesMichael/nagiosapi/archive/v%{version}.tar.gz
License:      Public Domain

BuildRequires: git
BuildRequires: golang

%if 0%{?rhel} != 7
BuildRequires: systemd-rpm-macros
%endif

%description
%{summary}

%package server
Summary: Nagios API Server

%description server
%{summary}

%prep
%setup

%build
go mod vendor
go build -mod vendor -o nagios-api main.go

%install

install -m 0755 -d %{buildroot}%{_bindir}
install -m 0755 -d %{buildroot}%{_libexecdir}
install -m 0755 -d %{buildroot}%{_sysconfdir}
install -m 0755 -d %{buildroot}%{_sysconfdir}
install -m 0755 -d %{buildroot}%{_unitdir}

cp -a nagios-api %{buildroot}%{_libexecdir}
cp -a etc/* %{buildroot}%{_sysconfdir}
cp -a usr/lib/systemd/system/* %{buildroot}%{_unitdir}

%check

%files server
%defattr(-,root,root,-)
%license LICENSE
%doc README.md

%config(noreplace) %{_sysconfdir}/nagiosapi/server.yaml
%{_libexecdir}/nagios-api-server
%{_unitdir}/nagios-api-server.service

%post server
%systemd_post nagios-api-server.service

%preun server
%systemd_preun nagios-api-server.service

%postun server
%systemd_postun_with_restart nagios-api-server.service

%changelog
* Sun Aug 2 2020 James Michael <jamesamichael@gmail.com>
- Initial package
