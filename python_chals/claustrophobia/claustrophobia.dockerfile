FROM python:3.10-alpine3.20

RUN apk add --no-cache shadow bash
WORKDIR /root

COPY requirements.txt /root/requirements.txt
RUN pip install -r requirements.txt --no-cache-dir

COPY shared_helpers.py /root/shared_helpers.py
COPY ./claustrophobia/claustrophobia.py /root/claustrophobia.py
RUN chmod 700 /root/claustrophobia.py 

RUN echo '$FLAG_PLACEHOLDER' > /flag.txt

RUN adduser -D -s /sbin/nologin ctfplayer

WORKDIR /home/ctfplayer
RUN mkdir -p /home/ctfplayer/.ssh && \
    echo "ssh-rsa FAKEKEY1234567890" > /home/ctfplayer/.ssh/authorized_keys && \
    mkdir -p /home/ctfplayer/Documents /home/ctfplayer/Downloads /home/ctfplayer/Desktop && \
    echo "Confidential Project - Do Not Share" > /home/ctfplayer/Documents/secret.txt && \
    echo "Welcome to our CTF event!" > /home/ctfplayer/Desktop/readme.txt && \
    echo "GALA{not_the_real_one}" > /home/ctfplayer/Documents/flag.txt && \
    mkdir -p /home/ctfplayer/.config && \
    echo "alias ll='ls -la'" > /home/ctfplayer/.bashrc && \
    touch /home/ctfplayer/.config/fake_config.cfg && \
    echo "password=123456" > /home/ctfplayer/.config/fake_config.cfg && \
    chown -R ctfplayer:ctfplayer /home/ctfplayer


ENTRYPOINT ["python", "/root/claustrophobia.py"]
