import asyncio
import random
import string
from starlette.applications import Starlette
from starlette.responses import HTMLResponse
from starlette.requests import Request
from jinja2 import Template
from shared_helpers import UserIDMiddleware, OriginalURLMiddleware

blacklisted_words = set(["import","eval","exec","execfile","subprocess","os","sys","time","sleep","input","print","open","read","write","close","exit","for","while","break","continue","pass","return","yield","range","quit","exit"])
app = Starlette()
app.add_middleware(UserIDMiddleware)
app.add_middleware(OriginalURLMiddleware)

def random_string(user_id: int, length: int) -> str:
    rng = random.Random(user_id)
    return ''.join(rng.choice(string.ascii_letters) for _ in range(length))

async def safe_eval(expression: str, flag: str, user_id: int):
    if any(word in expression.lower() for word in blacklisted_words):
        raise ValueError
    loop = asyncio.get_running_loop()
    return await loop.run_in_executor(None, lambda: eval(expression, {random_string(user_id, 8): lambda: flag}))

@app.route("/", methods=["GET", "POST"])
async def index(request: Request):
    FLAG = request.headers.get("X-GalaCTF-Flag")
    if request.method == "POST":
        data = await request.form()
        expression = f"{data['x']}{data['operator']}{data['y']}"
        try:
            value = await asyncio.wait_for(safe_eval(expression, FLAG, int(request.headers.get("X-User-ID", "0"))), timeout=1.0)
        except asyncio.TimeoutError:
            value = "Timed out"
        except ValueError:
            value = "Blocked by WAF"
        except Exception as e:
            value = f"Error: {e}"
    else:
        value = ""
    template = Template("""
    <!DOCTYPE html>
    <html lang="fr">
    <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>quantum calculator</title>
    <style>
        body {
        font-family: Arial, sans-serif;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
        background-color: #f2f2f2;
        }
        .calculator {
        background-color: #fff;
        padding: 20px;
        border-radius: 10px;
        box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        #display {
        width: 95%;
        height: 50px;
        font-size: 1.5em;
        margin-bottom: 10px;
        text-align: right;
        padding: 10px;
        border: 1px solid #ccc;
        border-radius: 5px;
        }
        .buttons {
        display: grid;
        grid-template-columns: repeat(4, 60px);
        gap: 10px;
        justify-content: center;
        }
        .buttons button {
        width: 60px;
        height: 60px;
        font-size: 1.2em;
        border: none;
        border-radius: 5px;
        background-color: #e0e0e0;
        cursor: pointer;
        transition: background-color 0.2s;
        }
        .buttons button:hover {
        background-color: #d4d4d4;
        }
        .buttons button.operator {
        background-color: #f9a825;
        color: #fff;
        }
        .buttons button.operator:hover {
        background-color: #f57f17;
        }
        .buttons button.equal {
        background-color: #388e3c;
        color: #fff;
        }
        .buttons button.equal:hover {
        background-color: #2e7d32;
        }
        .buttons button.clear {
        background-color: #d32f2f;
        color: #fff;
        grid-column: span 4;
        }
        .buttons button.clear:hover {
        background-color: #c62828;
        }
    </style>
    </head>
    <body>
    <div class="calculator">
        <form id="calcForm" action="{{original_url}}/" method="post">
        <input type="text" id="display" readonly value="{{value}}">
        <input type="hidden" name="x" id="x">
        <input type="hidden" name="operator" id="operator">
        <input type="hidden" name="y" id="y">
        
        <div class="buttons">
            <button type="button" onclick="appendNumber('7')">7</button>
            <button type="button" onclick="appendNumber('8')">8</button>
            <button type="button" onclick="appendNumber('9')">9</button>
            <button type="button" class="operator" onclick="appendOperator('/')">/</button>

            <button type="button" onclick="appendNumber('4')">4</button>
            <button type="button" onclick="appendNumber('5')">5</button>
            <button type="button" onclick="appendNumber('6')">6</button>
            <button type="button" class="operator" onclick="appendOperator('*')">*</button>

            <button type="button" onclick="appendNumber('1')">1</button>
            <button type="button" onclick="appendNumber('2')">2</button>
            <button type="button" onclick="appendNumber('3')">3</button>
            <button type="button" class="operator" onclick="appendOperator('-')">-</button>

            <button type="button" onclick="appendNumber('0')">0</button>
            <button type="button" onclick="appendNumber('.')">.</button>
            <button type="button" class="operator" onclick="appendOperator('+')">+</button>
            <button type="submit" class="equal">=</button>

            <button type="button" class="clear" onclick="clearAll()">C</button>
        </div>
        </form>
    </div>

    <script>
        let firstOperand = "";
        let operatorSelected = "";
        let secondOperand = "";

        function updateDisplay() {
        document.getElementById("display").value = firstOperand + operatorSelected + secondOperand;
        }

        function appendNumber(num) {
        if (operatorSelected === "") {
            firstOperand += num;
        } else {
            secondOperand += num;
        }
        updateDisplay();
        }

        function appendOperator(op) {
        if (firstOperand !== "" && operatorSelected === "") {
            operatorSelected = op;
        }
        updateDisplay();
        }

        function clearAll() {
        firstOperand = "";
        operatorSelected = "";
        secondOperand = "";
        updateDisplay();
        }

        document.getElementById("calcForm").addEventListener("submit", function(e) {
        if (firstOperand === "" || operatorSelected === "" || secondOperand === "") {
            e.preventDefault();
            alert("Invalid operation.");
            return false;
        }
        document.getElementById("x").value = firstOperand;
        document.getElementById("operator").value = operatorSelected;
        document.getElementById("y").value = secondOperand;
        });
    </script>
    </body>
    </html>

""")
    return HTMLResponse(content=template.render(original_url=request.state.original_url, value=value))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
