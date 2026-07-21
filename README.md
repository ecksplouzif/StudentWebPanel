# StudentWebPanel

A web panel for student groups: class schedule, homework deadlines, and announcements in one place.

## Project structure

```
.
├── main.go                  # entry point: config, routes, server startup
├── internal/
│   └── handlers/            # HTTP handlers
├── web/                     # templates and static assets
├── docs/                    # development plan and decision log
├── Makefile                 # run, build, lint, fmt
└── .env.example             # required environment variables
```