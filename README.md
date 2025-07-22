# Social Network Application

A modern social networking platform built with Vue.js and Go, containerized with Docker for easy deployment and scalability.

## Tech Stack

### Frontend
- Vue.js 3
- Vue Router for navigation
- Vuex for state management
- Vue Toast Notification for alerts
- Modern browser support

### Backend
- Go (Golang)
- RESTful API architecture
- Docker containerization

## Project Structure
```
.
├── frontend/          # Vue.js frontend application
├── backend/           # Go backend server
├── docker-compose.yml # Docker composition configuration
├── diagram.png        # Project architecture diagram
└── docs/             # Project documentation
    ├── tasks.md
    ├── audit_checklist.md
    ├── checklist.md
    └── docker_readme.md
```

## Prerequisites

- Docker and Docker Compose
- Node.js and npm (for local frontend development)
- Go 1.x (for local backend development)

## Quick Start with Docker

1. Clone the repository:
```bash
git clone <repository-url>
cd social-network
```

2. Start the application using Docker Compose:
```bash
docker-compose up
```

The application will be available at:
- Frontend: http://localhost:8080
- Backend: http://localhost:8001

## Local Development Setup

### Frontend (Vue.js)

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Run development server:
```bash
npm run serve
```

The frontend development server will be available at http://localhost:8080

### Backend (Go)

1. Navigate to the backend directory:
```bash
cd backend
```

2. Run the Go server:
```bash
go run main.go
```

The backend server will be available at http://localhost:8001

## Features

- User authentication
- Profile management
- Social interactions
- Real-time notifications
- Responsive design
- Docker containerization for easy deployment

## Project Documentation

- `tasks.md` - Project tasks and roadmap
- `audit_checklist.md` - Security and performance audit checklist
- `checklist.md` - Development and deployment checklists
- `docker_readme.md` - Docker-specific instructions
- `diagram.png` - Visual representation of the project architecture

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Notes

- The project uses Docker for containerization, making it easy to deploy and scale
- Frontend is built with Vue.js 3 for modern, reactive user interfaces
- Backend is implemented in Go for high performance and reliability
- All services are containerized and can be started with a single docker-compose command
