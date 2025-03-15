FROM python:3.10-alpine3.20

RUN mkdir -p /app
COPY ./just_a_moment/just_a_moment.py /app/just_a_moment.py
COPY shared_helpers.py /app/shared_helpers.py
COPY requirements.txt /app/requirements.txt

WORKDIR /app

RUN pip install -r requirements.txt --no-cache-dir

ENTRYPOINT ["python", "just_a_moment.py"]