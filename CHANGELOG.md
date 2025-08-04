# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Now able to skip TLS/SSL verification using the `-k`/`skip-verify-ssl` flag.
- Now gathers information on the target Ansible WorX/Tower instance using `/ping` before gathering resources.
- Now gathers Credentials used by Job Templates and creates an `ATUses` edge between them.


### Changed

- Moved `ToBHNode` to the objects themselves and now uses an interface to interact with it.
- Remodeled the object system to enable better code patterns.
- Changed Data Structure to handle nodes and edge from `[]AnsibleType` to `map[int]AnsibleType` to enable direct mapping using the Ansible ID.
- Renamed Edge generation function.
- OID for nodes is now based on the following data: `sha1(INSTALL_UUID + ID + RESOURCE_TYPE)`, allowing reproducible OIDs that stay different between Ansible instances.

[unreleased] https://github.com/TheSleekBoyCompany/AnsibleHound/?ref=main