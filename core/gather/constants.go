package gather

import "fmt"

const AAP_CONTROLLER_ENDPOINT = "/api/controller/v2/"
const AAP_TOKEN_ENDPOINT = "/api/gateway/v1/tokens/"
const TOWER_API_ENDPOINT = "/api/v2/"
const ME_ENDPOINT = "%s/me/"
const ORGANIZATIONS_ENDPOINT = "%s/organizations/"
const PROJECTS_ENDPOINT = "%s/projects/"
const INVENTORIES_ENDPOINT = "%s/inventories/"
const JOB_TEMPLATE_ENDPOINT = "%s/job_templates/"
const CREDENTIALS_ENDPOINT = "%s/credentials/"
const CREDENTIAL_TYPES_ENDPOINT = "%s/credential_types/"
const USERS_ENDPOINT = "%s/users/"
const USER_ROLES_ENDPOINT = "%s/users/%d/roles/"
const GROUPS_ENDPOINT = "%s/groups/"
const GROUP_HOSTS_ENDPOINT = "%s/groups/%d/hosts/"
const JOBS_ENDPOINT = "%s/jobs/"
const WORKFLOW_JOB_TEMPLATES_ENDPOINT = "%s/workflow_job_templates/"
const WORKFLOW_JOB_TEMPLATE_NODES_ENDPOINT = "%s/workflow_job_template_nodes/"
const HOSTS_ENDPOINT = "%s/hosts/"
const TEAMS_ENDPOINT = "%s/teams/"
const TEAM_ROLES_ENDPOINT = "%s/teams/%d/roles/"
const TEAM_USERS_ENDPOINT = "%s/teams/%d/users/"
const JOB_TEMPLATE_CREDENTIALS_ENDPOINT = "%s/job_templates/%d/credentials/"

const PING_ENDPOINT = "%s/ping"

const PAGE_SIZE = 200

var PAGE_SIZE_ARG = fmt.Sprintf("?page_size=%d", PAGE_SIZE)
var CURRENT_PAGE_ARG = "&page=%d"
