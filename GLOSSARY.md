# Glossary

Domain vocabulary used by the code, the API, and the interface. The rating
definitions follow [ADR-0001](docs/adr/0001-rating-contract.md).

## Domain nouns

**Movie** — A film in the graph, identified by its MovieLens `movieId`. It has a
title, release year, IMDb link, genres, and the people involved.

**Genre** — A category such as Action, Comedy, or Drama. A movie belongs to one
or more genres; onboarding saves at least three genres as the user's favourites.

**Cast / person** — An actor, director, or writer connected to a movie. The graph
stores all three under the `Cast` label with `ACTED_IN`, `DIRECTED`, or `WROTE`
relationships. Content-based similarity looks at this shared set of people.

**Rating** — The value on the `RATED` relationship between a user and a movie.
It is numeric but has two producers with different granularity:

- User input is binary: `1` means **liked**, `0` means **not liked**. The
  rate-movie service accepts only those two values and rejects anything else
  with `bad_request`.
- Seed ratings from the dataset are scaled decimals (`0.1`–`1`) and feed
  user-based collaborative filtering.

Only a rating of exactly `1` counts as a like for content-based similarity; a
decimal seed rating does not enter the liked set.

**Liked / Not liked** — The two outcomes of user rating. In the interface they
are **Loved it** (writes `1`) and **Not for me** (writes `0`).

**Watchlist** — Movies a user saved for later, stored as `WATCH_LATER`
relationships and shown under **Watch later** → **Movies waiting for you**.

**Watchlist history** — Watchlist movies the user has since rated, shown under
**Watch later** → **Your collection**. It is separate from the unrated
watchlist.

**Onboarding / Setup** — The two-step first-run flow: choose at least three
favourite genres, then rate at least fifteen movies. The API reports the current
step as `1`, `2`, or `3` (finished).

**Popular recommendations** — Movies with high ratings inside the user's
favourite genres. The onboarding swipe screen uses these, so a new user always
has something to rate.

**Content-based recommendation** — A recommendation scored by how much cast,
director, and writer overlap a movie has with movies the user rated `1`.
Surfaced in the **Our suggestions for you** tab.

**User-based recommendation** — A recommendation scored from the ratings of the
users most similar to the current user. Surfaced in the **Other people also
liked** tab.

**Nearest-neighbor circle** — The fixed number of most similar users (currently
25) that user-based recommendation keeps after ranking by Pearson correlation.
Everyone outside the circle is discarded before scores are calculated.

**Pearson correlation** — The similarity measure between two users' ratings on
movies they both rated. A higher coefficient means more similar taste; users
with too few ratings in common are excluded before ranking.

**Seed data** — The MovieLens small dataset imported into Neo4j on first start:
movies, genres, people, users, and decimal ratings. It is imported only when the
graph is empty and persists in `db/data`.

**Thumbnail / poster** — The image for a movie, served by `cp-thumbnails` at
`thumbnail_<movieId>.jpg`. Posters are optional and filled by
`scripts/backfill-posters.sh`; without them the interface renders a styled
posterless card.
