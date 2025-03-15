from starlette.applications import Starlette
from starlette.responses import JSONResponse, HTMLResponse, Response
from starlette.routing import Route
from starlette.requests import Request
from jinja2 import Template
from shared_helpers import UserIDMiddleware, OriginalURLMiddleware


app = Starlette()
app.add_middleware(UserIDMiddleware)
app.add_middleware(OriginalURLMiddleware)

@app.route("/", methods=["GET", "POST"])
async def index(request: Request):
    is_admin = request.cookies.get("isAdmin") == "true"
    if request.method == "POST" and is_admin:
        return JSONResponse({"flag": request.headers.get("X-GalaCTF-Flag")})
    template = Template("""
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Cookie Monster</title>
        <style>
            .container {
                display: flex;
                justify-content: center;
                align-items: center;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <form action="{{original_url}}/" method="post">
                {% if is_admin %}
                    <button type="submit">Get flag</button>
                {% else %}
                    <button type="button" disabled>You must be an admin to get the flag</button>
                {% endif %}
            </form>
        </div>
    </body>
    </html>
    """)
    response = HTMLResponse(
        content=template.render(
            original_url=request.state.original_url,
            is_admin=is_admin,
        )
    )
    if "isAdmin" not in request.cookies:
        response.set_cookie("isAdmin", "false")
    return response


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
