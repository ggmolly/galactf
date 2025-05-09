from starlette.applications import Starlette
from starlette.responses import HTMLResponse, PlainTextResponse, RedirectResponse
from starlette.requests import Request
from shared_helpers import UserIDMiddleware, OriginalURLMiddleware
import os
import pwd
import grp

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

    try:
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
    except PermissionError:
        return PlainTextResponse("Vous n'avez pas les droits pour accéder à ce dossier", status_code=403)
    except FileNotFoundError:
        return PlainTextResponse("Dossier ou fichier introuvable", status_code=404)
    except Exception:
        return PlainTextResponse("Une erreur est survenue", status_code=500)


if __name__ == "__main__":
    pw = pwd.getpwnam("ctfplayer")
    os.setgid(pw.pw_gid)
    os.setuid(pw.pw_uid)
    os.environ["HOME"] = pw.pw_dir

    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
