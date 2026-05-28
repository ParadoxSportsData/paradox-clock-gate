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

Install these before starting. Each has a verify command to confirm it's ready.

| Dependency | Required for | Install | Verify |
|---|---|---|---|
| Go 1.26.3+ | paradox-clock-gate | [go.dev/dl](https://go.dev/dl/) | `go version` |
| Python 3.12 | paradox-stats, game data fetch | [python.org/downloads](https://www.python.org/downloads/) | `python3.12 --version` |
| Python 3.14 | paradox-predict | [python.org/downloads](https://www.python.org/downloads/) | `python3.14 --version` |
| Node.js 20+ | paradox-ui | [nodejs.org](https://nodejs.org/) | `node --version` |
| libomp *(macOS only)* | paradox-predict (XGBoost) | `brew install libomp` | — |

> **macOS:** Install libomp before setting up paradox-predict. XGBoost will install without it but fail at runtime.

> **Python versions:** Both 3.12 and 3.14 are required. 3.12 is used for paradox-stats and the game data fetch script. 3.14 is used for paradox-predict. Using the wrong version will cause install failures.

---

## One-time setup

**1. Clone all repos into a common folder**
```bash
mkdir paradox && cd paradox
git clone https://github.com/ParadoxSportsData/paradox-clock-gate
git clone https://github.com/ParadoxSportsData/paradox-stats
git clone https://github.com/ParadoxSportsData/paradox-predict
git clone https://github.com/ParadoxSportsData/paradox-ui
```

**2. Download game data** (5–15 min first run; cached on subsequent runs)
```bash
cd paradox-clock-gate
python3.12 -m pip install -r scripts/requirements.txt
python3.12 scripts/fetch_games.py
cd ..
```

**3. Build clock-gate**
```bash
cd paradox-clock-gate
go build ./cmd/clock-gate/
cd ..
```

**4. Install Python dependencies**
```bash
# paradox-stats
cd paradox-stats && python3.12 -m venv .venv && .venv/bin/pip install -e . && cd ..

# paradox-predict (model.pkl is already included — no training required)
cd paradox-predict && python3.14 -m venv venv && venv/bin/pip install -e . && cd ..
```

**5. Install UI dependencies**
```bash
cd paradox-ui && npm install && cd ..
```

---

## Run (4 terminals)

Open four terminal tabs, each starting from inside the `paradox/` folder.

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
| `library not loaded: @rpath/libomp.dylib` | `brew install libomp`, then restart Terminal 3 |
| `nfl_data_py` install fails | Use `python3.12 -m pip` — this package does not support Python 3.13+ |
| Stats or Lab panels show empty | Confirm services on ports 8001 and 8002 are running |
| Port already in use | Another process is on 8080, 8001, 8002, or 5173 — stop it or restart the machine |
| `go build` fails with version error | Upgrade Go to 1.26.3+ from go.dev/dl |
| `python3.12` or `python3.14` not found | Install the specific version from python.org/downloads |
