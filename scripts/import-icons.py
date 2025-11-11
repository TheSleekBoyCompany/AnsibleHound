import requests
import urllib3
import json
import sys

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

CUSTOM_NODE_ENDPOINT = "/api/v2/custom-nodes"

def define_icon(url, token, icon_type, icon_name, icon_color):
    payload = {
        "custom_types": {
            icon_type: {
                "icon": {
                    "type": "font-awesome",
                    "name": icon_name,
                    "color": icon_color
                }
            }
        }
    }


    headers = {
        "Authorization": "Bearer %s" % token ,
        "Content-Type": "application/json"
    }

    response = requests.post(
        url,
        headers=headers,
        json=payload,
        verify=False  # Disables SSL verification
    )

    print(f"🔹 Sent icon for: {icon_type}")
    print("Status Code:", response.status_code)
    print("Response Body:", response.text)
    print("---")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("[!] Missing arguments.", file=sys.stderr)
        print("[?] Usage: import-icons.py <bloodhound-url> <jwt-token>", file=sys.stderr)
        print("[?] Example: import-icons.py http://127.0.0.1:8080 ey[...]", file=sys.stderr)
        sys.exit(1)

    url = sys.argv[1] + CUSTOM_NODE_ENDPOINT
    jwt_token = sys.argv[2]

    # Call function for each icon type you want to send
    define_icon(url, jwt_token, "ATAnsibleInstance", "sitemap", "#E43131")
    define_icon(url, jwt_token, "ATOrganization", "building", "#F59C36")
    define_icon(url, jwt_token, "ATInventory", "network-wired", "#FF78F2")
    define_icon(url, jwt_token, "ATUser", "user", "#7ADEE9")
    define_icon(url, jwt_token, "ATJob", "gears", "#7CAAFF")
    define_icon(url, jwt_token, "ATJobTemplate", "code", "#493EB0")
    define_icon(url, jwt_token, "ATProject", "folder-open", "#EC7589")
    define_icon(url, jwt_token, "ATCredential", "key", "#94E16A")
    define_icon(url, jwt_token, "ATCredentialType", "gear", "#94E16A")
    define_icon(url, jwt_token, "ATHost", "desktop", "#E9E350")
    define_icon(url, jwt_token, "ATTeam", "people-group", "#724752")
    define_icon(url, jwt_token, "ATGroup", "object-group", "#159b7c")
