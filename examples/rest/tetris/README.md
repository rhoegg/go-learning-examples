# Simple battleground.ai bot
Manual strategy for O, I, T, L, and J pieces. No support for S, Z pieces.

# Running
## Arguments
`<program-name> [switches] <API Endpoint>`
- `API Endpoint` the URL of the API endpoint on https://battleground.ai
### Switches
- `--slow` Makes the bot run slowly, good for watching the game run live
## Using go
In the `tetris` directory:

`go run ./... [switches] <API Endpoint>`
## Using docker
In the `tetris` directory:

`docker build -t tetris-bot .`

`docker run --rm tetris-bot [switches] <API Endpoint>`
