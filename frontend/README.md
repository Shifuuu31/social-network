# frontend

### Build and Run Docker Container (Only Frontend), you need Docker installed :)
Check you are in frontend directory.
Build image with 'sn-front' name:  `docker build -t sn-front . `
Run docker container and open ports: `docker run -p 8080:8080 sn-front`
Open Browser @ http://localhost:8080


## Project setup(for developing, needs once)
```
npm install
```

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```