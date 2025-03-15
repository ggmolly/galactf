FROM python:3.10-alpine3.20

RUN mkdir -p /app
COPY ./unsecure_notes/unsecure_notes.py /app/unsecure_notes.py
COPY shared_helpers.py /app/shared_helpers.py
COPY requirements.txt /app/requirements.txt

WORKDIR /app

RUN pip install -r requirements.txt --no-cache-dir

ENTRYPOINT ["python", "unsecure_notes.py"]