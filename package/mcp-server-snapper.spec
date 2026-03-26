#
# spec file for package mcp-server-snapper
#
# Copyright (c) 2025 SUSE LLC
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
Version:        0.1.0
Release:        0
Summary:        MCP Server for Snapper
License:        MIT
Source:         %{name}-%{version}.tar.xz
BuildArch:      noarch
BuildRequires:  python313-devel

%description
An MCP server for Snapper.

%prep
%autosetup

%build

%install
install -d -m 0755 %{buildroot}%{_bindir}
install -m 0755 src/mcp-server-snapper %{buildroot}%{_bindir}/mcp-server-snapper

%files
%license LICENSE
%doc README.md
%{_bindir}/mcp-server-snapper

%changelog
