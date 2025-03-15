import random
import hashlib
from threading import Lock
from starlette.applications import Starlette
from starlette.responses import JSONResponse, HTMLResponse
from starlette.requests import Request
from jinja2 import Template
from shared_helpers import UserIDMiddleware, OriginalURLMiddleware

FLAG_CHALS_LOCK = Lock()
FLAG_CHALS = {}
LOW = 0
HIGH = 1

app = Starlette()

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
            async function verify(y, C = 0) {
                const res = await fetch("{{original_url}}/verify", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        y: y
                    })
                });
                const data = await res.json();
                if (data.success) {
                    if (C === 0) {
                        $("#title").text("Captcha Solved - Get flag!");
                        $("#description").text("To get the flag, click on the 'Get flag' button below.");
                        $("#verify-btn").text("Get flag");
                        $("#verify-btn").attr("disabled", false);
                        $("#alert").removeClass("visually-hidden");
                        $("#verify-btn").click(async function() {
                            solveCaptcha("0000000")
                        });

                        $("#manual-btn").removeClass("visually-hidden");
                        $("#manual-btn").click(function() {
                            fetch("{{original_url}}/flag_x")
                                .then((res) => res.json())
                                .then((data) => {
                                    $("#x-ctn").text(`Solve y for x = ${data.x}, g = '0000000'.`);
                                    $("#x-ctn").parent().removeClass("visually-hidden");
                                    $("#manual-btn").text("Submit manual solution");
                                    $("#manual-btn").unbind("click");
                                    $("#manual-btn").click(function() {
                                        const y = prompt("y=?");
                                        verify(y, 1);
                                    });
                                });
                        });
                    } else if (C === 1 && data.flag) {
                        alert(`Congratulations! Here's your flag: ${data.flag}`);
                    }
                } else {
                    alert(data.message);
                }
            }
            const solveCaptcha = async (g = "0000") => {
                const x = await fetch("{{original_url}}/x")
                    .then((res) => res.json())
                    .then((data) => data.x);
                
                const maxY = 1e9;
                console.log(`[pow] solving y for x=${x}`);
                
                for (let y = 0; y < maxY; y++) {
                    const hash = forge.md.sha256.create().update(`${x}${y}`).digest().toHex();
                    if (hash.startsWith(g)) {
                        console.log(`[pow] Found y: ${y}`);
                        verify(y, g === "0000000" ? 1 : 0);
                        return;
                    }
                }
            };

            $("#verify-btn").click(function() {
                solveCaptcha();
                $("#verify-btn").attr("disabled", true);
                $("#verify-btn").text("Verifying...");
            });
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
async def verify(request: Request):
    data = await request.json()
    with FLAG_CHALS_LOCK:
        if request.state.user_id not in FLAG_CHALS:
            return JSONResponse({"error": "You have not generated a challenge yet!"}, status_code=400)
        chal = FLAG_CHALS[request.state.user_id]
        sol = data["y"]
        h = hashlib.sha256()
        h.update((str(chal["x"]) + str(sol)).encode())
        h = h.hexdigest()
        if chal["C"] == LOW and h.startswith("0000"):
            del FLAG_CHALS[request.state.user_id]
            return JSONResponse({"success": True})
        elif chal["C"] == HIGH and h.startswith("0000000"):
            del FLAG_CHALS[request.state.user_id]
            return JSONResponse({"success": True, "flag": request.headers.get("X-GalaCTF-Flag")})
    return JSONResponse({"success": False, "message": "An unexpected error occurred. Please try again later. (is the solution invalid?)"}, status_code=500)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
