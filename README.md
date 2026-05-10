# godis

A Redis-compatible in-memory data store written in Go from scratch. Implements the [RESP2](https://redis.io/docs/latest/develop/reference/protocol-spec/) for compatibility with existing Redis clients.


## Deployment (Docker)

**Build the image:**
```bash
docker build -t godis .
```

**Run:**
```bash
docker run -p 6379:6379 godis
```

godis listens on port `6379` by default. You can connect with any Redis client:

```bash
redis-cli ping
# PONG

redis-cli echo "hello"
# "hello"
```
