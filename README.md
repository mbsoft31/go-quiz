# Quiz Go Application

This project is a Go-based quiz management system with a TailwindCSS frontend, templated views, and a SQLite database. The project is set up to run in both development and production environments using Docker.

## Prerequisites

- Docker and Docker Compose installed on your system.
- WSL2 or Windows setup (if running on Windows).

## Development Setup

1. **Build and run the development environment:**

```bash
docker-compose up --build
```

This command will:
- Start the Go application with hot reloading.
- Watch and compile TailwindCSS.
- Watch and generate Templ views.
- Watch and compile JavaScript with esbuild.

2. **Access the application:**

The application will be available at `http://localhost:4000`.

3. **Stopping the development environment:**

Press `Ctrl+C` or run:

```bash
docker-compose down
```

## Production Setup

1. **Build the production Docker image:**

```bash
docker build -t quiz-go-app .
```

2. **Run the production container:**

```bash
docker run -p 4000:4000 quiz-go-app
```

3. **Access the production application:**

The application will be available at `http://localhost:4000`.

## Project Structure

- **Dockerfile:** Production Docker setup.
- **Dockerfile.dev:** Development Docker setup with hot reloading.
- **docker-compose.yml:** Configuration for running the development environment with Docker Compose.
- **Makefile:** Local development tasks for CSS, JS, and Templ.
- **public/:** Static files served by the application.
- **views/:** Templ views and CSS files.
- **database/:** SQLite database file.

## Notes

- If you are running on WSL or Windows, make sure Docker is configured to use WSL2 for better performance.
- If you encounter permission issues with mounted volumes, ensure that your WSL or Docker Desktop has the correct permissions.
