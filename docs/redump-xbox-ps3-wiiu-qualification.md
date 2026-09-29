# Xbox, Xbox 360, PlayStation 3 and Wii U qualification

Recorded 2026-09-20 UTC. The operator approved this batch and investigation of
NES acquisition. Standing redistribution approval is recorded in AGENTS.md.

## Identity and representations

[Official Redump downloads](https://redump.info/downloads) bind Xbox to XBOX,
Xbox 360 to XBOX360, PlayStation 3 to PS3 and Wii U to WIIU. The reviewed lowercase
HTTPS routes are `/datfile/xbox`, `/datfile/xbox360`, `/datfile/ps3`, `/datfile/wiiu`.
They return single-DAT ZIPs without redirects. ROMD keys are `xbox`, `x360`,
`ps3`, `wiiu`; catalog representation is `discs`. Exact original image metadata
and DAT bytes are preserved, including ancillary CD entries, cue files, prototypes
and separately named discs. There is no 1G1R, region selection, image trimming,
decryption or conversion. These DATs do not establish support for hashing
extracted directories, compressed images or other transformed game-file formats.
Separate BIOS and digital-distribution catalogs are outside this batch.

## Measured exports

| System | Header | Version | Entries / files | ZIP / DAT bytes | Floors |
| --- | --- | --- | --- | --- | --- |
| xbox | Microsoft - Xbox | 2026-09-19 20-42-35 | 2688 / 2694 | 216548 / 949755 | 2500 / 2500 |
| x360 | Microsoft - Xbox 360 | 2026-09-19 21-12-49 | 3709 / 3718 | 305372 / 1379853 | 3400 / 3400 |
| ps3 | Sony - PlayStation 3 | 2026-09-18 06-03-50 | 4518 / 4531 | 397009 / 1680338 | 4200 / 4200 |
| wiiu | Nintendo - Wii U | 2026-09-01 01-39-02 | 542 / 542 | 45146 / 205894 | 500 / 500 |

| System | DAT SHA-256 | Largest expected game-file size | Files above 4 GiB |
| --- | --- | --- | --- |
| xbox | `bb4b2fe70d52fb6b5acea2f268bee70f09a161e17d71dc216809954a1f2630fc` | 7838154752 | 2632 |
| x360 | `54504a55fe4d1e44e333c96b720bde5e3b2d6cd1c97f99199b83917cc67ffdd4` | 8738854912 | 3687 |
| ps3 | `cae71727f1d4405f2b83c86b02648b40eea4ca66aa979c78fec982c76a773103` | 49747525632 | 3459 |
| wiiu | `095a5742dd8497b97f7641bd06cd7c591def4f3e26d29d0b37b8358c5cc17073` | 25025314816 | 542 |

All DATs fit the unchanged 16 MiB input / 32 MiB expanded bounds. Expected game-file
sizes are metadata values, not files downloaded by this publisher. All records
have size, CRC32, MD5 and SHA1. Four explicit workflow opt-ins control publication.

## Complete-public reconciliation

- Xbox: all 2,688 public disc IDs across 27 pages match the standard DAT.
- PS3: all 4,518 public disc IDs across 46 pages match the standard DAT.
- Wii U: all 542 public disc IDs across six pages match the standard DAT.
- Xbox 360: 3,713 public IDs across 38 pages versus 3,709 in both official
  standard and serial/version DATs. The four absent IDs are
  [8654](https://redump.info/disc/8654), [8660](https://redump.info/disc/8660),
  [8663](https://redump.info/disc/8663), [8665](https://redump.info/disc/8665),
  German Xbox magazine demo discs. Every detail page labels its dump status
  Questionable. These are observed upstream export exclusions; this is not a
  general claim about all Questionable records. ROMD adds no exclusion and does
  not modify the official export to synthesize these records. No other public
  IDs are absent and no DAT IDs fall outside the public inventory.
- Each serial/version export has exactly the standard DAT's IDs, matching every
  binary size/CRC32/MD5/SHA1 tuple by disc ID and multiplicity. Zero mismatches.
- The six Xbox, nine Xbox 360 and thirteen PS3 cues match the standard DAT's CD
  subset by name, byte length/SHA1, and binary filenames/order. Wii U has no cues.

These comparisons check upstream consistency, not independent rehashing of game
images. Original exports, public listing snapshots and temporary scripts remain
outside the tooling repository.

## NES availability diagnosis

The previous batch's `form_changed` diagnostic did not establish an upstream
form-contract change. The current full NES selection form passes the unchanged
strict parser. A captured anonymous response instead contains Datomatic's explicit
notice that the requested export is temporarily unavailable and queued. This can
replace the selection form, preparation response or manager page.

The adapter now classifies that reviewed notice inside `main_form` as
`upstream_pending`. It stops acquisition, retains the last working artifact and
adds no automatic retries or invented retry deadline. Normal provider pacing and
all successful-form/DAT validation remain unchanged. Unrelated shoutbox content
and incomplete notices do not trigger the classification. The publisher's signed
public health remains the existing `failed` / `acquisition_failed` contract;
`upstream_pending` is the specific acquisition diagnostic in workflow logs.
Actual NES recovery still depends on upstream completing generation.

## Validation and rollout

`mise run check` passed race-enabled tests, vet, formatting, workflow validation
and builds. Tests exercise the queued notice at all three acquisition stages,
prove no follow-up request or candidate is produced, and preserve the original
artifact/event history. Signed candidate tests cover all four new mappings with
exact byte preservation, immutable downloads, feeds and restoration.
`mise run smoke`, definition validation and all four live production-adapter
acquisitions also passed with the exact qualified hashes. NES returned
`upstream_pending`, as expected from the captured provider response.
Local qualification alone does not imply production or application acceptance;
the observed rollout results are recorded below. No ROMD source, schema, API or runtime changes
are required for this catalog batch.


## Full-registry release guard

The initial release run 35478341726 stopped before acquisition or publication:
the CLI's historical combined selection limit was 25, while this registry has
26 catalogs. The existing live site was unchanged. The CLI limit is now bounded
at 32 combined active/paused selections; per-provider request pacing and adapter
limits are unchanged. A new offline test runs every current definition through
both active provider routing and all-paused mode, so future registry growth
cannot silently outgrow the CLI. Oversized mixed selections fail before registry
or network I/O. Full `mise run check` passed after the correction. The corrected
implementation was subsequently pinned before rerunning publication.


## Signed publication

Implementation [PR 35](https://github.com/MoonlarkStudios/romd-dat-catalogs/pull/35)
and capacity correction [PR 36](https://github.com/MoonlarkStudios/romd-dat-catalogs/pull/36)
passed CI and merged. The final data workflow pins both references to
`02ab4e05d556c794d14b73061723c67d02ce14f0` via
[data PR 12](https://github.com/MoonlarkStudios/romd-dat-data/pull/12).
[Run 35478946411](https://github.com/MoonlarkStudios/romd-dat-data/actions/runs/35478946411)
completed signed publication 3037 at 2026-09-20T00:36:00.817915095Z. All 26 real
catalogs are healthy. The original trust root and retained reader cache verified
discovery, and all four immutable candidate downloads match the measured DAT
bytes and hashes exactly. Each has a content event at sequence 84 and a full-commit
URL under data commit `a1c61129364ae9b2bcbac52b024b16615f239765`.

NES generation completed upstream. Its recovered Headered DAT has 7,654 entries,
7,656 files and 3,916,757 bytes, SHA-256
`b540371878284bebf631299dc79095c293c1020fdfffb15aa906ca30bb50ce6c`.
The independent signed candidate download matches that length/hash. This is a
genuine upstream revision, not a parser relaxation or locally fabricated change.
The original public database membership audit was not repeated for this revision.
The queue diagnostic fix remains in place for future temporary unavailability.

## Application acceptance

Wii U was discovered, reviewed and activated in the isolated ROMD environment,
with 542 entries and 542 files. Its unchanged check returned UpToDate and its
active download matches the signed DAT bytes. All 542 imported filenames, sizes
and CRC32/MD5/SHA1 tuples also match, including 25,025,314,816-byte size values.
The Xbox, Xbox 360 and PS3 systems
are prepared; their pending-subscription/import actions await the browser approval
requirement. No completed import, field comparison or outage test is claimed for
those three yet.
