#!/usr/bin/env bash
# Backfill movie posters from TMDB into the directory served by cp-thumbnails.
#
# The client requests thumbnail_<MovieLens id>.jpg, but TMDB needs a TMDB id,
# so this joins db/import/movies.csv to the links.csv mapping. Completed files
# are skipped and downloads are moved into place atomically, so an interrupted
# or rate-limited run can simply be re-run.
set -euo pipefail

if [ -z "${TMDB_API_KEY:-}" ]; then
  printf 'error: TMDB_API_KEY is not set.\n' >&2
  printf 'Export a TMDB v3 API key, for example:\n' >&2
  printf '  TMDB_API_KEY=your-key scripts/backfill-posters.sh\n' >&2
  exit 1
fi

for command_name in curl jq; do
  if ! command -v "$command_name" >/dev/null 2>&1; then
    printf 'error: %s is required but was not found on PATH.\n' "$command_name" >&2
    exit 1
  fi
done

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
repo_root=$(dirname -- "$script_dir")
movies_file="$repo_root/db/import/movies.csv"
links_file="$repo_root/dataset/01-ml-latest-small/links.csv"
out_dir="$repo_root/thumbnails/img"

for required_file in "$movies_file" "$links_file"; do
  if [ ! -r "$required_file" ]; then
    printf 'error: required file not found: %s\n' "$required_file" >&2
    exit 1
  fi
done

mkdir -p "$out_dir"

fetched=0
skipped=0
failed=0

warn() {
  printf 'warning: %s\n' "$*" >&2
}

while IFS=$'\t' read -r movie_id tmdb_id; do
  [ -n "$movie_id" ] || continue

  destination="$out_dir/thumbnail_${movie_id}.jpg"
  if [ -s "$destination" ]; then
    skipped=$((skipped + 1))
    continue
  fi

  if [ -z "$tmdb_id" ]; then
    warn "skip movie ${movie_id}: no tmdbId in links.csv"
    skipped=$((skipped + 1))
    continue
  fi

  details_url="https://api.themoviedb.org/3/movie/${tmdb_id}?api_key=${TMDB_API_KEY}"
  if details=$(curl --silent --show-error --location --fail \
      --retry 3 --retry-delay 2 --retry-connrefused --max-time 30 \
      "$details_url"); then
    :
  else
    warn "movie ${movie_id} (tmdb ${tmdb_id}): metadata request failed"
    failed=$((failed + 1))
    continue
  fi

  if ! poster_path=$(printf '%s' "$details" | jq -r 'if .poster_path == null then empty else .poster_path end' 2>/dev/null); then
    warn "movie ${movie_id} (tmdb ${tmdb_id}): invalid metadata response"
    failed=$((failed + 1))
    continue
  fi

  if [ -z "$poster_path" ]; then
    warn "skip movie ${movie_id}: no poster on TMDB ${tmdb_id}"
    skipped=$((skipped + 1))
    continue
  fi

  # A .part file is never considered complete, so an interrupted download is
  # retried on the next run and the final name only appears after a full body.
  partial="${destination}.part"
  if curl --silent --show-error --location --fail \
      --retry 3 --retry-delay 2 --retry-connrefused --max-time 60 \
      -o "$partial" "https://image.tmdb.org/t/p/w500${poster_path}"; then
    mv -f "$partial" "$destination"
    fetched=$((fetched + 1))
  else
    rm -f "$partial"
    warn "movie ${movie_id} (tmdb ${tmdb_id}): poster download failed"
    failed=$((failed + 1))
  fi
done < <(awk -F',' '
  NR == FNR {
    if (FNR == 1) next
    sub(/\r$/, "", $3)
    tmdb[$1] = $3
    next
  }
  {
    sub(/\r$/, "", $1)
    printf "%s\t%s\n", $1, ($1 in tmdb ? tmdb[$1] : "")
  }
' "$links_file" <(tail -n +2 "$movies_file"))

printf 'Poster backfill complete: %d fetched, %d skipped, %d failed.\n' \
  "$fetched" "$skipped" "$failed"
