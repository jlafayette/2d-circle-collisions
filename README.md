# collisions
Test project to practice implementing 2d collisions with the [ebiten game library](https://ebiten.org/).

Inspired by One Lone Coder's video: [Programming Balls #1 Circle Vs Circle Collisions C++](https://www.youtube.com/watch?v=LPzyNOHY3A4).

## Run Locally

```sh
go run main.go
```

## Run Locally in WebBrowser

```sh
go get github.com/hajimehoshi/wasmserve
wasmserve .
```

Open `http://localhost:8080/` on your browser.

## Rebuild shaders

```sh
go generate resources/shader/generate.go
```
