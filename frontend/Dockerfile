FROM node:lts-alpine
WORKDIR /frontend
# Copy only the necessary files initially (package.json and package-lock.json)
COPY package*.json ./
# Install dependencies
RUN npm install
COPY . .
RUN npm run build
EXPOSE 3000
CMD ["npm", "run", "serve"]
# ROM node:16.15.0
# RUN mkdir /frontend
# ADD . /frontend
# WORKDIR /frontend
# COPY ../../package*.json /frontend
# RUN npm install
# COPY . /frontend
# EXPOSE 8080
# CMD ["npm","run","serve"]