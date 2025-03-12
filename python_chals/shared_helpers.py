from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request
from starlette.responses import JSONResponse

class UserIDMiddleware(BaseHTTPMiddleware):
    """
    This middleware is used to isolate per user data inside the challenges.
    """
    async def dispatch(self, request: Request, call_next):
        user_id = request.headers.get("X-User-ID")
        if user_id is None:
            return JSONResponse({"error": "Invalid X-User-ID"}, status_code=401)
        try:
            request.state.user_id = int(user_id)
        except ValueError:
            return JSONResponse({"error": "Invalid X-User-ID"}, status_code=401)
        return await call_next(request)
