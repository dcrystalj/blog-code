FROM python:3.11-slim-bullseye

RUN pip install "celery[redis]" flower

RUN apt update && apt install -y curl

WORKDIR /app

CMD ["python3", "main.py"]
