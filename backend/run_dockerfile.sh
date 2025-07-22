docker build --tag backend-image .

docker run --name backend-container --publish 8000:8000 backend-image