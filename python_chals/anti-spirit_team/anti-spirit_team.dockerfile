FROM python:3.10-alpine3.20

RUN mkdir -p /app
WORKDIR /app

COPY requirements.txt /app/requirements.txt

RUN pip install -r requirements.txt --no-cache-dir

COPY shared_helpers.py /app/shared_helpers.py
COPY ./anti-spirit_team/anti-spirit_team.py /app/anti-spirit_team.py


ENTRYPOINT ["python", "anti-spirit_team.py"]