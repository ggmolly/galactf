import sqlite3
import threading
from starlette.applications import Starlette
from starlette.responses import JSONResponse, HTMLResponse
from starlette.routing import Route
from starlette.requests import Request
from jinja2 import Template
from shared_helpers import UserIDMiddleware

user_data_lock = threading.Lock()
user_databases = {}

BOOKS = [
    ("The Catcher in the Rye", "J.D. Salinger"),
    ("To Kill a Mockingbird", "Harper Lee"),
    ("1984", "George Orwell"),
    ("Pride and Prejudice", "Jane Austen"),
    ("The Great Gatsby", "F. Scott Fitzgerald"),
    ("Moby Dick", "Herman Melville"),
    ("War and Peace", "Leo Tolstoy"),
    ("The Odyssey", "Homer"),
    ("Crime and Punishment", "Fyodor Dostoevsky"),
    ("Brave New World", "Aldous Huxley"),
]

def get_user_db(user_id: int, flag: str = None):
    with user_data_lock:
        if user_id not in user_databases:
            conn = sqlite3.connect(":memory:", check_same_thread=False)
            cursor = conn.cursor()

            cursor.execute("""
                CREATE TABLE books (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    title TEXT NOT NULL,
                    author TEXT NOT NULL
                )
            """)

            cursor.executemany("INSERT INTO books (title, author) VALUES (?, ?)", BOOKS)

            cursor.execute("""
                CREATE TABLE flag (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    value TEXT NOT NULL
                )
            """)

            if flag:
                cursor.execute("INSERT INTO flag (value) VALUES (?)", (flag,))
            
            conn.commit()
            user_databases[user_id] = (conn, cursor)

        return user_databases[user_id]

async def list_books(request: Request):
    user_id = request.state.user_id
    flag = request.headers.get("X-GalaCTF-Flag")
    if flag is None:
        return JSONResponse({"error": "Invalid X-GalaCTF-Flag"}, status_code=401)
    conn, cursor = get_user_db(user_id, flag)

    search_query = request.query_params.get("q", "").strip()
    error_message = None
    books = []

    try:
        if search_query:
            query = f"SELECT id, title, author FROM books WHERE title LIKE '%{search_query}%'"
        else:
            query = "SELECT id, title, author FROM books"
        
        cursor.execute(query)
        books = [{"id": row[0], "title": row[1], "author": row[2]} for row in cursor.fetchall()]
    except Exception as e:
        error_message = str(e)

    template = Template("""
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>bobby's book library</title>
        <style>
            body { font-family: Arial, sans-serif; margin: 40px; text-align: center; }
            h1 { margin-bottom: 20px; }
            form { margin-bottom: 20px; }
            input[type="text"] { padding: 8px; width: 200px; }
            button { padding: 8px; cursor: pointer; }
            table { width: 100%; border-collapse: collapse; margin-top: 20px; }
            th, td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
            th { background-color: #f4f4f4; }
            .error { color: red; margin-top: 20px; }
        </style>
    </head>
    <body>
        <h1>bobby's book library</h1>
        <form method="get">
            <input type="text" name="q" placeholder="Search by title" value="{{ search_query }}">
            <button type="submit">Search</button>
        </form>
        {% if error_message %}
            <p class="error">SQL Error: {{ error_message }}</p>
        {% endif %}
        <table>
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Title</th>
                    <th>Author</th>
                </tr>
            </thead>
            <tbody>
                {% for book in books %}
                <tr>
                    <td>{{ book.id }}</td>
                    <td>{{ book.title }}</td>
                    <td>{{ book.author }}</td>
                </tr>
                {% endfor %}
            </tbody>
        </table>
    </body>
    </html>
    """)

    return HTMLResponse(template.render(books=books, search_query=search_query, error_message=error_message))

app = Starlette(routes=[Route("/", list_books)])
app.add_middleware(UserIDMiddleware)

if __name__ == '__main__':
    import uvicorn
    uvicorn.run(app, host='0.0.0.0', port=8080)
