FROM python:3.10-alpine3.20

RUN mkdir -p /app
WORKDIR /app

COPY requirements.txt /app/requirements.txt

RUN pip install -r requirements.txt --no-cache-dir

COPY shared_helpers.py /app/shared_helpers.py
COPY ./just_a_moment/just_a_moment.py /app/just_a_moment.py


ENTRYPOINT ["python", "just_a_moment.py"]