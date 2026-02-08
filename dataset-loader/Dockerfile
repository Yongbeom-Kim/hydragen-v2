FROM python:3.13-alpine3.23

WORKDIR /app/dataset-loader

COPY ./requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

CMD ["/bin/bash"]