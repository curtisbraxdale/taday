# Taday API

**Taday** is a modern RESTful API backend for a productivity web app that delivers daily agendas to users via SMS. Built with Go, PostgreSQL, and Fly.io deployment, the backend features secure authentication, user management, scheduling, tagging, and token-based authorization using HTTP-only cookies.

I built this web app because I have gotten into the digital-minimalist / dumbphone community. As such I try to have as few applications and as little-to-no notifications on my phone. I wanted to be able to stay productive and be reminded of my tasks and events without feeling like I was constantly reaching for my phone again and again. I still have text and call notifications enabled on my phone for emergencies and that's when the idea popped in my head. I could send myself a daily agenda via sms to my phone every morning at 8:00 am, in it there would be events that I had scheduled for the day and some todos that I had logged. This might seem like a regression in technology but for someone like me trying to reduce their screen-time as much as possible I thought that this would be a good use-case for a web app like mind. No app required, no persistent notifications. Just one daily text with the information that you have created and asked for.

---

## 🚀 Features

* **JWT-based Authentication** with Refresh Tokens
* **Secure Login, Registration, and Logout Flows**
* **Persistent Cookie-based Sessions** (HTTP-only, `SameSite=None`, `Secure`)
* **CRUD operations for:**

  * Users
  * Todos (day-based tasks)
  * Events (calendar blocks with recurrence)
  * Tags (labeling for events)
* **Event-Tag Relationships** (many-to-many)
* **Fly.io Deployment** (production-ready Docker container)

---

## 🔧 Tech Stack

| Layer      | Tech                      |
| ---------- | ------------------------- |
| Language   | Go                        |
| Database   | PostgreSQL (sqlc + pgx)   |
| Auth       | JWT + bcrypt              |
| Deployment | Fly.io                    |
| CORS/Auth  | Secure Cookies + CORS     |

---

## 🛡 Authentication

The API uses a combination of access and refresh JWTs:

* **Access Token**: Short-lived (1 hour), stored in an HTTP-only cookie
* **Refresh Token**: Long-lived (60 days), stored in HTTP-only cookie and validated against DB

### Endpoints

* `POST /api/login` — sets `access_token` and `refresh_token` cookies on success
* `POST /api/logout` — clears both cookies and revokes token in DB
* `POST /api/refresh` — validates refresh token and issues a new access token

All cookies are:

* `HttpOnly`
* `Secure`
* `SameSite=None`

> ✅ Compatible with cross-site frontends like `https://taday.io`

---

## 🧑‍💻 API Endpoints

### Users

| Method | Endpoint     | Description             |
| ------ | ------------ | ----------------------- |
| POST   | `/api/users` | Create a new user       |
| GET    | `/api/users` | Get current user (auth) |
| PUT    | `/api/users` | Update current user     |
| DELETE | `/api/users` | Delete account          |

### Todos

| Method | Endpoint         | Description         |
| ------ | ---------------- | ------------------- |
| GET    | `/api/todos`     | Get all todos       |
| POST   | `/api/todos`     | Create a todo       |
| GET    | `/api/todos/:id` | Get a specific todo |
| PUT    | `/api/todos/:id` | Update a todo       |
| DELETE | `/api/todos/:id` | Delete a todo       |

### Events

| Method | Endpoint          | Description             |
| ------ | ----------------- | ----------------------- |
| GET    | `/api/events`     | Get events (filterable) |
| POST   | `/api/events`     | Create event            |
| GET    | `/api/events/:id` | Get specific event      |
| PUT    | `/api/events/:id` | Update event            |
| DELETE | `/api/events/:id` | Delete event            |

### Tags

| Method | Endpoint        | Description    |
| ------ | --------------- | -------------- |
| GET    | `/api/tags`     | List all tags  |
| POST   | `/api/tags`     | Create new tag |
| PUT    | `/api/tags/:id` | Update tag     |
| DELETE | `/api/tags/:id` | Delete tag     |

### Event-Tag Relationship

| Method | Endpoint                      | Description           |
| ------ | ----------------------------- | --------------------- |
| GET    | `/api/events/:id/tags`        | List tags on event    |
| POST   | `/api/events/:id/tags`        | Add tag to event      |
| DELETE | `/api/events/:id/tags/:tagId` | Remove tag from event |

---

## 🧩 Frontend Integration

The frontend is hosted on [`https://taday.io`](https://taday.io). It communicates with this API by including credentials in fetch calls:

```ts
fetch("https://taday-api.fly.dev/api/users", {
  method: "GET",
  credentials: "include",
});
```

This ensures cookies (access & refresh tokens) are sent and handled securely across domains.

---

## ⚙️ Setup (Dev)

1. **Start PostgreSQL** (or use a Fly.io volume)
2. **Run migrations** with `sqlc`
3. **Build the Go server**

```bash
go run ./cmd/api
```

---

## 📦 Deployment

App is deployed to Fly.io using Docker. On login, cookies are set with production flags:

* `Secure: true`
* `SameSite=None`
* `HttpOnly: true`

TLS termination is handled by Fly.io.

---

## 📄 License

MIT — © Curtis Braxdale
