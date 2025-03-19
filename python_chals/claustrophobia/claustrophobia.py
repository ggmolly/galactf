from starlette.applications import Starlette
from starlette.responses import HTMLResponse, PlainTextResponse, RedirectResponse
from starlette.requests import Request
from shared_helpers import UserIDMiddleware, OriginalURLMiddleware
import os

app = Starlette()
app.add_middleware(UserIDMiddleware)
app.add_middleware(OriginalURLMiddleware)

BASE_DIR = os.path.join('/home/ctfplayer')

@app.route("/", methods=["GET"])
async def index(request: Request):
    root_dir = request.query_params.get("rootDir")

    if root_dir is None:
        return RedirectResponse(url=request.url.path + "api/v1/factories/claustrophobia?rootDir=.")

    target_path = os.path.join(BASE_DIR, root_dir)

    if os.path.exists(target_path):
        if os.path.isdir(target_path):
            files = os.listdir(target_path)
            file_list = "".join(
                f'<li><a href="?rootDir={os.path.join(root_dir, f)}">{f}</a></li>' for f in files
            )
            return HTMLResponse(f"<h1>Fichiers dans {target_path}</h1><ul>{file_list}</ul>")
        else:
            with open(target_path, "r") as f:
                content = f.read()
            return PlainTextResponse(content.replace("$FLAG_PLACEHOLDER", request.headers.get("X-GalaCTF-Flag")))
    else:
        return PlainTextResponse("Dossier ou fichier introuvable", status_code=404)


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
