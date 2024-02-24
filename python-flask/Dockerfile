# Build and Run
FROM python:3.13.0a4-alpine3.19
WORKDIR /app
COPY . .
RUN pip install gunicorn
RUN pip install -r requirements.txt

EXPOSE 3030
CMD ["gunicorn", "-b", "0.0.0.0:3030", "app:app"]
