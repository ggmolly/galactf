FROM python:3.10-alpine3.20

RUN mkdir -p /app
COPY ./cookie_monster/cookie_monster.py /app/cookie_monster.py
COPY shared_helpers.py /app/shared_helpers.py
COPY requirements.txt /app/requirements.txt

WORKDIR /app

RUN pip install -r requirements.txt --no-cache-dir

ENTRYPOINT ["python", "cookie_monster.py"]