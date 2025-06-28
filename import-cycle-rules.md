| Folder       | Contains                               | Can import from          | Must NOT import from     |
| ------------ | -------------------------------------- | ------------------------ | ------------------------ |
| `models`     | `User`, `Post`, etc. (structs only)    | `tools`                  | `db`, `handlers`         |
| `tools`      | `DecodeJSON`, `HashPassword`, etc.     | nothing (pure utils)     | anything in app          |
| `db`         | SQL queries, SQLite connection         | `models`, `tools`        | `handlers`, `middleware` |
| `handlers`   | All business logic (auth, posts, etc.) | `db`, `models`, `tools`  | nothing forbidden        |
| `middleware` | Auth, logging, cors, etc.              | `tools`, `models`        | `handlers`, `db`         |
| `router`     | Sets routes and groups using handlers  | `handlers`, `middleware` | `db`, `models`, `tools`  |
| `server.go`  | Main entry point                       | everything above         |                          |
