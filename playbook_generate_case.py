import limacharlie
import json
import ssl
import urllib.request

# LimaCharlie platform integer to human-readable name mapping
LC_PLATFORM_MAP = {
    "268435456": "windows",
    "536870912": "linux",
    "1073741824": "macos",
}

def json_to_html(json_data, indent_level=0):
    """
    Converts a Python object (derived from JSON) into a readable HTML string.

    Args:
        json_data: The Python object (dict, list, str, int, float, bool, None)
                   to convert to HTML.
        indent_level (int): Current indentation level for formatting (internal use).

    Returns:
        str: An HTML string representing the JSON data.
    """
    indent = "    " * indent_level
    html_output = ""

    if isinstance(json_data, dict):
        html_output += f'{indent}<div class="bg-gray-700 pl-4 rounded-lg shadow-md mb-1 text-white">\n'
        for key, value in json_data.items():
            html_output += f'{indent}  <div class="flex items-baseline mb-1">\n'
            html_output += f'{indent}    <strong class="text-blue-300 mr-2">{key}:</strong>\n'
            html_output += f'{indent}    {json_to_html(value, indent_level + 1)}\n'
            html_output += f'{indent}  </div>\n'
        html_output += f'{indent}</div>\n'
    elif isinstance(json_data, list):
        html_output += f'{indent}<ul class="list-disc list-inside bg-gray-700 pl-4 rounded-lg shadow-md mb-1 text-white">\n'
        if not json_data:
            html_output += f'{indent}  <li class="text-gray-400"><em>(empty list)</em></li>\n'
        else:
            for item in json_data:
                html_output += f'{indent}  <li class="mb-0">{json_to_html(item, indent_level + 1)}</li>\n'
        html_output += f'{indent}</ul>\n'
    else:
        value_str = json.dumps(json_data)
        if isinstance(json_data, str):
            html_output += f'<span class="text-green-300">{value_str}</span>'
        elif isinstance(json_data, (int, float)):
            html_output += f'<span class="text-purple-300">{value_str}</span>'
        elif isinstance(json_data, bool):
            html_output += f'<span class="text-orange-300">{value_str}</span>'
        elif json_data is None:
            html_output += f'<span class="text-red-300">null</span>'
        else:
            html_output += f'<span class="text-gray-300">{value_str}</span>'

    return html_output

def create_html(json_object, permalink):
    """
    Creates a complete HTML page with the converted JSON data.
    """
    html_content = ""
    for k in json_object.keys():
      heading = k.capitalize()
      html_content += f'<span class="detail-heading">{heading}</span>'
      html_content += json_to_html(json_object[k])

    full_html = f"""
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600&display=swap" rel="stylesheet">
    <style>
        .permalink-heading {{
            color: #63b3ed;
            font-size: 1.5rem;
            font-weight: bold;
            text-align: left;
        }}
        .lc-container {{
            width: 100%;
            padding: 24px;
            border-radius: 12px;
            box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
            background-color: #1a202c;
        }}
        .detail-heading {{
            color: #63b3ed;
            margin-bottom: 12px;
            font-size: 1.5rem;
            font-weight: bold;
            text-align: left;
        }}
        .permalink-link {{
            color: #90cdf4;
            margin-bottom: 12px;
            font-weight: normal;
            text-align: left;
        }}
        .text-green-300 {{ color: #68d391; }}
        .text-purple-300 {{ color: #b794f4; }}
        .text-orange-300 {{ color: #f6ad55; }}
        .text-red-300 {{ color: #fc8181; }}
        .text-blue-300 {{ color: #90cdf4; }}
    </style>
    <div class="lc-container">
        {html_content}
        <span class="permalink-heading">Permalink</span><br>
        <a class="permalink-link" href="{permalink}" target="_blank">{permalink}</a><br>
    </div>
"""
    return full_html


def resolve_platform(raw_platform):
    """Resolve LC platform integer to human-readable name."""
    plat_str = str(raw_platform).strip()
    return LC_PLATFORM_MAP.get(plat_str, plat_str)


