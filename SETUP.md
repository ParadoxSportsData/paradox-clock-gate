# ParadoxSportsData — Full Stack Setup

This guide brings up all five services end-to-end: the Go temporal engine, the Python stats and prediction services, and the React UI.

---

## Repos

| Service | Repo | Port |
|---------|------|------|
| [paradox-clock-gate](https://github.com/ParadoxSportsData/paradox-clock-gate) | Go HTTP server + CLI | 8080 |
| [paradox-stats](https://github.com/ParadoxSportsData/paradox-stats) | Python box-score service | 8001 |
| [paradox-predict](https://github.com/ParadoxSportsData/paradox-predict) | Python WP prediction service | 8002 |
| [paradox-ui](https://github.com/ParadoxSportsData/paradox-ui) | React + TypeScript frontend | 5173 |

---

## Prerequisites

| Tool | Required | Check |
|------|----------|-------|
| Go 1.26.3+ | paradox-clock-gate | `go version` |
| Python 3.14 | paradox-predict | `python3.14 --version` |
| Python 3.12 | paradox-stats | `python3.12 --version` |
| Node.js 20+ | paradox-ui | `node --version` |
| libomp (macOS) | paradox-predict (XGBoost) | `brew install libomp` |

---

## Step 1 — Clone all repos

```bash
git clone https://github.com/ParadoxSportsData/paradox-clock-gate
git clone https://github.com/ParadoxSportsData/paradox-stats
git clone https://github.com/ParadoxSportsData/paradox-predict
git clone https://github.com/ParadoxSportsData/paradox-ui
```

---

## Step 2 — Download NFL game data

Game data is fetched from nflverse and is not included in the repo (~4,100 files, ~400MB for all seasons).

```bash
cd paradox-clock-gate

# Install Python dependencies (Python 3.12 required for the fetch script)
pip install -r scripts/requirements.txt

# Download all seasons 2011–2025 (5–15 min on first run; cached on subsequent runs)
python3 scripts/fetch_games.py

# Or a single season for a quick start
python3 scripts/fetch_games.py --seasons 2024 2024
```

This writes one JSON file per game to `data/raw/`.

---

## Step 3 — Build and start paradox-clock-gate (port 8080)

```bash
cd paradox-clock-gate
go build ./cmd/clock-gate/
./clock-gate serve
```

Verify:
```bash
curl http://localhost:8080/health
# → {"status":"ok"}

curl http://localhost:8080/games | head -c 200
# → JSON array of game objects
```

---

## Step 4 — Start paradox-stats (port 8001)

```bash
cd paradox-stats
python3.12 -m venv .venv
.venv/bin/pip install -e .
.venv/bin/uvicorn main:app --port 8001
```

Verify:
```bash
curl http://localhost:8001/health
# → {"status":"ok","service":"paradox-stats"}
```

---

## Step 5 — Start paradox-predict (port 8002)

The trained model (`ml/model.pkl`) is included in the repo — no training step required.

```bash
# macOS only: install OpenMP runtime required by XGBoost
brew install libomp

cd paradox-predict
python3.14 -m venv venv
venv/bin/pip install -e .
venv/bin/uvicorn main:app --port 8002
```

Verify:
```bash
curl http://localhost:8002/health
# → {"status":"ok","service":"paradox-predict","model_loaded":true}

# 4th & 1 from the 1, down 4, 10s left, Q4 — expect ~15% WP
curl -s -X POST http://localhost:8002/predict/scenario \
  -H "Content-Type: application/json" \
  -d '{"down":4,"distance":1,"yardline_100":1,"quarter":4,"seconds_remaining_quarter":10,"score_differential":-4,"is_home_possession":true,"era_season":2024,"era_week":1}'
```

---

## Step 6 — Start paradox-ui (port 5173)

```bash
cd paradox-ui
npm install
npm run dev
```

Open [http://localhost:5173](http://localhost:5173).

Select a game → scrub the timeline → stats panels and win probability update live.

---

## Running all services at once

Each service runs in its own terminal. Open four tabs:

```bash
# Tab 1
cd paradox-clock-gate && ./clock-gate serve

# Tab 2
cd paradox-stats && .venv/bin/uvicorn main:app --port 8001

# Tab 3
cd paradox-predict && venv/bin/uvicorn main:app --port 8002

# Tab 4
cd paradox-ui && npm run dev
```

---

## UI without a backend (mock mode)

paradox-ui ships with fixture data for offline development:

```bash
cd paradox-ui
VITE_MOCK_MODE=true npm run dev
```

---

## Troubleshooting

**`ml/model.pkl: No such file or directory` in paradox-predict**
The model is committed to the repo — run `git pull` to get it.

**`library not loaded: @rpath/libomp.dylib` (macOS)**
Run `brew install libomp` and restart the uvicorn process.

**`nfl_data_py` install fails on Python 3.13+**
The fetch script requires Python 3.12. Use `python3.12 -m pip install -r scripts/requirements.txt`. paradox-predict uses 3.14 but does not depend on nfl_data_py.

**`go: command not found`**
Install Go 1.26.3 from [go.dev/dl](https://go.dev/dl/).

**Clock-gate returns `tick out of range`**
The requested tick exceeds the game's last play. Use `--list data/raw/` to see available games and query within a valid range.
