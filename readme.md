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
📂
 ├── app/        - Go source code for server and app
 ├── build/      - Docker and deployment scripts
 └── templates/  - HTML templates & fragments used by the app
```

[![CI Workflow](https://github.com/benc-uk/htmx-go-chat/actions/workflows/ci-workflow.yaml/badge.svg)](https://github.com/benc-uk/htmx-go-chat/actions/workflows/ci-workflow.yaml)

## 📝 Design Notes

The Go code resides in `app/` directory and, comprises a single `main` package, breaking it up over multiple packages was deemed unnecessary.

- `server.go` Main entry point and HTTP server, using Echo.
- `routes.go` All HTTP routes and endpoints, most of the app logic is here, and mostly returns rendered HTML templates.
- `renderer.go` Implements a HTML template renderer using the [html/template](https://pkg.go.dev/html/template) package, part of the Go standard library.
- `chat.go` See below.

All the HTML served by the app is held within the `templates/` folder. This is a mixture of full pages like `index.html` and HTML fragments used for various parts of the app, as well as any custom CSS.

The main views are the `login` template and the `chat` template which is only shown after users login. The term login is a misnomer here, all users have to do is enter their name to enter the chat, there is no formal login process or actual usernames & auth.

## 🎭 Chat Broker

The chat broker is the core part of the app that handles the multi-user chat using Server Side Events (SSE).

The SSE implementation has been factored out into this repo [benc-uk/go-rest-api](https://github.com/benc-uk/go-rest-api). This generic broker provides the SSE stream handler, which holds open the HTTP connection and streams events as they arrive, plus a connection registry, which handles multiple client connections, using Go channels. 

The broker uses generics and exposes callbacks/handlers for message formatting, and hooking into connection/disconnection events. As such, the code in `chat.go` uses these hooks as follows:

- Defines a `ChatMessage` struct type, used to send chat messages as SSE events.
- An in-memory message store so we can send some previous history of messages when users connect.
- Message formatter, converting `ChatMessage`'s to HTML in a HTMX friendly way.
- Handlers for connections/disconnections in order to notify other users.

One interesting thing about SSE is you can access the stream of events over a regular HTTP connection. So debugging and viewing the chat stream can be done by connecting to the `/chat-stream` URL directly in your browser.

## 🧑‍💻 Developer Guide

Pre-reqs

- Go (v1.20+)
- A Linux compatible system with bash, make, curl etc

Makefile reference:

```text
help                 💬 This help message :)
install-tools        🔧 Install dev tools into local project directory
watch                👀 Run the server with reloading
run                  🚀 Run the server
run-container        📦 Run from container
build                🔨 Build the server binary only
lint                 🔍 Lint & format check only, sets exit code on error for CI
lint-fix             📝 Lint & format, attempts to fix errors & modify code
image                🐳 Build container image
push                 📤 Push container image to the image registry
deploy               ⛅ Deploy to Azure
```

### Running Locally

Quickly run the server

```bash
make run
```

Open http://localhost:8000 in a browser

To run with reloading/watching on code changes

```bash
make install-tools
make watch
```

## 🐋 Building your own image

The makefile has two targets `image` and `push` which can be run to build and push an image. Set the variables `IMAGE_NAME` and `VERSION` to change the image name and tag. The image name should be fully qualified and include the registry if you are pushing it to one.

For example to build an image named `bob/my-chatapp` tagged with `dev` and pushed to the `myreg.io` registry, you would run:

```bash
make image push IMAGE_NAME=myreg.io/bob/my-chatapp VERSION=dev 
```

## ⛅ Deploying to Azure

You probably don't want to bother with this, but you can if you wish

```
make deploy
```