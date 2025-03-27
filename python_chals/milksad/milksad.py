import hashlib
import secrets
import base64
import json
import hmac
import random
from starlette.applications import Starlette
from starlette.responses import HTMLResponse, RedirectResponse
from starlette.requests import Request
from starlette.middleware.base import BaseHTTPMiddleware
from jinja2 import Template
from shared_helpers import UserIDMiddleware, OriginalURLMiddleware
from Crypto.Cipher import AES 
from time import time

LOGIN_TEMPLATE = Template("""<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Milksad Inc. - Admin Panel</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
</head>
<body class="d-flex justify-content-center align-items-center" style="height: 100vh; background-color: #f4f6f9;">
    <div class="card shadow-lg" style="width: 100%; max-width: 400px;">
        <div class="card-body">
            <h4 class="card-title text-center mb-4">Admin Login 🥛</h4>

            {% if error %}
            <div class="alert alert-danger" role="alert">
                {{ error }}
            </div>
            {% endif %}

            <form action="{{original_url}}/login" method="POST">
                <div class="mb-3">
                    <label for="username" class="form-label">Username</label>
                    <input type="text" class="form-control" id="username" placeholder="Enter username" required name="username">
                </div>
                <div class="mb-3">
                    <label for="password" class="form-label">Password</label>
                    <input type="password" class="form-control" id="password" placeholder="Enter password" required name="password">
                </div>
                <div class="d-flex justify-content-between mb-3">
                    <div class="form-check">
                        <input type="checkbox" class="form-check-input" id="rememberMe">
                        <label class="form-check-label" for="rememberMe">Remember me</label>
                    </div>
                    <a href="#" class="text-decoration-none small">Forgot password?</a>
                </div>
                <button type="submit" class="btn btn-primary w-100">Login</button>
                <p class="text-muted text-center mt-1">Milksad Inc. © 2025</p>
            </form>
        </div>
    </div>

    <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.7.1/jquery.slim.min.js" integrity="sha512-sNylduh9fqpYUK5OYXWcBleGzbZInWj8yCJAU57r1dpSK9tP2ghf/SRYCMj+KsslFkCOt3TvJrX2AV/Gc3wOqA==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" crossorigin="anonymous"></script>
</body>
</html>
""")
MFA_TEMPLATE   = Template("""<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Milksad Inc. - Admin Panel</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
</head>
<body class="d-flex justify-content-center align-items-center" style="height: 100vh; background-color: #f4f6f9;">
    <div class="card shadow-lg" style="width: 100%; max-width: 400px;">
        <div class="card-body">
            <h4 class="card-title text-center mb-4">Admin Login 🥛</h4>

            {% if error %}
            <div class="alert alert-danger" role="alert">
                {{ error }}
            </div>
            {% endif %}

            <form action="{{original_url}}/mfa" method="POST">
                <div class="mb-3">
                    <label for="mfa" class="form-label">MFA Code (rotates ever 30 seconds)</label>
                    <input type="number" class="form-control" id="mfa" name="mfa" placeholder="Enter 6-digit MFA code" required min="100000" max="999999" step="1" name="mfa">
                </div>
                <div class="d-flex justify-content-between mb-3">
                    <a href="#" class="text-decoration-none small">Forgot password?</a>
                </div>
                <button type="submit" class="btn btn-primary w-100">Login</button>
                <p class="text-muted text-center mt-1">Milksad Inc. © 2025 - Server time: {{ server_time }}</p>
            </form>
        </div>
    </div>

    <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.7.1/jquery.slim.min.js" integrity="sha512-sNylduh9fqpYUK5OYXWcBleGzbZInWj8yCJAU57r1dpSK9tP2ghf/SRYCMj+KsslFkCOt3TvJrX2AV/Gc3wOqA==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" crossorigin="anonymous"></script>
</body>
</html>
""")

# ====== Cryptographic stuff, this is used to authenticate user and prevent them from bypassing the real challenge =====
HMAC_KEY = "W5XO4VeBe17u59z8vbTKkjKr5xzt8VmV"
CIPHER_KEY = "86LedUgJ5FnzgRSy6TOhCRvvIdTHbOAG"
COOKIE_NAME = "milksad-admin-sess"

# HMAC to make sure users cannot modify the cookie (if they succeed to decipher it)
def generate_hmac(data: dict) -> str:
    return hmac.new(HMAC_KEY.encode(), json.dumps(data).encode(), hashlib.sha256).hexdigest()

def verify_hmac(data: dict) -> bool:
    h = data['sig']
    del data['sig']
    return hmac.new(HMAC_KEY.encode(), json.dumps(data).encode(), hashlib.sha256).hexdigest() == h

def get_iv() -> bytes:
    return secrets.token_bytes(16)

