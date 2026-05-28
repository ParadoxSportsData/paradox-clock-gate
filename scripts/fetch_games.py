"""
fetch_games.py — Download NFL play-by-play data and convert to clock-gate JSON format.

Downloads from nflverse via nfl_data_py, transforms each game into the format
clock-gate expects, and writes one JSON file per game to data/raw/.

Usage:
    python3 scripts/fetch_games.py                        # 2011–2025 (all seasons)
    python3 scripts/fetch_games.py --seasons 2024 2024    # single season
    python3 scripts/fetch_games.py --seasons 2020 2025    # range

Requirements:
    Python 3.12 (nfl_data_py pins pandas < 2.0, incompatible with 3.13+)
    pip install -r scripts/requirements.txt
"""

import argparse
import json
import logging
import sys
from datetime import datetime
from pathlib import Path

import nfl_data_py as nfl
import pandas as pd

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(message)s",
    handlers=[logging.StreamHandler(sys.stdout)],
)
logger = logging.getLogger(__name__)


def transform_game(game_df: pd.DataFrame) -> dict | None:
    """Convert one game's nflverse DataFrame to clock-gate JSON format."""
    try:
        row0 = game_df.iloc[0]

        game_id   = str(row0["game_id"])
        season    = int(row0["season"])
        week      = int(row0["week"])
        home_team = str(row0["home_team"])
        away_team = str(row0["away_team"])

        game_date_str = str(row0.get("game_date", ""))
        try:
            game_date = datetime.fromisoformat(game_date_str.replace("Z", "+00:00"))
        except (ValueError, AttributeError):
            game_date = datetime(season, 9, 1, 13, 0, 0)

        scores = game_df["total_home_score"].dropna()
        home_score = int(scores.iloc[-1]) if not scores.empty else None
        scores = game_df["total_away_score"].dropna()
        away_score = int(scores.iloc[-1]) if not scores.empty else None

        plays = []
        for _, r in game_df.iterrows():
            if pd.isna(r.get("play_id")) or pd.isna(r.get("game_seconds_remaining")):
                continue

            seconds_remaining = float(r["game_seconds_remaining"])
            quarter = int(r["qtr"]) if not pd.isna(r.get("qtr")) else 1

            if quarter <= 4:
                game_clock_total_seconds = int(3600 - seconds_remaining)
            else:
                ot_periods = quarter - 4
                ot_start = 3600 + ((ot_periods - 1) * 900)
                game_clock_total_seconds = ot_start + int(900 - seconds_remaining)

            mins = int(seconds_remaining // 60)
            secs = int(seconds_remaining % 60)

            plays.append({
                "play_id":                 int(r["play_id"]),
                "quarter":                 quarter,
                "game_clock":              f"{mins}:{secs:02d}",
                "game_clock_total_seconds": game_clock_total_seconds,
                "down":        int(r["down"])        if not pd.isna(r.get("down"))        else None,
                "ydstogo":     int(r["ydstogo"])     if not pd.isna(r.get("ydstogo"))     else None,
                "yardline_100": int(r["yardline_100"]) if not pd.isna(r.get("yardline_100")) else None,
                "play_type":   str(r["play_type"])   if not pd.isna(r.get("play_type"))   else None,
                "yards_gained": int(r["yards_gained"]) if not pd.isna(r.get("yards_gained")) else None,
                "description": str(r["desc"])        if not pd.isna(r.get("desc"))        else "",
                "posteam":     str(r["posteam"])     if not pd.isna(r.get("posteam"))     else None,
                "defteam":     str(r["defteam"])     if not pd.isna(r.get("defteam"))     else None,
                "posteam_score": int(r["posteam_score"]) if not pd.isna(r.get("posteam_score")) else None,
                "defteam_score": int(r["defteam_score"]) if not pd.isna(r.get("defteam_score")) else None,
                "wp":          float(r["wp"])         if not pd.isna(r.get("wp"))          else None,
                "epa":         float(r["epa"])        if not pd.isna(r.get("epa"))         else None,
            })

        if not plays:
            logger.warning(f"  no valid plays for {game_id}, skipping")
            return None

        return {
            "game_id":   game_id,
            "season":    season,
            "week":      week,
            "game_date": game_date.isoformat(),
            "home_team": home_team,
            "away_team": away_team,
            "home_score": home_score,
            "away_score": away_score,
            "plays":     plays,
        }

    except Exception as e:
        logger.error(f"  transform failed: {e}")
        return None


def fetch_season(season: int, output_dir: Path) -> tuple[int, int]:
    """Download one season and write JSON files. Returns (saved, failed)."""
    logger.info(f"Downloading {season} season from nflverse...")
    try:
        pbp = nfl.import_pbp_data(years=[season])
    except Exception as e:
        logger.error(f"Download failed for {season}: {e}")
        return 0, 0

    game_ids = pbp["game_id"].unique()
    logger.info(f"  {len(game_ids)} games found")

    saved, failed = 0, 0
    for gid in game_ids:
        game_df = pbp[pbp["game_id"] == gid].sort_values("play_id")
        data = transform_game(game_df)
        if data is None:
            failed += 1
            continue
        out = output_dir / f"{gid}.json"
        out.write_text(json.dumps(data, indent=2, ensure_ascii=False), encoding="utf-8")
        saved += 1

    logger.info(f"  saved {saved}, failed {failed}")
    return saved, failed


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Download NFL play-by-play data for clock-gate serve mode.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__,
    )
    parser.add_argument(
        "--seasons", nargs=2, type=int, metavar=("START", "END"),
        default=[2011, 2025],
        help="First and last season to download (inclusive). Default: 2011 2025",
    )
    parser.add_argument(
        "--output", default="data/raw",
        help="Directory to write JSON files. Default: data/raw",
    )
    args = parser.parse_args()

    start, end = args.seasons
    if start > end:
        print("error: START must be <= END", file=sys.stderr)
        sys.exit(1)

    output_dir = Path(args.output)
    output_dir.mkdir(parents=True, exist_ok=True)

    total_saved = total_failed = 0
    for season in range(start, end + 1):
        s, f = fetch_season(season, output_dir)
        total_saved  += s
        total_failed += f

    print(f"\nDone. {total_saved} games saved to {output_dir}/")
    if total_failed:
        print(f"       {total_failed} games failed (see log above)")


if __name__ == "__main__":
    main()
