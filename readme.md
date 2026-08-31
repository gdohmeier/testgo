Here’s a complete, first-try setup: a small Go HTTP server that serves a playable Snake game and returns **200 OK** on `/up`, plus a Dockerfile and the GitHub / Docker Hub commands.

### Project layout

```
testgo/
  main.go
  go.mod
  Dockerfile
  .dockerignore
```

---

### `go.mod`

```go
module testgo

go 1.23
```

---

### `main.go`

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const gameHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>testgo — Snake</title>
  <style>
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      font-family: system-ui, sans-serif;
      background: #111;
      color: #eee;
    }
    h1 { margin: 0 0 8px; font-size: 1.4rem; }
    p { margin: 0 0 12px; color: #aaa; font-size: 0.9rem; }
    canvas {
      background: #1a1a1a;
      border: 2px solid #3a3;
      image-rendering: pixelated;
    }
    .bar { margin-top: 10px; font-variant-numeric: tabular-nums; }
    button {
      margin-top: 10px;
      padding: 8px 16px;
      background: #2a2;
      color: #111;
      border: 0;
      border-radius: 4px;
      font-weight: 700;
      cursor: pointer;
    }
  </style>
</head>
<body>
  <h1>testgo snake</h1>
  <p>Arrow keys or WASD. Eat food. Don't hit walls or yourself.</p>
  <canvas id="c" width="400" height="400"></canvas>
  <div class="bar">Score: <span id="score">0</span></div>
  <button id="restart" type="button">Restart</button>
  <script>
    const SIZE = 20, CELLS = 20;
    const canvas = document.getElementById("c");
    const ctx = canvas.getContext("2d");
    const scoreEl = document.getElementById("score");

    let snake, dir, nextDir, food, score, alive, timer;

    function randCell() {
      return Math.floor(Math.random() * CELLS);
    }

    function placeFood() {
      let x, y, ok;
      do {
        x = randCell(); y = randCell();
        ok = !snake.some(s => s.x === x && s.y === y);
      } while (!ok);
      food = { x, y };
    }

    function reset() {
      snake = [{ x: 10, y: 10 }, { x: 9, y: 10 }, { x: 8, y: 10 }];
      dir = nextDir = { x: 1, y: 0 };
      score = 0;
      alive = true;
      scoreEl.textContent = "0";
      placeFood();
      if (timer) clearInterval(timer);
      timer = setInterval(tick, 120);
      draw();
    }

    function tick() {
      if (!alive) return;
      dir = nextDir;
      const head = { x: snake[0].x + dir.x, y: snake[0].y + dir.y };
      if (head.x < 0 || head.y < 0 || head.x >= CELLS || head.y >= CELLS) {
        alive = false; draw(); return;
      }
      if (snake.some(s => s.x === head.x && s.y === head.y)) {
        alive = false; draw(); return;
      }
      snake.unshift(head);
      if (head.x === food.x && head.y === food.y) {
        score += 1;
        scoreEl.textContent = String(score);
        placeFood();
      } else {
        snake.pop();
      }
      draw();
    }

    function draw() {
      ctx.fillStyle = "#1a1a1a";
      ctx.fillRect(0, 0, canvas.width, canvas.height);
      ctx.fillStyle = "#c33";
      ctx.fillRect(food.x * SIZE, food.y * SIZE, SIZE - 1, SIZE - 1);
      snake.forEach((s, i) => {
        ctx.fillStyle = i === 0 ? "#6f6" : "#3a3";
        ctx.fillRect(s.x * SIZE, s.y * SIZE, SIZE - 1, SIZE - 1);
      });
      if (!alive) {
        ctx.fillStyle = "rgba(0,0,0,0.55)";
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = "#fff";
        ctx.font = "24px system-ui";
        ctx.textAlign = "center";
        ctx.fillText("Game over", canvas.width / 2, canvas.height / 2);
      }
    }

    const keymap = {
      ArrowUp: { x: 0, y: -1 }, ArrowDown: { x: 0, y: 1 },
      ArrowLeft: { x: -1, y: 0 }, ArrowRight: { x: 1, y: 0 },
      w: { x: 0, y: -1 }, s: { x: 0, y: 1 },
      a: { x: -1, y: 0 }, d: { x: 1, y: 0 },
      W: { x: 0, y: -1 }, S: { x: 0, y: 1 },
      A: { x: -1, y: 0 }, D: { x: 1, y: 0 },
    };
    window.addEventListener("keydown", (e) => {
      const nd = keymap[e.key];
      if (!nd) return;
      e.preventDefault();
      if (nd.x === -dir.x && nd.y === -dir.y) return;
      nextDir = nd;
    });
    document.getElementById("restart").addEventListener("click", reset);
    reset();
  </script>
</body>
</html>
`

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(gameHTML))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Println("testgo listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
```

---

### `Dockerfile`

```dockerfile
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/testgo .

FROM alpine:3.20
RUN adduser -D -H -u 10001 app
USER app
WORKDIR /home/app
COPY --from=build /out/testgo /usr/local/bin/testgo
EXPOSE 8080
ENV PORT=8080
HEALTHCHECK --interval=15s --timeout=3s --retries=3 CMD wget -qO- http://127.0.0.1:8080/up || exit 1
CMD ["testgo"]
```

Alpine’s `wget` is available by default in many tags; to be extra safe on first try, use this slightly more portable runtime if `wget` is missing in your base:

```dockerfile
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/testgo .

FROM alpine:3.20
RUN apk add --no-cache wget && adduser -D -H -u 10001 app
USER app
WORKDIR /home/app
COPY --from=build /out/testgo /usr/local/bin/testgo
EXPOSE 8080
ENV PORT=8080
HEALTHCHECK --interval=15s --timeout=3s --retries=3 CMD wget -qO- http://127.0.0.1:8080/up || exit 1
CMD ["testgo"]
```

Use the second Dockerfile (the one with `apk add wget`) so the healthcheck works on first try.

---

### `.dockerignore`

```
.git
.gitignore
README.md
*.md
```

---

### Run locally (no Docker)

```bash
cd testgo
go run .
```

Open http://localhost:8080 — play Snake.  
Check health: `curl -i http://localhost:8080/up` → `200 OK`.

### Run with Docker

```bash
docker build -t testgo:local .
docker run --rm -p 8080:8080 testgo:local
```

---

### Push the repo to GitHub

Replace `YOUR_GITHUB_USER` with your GitHub username. Create an empty repo named `testgo` on GitHub first (no README).

```bash
cd testgo
git init
git add main.go go.mod Dockerfile .dockerignore
git commit -m "Initial commit: testgo snake web app"
git branch -M main
git remote add origin https://github.com/YOUR_GITHUB_USER/testgo.git
git push -u origin main
```

SSH variant:

```bash
git remote add origin git@github.com:YOUR_GITHUB_USER/testgo.git
git push -u origin main
```

---

### Push the image to Docker Hub

Replace `YOUR_DOCKERHUB_USER` with your Docker Hub username.

```bash
docker login
docker build -t YOUR_DOCKERHUB_USER/testgo:latest .
docker push YOUR_DOCKERHUB_USER/testgo:latest
```

Then others can run:

```bash
docker pull YOUR_DOCKERHUB_USER/testgo:latest
docker run --rm -p 8080:8080 YOUR_DOCKERHUB_USER/testgo:latest
```

No extra Go modules, no frontend build step — `go build` and `docker build` are enough for a clean first run.
