#
# spec file for package mcp-server-snapper
#
# Copyright (c) 2026 SUSE LLC
#
# All modifications and additions to the file contributed by third parties
# remain the property of their copyright owners, unless otherwise agreed
# upon. The license for this file, and modifications and additions to the
# file, is the same license as for the pristine package itself (unless the
# license for the pristine package is not an Open Source License, in which
# case the license is the MIT License). An "Open Source License" is a
# license that conforms to the Open Source Definition (Version 1.9)
# published by the Open Source Initiative.

# Please submit bugfixes or comments via https://bugs.opensuse.org/

Name:           mcp-server-snapper
Version:        0.3.0
Release:        0
Summary:        MCP Server for Snapper
License:        MIT
URL:            https://github.com/aschnell/mcp-server-snapper
Source:         %{name}-%{version}.tar.xz
BuildRequires:  go >= 1.22
Requires:       snapper

%description
An MCP server for Snapper.

%prep
%setup -q

%build
./build.sh

%check
for test in tools/tools ; do
    echo "Running $test..."
    MCPSERVER=src/mcp-server-snapper "testsuite/$test" || { echo "Test $test failed!" ; exit 1; }
done

%install
install -d -m 0755 %{buildroot}%{_bindir}
install -m 0755 src/mcp-server-snapper %{buildroot}%{_bindir}/mcp-server-snapper

install -d -m 0755 %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite
install -m 0644 testsuite/README %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/README
install -d -m 0755 %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/create-snapshot-1
install -m 0755 testsuite/create-snapshot-1/create-snapshot-1 %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/create-snapshot-1/create-snapshot-1
install -d -m 0755 %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/create-snapshot-2
install -m 0755 testsuite/create-snapshot-2/create-snapshot-2 %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/create-snapshot-2/create-snapshot-2
install -d -m 0755 %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/get-config
install -m 0755 testsuite/get-config/get-config %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/get-config/get-config
install -d -m 0755 %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/list-configs
install -m 0755 testsuite/list-configs/list-configs %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/list-configs/list-configs
install -d -m 0755 %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/list-snapshots
install -m 0755 testsuite/list-snapshots/list-snapshots %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/list-snapshots/list-snapshots
install -d -m 0755 %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/rollback
install -m 0755 testsuite/rollback/rollback %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/rollback/rollback
install -d -m 0755 %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/tools
install -m 0755 testsuite/tools/tools %{buildroot}%{_prefix}/lib/mcp-server-snapper/testsuite/tools/tools

%files
%license LICENSE
%doc README.md
%{_bindir}/mcp-server-snapper

%package testsuite
Summary:        Testsuite for package %{name}
Requires:       %{name}

%description testsuite
Testsuite for package %{name}

Note: This package is for testing purposes only. It is intended for
use by quality assurance and requires a dedicated testing environment.

Do not install on a production system!

%files testsuite
%{_prefix}/lib/mcp-server-snapper/
%{_prefix}/lib/mcp-server-snapper/testsuite/

%changelog
