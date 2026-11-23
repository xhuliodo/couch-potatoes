# 0001 — Rating contract: binary user input, decimal seed data

## Status

Accepted

## Context

The graph stores a single `rating` property on the `RATED` relationship, but that property is written by two different producers:

- **Users**, through the rate-movie endpoint. The interface only offers two actions — "loved it" and "not for me" — so the meaningful contract is binary: `1` means liked, `0` means not liked.
- **The dataset import**, through `db/import/ratings.csv`. MovieLens ratings of `0.5`–`5` stars are scaled to `0.1`–`1` decimal values and loaded with `toFloat`. These decimals feed collaborative filtering, where the distance between two users' ratings carries information that a binary value would discard.

The two producers therefore put values with different granularity into the same edge property. Without a written decision, a later change could accept `2` or `4` from a user, store a value that user-based similarity treats as a strong signal while content-based similarity ignores it, or "normalize" the seed data and silently change recommendation output.

## Decision

- User input is binary. The rate-movie service accepts only `0` (`not liked`) and `1` (`liked`) and rejects every other integer with a `bad_request` client error.
- Validation belongs to the rate-movie service, not the HTTP handler. The service is the seam the tests drive, and every caller of the service gets the same guarantee.
- Rejected input persists nothing: validation runs before the user and movie lookups and before any write.
- Seed ratings imported from the dataset remain scaled decimals (`0.1`–`1`) and continue to feed user-based collaborative filtering unchanged.
- Only a rating of exactly `1` counts as a like for content-based similarity. A decimal seed rating such as `0.8` does not put a movie into the liked set.

## Consequences

- The database can contain decimals that no user request can write. This is intentional, not a data-integrity bug: the values come from a trusted import, not from the API.
- Content-based recommendations only grow from explicit `1` ratings. A user who has rated many movies but liked none of them gets no content-based neighbours, even though collaborative filtering can still use their decimal history.
- Adding a finer-grained user rating scale later (for example a five-star control) requires revisiting this ADR, the service validation, content-based similarity, and possibly the import normalization together.
- The service-level test suite pins the contract: invalid values are rejected without writing, both valid values persist, and re-rating a movie updates the existing edge.
