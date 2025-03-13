import sqlite3
import asyncio
import random
from starlette.applications import Starlette
from starlette.responses import JSONResponse, HTMLResponse, Response
from starlette.routing import Route
from starlette.requests import Request
from jinja2 import Template
from shared_helpers import UserIDMiddleware, OriginalURLMiddleware

STICKY_NOTES = [
    (1, "Hello, world!", "This is a sample sticky note.", "Admin"),
    (2, "This is a test", "Here's a second sample sticky note.", "Admin"),
    (3, "Wow!", "Great job sweetie, I'm really proud of you! Did you do your homework today? :) -mom", "Mom"),
    (5, "Congratulations", "Here's your flag: $FLAG_TEMPLATE", "m"),
    (8, "SQL Injection", "Always sanitize user input to avoid SQL injection attacks. Use parameterized queries.", "Alice"),
    (9, "Buffer Overflow", "Careful with buffer sizes! Overwriting memory can lead to serious vulnerabilities.", "Bob"),
    (10, "Cross-Site Scripting", "Never trust user input! Always escape HTML special characters before rendering.", "Charlie"),
    (11, "Password Hashing", "Use bcrypt or Argon2 for password hashing, never store plaintext passwords.", "Dave"),
    (16, "Two-Factor Authentication", "Implement 2FA where possible to add an extra layer of security.", "Ivy"),
    (19, "API Rate Limiting", "Implement rate limiting on your APIs to protect against DoS attacks.", "Mia"),
    (20, "JWT", "Never store sensitive information in JWT payloads; always use proper encryption.", "Noah"),
    (24, "Session Management", "Ensure proper session expiration and token invalidation to protect user sessions.", "Riley"),
    (25, "Input Validation", "Never trust user input. Validate and sanitize inputs before processing them.", "Sam"),
    (26, "WebSockets", "Use WebSockets for real-time communication, but make sure to authenticate and encrypt.", "Grégory"),
    (27, "HTTP Headers", "Always set security-related HTTP headers like Content-Security-Policy and Strict-Transport-Security.", "Uma")
]

db = sqlite3.connect(":memory:", check_same_thread=False)
cursor = db.cursor()

cursor.execute("""
    CREATE TABLE sticky_notes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        author TEXT NOT NULL
    )
""")

cursor.executemany("""
    INSERT INTO sticky_notes (id, title, content, author)
    VALUES (?, ?, ?, ?)
""", STICKY_NOTES)
db.commit()

async def serve_css(request):
    return Response(
        content="""
        .post-it {
            background-color: #f7f3a5;
            border: 3px solid #e4d836;
            padding: 20px;
            max-width: 300px;
            box-shadow: 5px 5px 15px rgba(0, 0, 0, 0.2);
            font-family: sans-serif;
            margin-top: 50px;
            position: relative;
            overflow-wrap: break-word;
            word-wrap: break-word;
            hyphens: auto;
        }
        .container {
            display: flex;
            justify-content: center;
            align-items: center;
            gap: 16px;
            flex-wrap: wrap;
        }
        .center {
            text-align: center;
        }
        .post-it::after {
            content: '';
            position: absolute;
            bottom: -10px;
            right: 10px;
            width: 40px;
            height: 40px;
            background-color: #e4d836;
            transform: rotate(45deg);
            z-index: 1;
        }
        """, 
        media_type="text/css", 
        status_code=200
    )

async def index(request):
    cursor.execute("SELECT * FROM sticky_notes ORDER BY id DESC LIMIT 4")
    sticky_notes = cursor.fetchall()

    sticky_notes.reverse()
    if sticky_notes:
        html_template = """
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Sticky Notes</title>
            <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/css/bootstrap.min.css" rel="stylesheet">
            <link rel="stylesheet" href="{{original_url}}/main.css">
        </head>
        <body>
            <h1 class="mt-5 center">Latest Public Sticky Notes</h1>
            <p class="center">Click on the titles for more details!</p>
            <div class="container text-center">
                {% for sticky_note in sticky_notes %}
                <div class="post-it">
                    <a href="{{original_url}}/note/{{ sticky_note[0] }}">
                        <h5>{{ sticky_note[1] }}</h5>
                    </a>
                    <p>{{ sticky_note[2] }}</p>
                    <small>Author: {{ sticky_note[3] }}</small>
                </div>
                {% endfor %}
            </div>
            <hr class="my-4" />
            <div class="container text-center mt-4">
                <form action="login" method="post">
                    <div class="mb-3">
                        <label for="username" class="form-label">Username</label>
                        <input type="text" id="username" name="username" class="form-control" required>
                    </div>
                    <div class="mb-3">
                        <label for="password" class="form-label">Password</label>
                        <input type="password" id="password" name="password" class="form-control" required>
                    </div>
                    <p class="text-danger">
                        You must be logged in to create a sticky note!
                        <br />
                        They can be public, hidden or private.
                    </p>
                    <button type="submit" class="btn btn-primary w-100 mt-2">Login</button>
                </form>
            </div>
        </body>
        </html>
        """
        template = Template(html_template)
        html_content = template.render(
            sticky_notes=sticky_notes,
            original_url=request.state.original_url
        )

        return HTMLResponse(content=html_content)
    else:
        return HTMLResponse(content="<h1>No sticky notes found</h1>")

async def note(request):
    note_id = request.path_params['id']
    cursor.execute("SELECT * FROM sticky_notes WHERE id = ?", (note_id,))
    sticky_note = cursor.fetchone()

    if sticky_note:
        title, content, author = sticky_note[1], sticky_note[2], sticky_note[3]
        content = content.replace("$FLAG_TEMPLATE", request.headers.get("X-GalaCTF-Flag", "Error. Contact admins."))
        html_template = """
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Sticky Note</title>
            <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/css/bootstrap.min.css" rel="stylesheet">
            <link rel="stylesheet" href="{{original_url}}/main.css">
        </head>
        <body>
            <div class="container d-flex justify-content-center align-items-center" style="height: 100vh;">
                <div class="post-it">
                    <h5>{{ title }}</h5>
                    <p>{{ content }}</p>
                    <small>Author: {{ author }}</small>
                </div>
            </div>
            <a href="{{original_url}}" class="btn btn-primary position-fixed bottom-0 end-0 m-3">Back to Home</a>
        </body>
        </html>
        """
        template = Template(html_template)
        html_content = template.render(
            title=title,
            content=content,
            author=author,
            original_url=request.state.original_url
        )
        return HTMLResponse(content=html_content)
    else:
        return HTMLResponse(content="<h1>Sticky note not found, it might be private.</h1>")


async def dummy_login(request):
    await asyncio.sleep(random.random() * 2)
    return Response(content="Wrong credentials. Try again later.", status_code=401)

app = Starlette(routes=[
    Route("/", index),
    Route("/login", dummy_login, methods=["POST"]),
    Route("/note/{id:int}", note),
    Route("/main.css", serve_css)
])
app.add_middleware(UserIDMiddleware)
app.add_middleware(OriginalURLMiddleware)

if __name__ == '__main__':
    import uvicorn
    uvicorn.run(app, host='0.0.0.0', port=8080)
