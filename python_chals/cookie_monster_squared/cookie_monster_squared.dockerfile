FROM python:3.10-alpine3.20

RUN mkdir -p /app
COPY ./cookie_monster_squared/cookie_monster_squared.py /app/cookie_monster_squared.py
COPY shared_helpers.py /app/shared_helpers.py
COPY requirements.txt /app/requirements.txt

WORKDIR /app

RUN pip install -r requirements.txt --no-cache-dir

ENTRYPOINT ["python", "cookie_monster_squared.py"]
