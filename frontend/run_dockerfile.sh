docker build --tag frontend-image .

# run backend+frontend images in one container "social-network" with docker-compose

docker run --name frontend-container --publish 8080:8080 frontend-image