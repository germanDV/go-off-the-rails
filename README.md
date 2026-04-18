# Go Off The Rails (GOTR)

A web framework for Go. Heavily inspired by Ruby on Rails.

It is heavily opinionated and adapted to particular needs.
You are probably better of using Ruby on Rails, Phoenix or Loco.rs instead.
But if the opinions fit your needs as well, welcome!

## Auth

It scaffolds a multi-tenant system with basic RBAC (`ADMIN` and `USER` roles).
When a user signs up, an org is created for them and they are assigned to the `ADMIN` role.
To join an existing org, an `ADMIN` needs to send an invitation.

The auth token is a JWT stored in a secured http-only cookie.
The token is long-lived (7 days).

A `SUPERADMIN` role is available for the maintainers of the system to perform administrative tasks.

### Data Models

```
orgs:   id, name, created_at, updated_at
users:  id, email, password_hash, role, org_id, created_at, updated_at
```

### JWT

```json
{
  "sub": "user_id",
  "org_id": "org_id",
  "role": "ADMIN",
   ...
}
```

## Database

SQLite. IDs are UUIDv7.

SQL is a firt-class citizen in GOTR. We don't use ORMs. You have to write your migrations and queries. Which then produced Go code using `sqlc`.

The CLI can generate migrations and CRUD queries for you to have a starting point:

```
gotr generate scaffold movies title:string! rating:int
```

Produces:
```
db/migrations/002_create_movies.sql
db/queries/movies_queries.sql
```

```sql
-- db/migrations/002_create_movies.sql
CREATE TABLE movies (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  rating INTEGER,
  version INTEGER NOT NULL DEFAULT 1,
  org_id TEXT NOT NULL REFERENCES orgs(id),
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

```sql
-- db/queries/movies_queries.sql

-- name: GetMovie :one
SELECT * FROM movies WHERE id = ? AND org_id = ?;

