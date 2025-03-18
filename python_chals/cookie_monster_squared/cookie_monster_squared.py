from starlette.applications import Starlette
from starlette.responses import JSONResponse, HTMLResponse
from starlette.requests import Request
from jinja2 import Template
import jwt
from shared_helpers import UserIDMiddleware, OriginalURLMiddleware

app = Starlette()
app.add_middleware(UserIDMiddleware)
app.add_middleware(OriginalURLMiddleware)

JWT_WEAK_SECRET = 'undefined'

@app.route("/", methods=["GET", "POST"])
async def index(request: Request):
    privileges = request.cookies.get("privileges")
    error = "You must be an admin to get the flag."
    decoded_privileges = {"isAdmin": False}
    try:
        decoded_privileges = jwt.decode(privileges, JWT_WEAK_SECRET, algorithms=["HS256"])
    except jwt.exceptions.InvalidSignatureError:
        error = "jwt error: Invalid signature"
    except jwt.exceptions.InvalidAlgorithmError:
        error = "jwt error: Invalid algorithm"
    except Exception:
        pass
    is_admin = decoded_privileges.get("isAdmin") == True
    if request.method == "POST" and is_admin:
        return JSONResponse({"flag": request.headers.get("X-GalaCTF-Flag")})
    template = Template("""
    <!DOCTYPE html>
    <html lang="en">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <title>Cookie Monster² Challenge</title>
      <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
      <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" crossorigin="anonymous"></script>
    </head>
    <body>
      <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container-fluid">
          <a class="navbar-brand" href="#">CM Inc. LLC</a>
          <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav"
            aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
          </button>
          <div class="collapse navbar-collapse" id="navbarNav">
            <ul class="navbar-nav ms-auto">
              <li class="nav-item">
                <a class="nav-link" href="#">Home</a>
              </li>
              {% if is_admin %}
              <li class="nav-item">
                <a class="nav-link" href="#">Admin Panel</a>
              </li>
              {% endif %}
            </ul>
          </div>
        </div>
      </nav>
      <div class="container mt-5">
        <div class="row">
          <div class="col-md-8 offset-md-2">
            <div class="card">
              <div class="card-header bg-primary text-white">
                <h3>Welcome to the Cookie Monster² Challenge</h3>
              </div>
              <div class="card-body">
                <p>
                  This CTF challenge is designed to test your wits and skills. Explore the page, follow the clues,
                  and see if you can uncover the hidden secrets. Remember, not everything is as it seems.
                </p>
                {% if is_admin %}
                  <form action="{{ original_url }}/" method="post">
                    <button type="submit" class="btn btn-success">Reveal the Flag</button>
                  </form>
                {% else %}
                  <div class="alert alert-warning" role="alert">
                    {{ error }}
                  </div>
                {% endif %}
              </div>
              <div class="card-footer text-muted">
                &copy; 2025 CTF Challenge. All rights reserved.
              </div>
            </div>
          </div>
        </div>
      </div>
      <footer class="bg-dark text-white text-center py-3 mt-5">
        <div class="container">
          <p>🤤🤤🤤🤤🤤🤤🤤</p>
          <p>🤤 cookies² 🤤</p>
          <p>🤤🤤🤤🤤🤤🤤🤤</p>
        </div>
      </footer>
    </body>
    </html>
    """)
    response = HTMLResponse(
        content=template.render(
            original_url=request.state.original_url,
            is_admin=is_admin,
            error=error,
        )
    )
    if not privileges or not request.cookies.get("privileges").startswith("ey"):
        response.set_cookie("privileges", jwt.encode({"isAdmin": False}, JWT_WEAK_SECRET, algorithm="HS256"))
    return response

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
