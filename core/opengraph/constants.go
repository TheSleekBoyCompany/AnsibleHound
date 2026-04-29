package opengraph

const LDAP_VALUE = "ldap"

const GIT_SCM_TYPE = "git"
const DOT_GIT_SCM_TYPE = "." + GIT_SCM_TYPE

const MATCH_BY_ID = "id"
const MATCH_BY_NAME = "name"

const ACTIVE_DIRECTORY_BASE = "Base"
const ANSIBLE_BASE = "AnsibleBase"
const GITHUB_BASE = "GHBase"

const CREDENTIAL_USERNAME = "username"
const CREDENTIAL_KIND = "scm"

const ORGANIZATION_RESOURCE_TYPE = "organization"
const INVENTORY_RESOURCE_TYPE = "inventory"
const TEAM_RESOURCE_TYPE = "team"
const CREDENTIAL_RESOURCE_TYPE = "credential"
const JOB_TEMPLATE_RESOURCE_TYPE = "job_template"
const WORKFLOW_JOB_TEMPLATE_RESOURCE_TYPE = "workflow_job_template"

const SCM_CREDENTIAL_TYPE = "Source Control"
const SECRET_SERVER_CREDENTIAL_TYPE = "Thycotic Secret Server"
const HASHICORP_VAULT_CREDENTIAL_TYPE = "HashiCorp Vault Secret Lookup"
const MACHINE_CREDENTIAL_TYPE = "Machine"

const MACHINE_CREDENTIAL_SSH_SUBTYPE = "ssh"
const MACHINE_CREDENTIAL_PASSWORD_SUBTYPE = "password"

const ADMIN_ROLE_EDGE = "ATAdmin"
const AUTIDOR_ROLE_EDGE = "ATAuditor"
const USE_ROLE_EDGE = "ATUse"
const INVENTORY_ADMIN_ROLE_EDGE = "ATInventoryAdmin"

const CONTAINS_EDGE = "ATContains"
const USES_EDGE = "ATUses"
const USES_TYPE_EDGE = "ATUsesType"

const GROUP_NODE = "ATGroup"
const HOST_NODE = "ATHost"
const CREDENTIAL_NODE = "ATCredential"
const PROJECT_NODE = "ATProject"
const ORGANIZATION_NODE = "ATOrganization"
const INVENTORY_NODE = "ATInventory"
const JOB_TEMPLATE_NODE = "ATJobTemplate"

const SYNCED_TO_USER_AD_EDGE = "SyncedToATUser"
const HAS_SOURCE_CONTROL_URL_GITHUB_EDGE = "ATHasSourceControlUrl"
const IS_CREDENTIAL_OF_GITHUB_EDGE = "ATIsCredentialOf"

const COMPROMISE_WITH_PLAYBOOK_POST_PROCESSING_EDGE = "ATCompromiseWithMaliciousPlaybook"
const CAN_USE_IN_ADHOC_COMMANDS_POST_PROCESSING_EDGE = "ATCanUseInADHOCCommands"
const COMPROMISE_WITH_HONEYPOT_POST_PROCESSING_EDGE = "ATCompromiseWithHoneypot"
const SSH_HIJACK_AGENT_POST_PROCESSING_EDGE = "ATSSHHijackAgent"
const VALID_FOR_POST_PROCESSING_EDGE = "ATValidFor"
const COMPROMISE_WITH_REQUESTBIN_POST_PROCESSING_EDGE = "ATCompromiseWithRequestbin"
const COMPROMISE_WITH_FAKE_SCM_SERVER_POST_PROCESSING_EDGE = "ATCompromiseWithFakeSCM"
