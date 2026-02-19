"""
LimaCharlie Actions Microservice for DFIR-IRIS.

Provides one-click sensor isolation/rejoin from DFIR-IRIS asset custom attributes.
Runs alongside DFIR-IRIS on port 4443 with self-signed TLS.
"""
import os
import json
import ssl
import urllib.request
from flask import Flask, request, jsonify, render_template_string

app = Flask(__name__)

LC_API_KEY = os.environ.get("LC_API_KEY", "")
LC_OID = os.environ["LC_OID"]  # Required: LimaCharlie Organization ID
LC_API_BASE = "https://api.limacharlie.io"

# -- LimaCharlie REST API helpers --

def lc_headers():
    return {
        "Authorization": f"Bearer {LC_API_KEY}",
        "Content-Type": "application/json",
    }


def lc_get(path):
    req = urllib.request.Request(f"{LC_API_BASE}{path}", headers=lc_headers())
    with urllib.request.urlopen(req) as resp:
        return json.loads(resp.read().decode("utf-8"))


def lc_post(path, body=None):
    data = json.dumps(body).encode("utf-8") if body else None
    req = urllib.request.Request(f"{LC_API_BASE}{path}", data=data, headers=lc_headers(), method="POST")
    with urllib.request.urlopen(req) as resp:
        return json.loads(resp.read().decode("utf-8"))


def lc_sensor_info(sid):
    return lc_get(f"/v1/{LC_OID}/sensor/{sid}")


def lc_is_isolated(sid):
    resp = lc_get(f"/v1/{LC_OID}/isolation?sid={sid}")
    return resp.get("is_isolated", False)


def lc_isolate(sid):
    return lc_post(f"/v1/{LC_OID}/isolation", {"sid": sid})


def lc_rejoin(sid):
    req = urllib.request.Request(
        f"{LC_API_BASE}/v1/{LC_OID}/isolation",
        data=json.dumps({"sid": sid}).encode("utf-8"),
        headers=lc_headers(),
        method="DELETE",
    )
    with urllib.request.urlopen(req) as resp:
        return json.loads(resp.read().decode("utf-8"))


# -- Templates --

PAGE_TEMPLATE = """
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LimaCharlie - {{ hostname }}</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #0f172a; color: #e2e8f0; min-height: 100vh; display: flex; align-items: center; justify-content: center; }
        .card { background: #1e293b; border-radius: 12px; padding: 32px; max-width: 480px; width: 100%; box-shadow: 0 25px 50px rgba(0,0,0,.4); }
        .header { display: flex; align-items: center; gap: 12px; margin-bottom: 24px; }
        .status-dot { width: 12px; height: 12px; border-radius: 50%; }
        .status-dot.online { background: #22c55e; box-shadow: 0 0 8px #22c55e; }
        .status-dot.isolated { background: #ef4444; box-shadow: 0 0 8px #ef4444; }
        .status-dot.offline { background: #6b7280; }
        h1 { font-size: 18px; font-weight: 600; word-break: break-all; }
        .meta { color: #94a3b8; font-size: 13px; margin-bottom: 20px; }
        .meta div { margin-bottom: 4px; }
        .meta span { color: #cbd5e1; }
        .badge { display: inline-block; padding: 4px 10px; border-radius: 6px; font-size: 12px; font-weight: 600; }
        .badge-isolated { background: #7f1d1d; color: #fca5a5; }
        .badge-connected { background: #14532d; color: #86efac; }
        .actions { display: flex; gap: 12px; margin-top: 24px; }
        .btn { flex: 1; padding: 12px; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; transition: all .2s; }
        .btn:disabled { opacity: .5; cursor: not-allowed; }
        .btn-danger { background: #dc2626; color: white; }
        .btn-danger:hover:not(:disabled) { background: #b91c1c; }
        .btn-success { background: #16a34a; color: white; }
        .btn-success:hover:not(:disabled) { background: #15803d; }
        .btn-secondary { background: #334155; color: #94a3b8; }
        .btn-secondary:hover { background: #475569; }
        .msg { margin-top: 16px; padding: 12px; border-radius: 8px; font-size: 13px; display: none; }
        .msg-ok { background: #14532d; color: #86efac; }
        .msg-err { background: #7f1d1d; color: #fca5a5; }
        .spinner { display: inline-block; width: 16px; height: 16px; border: 2px solid rgba(255,255,255,.3); border-top-color: white; border-radius: 50%; animation: spin .6s linear infinite; }
        @keyframes spin { to { transform: rotate(360deg); } }
    </style>
</head>
<body>
    <div class="card">
        <div class="header">
            <div class="status-dot {{ 'isolated' if is_isolated else 'online' }}"></div>
            <h1>{{ hostname }}</h1>
        </div>
        <div class="meta">
            <div>SID: <span>{{ sid }}</span></div>
            <div>Platform: <span>{{ platform }}</span></div>
            <div>External IP: <span>{{ ext_ip }}</span></div>
            <div>Internal IP: <span>{{ int_ip }}</span></div>
            <div>Status:
                {% if is_isolated %}
                    <span class="badge badge-isolated">ISOLATED</span>
                {% else %}
                    <span class="badge badge-connected">CONNECTED</span>
                {% endif %}
            </div>
        </div>
        <div class="actions">
            {% if is_isolated %}
                <button class="btn btn-success" id="btn-action" onclick="doAction('rejoin')">Rejoin Network</button>
            {% else %}
                <button class="btn btn-danger" id="btn-action" onclick="doAction('isolate')">Isolate Host</button>
            {% endif %}
            <button class="btn btn-secondary" onclick="window.close()">Close</button>
        </div>
        <div class="msg" id="msg"></div>
    </div>
    <script>
        function doAction(action) {
            var btn = document.getElementById('btn-action');
            var msg = document.getElementById('msg');
            var label = action === 'isolate' ? 'Isolating...' : 'Rejoining...';
            if (!confirm(action === 'isolate'
                ? 'Isolate this host from the network? It will only be able to communicate with LimaCharlie.'
                : 'Rejoin this host to the network? Normal connectivity will be restored.'
            )) return;
            btn.disabled = true;
            btn.innerHTML = '<span class="spinner"></span> ' + label;
            msg.style.display = 'none';
            fetch('/api/' + action, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({sid: '{{ sid }}', oid: '{{ oid }}'})
            })
            .then(function(r) { return r.json(); })
            .then(function(d) {
                if (d.error) {
                    msg.className = 'msg msg-err';
                    msg.textContent = 'Error: ' + d.error;
                    msg.style.display = 'block';
                    btn.disabled = false;
                    btn.textContent = action === 'isolate' ? 'Isolate Host' : 'Rejoin Network';
                } else {
                    msg.className = 'msg msg-ok';
                    msg.textContent = d.message;
                    msg.style.display = 'block';
                    setTimeout(function() { location.reload(); }, 1500);
                }
            })
            .catch(function(e) {
                msg.className = 'msg msg-err';
                msg.textContent = 'Request failed: ' + e;
                msg.style.display = 'block';
                btn.disabled = false;
                btn.textContent = action === 'isolate' ? 'Isolate Host' : 'Rejoin Network';
            });
        }
    </script>
</body>
</html>
"""

