# HTMX Chat App - Go and SSE

A basic multi user chat web app built using Go, and using HTMX for all frontend and UI interactions. The chat system is based on [Server Sent Events (SSE) ](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) rather than Websockets.

The app looks, feels and interacts like a SPA, without full page reloads, but has zero lines of JavaScript

Built using:

- [Go](https://go.dev/)
- [Echo](https://echo.labstack.com/) - Minimal web framework and router
- [HTMX](https://htmx.org/) - High power tools for HTML
- [Bulma](https://bulma.io/) - CSS framework
- [Font Awesome](https://fontawesome.com/) - Icons

```
.
├── 📂 app        - Go source code for server and app
├── 📂 build      - Docker and deployment scripts
└── 📂 templates  - HTML templates & fragments used by the app
```

[![CI Workflow](https://github.com/benc-uk/htmx-go-chat/actions/workflows/ci-workflow.yaml/badge.svg)](https://github.com/benc-uk/htmx-go-chat/actions/workflows/ci-workflow.yaml)

## Developer Guide

Pre-reqs

- Go
- A Linux compatible system with bash, make, curl etc

Makefile reference:

```text
help                 💬 This help message :)
install-tools        🔧 Install dev tools into local project tools directory
watch                👀 Run the server with reloading
run                  🚀 Run the server
run-container        📦 Run the server from container
build                🔨 Build the server
lint                 🔍 Lint & format check only, sets exit code on error for CI
lint-fix             📝 Lint & format, attempts to fix errors & modify code
image                🐳 Build the container image
push                 📤 Push the container image to the image registry
deploy               ⛅ Deploy to Azure
```

### Running Locally

```bash
make run
```

Or for reloading/watching

```bash
make install-tools
make watch
```


## Deploying to Azure

Blah