-- name: ListMovies :many
SELECT * FROM movies
WHERE org_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateMovie :exec
INSERT INTO movies (id, title, rating, org_id, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateMovie :exec
UPDATE movies
SET title = ?, rating = ?, updated_at = ?, version = version + 1
WHERE id = ? AND org_id = ? AND version = ?;

-- name: DeleteMovie :exec
DELETE FROM movies WHERE id = ? AND org_id = ?;
```

## Project Structure

```
app/
  domain/
    movie.go
  controllers/
    movies_controller.go
  views/
    layout.templ
    sidebar.templ
    login.templ
    movies_index.templ    // views.MoviesIndex
    movies_show.templ     // views.MoviesShow
    movies_new.templ      // views.MoviesNew
    movies_edit.templ     // views.MoviesEdit
  db/
    queries/
      movies_queries.sql
    migrations/
      001_movies_migration.sql
    generated/
      movies.sql.go  // sqlc output
    orgs_repository.go
    movies_repository.go
```


- `controllers` package is an orchestration layer, it depends on the `db` repositories to read and write entities, it delegates all business logic to the domain entities
- `db` package owns all data access and mapping from sqlc structs to domain entities
- `domain` package defines the entities and owns all the business logic as well as validation rules
- `views` package is a templating layer that is in charge of rendering the domain, it uses `Templ` templates


```go
// domain/movie.go
type Movie struct {
    ID     uuid.UUID
    OrgID  uuid.UUID
    Title  string
    Rating int
    Version int
    CreatedAt time.Time
    UpdatedAt time.Time
}

// db/movies_repository.go
type MoviesRepository struct {
    queries *generated.Queries  // sqlc generated
}

func (r *MoviesRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Movie, error) {
    row, err := r.queries.GetMovie(ctx, id.String())
    if err != nil {
        return domain.Movie{}, err
    }
    return domain.Movie{          // mapping lives here
        ID:     uuid.MustParse(row.ID),
        Title:  row.Title,
        Rating: int(row.Rating),
    }, nil
}

// controllers/movies.go
type MoviesController struct {
    movies *db.MovieRepository
}
```

## Domain Services

When business logic grows beyond what a single entity can encapsulate — for example, computing a result that involves multiple entities — a **domain service** can be introduced in the `domain` package.

Domain services are kept **pure**: they receive entities, perform in-memory transformations or computations, and return the results. They have no dependencies on repositories, databases, or any infrastructure.

The **controller acts as a workflow script**: it reads entities from repositories, passes them to the domain service, and then takes care of persistence and any other side effects (notifications, logging, etc.).

```go
// domain/movie_pricing.go — pure domain service, no dependencies
func ApplySeasonalDiscount(movie Movie, season string) (Movie, float64) {
    discount := 0.0
    if movie.Rating >= 8 && season == "summer" {
        discount = 0.15
    }
    return movie, discount
}

// controllers/movies_controller.go — orchestrates the workflow
func (c *MoviesController) ApplyDiscount(w http.ResponseWriter, r *http.Request) {
    movie, err := c.moviesRepo.GetByID(r.Context(), orgID, movieID)
    // ...
    updated, discount := domain.ApplySeasonalDiscount(movie, "summer")
    err = c.moviesRepo.Update(r.Context(), updated)
    // ... send notification, render response, etc.
}
```

This keeps domain logic easy to test (no mocks needed) while avoiding unnecessary layers of abstraction. If a workflow later needs to be reused from multiple entry points (HTTP, CLI, background jobs), it can be promoted into a service that accepts repository interfaces — but that step is deferred until actually needed.

## Controllers

### Pagination

Listing endpoints are paginated by default. The `limit` and `offset` query params are optional, with useful defaults.

```go
// domain/pagination.go
type PaginationParams struct {
    Limit  int
    Offset int
}

func NewPaginationParams(limit string, offset string) PaginationParams {
	return PaginationParams{
		limit:  cap(parseIntOr(limit, DefaultLimit), MaxLimit),
		offset: parseIntOr(offset, DefaultOffset),
	}
}

// controllers/movies.go
func (c *MoviesController) Index(w http.ResponseWriter, r *http.Request) {
    page := domain.NewPaginationParams(
        r.URL.Query().Get("limit"),
        r.URL.Query().Get("offset"),
    )

    movies, err := c.moviesRepo.List(r.Context(), page)
    ...
}
```

`views/*_index.templ` automatically renders the pagination links.

### Validation

```go
// domain/movies.go
func NewMovie(orgID uuid.UUID, titleInput string, ratingInput string) (Movie, error) {
	title := strings.TrimSpace(titleInput)

	rating, err := strconv.Atoi(strings.TrimSpace(ratingInput))
	if err != nil {
		return Movie{}, err
	}

	err = validate(title, rating)
	if err != nil {
		return Movie{}, err
	}

	return Movie{
		ID:        uuid.Must(uuid.NewV7()),
		OrgID:     orgID,
		Title:     title,
		Rating:    rating,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (m *Movie) Update(title string, rating string) error {
   if err := validate(title, rating); err != nil {
        return Movie{}, err
    }

    m.Title = title
    m.Rating = rating
    m.UpdatedAt = time.Now().UTC()

    return nil
}

func validate(title string, rating int) error {
	if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 100 {
		return errors.New("title must be between 3 and 100 characters")
	}

	if rating < 1 || rating > 10 {
		return errors.New("rating must be between 1 and 10")
	}

    return nil
}
```

### Middleware

- auth - checks if there is an auth token in the request and if so, decodes it and sets the user in the context
- recover - recovers from panics and renders a nice error page
- logger - logs the request
- helmet - sets some security headers
- realip - sets the real ip address of the request
- ratelimit - limits the number of requests per IP per second

## Views

`gotr` provides a basic UI which includes a sidebar with auth-related actions and pages for CRUD operations on wahtever new things you generate.

GOTR uses [HTMX](https://htmx.org/) to enhance requests to the backend. Among other things, this allows forms to use HTTP verbs like `PUT` and `DELETE` instead of relying solely on `POST`. Data is still communicated to the backend as form values — not JSON.

```go
templ Layout(title string, ctx context.Context) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="UTF-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
            <title>{ title }</title>
            <link rel="stylesheet" href="/static/styles.css"/>
            <script src="https://unpkg.com/htmx.org@2" defer></script>
        </head>
        <body>
            @maybeSidebar(ctx)
            <main>
                { children... }
            </main>
        </body>
    </html>
}

templ maybeSidebar(ctx context.Context) {
    if user, ok := userFromContext(ctx); ok {
        @sidebar(user)
    }
}

templ sidebar(user domain.User) {
    <nav>
        <span>{ user.Email }</span>
        <form action="/auth/signout" method="POST">
            <button type="submit">Sign out</button>
        </form>
    </nav>
}
```

## Routes And Views

| Mth    | Route             | View                 | Description                  |
|--------|-------------------|----------------------|------------------------------|
| GET    | /movies           | movies_index.templ   | List all movies              |
| GET    | /movies/new       | movies_new.templ     | Show create form             |
| POST   | /movies           | n/a                  | Create movie form submission |
| GET    | /movies/{id}      | movies_show.templ    | Show movie details           |
| GET    | /movies/{id}/edit | movies_edit.templ    | Show edit form               |
| PUT    | /movies/{id}      | n/a                  | Handle edit form submission  |
| DELETE | /movies/{id}      | n/a                  | Delete movie                 |

## Optimistic Concurrency

The `version` column is used to implement optimistic concurrency.
When an entity is updated, the `version` is incremented.
The last read version is passed to the update query. If the update does not affect any rows, it means the entity was updated by another process and the controller returns a 409.

## CLI

- Start a new project with `gotr new <project-name>`.
- Run the development server with `gotr dev`. It uses `air` for hot reloading.
- Scaffold a new resource with full CRUD functionality with `gotr generate scaffold <resource-name> <field-name>:<type> [<field-name>:<type> ...]`.
- Run DB migrations with `gotr db migrate`.