ERROR_TEMPLATE = """
<!DOCTYPE html>
<html><head><title>Error</title>
<style>body{font-family:sans-serif;background:#0f172a;color:#e2e8f0;display:flex;align-items:center;justify-content:center;min-height:100vh;}
.card{background:#1e293b;border-radius:12px;padding:32px;max-width:400px;text-align:center;}
h1{color:#ef4444;margin-bottom:12px;}</style></head>
<body><div class="card"><h1>Error</h1><p>{{ error }}</p></div></body></html>
"""

LC_PLATFORM_MAP = {"268435456": "Windows", "536870912": "Linux", "1073741824": "macOS"}

# -- Routes --

@app.route("/sensor/<sid>")
def sensor_page(sid):
    try:
        info = lc_sensor_info(sid)
        isolated = lc_is_isolated(sid)
        hostname = info.get("hostname", sid)
        plat_raw = str(info.get("plat", ""))
        platform = LC_PLATFORM_MAP.get(plat_raw, plat_raw)
        ext_ip = info.get("ext_ip", "N/A")
        int_ip = info.get("int_ip", "N/A")
        return render_template_string(PAGE_TEMPLATE,
            sid=sid, oid=LC_OID, hostname=hostname, platform=platform,
            ext_ip=ext_ip, int_ip=int_ip, is_isolated=isolated)
    except Exception as e:
        return render_template_string(ERROR_TEMPLATE, error=str(e)), 500


@app.route("/api/isolate", methods=["POST"])
def api_isolate():
    data = request.get_json() or {}
    sid = data.get("sid")
    if not sid:
        return jsonify({"error": "Missing sid"}), 400
    try:
        lc_isolate(sid)
        return jsonify({"message": f"Host {sid} has been isolated from the network.", "isolated": True})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@app.route("/api/rejoin", methods=["POST"])
def api_rejoin():
    data = request.get_json() or {}
    sid = data.get("sid")
    if not sid:
        return jsonify({"error": "Missing sid"}), 400
    try:
        lc_rejoin(sid)
        return jsonify({"message": f"Host {sid} has rejoined the network.", "isolated": False})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@app.route("/api/status/<sid>")
def api_status(sid):
    try:
        isolated = lc_is_isolated(sid)
        return jsonify({"sid": sid, "is_isolated": isolated})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


if __name__ == "__main__":
    import sys
    cert = sys.argv[1] if len(sys.argv) > 1 else None
    key = sys.argv[2] if len(sys.argv) > 2 else None
    ssl_ctx = None
    if cert and key:
        ssl_ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
        ssl_ctx.load_cert_chain(cert, key)
    app.run(host="0.0.0.0", port=4443, ssl_context=ssl_ctx)
