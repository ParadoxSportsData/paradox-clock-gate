# ParadoxSportsData — Setup Guide

Runs the full stack: temporal engine, stats, win-probability prediction, and the React UI.

---

## Repos

| Repo | Port |
|------|------|
| [paradox-clock-gate](https://github.com/ParadoxSportsData/paradox-clock-gate) | 8080 |
| [paradox-stats](https://github.com/ParadoxSportsData/paradox-stats) | 8001 |
| [paradox-predict](https://github.com/ParadoxSportsData/paradox-predict) | 8002 |
| [paradox-ui](https://github.com/ParadoxSportsData/paradox-ui) | 5173 |

---

## Prerequisites

- Go 1.26.3+ — [go.dev/dl](https://go.dev/dl/)
- Python 3.14 — for paradox-predict
- Python 3.12 — for paradox-stats
- Node.js 20+
- macOS only: `brew install libomp` (required by XGBoost)

---

## One-time setup

**1. Clone**
```bash
git clone https://github.com/ParadoxSportsData/paradox-clock-gate
git clone https://github.com/ParadoxSportsData/paradox-stats
git clone https://github.com/ParadoxSportsData/paradox-predict
git clone https://github.com/ParadoxSportsData/paradox-ui
```

**2. Download game data** (5–15 min first run; cached after)
```bash
cd paradox-clock-gate
pip install -r scripts/requirements.txt
python3 scripts/fetch_games.py
```

**3. Build clock-gate**
```bash
go build ./cmd/clock-gate/
```

**4. Install Python dependencies**
```bash
# paradox-stats
cd ../paradox-stats && python3.12 -m venv .venv && .venv/bin/pip install -e .

# paradox-predict
cd ../paradox-predict && python3.14 -m venv venv && venv/bin/pip install -e .
```

**5. Install UI dependencies**
```bash
cd ../paradox-ui && npm install
```

---

## Run (4 terminals)

```bash
# Terminal 1 — clock-gate
cd paradox-clock-gate && ./clock-gate serve

# Terminal 2 — stats
cd paradox-stats && .venv/bin/uvicorn main:app --port 8001

# Terminal 3 — prediction
cd paradox-predict && venv/bin/uvicorn main:app --port 8002

# Terminal 4 — UI
cd paradox-ui && npm run dev
```

Open **[http://localhost:5173](http://localhost:5173)**.

---

## Troubleshooting

| Symptom | Fix |
|---------|-----|
| `library not loaded: @rpath/libomp.dylib` | `brew install libomp` |
| `nfl_data_py` install fails | Use Python 3.12 for the fetch script only |
| Stats or Lab panels show empty | Check that services on 8001 and 8002 are running |