def generate_cookie(data: dict, iv: bytes) -> str:
    cipher = AES.new(CIPHER_KEY.encode(), AES.MODE_GCM, iv)
    data['sig'] = generate_hmac(data)
    return base64.b64encode(iv + cipher.encrypt(json.dumps(data).encode())).decode()

def decrypt_cookie(cookie: str) -> dict:
    cookie = base64.b64decode(cookie)
    iv = cookie[:16]
    cipher = AES.new(CIPHER_KEY.encode(), AES.MODE_GCM, iv)
    data = json.loads(cipher.decrypt(cookie[16:]).decode())
    if not verify_hmac(data):
        raise ValueError("Invalid signature")
    return data
# ====== End of crypto stuff =====


app = Starlette()

app.add_middleware(UserIDMiddleware)
app.add_middleware(OriginalURLMiddleware)

class CheckCookieMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        if request.method != "GET":
            return await call_next(request)
        cookie = request.cookies.get(COOKIE_NAME)
        if cookie is None and request.url.path != "/":
            return RedirectResponse(url=request.headers.get("X-Root-Uri", "") + "/")
        elif cookie is None and request.url.path == "/":
            return await call_next(request)

        try:
            data = decrypt_cookie(cookie)
            if data.get("mfa_valid"):
                return HTMLResponse(
                    "Welcome! But the admin panel is in another server. Here's the flag instead: " + request.headers.get("X-GalaCTF-Flag", "wait.. there is no flag? (this is a bug)")
                )
            elif not request.url.path.startswith("/mfa"):
                return RedirectResponse(url=request.headers.get("X-Root-Uri", "") + "/mfa")
        except Exception:
            return RedirectResponse(url=request.headers.get("X-Root-Uri", "") + "/")

        response = await call_next(request)
        return response

app.add_middleware(CheckCookieMiddleware)

@app.route("/", methods=["GET"])
async def index(request: Request):
    return HTMLResponse(LOGIN_TEMPLATE.render(original_url=request.state.original_url))

@app.route("/login", methods=["POST", "GET"])
async def do_login(request: Request):
    if request.method == "GET":
        return RedirectResponse(url=request.headers.get("X-Root-Uri", "") + "/", status_code=302)
    data = await request.form()
    if data.get("username") != "admin" or data.get("password") != "admin":
        return HTMLResponse(LOGIN_TEMPLATE.render(original_url=request.state.original_url, error="Invalid username or password"))

    # send a cookie to the user with mfa_valid set to false
    data = {
        "username": data.get("username"),
        "mfa_valid": False,
    }
    iv = get_iv()
    cookie = generate_cookie(data, iv)
    response = RedirectResponse(url=request.headers.get("X-Root-Uri", "") + "/mfa", status_code=303)
    response.set_cookie(COOKIE_NAME, cookie, httponly=True, secure=True, samesite="lax", path="/")
    return response

def verify_mfa(request: Request, mfa: str):
    try:
        x = int(mfa)
    except ValueError:
        return HTMLResponse(MFA_TEMPLATE.render(original_url=request.state.original_url, error="🤨", server_time=int(time())))
    
    if x < 100000 or x > 999999:
        return HTMLResponse(MFA_TEMPLATE.render(original_url=request.state.original_url, error="🤨", server_time=int(time())))
    
    seed = int(time()) // 30 # rotate every 30 seconds
    allowed_mfas = []

    for s in [seed-1, seed, seed+1]: # tolerance of -30, 0, 30 seconds
        mt = random.Random(s)
        allowed_mfas.append(mt.randint(100000, 999999))

    if x not in allowed_mfas:
        return HTMLResponse(MFA_TEMPLATE.render(original_url=request.state.original_url, error="Wrong MFA code :(. Is your clock synced?", server_time=int(time())))

    # forge a valid cookie
    data = {
        "username": request.state.user_id,
        "mfa_valid": True,
    }
    iv = get_iv()
    cookie = generate_cookie(data, iv)
    response = RedirectResponse(url=request.headers.get("X-Root-Uri", "") + "/mfa", status_code=303)
    response.set_cookie(COOKIE_NAME, cookie, httponly=True, secure=True, samesite="lax", path="/")
    return response

@app.route("/mfa", methods=["GET", "POST"])
async def mfa(request: Request):
    cookie = request.cookies.get(COOKIE_NAME)
    try:
        data = decrypt_cookie(cookie)
        if data.get("mfa_valid") == True:
            return HTMLResponse(request.headers.get("X-GalaCTF-Flag"))
        elif request.method == "POST":
            data = await request.form()
            return verify_mfa(request, data.get("mfa"))
        else:
            return HTMLResponse(MFA_TEMPLATE.render(original_url=request.state.original_url, server_time=int(time())))
    except:
        pass
    # user has no valid cookie, redirect to login so they can login again
    return RedirectResponse(url=request.headers.get("X-Root-Uri", "") + "/", status_code=302)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
