# No-Intro SNES acquisition and publication gate

Investigated 2026-09-06 UTC. Scope: `no-intro/snes/standard`, stable system
`snes`, Datomatic system `49`, expected header
`Nintendo - Super Nintendo Entertainment System`. This continues the earlier
SNES investigation. No ROMD application checkout was changed.

**Status: the approved SNES mirror is live through signed Git/Pages/RSS publication
as of 2026-09-07 UTC.** The public DAT selection is reconciled and
redistribution/retention is operator-approved. Synthetic tests, real production
evidence, and application acceptance are distinguished below. The
[September 18 closeout](catalog-expansion-plan.md) records current healthy
publication, a subsequent real document revision, and the separate ROMD
acceptance report. Statements below about unobserved changes and unverified
application acceptance describe the original September 7 evidence.

## Official acquisition

The [official guide](https://wiki.no-intro.org/index.php?title=Database_Navigation_Guide)
describes per-system and daily downloads. The selected
[Standard DAT form](https://datomatic.no-intro.org/index.php?page=download&op=dat&s=49)
works with ordinary HTTPS and an anonymous, acquisition-local cookie session:

1. GET the per-system Standard DAT form and check its selected system/name and
   reviewed inclusion controls.
2. POST its current Prepare control and the explicit complete-public selection.
3. Accept only the returned same-origin `/index.php?page=manager&s=49&download=…`
   redirect, then GET that manager page.
4. POST the visible Download control and validate the bounded ZIP/XML response.

This is a public form flow, **not a documented stable API**. The manager currently
contains a hidden submit and a separate visible submit with the same label.
Choosing the first matching label is unreliable. The adapter selects only the
reviewed visible form and rejects ambiguous or changed controls. It never logs
session cookies/download controls, persists them, logs in, runs a browser,
executes JavaScript, solves CAPTCHAs, or substitutes another mirror.
An authentication/challenge or form change is a failed acquisition requiring
review, not an invitation to bypass it. No upstream messages were sent.

The current [robots guidance](https://datomatic.no-intro.org/robots.txt) has an
empty Disallow, a five-second crawl delay, and one request per five seconds.
The adapter serializes requests with at least five seconds between starts,
30-second request timeouts, a two-minute total deadline, and at most four
requests per attempt. It retries transport failures and unexpected preparation responses up to three total attempts within the same two-minute deadline, with paced backoff. Queued exports and form or document validation failures are not retried. 429/503 terminate
the attempt, honor numeric/date Retry-After (one-minute default), and return
that deadline through the existing signed snapshot. The CLI honors restored
retry deadlines across daily runs. Reuse one adapter per provider/process.
No new scheduler, checkpoint, recovery archive, or GitHub Release is introduced.

Forms are limited to 1 MiB, downloaded input to 16 MiB, and extracted documents
to 32 MiB. Only a single regular DAT member is accepted. Unexpected redirects,
encoding, truncation, ambiguous archives, changed forms, HTML, wrong header/name,
wrong header ID, malformed ROM records, and count-floor violations fail closed.
Acquisition failures preserve the last working catalog via the existing publisher.
Filesystem/configuration errors abort publication rather than replacing state.

## Inclusion policy v1 and completeness limit

The Standard DAT request explicitly chooses default representation, non-merged
collection, and canonical No-Intro names. It includes BIOS, full titles and
preproduction/demo entries, licensed/pirate/other-unlicensed records, original
lifespan and aftermarket records, adult records, physical/digital storage, and
nodump entries. All regions, languages, and both special-flag groups are selected.
The current Aftermarket checkbox submits the literal value `0`; its presence
selects it. No client-side filtering, 1G1R, Retool, or document rewriting occurs.

`nodump` records may omit size and hashes. Preserve that unknown status and the
original bytes; never invent values. Other ROMs need an unsigned numeric size
and correctly shaped CRC32/MD5/SHA1 fields; SHA256 is checked when supplied.
These checks establish structure, not authenticity of game dumps.

The official export contains 4,358 games/ROMs, including 202 names with the
Aftermarket tag, 28 nodump ROMs, and 18 baddump ROMs. A follow-up reconciliation
on 2026-09-06 compared it with two independent official public exports at the
same version, `20260818-050713`:

| Official surface | Records | Meaning |
| --- | ---: | --- |
| Download overview | 4,375 | Summary counter; does not match public database/search |
| Public search | 4,361 | Unfiltered public search result total |
| Dumplog CSV | 4,361 | Unique public archive IDs |
| Database XML export | 4,361 | Same archive-ID set as the dumplog |
| Database XML with `dat="0"` | 3 | Explicit upstream DAT exclusions |
| Database XML remaining IDs | 4,358 | Exact ID-set match with the downloaded DAT |

The three excluded records are:

- `4133`: Jeopardy! (Unknown) (Beta), no dump.
- `4220`: Chrono Trigger (Japan) (Final Fantasy Chronicles).
- `4221`: Chrono Trigger (Japan) (Beta) (Final Fantasy Chronicles).

All three have `dat="0"` on their database `<archive>` element. The two Chrono
Trigger records also carry the upstream sticky note that in-game dialogs appear
corrupted. We preserve the official export selection; we do not synthesize or
insert these intentionally excluded database records into the DAT.

Every eligible public database ID appears in the DAT, with no additional DAT
IDs. All 4,330 dumped ROMs match size and every supplied CRC32/MD5/SHA1/SHA256
against a source or release file in the corresponding database archive. The
other 28 entries remain nodump, with unknown sizes/hashes preserved. Names
match except for the expected `[UNDUMPED]` suffix on those 28 DAT entries.
This establishes completeness against the public DAT-enabled database snapshot,
not completeness of No-Intro's non-public data or independently verified ROMs.

The original 17-record difference therefore decomposes into **three deliberate
DAT exclusions and a 14-record overview-versus-public-database difference**.
Public search, dumplog, and database XML agree on the lower total. The reason for
those remaining 14 summary-counter entries is not exposed by the inspected
public surfaces. A stale counter or non-public records are possible explanations,
not established facts. Do not claim all 17 are private records. The overview
counter is unsuitable as an exact DAT acceptance count; the complete ID-set
reconciliation is stronger evidence that our acquisition loses no eligible
public records. No further mirroring-completeness gate is imposed solely by
that unexplained aggregate counter.

Evidence sources:

- [Download overview](https://datomatic.no-intro.org/index.php?page=download&s=49).
- [Public search](https://datomatic.no-intro.org/index.php?page=search&s=49),
  ordinary blank full-name search: `Showing: 200/4,361 item(s).`
- [Dumplog](https://datomatic.no-intro.org/index.php?page=download&op=dumplog&s=49):
  803,608 extracted CSV bytes, SHA256
  `3317fba6714269dee26f3bf40024a1c7d0cf2ae5ea48daf865b8f7941997d5a3`.
- [Database export](https://datomatic.no-intro.org/index.php?page=download&op=db&s=49):
  7,953,728 extracted XML bytes, SHA256
  `bb622ba7734557fb07b02a8d1809fbbad6178bc544bbbabeedcfa4d8f3c157b3`.
- Public records [4133](https://datomatic.no-intro.org/index.php?page=show_record&s=49&n=4133),
  [4220](https://datomatic.no-intro.org/index.php?page=show_record&s=49&n=4220),
  and [4221](https://datomatic.no-intro.org/index.php?page=show_record&s=49&n=4221).

The reconciliation used bounded anonymous downloads with five-second pacing,
then local ID and file-field comparisons. Upstream exports and the disposable
probe are not committed. It does not add daily database downloads, scraping,
platform enrollment, or scheduling infrastructure. The registry's 4,000-game/ROM
floors remain anomaly checks, not exact expected counts.

## Redistribution and rollout authorization

The operator confirmed redistribution and retention approval for this scoped
No-Intro mirror on 2026-09-06. That gate is accepted; it is not inferred from
robots guidance or Redump's separate terms. The inspected
[terms](https://datomatic.no-intro.org/stuff/terms.txt) and original DAT author,
trademark, and piracy notices remain preserved. This record does not invent an
upstream license grant or publish private operator qualification records.

The operator subsequently approved merge and deployment. Tooling
[PR #19](https://github.com/MoonlarkStudios/romd-dat-catalogs/pull/19) and data
[PR #3](https://github.com/MoonlarkStudios/romd-dat-data/pull/3) are merged;
both data workflow pins use tooling merge
`069bff51587ead4bcc62e6942af2bc8364ce3af0`.
`NOINTRO_SNES_PUBLISH_ENABLED=true` now selects SNES in the existing daily workflow.
The existing trust root, online signing conventions, and schedule were retained.
No credentials, session material, or private keys are committed. Application-side
reader deployment remains owned by the ROMD application work and was not verified
in this repository; see [the reader contract](deployment.md#no-intro-rollout-walkthrough).

## Evidence and acceptance boundaries

Real source evidence, local only:

- Anonymous form download: ZIP 542,753 bytes; DAT 1,911,272 bytes.
- Header ID `49`; version `20260818-050713`; 4,358 games and 4,358 ROMs.
- Extracted SHA256:
  `307fc5258970736d511c05b87d16268a9fbfff38d5bf6aa39dafd500e6500b19`.
- Production Go CLI acquisition/local publication succeeded at
  `2026-09-06T21:05:32.228022Z` and `2026-09-06T21:07:48.529476Z`.
  The second check retained the same artifact and first change timestamp.
  These are repeated same-version observations, not two actual upstream versions
  or evidence of sustained daily reliability.
- The same real document passed local TUF signing, authenticated candidate
  reading, asset/hash verification, and restoration with a disposable test root.
  Git asset resolution was simulated; this proves local signed-path compatibility,
  not a published immutable Git asset. The disposable probe was removed.
- Existing deployed Pages metadata and PSX documents restored successfully using
  the existing `trust/1.root.json`, without resetting the root or versions.

Real public deployment evidence (2026-09-07 UTC):

- [Initial production run 34069503517](https://github.com/MoonlarkStudios/romd-dat-data/actions/runs/34069503517)
  acquired the official DAT and published metadata version 3008, advancing 3007.
- The complete uncompressed document is committed at stable path
  `no-intro/snes/standard.dat`, with the signed index and RSS enclosure bound to
  [immutable commit 4f18544](https://raw.githubusercontent.com/MoonlarkStudios/romd-dat-data/4f18544db79f5a8f31e5cec937398d56adfa4908/no-intro/snes/standard.dat).
  Its hash, size, and 4,358 game/ROM counts match the real source evidence above.
- Independent `reference-data`, SNES and PSX `candidate`, and `restore` commands
  passed against deployed Pages using the original pinned root and retained TUF
  cache. The downloaded RSS matched its authenticated TUF target digest and
  included the SNES event and exact immutable enclosure URL.
- Existing PSX document hash remained
  `62572360b7abe18df80886283e39782c2c4a325e48c6134fc179d589936fe41d`.

- [Repeated production run 34069754217](https://github.com/MoonlarkStudios/romd-dat-data/actions/runs/34069754217)
  advanced signed metadata from 3008 to 3009 with fresh healthy acquisition checks.
  Independent authenticated reading confirmed identical DAT hashes, sizes,
  immutable URLs, change timestamps, and RSS bytes. Data `main` remained at
  `4f18544db79f5a8f31e5cec937398d56adfa4908`: no new document commit or content event.
  This is real same-document evidence; ZIP repackaging and genuine document
  changes remain covered by synthetic tests, not a newly observed upstream pair.

Synthetic/offline evidence:

- Explicit system/provider routing, policy controls, rejected form/redirect
  changes, header ID and ROM structure, nodump preservation, request pacing,
  size/encoding/truncation failures, 429/503 backoff, cancellation, and last-working
  artifact retention.
- Repackaged ZIPs retain document identity and RSS GUID; changed synthetic bytes
  produce new content identity and RSS. Real local Git commits prove stable
  `no-intro/snes/standard.dat` paths and full-commit URL history behavior.
- A hand-authored No-Intro-shaped fixture passes provider validation, export,
  TUF signing/verification, shared system binding, exact-byte candidate reading,
  immutable asset verification, RSS enclosure checks, and authenticated restore.
- `mise install`, `mise run check` (race tests, vet, gofmt, actionlint, three
  builds), and `mise run smoke` all passed on 2026-09-06 UTC.

No genuine changed upstream pair was observed, and ROMD activation was not
exercised. Consumers
must use publisher reader code that accepts the `no-intro` provider and its
hyphenated immutable paths; older strict readers can reject the added registry.
The signed contract remains schema 2 and Git format 2; there is no new schema or
ROMD database identity. Availability comes from the published catalog snapshot,
not merely the presence of this gated registry entry.

## Fresh1G1R comparison

Rechecked [Fresh1G1R commit 86e21ac](https://github.com/UnluckyForSome/Fresh1G1R/tree/86e21acf76cffd20051eb5bb0fd9fcbddba3eb0b).
Its No-Intro path uses headless Chromium on the daily pack and the first matching
Download label, then processes data with Retool. This informed investigation of
the form sequence, not ROMD's acquisition implementation or completeness policy.
No third-party code was copied or executed. The official per-system form was
independently exercised instead of downloading/filtering an entire daily pack.

If the form stops working, pause this catalog and retain its last version.
The smallest fallback is an operator-downloaded, all-category official Standard
DAT passed through the same No-Intro validator and publisher/signing/export path
in an explicit reviewed local run. The existing generic manifest path alone does
not enforce provider-specific completeness/ID checks. Do not silently call manual
ingestion automatic support or enable a fallback mirror. No manual upload service
or credentials-based fallback is added in this slice.
