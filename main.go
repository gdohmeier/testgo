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
  <meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover" />
  <title>testgo — Snake</title>
  <style>
    * { box-sizing: border-box; }
    html, body { height: 100%; }
    body {
      margin: 0;
      min-height: 100dvh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 8px;
      padding: max(12px, env(safe-area-inset-top))
               max(12px, env(safe-area-inset-right))
               max(12px, env(safe-area-inset-bottom))
               max(12px, env(safe-area-inset-left));
      font-family: system-ui, sans-serif;
      background: #111;
      color: #eee;
      touch-action: manipulation;
    }
    h1 { margin: 0; font-size: clamp(1.1rem, 4vw, 1.4rem); }
    p { margin: 0; color: #aaa; font-size: clamp(0.75rem, 3vw, 0.9rem); text-align: center; }
    canvas {
      width: min(92vw, 72vh, 520px);
      height: auto;
      aspect-ratio: 1;
      background: #1a1a1a;
      border: 2px solid #3a3;
      image-rendering: pixelated;
      touch-action: none;
      max-width: 100%;
    }
    .bar { font-variant-numeric: tabular-nums; }
    button {
      min-height: 44px;
      min-width: 44px;
      padding: 8px 16px;
      background: #2a2;
      color: #111;
      border: 0;
      border-radius: 4px;
      font-weight: 700;
    }
    .pad {
      display: none;
      grid-template-columns: repeat(3, 56px);
      grid-template-rows: repeat(3, 56px);
      gap: 8px;
      justify-content: center;
      user-select: none;
      -webkit-user-select: none;
    }
    .pad button { font-size: 1.2rem; background: #333; color: #eee; }
    .pad .up    { grid-column: 2; grid-row: 1; }
    .pad .left  { grid-column: 1; grid-row: 2; }
    .pad .right { grid-column: 3; grid-row: 2; }
    .pad .down  { grid-column: 2; grid-row: 3; }
    @media (hover: none), (pointer: coarse) {
      .pad { display: grid; }
      .hint-keys { display: none; }
    }
  </style>
</head>
<body>
  <h1>testgo snake</h1>
  <p class="hint-keys">Arrow keys or WASD. Eat food. Don't hit walls or yourself.</p>
  <p class="hint-touch" hidden>Swipe or use the pad. Don't hit walls or yourself.</p>
  <canvas id="c" width="400" height="400"></canvas>
  <div class="bar">Score: <span id="score">0</span></div>
  <button id="restart" type="button">Restart</button>
  <div class="pad" aria-label="Direction pad">
    <button type="button" class="up" data-dir="up">▲</button>
    <button type="button" class="left" data-dir="left">◀</button>
    <button type="button" class="right" data-dir="right">▶</button>
    <button type="button" class="down" data-dir="down">▼</button>
  </div>
  <script>
    const SIZE = 20, CELLS = 20;
    const canvas = document.getElementById("c");
    const ctx = canvas.getContext("2d");
    const scoreEl = document.getElementById("score");
    const dirs = {
      up: { x: 0, y: -1 }, down: { x: 0, y: 1 },
      left: { x: -1, y: 0 }, right: { x: 1, y: 0 },
    };
    const keymap = {
      ArrowUp: "up", ArrowDown: "down", ArrowLeft: "left", ArrowRight: "right",
      w: "up", s: "down", a: "left", d: "right",
      W: "up", S: "down", A: "left", D: "right",
    };

    let snake, dir, nextDir, food, score, alive, timer;

    if (window.matchMedia("(hover: none), (pointer: coarse)").matches) {
      document.querySelector(".hint-touch").hidden = false;
    }

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

    function setDir(nd) {
      if (!nd) return;
      if (nd.x === -dir.x && nd.y === -dir.y) return;
      nextDir = nd;
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

    window.addEventListener("keydown", (e) => {
      const name = keymap[e.key];
      if (!name) return;
      e.preventDefault();
      setDir(dirs[name]);
    });

    document.querySelector(".pad").addEventListener("pointerdown", (e) => {
      const btn = e.target.closest("[data-dir]");
      if (!btn) return;
      e.preventDefault();
      setDir(dirs[btn.dataset.dir]);
    });

    let touchStart = null;
    canvas.addEventListener("pointerdown", (e) => {
      touchStart = { x: e.clientX, y: e.clientY };
    });
    canvas.addEventListener("pointerup", (e) => {
      if (!touchStart) return;
      const dx = e.clientX - touchStart.x;
      const dy = e.clientY - touchStart.y;
      touchStart = null;
      if (Math.abs(dx) < 24 && Math.abs(dy) < 24) return;
      if (Math.abs(dx) > Math.abs(dy)) setDir(dx > 0 ? dirs.right : dirs.left);
      else setDir(dy > 0 ? dirs.down : dirs.up);
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
