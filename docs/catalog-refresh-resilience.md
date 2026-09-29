# Catalog refresh resilience

A failed upstream refresh does not revoke the last signed artifact. The candidate
reader accepts retained artifacts from a freshly authenticated directory, with
unchanged identity, size, hash, trust-root, expiration and rollback checks. ROMD
must ship the updated reader together with its backend/UI availability change.
Public health continues to describe the latest acquisition attempt; no schema or
trust-root migration is required. Missing artifacts remain unavailable.

No-Intro makes up to three attempts for `prepare_failed`, `request_failed` and
`incomplete_response` within one two-minute deadline. Fresh anonymous sessions
use the same complete selection policy and shared five-second admission gate;
retry backoffs are five and ten seconds. HTTP 429/503 retain their Retry-After
cooldown. Queued exports, changed forms and invalid documents are not retried.
Response diagnostics contain only attempt, fixed stage/page classification,
status, byte count, SHA-256 and fixed outcome code. Cookies, form values, response
text and redirect URLs are never logged. The workflow warns about partial
acquisition failures even when signing and publication succeed.

## N64 and SNES qualification, September 29, 2026 UTC

The public N64 form omits `inc_nodump`. Only N64 now permits that control to be
absent, selecting inclusion whenever it is present and rejecting exclusion-only
controls. All other reviewed selections and DAT validation remain enforced.
The live adapter acquired BigEndian version `20260928-175532`: 1,263 entries and
files, no nodump records, SHA-256
`18f6eb0181a137823b5eadad9996f52ee9ab705d03d5bcf38c1d5c1518c42dd6`.
Compared with the signed retained version `20260918-065107` (1,262 records, no
nodump), all prior SHA-1 identities remain and one is added. This validates the
observed form transition, not undisclosed upstream database content.

SNES acquisition returned version `20260928-123938`, 4,362 entries / 4,363 files,
SHA-256 `1ccc7450a5dc472811597941fc48621d2255c4633eb1c7899d80492c8a754142`.
Both acquisitions passed the existing identity, structure and count checks.
These candidates were written to temporary local state only, not published.
The new reader also authenticated and downloaded the currently failed-health
N64 retained artifact from production using the existing trust root.

## Rollout

Update both the data repository reusable-workflow and tooling SHA pins to the
reviewed publisher revision. Update ROMD's Docker reader source checksum and SHA
to that same revision. Deploy the matching ROMD backend/admin client; then run
the data publisher and verify the signed directory. Do not reset retained trust
caches, bootstrap the site, or bypass an expired signature to perform rollout.
