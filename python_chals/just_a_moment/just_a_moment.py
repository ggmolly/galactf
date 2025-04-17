import random
import hashlib
import time
from threading import Lock
from starlette.applications import Starlette
from starlette.responses import JSONResponse, HTMLResponse
from starlette.requests import Request
from jinja2 import Template
from shared_helpers import UserIDMiddleware, OriginalURLMiddleware
from slowapi import Limiter, _rate_limit_exceeded_handler
from slowapi.errors import RateLimitExceeded
from slowapi.util import get_remote_address

FLAG_CHALS_LOCK = Lock()
FLAG_CHALS = {}
LOW = 0
HIGH = 1

limiter = Limiter(key_func=lambda request: request.state.user_id)
app = Starlette()
app.state.limiter = limiter
app.add_exception_handler(RateLimitExceeded, _rate_limit_exceeded_handler)

app.add_middleware(UserIDMiddleware)
app.add_middleware(OriginalURLMiddleware)

@app.route("/", methods=["GET"])
async def index(request: Request):
    
    html_content = """
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Captcha Challenge</title>
        <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
    </head>
    <body>
        <div class="container">
            <div class="row justify-content-center mt-5">
                <div class="col-md-6">
                    <div class="card">
                        <div class="card-header">
                            <h4 id="title">Captcha Challenge</h4>
                        </div>
                        <div class="card-body">
                            <p class="text-center" id="description">This captcha uses an algorithm to verify that you are not a spammer, employing very complex cryptographic computations.</p>
                            <p class="visually-hidden text-danger text-center" id="alert">
                                <strong>Warning!</strong> This may crash your browser.
                            </p>
                            <div class="mb-3 text-center">
                                <button type="submit" class="btn btn-success" id="verify-btn">Verify me!</button>
                                <button type="button" class="btn btn-primary visually-hidden" id="manual-btn">Solve manually</button>
                            </div>
                            <div class="visually-hidden">
                                <p id="x-ctn"></p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.7.1/jquery.slim.min.js" integrity="sha512-sNylduh9fqpYUK5OYXWcBleGzbZInWj8yCJAU57r1dpSK9tP2ghf/SRYCMj+KsslFkCOt3TvJrX2AV/Gc3wOqA==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/forge/1.3.1/forge.all.min.js" integrity="sha512-y8oWKnULY59b/Ce+mlekagFu+2M1R4FCPoQvG1Gvgp5mpM3UiTAQZY/3ai91gE0IW8/yk76JojiUkRlUP59F0Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
        <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" crossorigin="anonymous"></script>
        <script defer>
            async function verify(t,e=0){const n=await fetch("{{original_url}}/verify",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({y:t})});if(429===n.status)return void alert("You're doing that too fast! Retry in a minute.");const a=await n.json();a.success?0===e?($("#title").text("Captcha Solved - Get flag!"),$("#description").text("To get the flag, click on the 'Get flag' button below."),$("#verify-btn").text("Get flag"),$("#verify-btn").attr("disabled",!1),$("#alert").removeClass("visually-hidden"),$("#verify-btn").click((async function(){solveCaptcha("0000000")})),$("#manual-btn").removeClass("visually-hidden"),$("#manual-btn").click((function(){fetch("{{original_url}}/flag_x").then((t=>t.json())).then((t=>{$("#x-ctn").text(`Solve y for x = ${t.x}, g = '0000000'.`),$("#x-ctn").parent().removeClass("visually-hidden"),$("#manual-btn").text("Submit manual solution"),$("#manual-btn").unbind("click"),$("#manual-btn").click((function(){verify(prompt("y=?"),1)}))}))}))):1===e&&a.flag&&alert(`Congratulations! Here's your flag: ${a.flag}`):alert(a.message)}const solveCaptcha=async(t="0000")=>{const e=await fetch("{{original_url}}/x").then((t=>t.json())).then((t=>t.x));console.log(`[pow] solving y for x=${e}`);for(let n=0;n<1e9;n++){if(forge.md.sha256.create().update(`${e}${n}`).digest().toHex().startsWith(t))return console.log(`[pow] Found y: ${n}`),void verify(n,"0000000"===t?1:0)}};$("#verify-btn").click((function(){solveCaptcha(),$("#verify-btn").attr("disabled",!0),$("#verify-btn").text("Verifying...")}));
        </script>
    </body>
    </html>
    """
    template = Template(html_content)
    return HTMLResponse(template.render(original_url=request.state.original_url))

@app.route("/x", methods=["GET"])
async def get_x(request: Request):
    with FLAG_CHALS_LOCK:
        x = random.randint(100000000000000, 9007199254740991)
        FLAG_CHALS[request.state.user_id] = {
            "x": x,
            "C": LOW, # complexity, 0 is for 4 leading zeros, 1 is for 8 leading zeros
        }
        return JSONResponse(FLAG_CHALS[request.state.user_id])

@app.route("/flag_x", methods=["GET"])
async def get_flag_x(request: Request):
    with FLAG_CHALS_LOCK:
        x = random.randint(100000000000000, 9007199254740991)
        FLAG_CHALS[request.state.user_id] = {
            "x": x,
            "C": HIGH,
        }
        return JSONResponse(FLAG_CHALS[request.state.user_id])

@app.route("/verify", methods=["POST"])
@limiter.limit(limit_value="5/minute")
async def verify(request: Request):
    time.sleep(random.randint(1, 5))
    data = await request.json()
    with FLAG_CHALS_LOCK:
        if request.state.user_id not in FLAG_CHALS:
            return JSONResponse({"error": "You have not generated a challenge yet!"}, status_code=400)
        chal = FLAG_CHALS[request.state.user_id]
        sol = str(data["y"])
        if len(sol) > 16: # max length for 9007199254740991, early exit to avoid DoS
            return JSONResponse({"success": False, "message": "Invalid solution. Please try again later."}, status_code=400)
        h = hashlib.sha256()
        h.update((str(chal["x"]) + sol).encode())
        h = h.hexdigest()
        if chal["C"] == LOW and h.startswith("0000"):
            del FLAG_CHALS[request.state.user_id]
            return JSONResponse({"success": True})
        elif chal["C"] == HIGH and h.startswith("0000000"):
            del FLAG_CHALS[request.state.user_id]
            return JSONResponse({"success": True, "flag": request.headers.get("X-GalaCTF-Flag")})
    return JSONResponse({"success": False, "message": "Invalid solution."}, status_code=409)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
