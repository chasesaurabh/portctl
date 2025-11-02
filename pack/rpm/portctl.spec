Name: portctl
Version: 0.0.0
Release: 1%{?dist}
Summary: portctl — discover and kill processes listening on a TCP port
License: MIT
URL: https://github.com/chasesaurabh/portctl
Source0: %{name}-%{version}.tar.gz
BuildArch: x86_64

%description
A small CLI to discover processes binding a TCP port and optionally send signals to them.

%prep
%setup -q

%build
# no compiled build here — we will place binary in %{buildroot}/usr/local/bin

%install
mkdir -p %{buildroot}/usr/local/bin
install -m 0755 portctl %{buildroot}/usr/local/bin/portctl

%files
/usr/local/bin/portctl

%changelog
* $(date +"%a %b %d %Y") Saurabh Chase <chasesaurabh@gmail.com> - 0.0.0-1
- initial package