def iris_request(url, iris_key, ctx, method="GET", body=None):
    """Helper to make DFIR-IRIS API requests."""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {iris_key}"
    }
    data = json.dumps(body).encode('utf-8') if body else None
    req = urllib.request.Request(url, data=data, headers=headers, method=method)
    with urllib.request.urlopen(req, context=ctx) as resp:
        return json.loads(resp.read().decode('utf-8'))


def playbook(sdk, data):
    # Get the secret we need from LimaCharlie.
    # Retrieve the DFIR-IRIS API key and base URL from LimaCharlie Hive secrets.
    # Create these secrets in the LimaCharlie web UI under Organization > Secrets.
    irisKey = limacharlie.Hive(sdk, "secret").get("iris-api-key").data["secret"]
    IRIS_BASE = limacharlie.Hive(sdk, "secret").get("iris-base-url").data["secret"]
    iris_data = data["data"]

    # Parse the details of the base event
    details = json.loads(iris_data["custom_attributes"]["LimaCharlie"]["Details"])
    permalink = iris_data["custom_attributes"]["LimaCharlie"]["Permalink"]

    my_html = create_html(details, permalink)
    iris_data["custom_attributes"]["LimaCharlie"]["Details"] = my_html

    # Extract and remove the debug flag before sending to DFIR-IRIS
    debug = iris_data.pop("debug", "false")

    # Extract sensor data before sending to DFIR-IRIS (not part of case payload)
    sensor = iris_data.pop("sensor", {})

    # Resolve LC platform integer to human-readable name
    platform_name = resolve_platform(sensor.get("platform", ""))

    # if debug is true, return the data as is and don't create a case
    if str(debug).lower() == "true":
        return {"data": iris_data}

    # Skip SSL verification (demo system)
    ctx = ssl.create_default_context()
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE

    # --- Step 1: Create the case ---
    try:
        parsed_response = iris_request(
            f"{IRIS_BASE}/manage/cases/add", irisKey, ctx,
            method="POST", body=iris_data
        )
    except Exception as e:
        return {"error": str(e)}

    case_id = parsed_response.get("data", {}).get("case_id")
    if not case_id:
        return {"data": parsed_response}

    # --- Step 2: Create an Alert and merge it into the case ---
    try:
        case_name = iris_data.get("case_name", "LimaCharlie Detection")
        severity_id = iris_data.get("severity_id", 2)  # Default: Unspecified

        alert_payload = {
            "alert_title": case_name,
            "alert_severity_id": int(severity_id),
            "alert_status_id": 8,  # Escalated (already linked to case)
            "alert_customer_id": int(iris_data.get("case_customer", 1)),
            "alert_source": "LimaCharlie",
            "alert_source_link": permalink,
            "alert_source_content": details,
            "alert_description": iris_data.get("case_description", ""),
            "alert_tags": "limacharlie,automated",
            "alert_note": f"Auto-generated from LimaCharlie detection. Case #{case_id}.",
        }

        alert_resp = iris_request(
            f"{IRIS_BASE}/alerts/add", irisKey, ctx,
            method="POST", body=alert_payload
        )
        alert_id = alert_resp.get("data", {}).get("alert_id")
        parsed_response["alert"] = alert_resp

        # Merge the alert into the case
        if alert_id:
            merge_resp = iris_request(
                f"{IRIS_BASE}/alerts/merge/{alert_id}", irisKey, ctx,
                method="POST", body={
                    "target_case_id": case_id,
                    "iocs_import_list": [],
                    "assets_import_list": [],
                    "note": "",
                    "import_as_event": False,
                    "case_tags": "",
                }
            )
            parsed_response["alert_merge"] = merge_resp
    except Exception as e:
        parsed_response["alert_error"] = str(e)

    # --- Step 3: Add the host as an asset (if SID exists) ---
    asset_id = None
    sid = sensor.get("sid", "")
    if sid and sid != "<no value>":
        try:
            # Resolve asset type from platform name
            platform_keyword_map = {
                "windows": "Windows - Computer",
                "linux": "Linux - Computer",
                "macos": "Mac - Computer",
            }
            target_type_name = platform_keyword_map.get(platform_name, "Linux - Computer")

            # Query DFIR-IRIS for available asset types
            asset_type_id = 1
            try:
                type_data = iris_request(
                    f"{IRIS_BASE}/manage/asset-type/list", irisKey, ctx
                )
                for at in type_data.get("data", []):
                    if at.get("asset_name") == target_type_name:
                        asset_type_id = at["asset_id"]
                        break
            except Exception:
                pass

            # Build isolation button HTML for the Actions custom attribute
            # Set LC_ACTIONS_URL to the host running lc-actions-svc (e.g. https://host:4443)
            actions_base = limacharlie.Hive(sdk, "secret").get("lc-actions-url").data["secret"]
            isolation_url = f"{actions_base}/sensor/{sid}"
            actions_html = (
                f'<a href="{isolation_url}" target="_blank" '
                f'style="display:inline-block;padding:8px 16px;background:#dc2626;'
                f'color:white;border-radius:6px;font-weight:600;text-decoration:none;'
                f'font-size:13px;">Isolate / Manage Host</a>'
            )

            asset_payload = {
                "asset_name": sensor.get("hostname", "Unknown Host"),
                "asset_type_id": asset_type_id,
                "asset_ip": sensor.get("ext_ip", ""),
                "asset_description": f"LimaCharlie EDR sensor ({platform_name})",
                "asset_tags": "limacharlie,edr",
                "analysis_status_id": 1,  # Unspecified - REQUIRED for assets to appear in list
                "custom_attributes": {
                    "LimaCharlie": {
                        "SID": sid,
                        "Actions": actions_html,
                    }
                },
            }

            asset_resp = iris_request(
                f"{IRIS_BASE}/case/assets/add?cid={case_id}", irisKey, ctx,
                method="POST", body=asset_payload
            )
            asset_id = asset_resp.get("data", {}).get("asset_id")
            parsed_response["asset"] = asset_resp
        except Exception as e:
            parsed_response["asset_error"] = str(e)

    # --- Step 4: Add detection base event as a timeline event ---
    try:
        routing = details.get("routing", {})
        event_data = details.get("event", details)
        event_type = routing.get("event_type", "detection")
        detect_name = iris_data.get("case_name", "LimaCharlie Detection")

        # Extract timestamp from the LC event (epoch seconds -> ISO format)
        import datetime
        raw_ts = routing.get("event_time", None)
        if raw_ts:
            # LC timestamps are in epoch microseconds
            try:
                ts = int(raw_ts)
                if ts > 1e15:  # microseconds
                    ts = ts / 1e6
                elif ts > 1e12:  # milliseconds
                    ts = ts / 1e3
                event_dt = datetime.datetime.utcfromtimestamp(ts)
            except (ValueError, OSError):
                event_dt = datetime.datetime.utcnow()
        else:
            event_dt = datetime.datetime.utcnow()
        event_date_str = event_dt.strftime("%Y-%m-%dT%H:%M:%S.000")

        # Build routing HTML for the custom attribute
        routing_html = json_to_html(routing)

        # Build permalink HTML for the custom attribute
        permalink_html = (
            f'<a href="{permalink}" target="_blank" '
            f'style="color:#90cdf4;font-weight:600;">'
            f'{permalink}</a>'
        )

        # Link the event to the asset if one was created
        event_assets = []
        if sid and sid != "<no value>" and asset_id:
            event_assets = [int(asset_id)]

        timeline_payload = {
            "event_title": f"[{event_type}] {detect_name}",
            "event_date": event_date_str,
            "event_tz": "+00:00",
            "event_category_id": 1,  # Unspecified (could map MITRE tactics later)
            "event_source": "LimaCharlie",
            "event_content": f"Detection: {detect_name}\nEvent type: {event_type}\nSensor: {sensor.get('hostname', 'N/A')} ({sid})",
            "event_raw": json.dumps(details, indent=2),
            "event_assets": event_assets,
            "event_iocs": [],
            "event_tags": "limacharlie,detection",
            "custom_attributes": {
                "LimaCharlie": {
                    "Routing": routing_html,
                    "Permalink": permalink_html,
                }
            },
        }

        timeline_resp = iris_request(
            f"{IRIS_BASE}/case/timeline/events/add?cid={case_id}", irisKey, ctx,
            method="POST", body=timeline_payload
        )
        parsed_response["timeline_event"] = timeline_resp
    except Exception as e:
        parsed_response["timeline_event_error"] = str(e)

    return {"data": parsed_response}
