# tanky

Two player Battle City in the terminal. Go, tcell, a tiny TCP relay server and room codes. Needs Go 1.25.

## Run

Start a server somewhere both players can reach:

    go run ./cmd/server

It prints the addresses people on your network can connect to. `-addr` changes the port, default `:8080`.

Then each player runs the client:

    go run ./cmd/client

Host gets a room code, joiner types it in. Skip the menu if you want:

    go run ./cmd/client host
    go run ./cmd/client join ABC123

`-server host:port` if the server isn't on localhost.

## Play

WASD or arrows to move, space to shoot, Esc to quit.
