# AnsibleHound

![](./images/ansiblehound.png)

## Overview

**AnsibleHound** is a BloodHound OpenGraph collector for **Ansible WorX** and **Ansible Tower**. The collector is designed to map the structure and permission of your organization into a navigable attack‑path graph.

Developped by [@Ramoreik](https://github.com/Ramoreik) and [@s_lck](https://github.com/s-lck).

## Collector Setup & Usage

### Building the tool

```bash
go build . -o build/collector
```

### Running the Collection

The collector can be run using any of the following authentication materials:

- `token` who is working only for local users
- `username/password` who is working for both local and LDAP users

The collector will then list the access permissions available to the user. The collection's results depend on the user's level of access used for collection.

> Note : If you have multiple instances of Ansible you need to run the collector against each of them

#### Token

To obtain a valid token for **Ansible WorX** or **Ansible Tower**, you can navigate to the **User Details** of your current user.

![](./images/user-details.png)

Then the **tokens** tab.

![](./images/tokens-tab.png)

Finally, create a token and give it **Read** permissions.

![](./images/create-token.png)

To run the collector, provide it with a target and a token:

```bash
./collector -t '<ansible-url>' --token '<token>'

# Example
./collector -t 'http://localhost:8080/' --token '56KOmh...'
```

> Using local authentication will prevent you from connecting the Ansible and Active Directory graphs. In this way, the `SyncedToATUser` edge will not appear.

#### Username/Password

To run the collector, provide it with a username and a password:

```bash
./collector -u '<username>' -p '<password>' -t '<ansible-url>'

# Example
./collector -u 'admin' -p 'tcrA...' -t 'http://10.10.10.10:8080'
```

##### Active Directory accounts

For Active Directory / LDAP accounts accounts, provide the domain username and password:

```bash
./collector -u '<username>' -p '<password>' -t '<ansible-url>' --dc-ip <dc-ip> --domain <domain-name>

# Example
./collector -u 'admin' -p 'tcrA...' -t 'http://10.10.10.10:8080' --dc-ip '10.10.10.100' --domain 'sleekboy'
```

> Using an Active Directory account will allow you to connect Ansible and Active Directory graphs.

### Load Icons

A script is provided to import the icon for the custom nodes used by AnsibleHound.
You have to provide it the `bloodhound-url` and `jwt-token`.

```bash
python3 ./scripts/import-icons.py <bloodhound-url> <jwt-token>

# Example
python3 ./scripts/import-icons.py 'http://localhost:8080' 'ey[...]'
```

### Samples

If you don't have any Ansible WorX or Tower environment, you can just drop `./samples/example.json` on BHCE to enjoy the graph.

## Schema

![Schema](./images/cypher.png)

### Nodes

Nodes correspond to each object type.

| Node              | Description                                                                                                           | Icon          | Color   |
| ----------------- | --------------------------------------------------------------------------------------------------------------------- | ------------- | ------- |
| ATAnsibleInstance | Complete installation of Ansible                                                                                      | sitemap       | #F59C36 |
| ATOrganization    | Logical collection of users, teams, projects, and inventories. It is the highest-level object in the object hierarchy | building      | #F59C36 |
| ATInventory       | Collection of hosts and groups                                                                                        | network-wired | #FF78F2 |
| ATGroup           | Group of hosts                                                                                                        | object-group  | #159b7c |
| ATUser            | An individual user account                                                                                            | user          | #7ADEE9 |
| ATJob             | Instance launching a playbook against an inventory of hosts                                                           | gears         | #7CAAFF |
| ATJobTemplate     | Combines an Ansible playbook from a project and the settings required to launch it                                    | code          | #493EB0 |
| ATProject         | Logical collection of Ansible playbooks                                                                               | folder-open   | #EC7589 |
| ATCredential      | Authenticate the user to launch playbooks (passwords - SSH keys) against inventory hosts                              | key           | #94E16A |
| ATCredentialType  | Type of the Credential and information about this type.                                                               | key           | #94E16A |
| ATHost            | These are the target devices (servers, network appliances or any computer) you aim to manage                          | desktop       | #E9E350 |
| ATTeam            | A group of users                                                                                                      | people-group  | #724752 |

### Edges

All the edges are prefixed by `AT` to make it distinct from other collectors edges.

#### Ansible edges

Ansible edges only create relations between Ansible nodes:

| Edge Type    | Source              | Target                                                                                       |
| ------------ | ------------------- | -------------------------------------------------------------------------------------------- |
| `ATContains` | `ATAnsibleInstance` | `ATOrganization`                                                                             |
| `ATContains` | `ATOrganization`    | `ATInventory`                                                                                |
| `ATContains` | `ATInventory`       | `ATHost`                                                                                     |
| `ATContains` | `ATInventory`       | `ATGroup`                                                                                    |
| `ATContains` | `ATGroup`           | `ATHost`                                                                                     |
| `ATContains` | `ATJobTemplate`     | `ATJob`                                                                                      |
| `ATContains` | `ATOrganization`    | `ATJobTemplate`                                                                              |
| `ATContains` | `ATOrganization`    | `ATCredential`                                                                               |
| `ATContains` | `ATOrganization`    | `ATProject`                                                                                  |
| `ATUses`     | `ATJobTemplate`     | `ATProject`                                                                                  |
| `ATUses`     | `ATJobTemplate`     | `ATInventory`                                                                                |
| `ATUsesType` | `ATCredential`      | `ATCredentialType`                                                                           |
| `ATExecute`  | `ATUser`            | `ATJobTemplate`                                                                              |
| `ATExecute`  | `ATTeam`            | `ATJobTemplate`                                                                              |
| `ATMember`   | `ATUser`            | `ATOrganization` - `ATTeam`                                                                  |
| `ATRead`     | `ATUser`            | `ATOrganization` - `ATTeam` - `ATInventory` - `ATProject` - `ATJobTemplate`                  |
| `ATRead`     | `ATTeam`            | `ATOrganization` - `ATUser` - `ATInventory` - `ATProject` - `ATJobTemplate`                  |
| `ATAuditor`  | `ATUser`            | `ATOrganization` - `ATProject` - `ATInventory` - `ATJobTemplate`                             |
| `ATAdmin`    | `ATUser`            | `ATOrganization` - `ATTeam` - `ATInventory` - `ATProject` - `ATJobTemplate` - `ATCredential` |

#### Hybrid edges

Hybrid edges establish connections between Ansible and other technologies. AnsibleHound currently handles two types of hybrid edge:

| Edge Type               | Source Graph      | Target Graph | Source Node  | Target Node  |
| ----------------------- | ----------------- | ------------ | ------------ | ------------ |
| `SyncedToATUser`        | Active Directory  | Ansible      | User         | ATUser       |
| `ATHasSourceControlUrl` | Ansible           | GitHub       | ATProject    | GHRepository |
| `ATIsCredentialOf`      | Ansible           | GitHub       | ATCredential | GHUser       |

The following collectors must be used in order to use those hybrid graphs:

- [SharpHound](https://github.com/SpecterOps/SharpHound) for Active Directory
- [GitHound](https://github.com/SpecterOps/GitHound) for GitHub

The output results of these collectors must be uploaded on BloodHound **before** the Ansible output is uploaded. BloodHound will then automatically establish the connection between the graphs.

##### SyncedToATUser

The `SyncedToATUser` edge will allows you to connect Ansible and Active Directory graphs:

![SyncedToATUser hybrid edge](./images/SyncedToATUser.png)

This edge highlights the ability of an Active Directory user to authenticate on the Ansible instance.

The `security identifier` (SID) of the Active Directory user is used to link that user with the Ansible user. The `distinguished name` (DN) stored in the Ansible instance is used to recover the SID of the Active Directory user.

> Since the link is based on the `DN`, there is a risk of collision if two or more domains with the same name share a user with the same SID and name..

##### ATHasSourceControlUrl

The `ATHasSourceControlUrl` edge will allows you to connect Ansible and GitHub graphs:

![ATHasSourceControlUrl hybrid edge](./images/ATHasSourceControlUrl.png)

This edge highlights the link between an Ansible project and a Git Source Control Type.

The name of the repository is used to link the GitHub repository with the Ansible project.

> Since the link is name based, there is a risk of collision if two or more repositories with the same name, hosted on different GitHub accounts, are used in the same Ansible instance.

##### ATIsCredentialOf

The `ATIsCredentialOf` edge will allows you to connect Ansible and GitHub graphs:

![ATIsCredentialOf hybrid edge](./images/ATIsCredentialOf.png)

This edge highlights the link between an Ansible credential and a GitHub user.

The username is used to link the GitHub user with the Ansible credential.

> Because the link is username-based, no edge will be created if the username field of the Ansible credential object is left blank or contains an email address.

## Requirements

### Postgres

As AnsibleHound uses OpenGraph, which is only officially supported with Postgres as the backend, it is recommended to switch to Postgres to avoid ingestion bugs.

Documentation : <https://github.com/SpecterOps/BloodHound/tree/main/examples/docker-compose>

### Ingestion

A minimal version of `8.5.0` is required to handle properly hybrid graph ingestion.

## Licensing

```
                    GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <https://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.
```
