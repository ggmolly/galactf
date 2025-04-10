import random
import hashlib
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
            function _0x108f(_0x5e942e,_0x3efc20){const _0x422dc3=_0x247d();return _0x108f=function(_0xdd2728,_0x1bf993){_0xdd2728=_0xdd2728-0x1cf;let _0x1cb790=_0x422dc3[_0xdd2728];return _0x1cb790;},_0x108f(_0x5e942e,_0x3efc20);}const _0x442c0c=_0x108f;(function(_0x1ce551,_0x16b4da){const _0x149c82=_0x108f,_0x104bd4=_0x1ce551();while(!![]){try{const _0x24adb9=-parseInt(_0x149c82(0x1ef))/0x1+-parseInt(_0x149c82(0x1e7))/0x2+-parseInt(_0x149c82(0x1e8))/0x3+-parseInt(_0x149c82(0x1f9))/0x4*(parseInt(_0x149c82(0x1e5))/0x5)+parseInt(_0x149c82(0x204))/0x6+parseInt(_0x149c82(0x1dc))/0x7+parseInt(_0x149c82(0x203))/0x8;if(_0x24adb9===_0x16b4da)break;else _0x104bd4['push'](_0x104bd4['shift']());}catch(_0x775721){_0x104bd4['push'](_0x104bd4['shift']());}}}(_0x247d,0x4895b));const _0x32cb98=(function(){let _0x4be89a=!![];return function(_0x39dc26,_0x590476){const _0x1d95e9=_0x4be89a?function(){const _0x59b4b6=_0x108f;if(_0x590476){const _0xdabdc3=_0x590476[_0x59b4b6(0x1e6)](_0x39dc26,arguments);return _0x590476=null,_0xdabdc3;}}:function(){};return _0x4be89a=![],_0x1d95e9;};}()),_0x3fda2b=_0x32cb98(this,function(){const _0x174545=_0x108f;return _0x3fda2b['toString']()[_0x174545(0x200)](_0x174545(0x1ff))['toString']()[_0x174545(0x1ee)](_0x3fda2b)[_0x174545(0x200)](_0x174545(0x1ff));});function _0x247d(){const _0x2ffa91=['toString','3537135oRuQzH','Get\x20flag','log','json','sha256','[pow]\x20Found\x20y:\x20','table','{{original_url}}/verify','status','15cxPAPC','apply','388562iIvyug','1753737rDRcVz','Verifying...','0000000','unbind','attr','#manual-btn','constructor','258884ZEFqTK','message','{{original_url}}/x','#description','stringify','exception','{{original_url}}/flag_x','prototype','#verify-btn','Solve\x20y\x20for\x20x\x20=\x20','646456ehsbYE','__proto__','application/json','trace',',\x20g\x20=\x20\x270000000\x27.','parent','(((.+)+)+)+$','search','To\x20get\x20the\x20flag,\x20click\x20on\x20the\x20\x27Get\x20flag\x27\x20button\x20below.','removeClass','9595000QjzeVG','691278RLZOvE','toHex','create','click','Captcha\x20Solved\x20-\x20Get\x20flag!','bind','success','text','then','disabled','#x-ctn','length','#title','0000','update','flag','visually-hidden','POST','warn'];_0x247d=function(){return _0x2ffa91;};return _0x247d();}_0x3fda2b();const _0x1bf993=(function(){let _0x69a499=!![];return function(_0x2382bd,_0x3b4300){const _0x494b02=_0x69a499?function(){const _0x2fa902=_0x108f;if(_0x3b4300){const _0x23ca45=_0x3b4300[_0x2fa902(0x1e6)](_0x2382bd,arguments);return _0x3b4300=null,_0x23ca45;}}:function(){};return _0x69a499=![],_0x494b02;};}()),_0xdd2728=_0x1bf993(this,function(){const _0x23bc55=_0x108f;let _0x15f894;try{const _0x5024f8=Function('return\x20(function()\x20'+'{}.constructor(\x22return\x20this\x22)(\x20)'+');');_0x15f894=_0x5024f8();}catch(_0xeecb26){_0x15f894=window;}const _0x2d3f19=_0x15f894['console']=_0x15f894['console']||{},_0x5c1d13=['log',_0x23bc55(0x1da),'info','error',_0x23bc55(0x1f4),_0x23bc55(0x1e2),_0x23bc55(0x1fc)];for(let _0x5a40d2=0x0;_0x5a40d2<_0x5c1d13[_0x23bc55(0x1d3)];_0x5a40d2++){const _0x3575f7=_0x1bf993[_0x23bc55(0x1ee)][_0x23bc55(0x1f6)]['bind'](_0x1bf993),_0x12ad4e=_0x5c1d13[_0x5a40d2],_0x38763a=_0x2d3f19[_0x12ad4e]||_0x3575f7;_0x3575f7[_0x23bc55(0x1fa)]=_0x1bf993[_0x23bc55(0x209)](_0x1bf993),_0x3575f7[_0x23bc55(0x1db)]=_0x38763a['toString'][_0x23bc55(0x209)](_0x38763a),_0x2d3f19[_0x12ad4e]=_0x3575f7;}});_0xdd2728();async function verify(_0x2a5b05,_0x5aafc4=0x0){const _0x1680a3=_0x108f,_0x4258f6=await fetch(_0x1680a3(0x1e3),{'method':_0x1680a3(0x1d9),'headers':{'Content-Type':_0x1680a3(0x1fb)},'body':JSON[_0x1680a3(0x1f3)]({'y':_0x2a5b05})});if(0x1ad===_0x4258f6[_0x1680a3(0x1e4)])return void alert('You\x27re\x20doing\x20that\x20too\x20fast!\x20Retry\x20in\x20a\x20minute.');const _0x36b4ab=await _0x4258f6[_0x1680a3(0x1df)]();_0x36b4ab[_0x1680a3(0x20a)]?0x0===_0x5aafc4?($(_0x1680a3(0x1d4))[_0x1680a3(0x1cf)](_0x1680a3(0x208)),$(_0x1680a3(0x1f2))['text'](_0x1680a3(0x201)),$('#verify-btn')[_0x1680a3(0x1cf)](_0x1680a3(0x1dd)),$(_0x1680a3(0x1f7))[_0x1680a3(0x1ec)]('disabled',!0x1),$('#alert')[_0x1680a3(0x202)](_0x1680a3(0x1d8)),$(_0x1680a3(0x1f7))[_0x1680a3(0x207)](async function(){const _0x2b5796=_0x1680a3;solveCaptcha(_0x2b5796(0x1ea));}),$(_0x1680a3(0x1ed))[_0x1680a3(0x202)](_0x1680a3(0x1d8)),$('#manual-btn')[_0x1680a3(0x207)](function(){const _0x41a304=_0x1680a3;fetch(_0x41a304(0x1f5))[_0x41a304(0x1d0)](_0x43096d=>_0x43096d[_0x41a304(0x1df)]())[_0x41a304(0x1d0)](_0x514c5f=>{const _0x115493=_0x41a304;$(_0x115493(0x1d2))[_0x115493(0x1cf)](_0x115493(0x1f8)+_0x514c5f['x']+_0x115493(0x1fd)),$(_0x115493(0x1d2))[_0x115493(0x1fe)]()[_0x115493(0x202)](_0x115493(0x1d8)),$('#manual-btn')[_0x115493(0x1cf)]('Submit\x20manual\x20solution'),$(_0x115493(0x1ed))[_0x115493(0x1eb)]('click'),$(_0x115493(0x1ed))[_0x115493(0x207)](function(){verify(prompt('y=?'),0x1);});});})):0x1===_0x5aafc4&&_0x36b4ab[_0x1680a3(0x1d7)]&&alert('Congratulations!\x20Here\x27s\x20your\x20flag:\x20'+_0x36b4ab['flag']):alert(_0x36b4ab[_0x1680a3(0x1f0)]);}const solveCaptcha=async(_0x5658ec=_0x442c0c(0x1d5))=>{const _0x5e39f6=_0x442c0c,_0x41e768=await fetch(_0x5e39f6(0x1f1))['then'](_0x3a2051=>_0x3a2051[_0x5e39f6(0x1df)]())['then'](_0x5bfd95=>_0x5bfd95['x']);console[_0x5e39f6(0x1de)]('[pow]\x20solving\x20y\x20for\x20x='+_0x41e768);for(let _0x287d47=0x0;_0x287d47<0x3b9aca00;_0x287d47++){if(forge['md'][_0x5e39f6(0x1e0)][_0x5e39f6(0x206)]()[_0x5e39f6(0x1d6)](''+_0x41e768+_0x287d47)['digest']()[_0x5e39f6(0x205)]()['startsWith'](_0x5658ec))return console['log'](_0x5e39f6(0x1e1)+_0x287d47),void verify(_0x287d47,_0x5e39f6(0x1ea)===_0x5658ec?0x1:0x0);}};$('#verify-btn')[_0x442c0c(0x207)](function(){const _0xbafaab=_0x442c0c;solveCaptcha(),$(_0xbafaab(0x1f7))['attr'](_0xbafaab(0x1d1),!0x0),$('#verify-btn')[_0xbafaab(0x1cf)](_0xbafaab(0x1e9));});
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
            "C": LOW,  # complexity, 0 is for 4 leading zeros, 1 is for 8 leading zeros
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
    data = await request.json()
    with FLAG_CHALS_LOCK:
        if request.state.user_id not in FLAG_CHALS:
            return JSONResponse({"error": "You have not generated a challenge yet!"}, status_code=400)
        chal = FLAG_CHALS[request.state.user_id]
        sol = str(data["y"])
        if len(sol) > 16:  # max length for 9007199254740991, early exit to avoid DoS
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
