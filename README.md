# blogAggregator


(gator) A multi-user command line tool for aggregating RSS feeds and viewing the posts.

## Installation

Make sure you have the latest [Go toolchain](https://golang.org/dl/) installed as well as a local Postgres database. You can then install `gator` with:

```bash
go install ...
```

## Config

Create a `.gatorconfig.json` file in your home directory with the following structure:

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the values with your database connection string.

## Usage

Create a new user:

```bash
gator register <name>
```

Add a feed:

```bash
gator addfeed <url>
```

Start the aggregator:

```bash
gator agg 30s
```

View the posts:

```bash
gator browse [limit]
```

There are a few other commands you'll need as well:
# Command | Description
1. register <username> | Registers a new user
2. login <username> | Logs into an existing user
3. reset | Removes all users
4. users | Lists all registered users
5. addfeed <title> <url> | Adds a new RSS feed
6. feeds | Lists all feeds and who added them
7. follow <feed_url> | Follow a feed
8. unfollow <feed_url> | Unfollow a feed
9. following | Lists feeds followed by the current user
10. agg <duration> | Aggregates new posts from feeds (e.g. 1m, 1h)
11. browse [limit] | Displays posts from followed feeds(limit is 10 at a time)

![image](https://github.com/user-attachments/assets/cd8b69d9-70c5-419a-90d9-bf0882f6b0e4)
![image](https://github.com/user-attachments/assets/da07ef08-dc04-43ec-bffe-b30afcd2304b)



## Limitations 
Gator expects valid RSS/XML feeds. Some feeds may be malformed or use non-standard formats and may not work properly.

Use direct RSS URLs, not general blog URLs, for the best results.

Aggregation fetches are limited to feeds that return XML content via HTTP.


