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
