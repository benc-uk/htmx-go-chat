# HTMX Chat App - Go and SSE

A basic multi user chat web app built using Go, and using HTMX for all frontend and UI interactions. The chat system is based on [Server Sent Events (SSE) ](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) rather than Websockets.

The app looks, feels and interacts like a SPA, without full page reloads, but has zero lines of JavaScript

Stack / Libraries

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