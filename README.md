# Gator

Gator is a command-line RSS feed aggregator written in Go. It stores users,
feeds, follows, and posts in PostgreSQL. You can register users, add RSS feeds,
follow feeds, continuously collect new posts, and browse the posts collected
for the logged-in user.

## Requirements

Install the following before running Gator:

- [Go](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)

The project uses the Go version declared in `go.mod`. PostgreSQL must be
running and you must have a database and credentials that Gator can use.


## Database setup

The SQL migrations in `sql/schema` are Goose migrations. Install Goose if it is
not already available:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Create an empty PostgreSQL database, then apply the migrations from the
repository root. Replace the connection string with your PostgreSQL details:

```sh
goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
```

## Configuration

Gator reads its configuration from `~/.gatorconfig.json`. Create that file and
set `db_url` to the same PostgreSQL connection string used for the migrations:

```json
{
	"db_url": "Your_PostgreSQL_Connection_string",
	"current_user_name": ""
}
```

The `current_user_name` value is updated automatically by `register` and
`login`. Do not commit this file if it contains database credentials.

## Running Gator

From the repository root, run commands with either the installed binary or Go:

```sh
gator <command> [args...]
go run . <command> [args...]
```

Register a user and make that user current:

```sh
gator register alice
```

Useful commands include:

```text
login <name>                  Log in as an existing user
users                         List registered users
addfeed <name> <url>          Add an RSS feed and follow it
feeds                         List all feeds
follow <url>                  Follow an existing feed
following                     List feeds followed by the current user
unfollow <url>                Stop following a feed
agg <duration>                Continuously fetch feeds, e.g. 15s or 1m
browse [limit]               Show collected posts; defaults to 2
reset                         Delete all application data
```

For example, add a feed, start the collector in one terminal, and browse new
posts from another:

```sh
gator addfeed golang "https://blog.golang.org/feed.atom"
gator agg 1m
gator browse 10
```

`addfeed`, `follow`, `unfollow`, `following`, and `browse` require a current
user. Stop the long-running `agg` command with `Ctrl-C`.
