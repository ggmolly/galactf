FROM python:3.10-alpine3.20

RUN mkdir -p /app
WORKDIR /app

COPY requirements.txt /app/requirements.txt

RUN pip install -r requirements.txt --no-cache-dir

COPY shared_helpers.py /app/shared_helpers.py
COPY ./bobby_library/bobby_library.py /app/bobby_library.py


ENTRYPOINT ["python", "bobby_library.py"]