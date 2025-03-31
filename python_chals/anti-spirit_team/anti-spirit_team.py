import json
import random
import string
import re
from starlette.applications import Starlette
from starlette.responses import HTMLResponse, RedirectResponse, PlainTextResponse
from starlette.requests import Request
from starlette.middleware.base import BaseHTTPMiddleware
from jinja2 import Template
from shared_helpers import UserIDMiddleware, OriginalURLMiddleware
from time import time

CHARSET = string.ascii_letters + string.digits + "_-"
CHARS_CLASS = [
    ".",
    "\\d",
    "\\D",
    "\\w",
    "\\W",
]
LENGTH = 4096

app = Starlette()
app.add_middleware(UserIDMiddleware)
app.add_middleware(OriginalURLMiddleware)

def random_alphanum(length: int, rng: random.Random) -> str:
    return ''.join(rng.choice(CHARSET) for _ in range(length))

def random_repeat(rng: random.Random) -> str:
    min_val = rng.randint(1, 24)
    max_val = rng.randint(min_val, 24)
    return "{" + str(min_val) + "," + str(max_val) + "}"

def random_literal(rng: random.Random) -> str:
    char = rng.choice(string.ascii_letters + string.digits + "_-")
    if rng.randint(0, 1) == 0:
        char = re.escape(char)
    return char


def random_character_class(rng: random.Random) -> str:
    repeat = rng.randint(-8, 12)
    if repeat <= 0:
        repeat = 1
    return rng.choice(CHARS_CLASS) + random_repeat(rng)

def random_alternate(rng: random.Random) -> str:
    is_class = rng.randint(0, 1) == 0
    if is_class:
        return "(" + random_character_class(rng) + "|" + random_character_class(rng) + ")"
    else:
        return "(" + random_literal(rng) + "|" + random_literal(rng) + ")"

tokens = {
    "literal": random_literal,
    "character_class": random_character_class,
    "random_alternate": random_alternate,
}

def generate_regex(user_id: int) -> str:
    rng = random.Random(user_id)
    regex = "^"
    while len(regex) < LENGTH:
        factory = rng.choice(list(tokens.keys()))
        regex += tokens[factory](rng)
    regex += "$"
    return regex

template = Template('''
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AST Login: Regex Challenge</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/css/bootstrap.min.css" rel="stylesheet">
    <style>
        body {
            background-color: #121212;
            color: #eaeaea;
            font-family: 'Courier New', Courier, monospace;
        }
        .container {
            max-width: 600px;
            margin-top: 80px;
        }
        .banner {
            background-color: #343a40;
            color: #28a745;
            padding: 20px;
            text-align: center;
            border-radius: 10px;
            box-shadow: 0 2px 15px rgba(0, 0, 0, 0.5);
        }
        .banner h1 {
            font-size: 2.5rem;
            letter-spacing: 2px;
        }
        .card {
            background-color: #1a1a1a;
            border: 1px solid #444;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.7);
            border-radius: 10px;
            margin-top: 30px;
        }
        .card-header {
            background-color: #28a745;
            color: white;
            font-weight: bold;
            text-transform: uppercase;
        }
        .card-body {
            color: #eaeaea;
            font-size: 1rem;
        }
        .card-body blockquote {
            font-style: italic;
            color: #c7c7c7;
            border-left: 3px solid #28a745;
            padding-left: 15px;
            margin-left: 0;
        }
        .footer {
            text-align: center;
            margin-top: 40px;
            font-size: 0.9rem;
            color: #777;
        }
        .btn-primary {
            background-color: #28a745;
            border-color: #28a745;
        }
        .btn-success {
            background-color: #007bff;
            border-color: #007bff;
        }
    </style>
</head>
<body>
  <div class="container">
    <div class="banner">
        <h1>Welcome to the AST Portal</h1>
    </div>

    <div class="card mb-4">
      <div class="card-header">
        Mission Brief
      </div>
      <div class="card-body">
        <p>As a newly recruited member of the Anti-Spirit Team, you're tasked with decoding the cryptic flag hidden within the regex pattern.</p>
        <p>Before you can access the inner sanctum, you must correctly match the provided regex. <br/><br/>And remember:</p>
        <blockquote class="blockquote">
            <p class="mb-0">"I fear no man, but that thing: <code>/.*.*=.*;/</code> it scares me."<br/>- a Cloudflare WAF engineer</p>
        </blockquote>
        <p class="mb-0">Only those who can match the pattern will uncover the flag and be welcomed into the ranks of the AST.</p>
      </div>
    </div>
                    

    {% if message %}
      <div class="alert alert-{{ alert_type }}" role="alert">
        {{ message }}
      </div>
    {% endif %}

    <div class="mb-4">
      <a href="{{original_url}}/download" class="btn btn-primary">Download Regex</a>
    </div>

    <form action="{{original_url}}/" method="post">
      <div class="mb-3">
        <label for="user_input" class="form-label">Enter your input that matches the regex pattern</label>
        <textarea type="text" class="form-control" id="user_input" name="user_input" cols="30" rows="10" maxlength="10000" required></textarea>
      </div>
      <button type="submit" class="btn btn-success">Submit</button>
    </form>
    <div class="footer">
        <p>note: this challenge relies on python3's regex engine.</p>
    </div>
  </div>
</body>
</html>


''')

@app.route("/", methods=["GET", "POST"])
async def index(request: Request):
    if request.method == "POST":
        return await submit(request)
    content = template.render(message=None, alert_type=None, original_url=request.state.original_url)
    return HTMLResponse(content)

@app.route("/download", methods=["GET"])
async def download(request: Request):
    regex = generate_regex(request.state.user_id)
    response = PlainTextResponse(regex)
    response.headers["Content-Disposition"] = "attachment; filename=regex.txt"
    return response

async def submit(request: Request):
    form = await request.form()
    user_input = form.get("user_input").replace("\n", " ").replace("\t", " ").strip()
    if len(user_input) > 10000:
        message = "Failure. Your input is too long."
        alert_type = "danger"
        return HTMLResponse(template.render(message=message, alert_type=alert_type, original_url=request.state.original_url))
    regex = generate_regex(request.state.user_id)
    if re.fullmatch(regex, user_input):
        message = "Success! Your input matches the regex. <br/> Here's your flag: " + request.headers.get("X-GalaCTF-Flag", "there is no flag? (this is a bug)")
        alert_type = "success"
    else:
        message = "Failure. Your input does not match the regex."
        alert_type = "danger"
    content = template.render(message=message, alert_type=alert_type, original_url=request.state.original_url)
    return HTMLResponse(content)


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
