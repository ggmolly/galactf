FROM python:3.10-alpine3.20

RUN mkdir -p /app
WORKDIR /app

COPY requirements.txt /app/requirements.txt

RUN pip install -r requirements.txt --no-cache-dir

COPY shared_helpers.py /app/shared_helpers.py
COPY ./unsecure_notes/unsecure_notes.py /app/unsecure_notes.py


ENTRYPOINT ["python", "unsecure_notes.py"]