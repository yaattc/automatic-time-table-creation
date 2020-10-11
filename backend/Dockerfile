FROM python:3.8

COPY . /srv

RUN pip3 install -r /srv/requirements.txt

WORKDIR /srv

EXPOSE 8083
ENTRYPOINT ["python3", "/srv/main.py"]