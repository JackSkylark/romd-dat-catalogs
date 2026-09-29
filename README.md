# ROMD DAT catalogs

Companion publisher for ROMD's planned one-click No-Intro and Redump DAT
subscriptions. Complete upstream catalogs come first; future 1G1R filtering
belongs in ROMD's library layer.

**Current status: 26 signed catalogs are live: PSX, SNES, Game Boy, Game Boy Color, Game Boy Advance, NES (Headered), Genesis, N64 (BigEndian), Master System, Game Gear, PC Engine/TurboGrafx-16, 32X, Saturn, Sega CD, Dreamcast, GameCube, PlayStation 2, Wii, Nintendo DS (Decrypted), Nintendo 3DS (Decrypted), PSP, PC Engine CD/TurboGrafx-CD, Xbox, Xbox 360, PlayStation 3 and Wii U.**
Signed publication 3037 reports all 26 real catalogs healthy, including NES
after Datomatic completed its queued export.
See the [SNES closeout and expansion plan](docs/catalog-expansion-plan.md).
The Game Boy family has [passed acquisition, public membership and isolated
ROMD catalog acceptance](docs/nointro-handheld-qualification.md).
[Public RSS](https://moonlarkstudios.github.io/romd-dat-data/feed.xml) ·
[Deployment runbook](docs/deployment.md)

The local publisher still emits explicitly unsigned intermediate state. The
distribution command wraps it with TUF-authenticated catalog/feed targets for
GitHub Pages. The live workflow stores complete uncompressed DATs in a
separate public data repository and signs full-commit raw URLs. That repository
also hosts Pages/RSS, calling a pinned reusable workflow from this tooling repo; routine Releases
and recovery archives are no longer the storage model. The [Redump acquisition
library](docs/redump-adapter.md) has a qualified PSX workflow path. The
[No-Intro SNES adapter](docs/nointro-snes-qualification.md) is also deployed,
following public-catalog reconciliation and operator redistribution approval.
The Game Boy family is enabled after complete-public qualification and application
acceptance. [NES, Genesis and N64](docs/nointro-home-console-qualification.md)
are also live after complete-public qualification. [Master System, Game Gear,
PC Engine/TurboGrafx-16 and 32X](docs/nointro-batch-three-qualification.md) passed
qualification, publication and isolated ROMD acceptance. [Saturn, Sega CD and
Dreamcast](docs/redump-sega-qualification.md) also passed those gates using the
standard Redump disc catalogs. [GameCube, PlayStation 2 and Wii](docs/redump-gamecube-ps2-wii-qualification.md)
also passed qualification, signed publication and isolated ROMD acceptance.
[Nintendo DS, Nintendo 3DS, PSP and PC Engine CD](docs/nds-3ds-psp-pcecd-qualification.md)
have now passed the same release gates.
[Xbox, Xbox 360, PS3 and Wii U](docs/redump-xbox-ps3-wiiu-qualification.md)
are published and healthy; their application acceptance status is recorded in
that qualification report. These are the 26 enabled
real catalogs; synthetic fixtures remain
separate. An operator-driven [candidate reader](docs/romd-candidate.md)
is available for the first ROMD review integration. Do not enable production
AutoApply against synthetic catalogs.

ROMD owns application reference identities and presentation metadata. This repository
owns [DAT acquisition definitions](docs/system-definitions.md) and vendors a versioned
ROMD system-key export to validate explicit catalog mappings. Signed publications
carry DAT catalog definitions; they do not update application reference data.

## Run

Use mise on macOS or Linux (POSIX file locking required). `.mise.toml` pins
Go 1.27.1, actionlint, and ShellCheck for development and CI. `go.mod` records
the minimum language version; `GOTOOLCHAIN=local` prevents silent toolchain
switching. Each CLI builds into a standalone executable:

```sh
mise trust
mise install
mise run check
mise run smoke
```

CLI flags precede the manifest argument. The example hostname is deliberately non-routable. For local serving tests,
the objects remain accessible under `output/`; HTTPS hosting is a later step.

The manifest is a trusted local registry of acquisition attempts. Each entry
contains a stable `catalogId`, exact expected DAT header name, public source
URL, and either a path to a DAT/single-DAT ZIP or a `failure` marker. Paths
are relative to the manifest. Never put credentials in URLs or manifests.
This prototype's limits are intentionally suitable for small fixtures, not
yet qualified for the entire upstream daily pack.

## What works

- Exact extracted-document hashes, preserving upstream bytes and notices.
- Bounded DAT/ZIP handling, safe XML parsing, expected-identity checks.
- Per-catalog failure state retaining the prior working artifact.
- Immutable content-addressed artifacts, RSS, and complete catalog snapshots.
- Stable feed event identities; an upstream restoration can reuse old content
  with a new event. Feed truncation does not limit index-based catch-up.
- Atomic replacement of `current.json` after its objects are written and
  verified, with a persistent OS lock serializing local publisher processes.
- Rejected-version exclusion in a small reference catch-up function.

The Go implementation resumes existing Python prototype snapshots; a committed
synthetic snapshot from Python commit `02d300c` exercises compatibility without
a Python runtime. JSON/XML serialization can differ while DAT hashes and event
identities remain stable. DAT comparison remains byte-based: semantic handling
of header-only or formatting-only changes is a separate future feature.

XML validation reads tokens without a full document tree, although the bounded
original document and a set of game names remain in memory. UTF-8 XML is
supported; other declared encodings fail closed pending explicit qualification.
The full-source acquisition path still needs large-catalog qualification. To measure the synthetic 10,000-game parser workload:

```sh
mise run bench
```

This benchmark is a reproducible measurement tool, not a full upstream run or
evidence of a speedup over Python.

`current.json` points to an exact snapshot by path, length, and SHA-256 and
explicitly declares `trust: unsigned-development`. The snapshot references
an immutable RSS file and every current catalog. RSS is a notification hint;
the complete index is the basis for reconciliation. This local intermediate state
has no mutable `feed.xml`; signed distribution adds the stable RSS alias and
TUF-authenticated targets described in the deployment runbook.

Publisher tests exercise local process interruption before pointer commit;
distribution tests exercise signed metadata and client verification over local
HTTP fixtures. They do not prove host power-loss durability, GitHub deployment,
upstream authenticity, or ROMD activation safety. DAT checks here are structural;
full ROM/hash semantics and anomaly policies remain pending.

## Next work

SNES acquisition, signed publication, and initial ROMD catalog acceptance are
complete within the recorded scope. The September 2026 closeout verifies a
subsequent real document change and current healthy acquisition. ROM import,
hardware launch, and activation of that latest real revision were not exercised
in the closeout. See [qualification evidence](docs/nointro-snes-qualification.md)
and the [staged expansion plan](docs/catalog-expansion-plan.md).

Qualify each additional catalog before enabling it. Reuse the existing adapters,
daily publication, signed distribution, and reviewed application activation.
Existing [source-state tools](docs/source-state.md) remain available for explicit
operator use; they are separate from the deployed daily acquisition workflow.
Readers do not need a GitHub API token. Operational outage recovery remains a
separate acceptance area.

Only hand-authored synthetic DATs belong in current fixtures. No upstream
content is distributed by this repository. No-Intro and Redump remain the
authorities for their catalogs; this project is not affiliated with them.
