package core

import "fmt"

const VERSION = "0.0.1"

const API_ENDPOINT = "/api/v2/"
const ORGANIZATIONS_ENDPOINT = API_ENDPOINT + "organizations/"
const PROJECTS_ENDPOINT = API_ENDPOINT + "projects/"
const INVENTORIES_ENDPOINT = API_ENDPOINT + "inventories/"
const JOB_TEMPLATE_ENDPOINT = API_ENDPOINT + "job_templates/"
const CREDENTIALS_ENDPOINT = API_ENDPOINT + "credentials/"
const CREDENTIAL_TYPES_ENDPOINT = API_ENDPOINT + "credential_types/"
const USERS_ENDPOINT = API_ENDPOINT + "users/"
const USER_ROLES_ENDPOINT = API_ENDPOINT + "users/%d/roles/"
const GROUPS_ENDPOINT = API_ENDPOINT + "groups/"
const GROUP_HOSTS_ENDPOINT = API_ENDPOINT + "groups/%d/hosts/"
const JOBS_ENDPOINT = API_ENDPOINT + "jobs/"
const HOSTS_ENDPOINT = API_ENDPOINT + "hosts/"
const TEAMS_ENDPOINT = API_ENDPOINT + "teams/"
const TEAM_ROLES_ENDPOINT = API_ENDPOINT + "teams/%d/roles/"
const TEAM_USERS_ENDPOINT = API_ENDPOINT + "teams/%d/users/"
const JOB_TEMPLATE_CREDENTIALS_ENDPOINT = API_ENDPOINT + "job_templates/%d/credentials/"
const PING_ENDPOINT = API_ENDPOINT + "ping"
const PAGE_SIZE = 200

var PAGE_SIZE_ARG = fmt.Sprintf("?page_size=%d", PAGE_SIZE)
var CURRENT_PAGE_ARG = "&page=%d"
