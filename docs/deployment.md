> Current reader contract: new tooling publishes DAT-only schema 3 definitions.
> ROMD now owns application reference data. Use `distribution catalogs` for
> discovery and coordinate ROMD's pinned reader before updating the live workflow.
> Historical publication evidence below describes the earlier reference-data format.

# Signed catalog deployment

The companion repository contains Go tooling, workflows, fixtures, and platform
identities. The separate public data repository
`MoonlarkStudios/romd-dat-data` contains complete, uncompressed DAT documents at
stable paths, for example `redump/psx/discs.dat`. Git retains older versions.
There is no automated trimming or history rewriting.

Status (2026-09-06): the data repository's Pages site is live at
https://moonlarkstudios.github.io/romd-dat-data/, with RSS at
https://moonlarkstudios.github.io/romd-dat-data/feed.xml. It publishes only the
hand-authored synthetic catalog. Daily metadata renewal is enabled; public
PSX mirroring is enabled separately after deploying the reviewed tooling pin;
see [the scoped redistribution qualification](redump-psx-qualification.md).

The [migration run](https://github.com/MoonlarkStudios/romd-dat-data/actions/runs/34016250699)
and [ordinary unchanged run](https://github.com/MoonlarkStudios/romd-dat-data/actions/runs/34016317753)
both passed. Independent local verification reused the old site's pinned root
and TUF cache: metadata advanced from 3001 to 3002 to 3003. The second publication
kept data commit `867713188e3360334785dd13ef76aea4ade0cb57` unchanged. Its only
DAT is `synthetic/console/standard.dat` (411 bytes, SHA-256
`86ac64a0c11c1c8f2c864fc57589de00d5dba85821c8a495034e463248bd68e9`).
The candidate reader verified the new public URL; RSS bytes matched the verified
TUF feed target. No Releases or recovery archives were created in the data repo.

This proves synthetic Git-backed publication and site migration, not public
Redump acquisition or a new ROMD UI/activation/hardware acceptance run. The
existing private ROMD demo configuration was not changed.

## Shared reference-data publication

The initial data workflow pinned tooling revision
`19f983c0e58d62e7f97ba8bf860c6d1a8c289fde`. The
[initial reference-data run](https://github.com/MoonlarkStudios/romd-dat-data/actions/runs/34039634452)
and [restoration/unchanged run](https://github.com/MoonlarkStudios/romd-dat-data/actions/runs/34039736902)
both passed, advancing signed metadata to versions 3004 and 3005. Independent
verification reused the existing root/cache and confirmed that public
[reference-data.json](https://moonlarkstudios.github.io/romd-dat-data/reference-data.json)
matched its TUF target, the signed index definitions, and restored local state.

That initial snapshot contained one system (PSX), its catalog definition, and 15 company
identities. Its SHA-256 is
`95268b885395daa596a7561f7d714b39916e65a42bac25d106c47f1289aa86bd`.
It stayed identical on the second run, with no additional data-repository
commit. The synthetic candidate also passed independent reader verification.
PSX acquisition remains disabled; this publishes reference data, not a PSX DAT.

The [complete seed publication](https://github.com/MoonlarkStudios/romd-dat-data/actions/runs/34040934356)
passed at metadata version 3006, using tooling
`249ac9dd0bb25fd9ba2faf9eb9795ab963a16226`. The schema-2 snapshot has 56 systems,
15 companies, 20 regions, and 16 languages, including ROMD's existing aliases
and IGDB mappings. It is 9,651 bytes with SHA-256
`a0f539d85ed9db1619ce1e44574a1973436070769e23afcad3027ba0cf0cbd4f`.
The metadata-only reader reused the existing bootstrap root and client cache,
verified the new signed index, and confirmed byte equality with the public
reference-data alias. The workflow also verified deployed signatures and full
restoration. Only the synthetic catalog is available; the PSX definition is not
a published PSX document. ROMD's shared-data consumer/UI remains a separate
acceptance checkpoint.

## Development and verification

Run `mise trust`, `mise install`, then `mise run check` (race tests, vet,
formatting, actionlint/ShellCheck, and builds). Git must also be available;
new lifecycle tests exercise real local Git commits without network access.
Go and validation tool versions are pinned in `.mise.toml`. `mise run smoke`
exercises the unsigned local publisher. No Python runtime is needed.

## Data and trust contracts

Only changed extracted DAT bytes produce a data commit. ZIP repackaging, check
timestamps, signed metadata renewal, and failed acquisition produce no DAT
change. The workflow exports only current documents, compares the Git index,
and commits/pushes only if that index changed. Store DATs without Git LFS or
text conversion (`*.dat -text -filter` in the data repository's `.gitattributes`).
Git's own internal object compression is independent of the raw documents served
to clients.

Signed catalog format `romd-signed-catalog-2` binds every current document's
length and SHA-256 to a raw GitHub URL containing a full 40-character commit ID:
`https://raw.githubusercontent.com/OWNER/DATA_REPO/COMMIT/redump/psx/discs.dat`.
Branch URLs, short commits, query strings, and fragments are rejected. The data
commit is pushed and all referenced raw documents are downloaded/hash-verified
before Pages advances. Never force-push or rewrite data history: that could
break URLs in previously signed indexes. No historical browsing or publisher
rollback service is a beta requirement. ROMD owns its installed active document
locally and validates/reviews each update before activation.

The data repository's Pages site publishes the signed index and RSS using the existing pinned trust root
and [go-tuf/v2](https://github.com/theupdateframework/go-tuf). Clients verify the
root chain, signatures, expiry, cached versions, and exact document hash/size.
The root is supplied independently; never use trust-on-first-use. Preserve the
client cache for rollback detection. Old `ae48371` candidate readers understand
only format 1 and must be rebuilt from this Git-format implementation before
cutover. The `candidate.json` contract consumed by ROMD does not change. Older
readers fail closed; they do not activate an unverified update.

Root metadata expires after one year; targets after 30 days; snapshot after
seven days; timestamp after 48 hours. Daily runs renew signatures even when DAT
bytes are unchanged. Latest-only RSS keeps one content-change event per catalog
with a stable GUID across unchanged checks. RSS is a notification hint; clients
reconcile against the verified full index. The plain `/feed.xml` is not a trust
entry point. Historical feed browsing is not provided.

## Configure the Git-backed workflow

The implementation is a reusable `publish-synthetic.yml` workflow in the tooling
repository. A small caller in the data repository owns manual/daily triggers.
Both the workflow reference and tooling checkout are pinned to the same full
reviewed commit. [GitHub reusable workflows](https://docs.github.com/en/actions/reference/workflows-and-actions/reusing-workflow-configurations)
run in the caller's repository context: data commits, Pages artifacts/deployment,
RSS, variables, and the online signing secret all belong to the data repository.
The ordinary `GITHUB_TOKEN` writes that repository; no cross-repository token or
`DAT_DATA_REPOSITORY` variable is needed. Tooling checkout credentials are not
persisted. The offline root key remains outside CI.

The canonical site is
`https://moonlarkstudios.github.io/romd-dat-data/`, with RSS at `/feed.xml`.
The tooling-repository site is the old prototype and no longer renews its
metadata. Its old Releases remain intact; the old RSS URL is not redirected.
Use the data site for new clients. The steps below document the completed
cutover and remain the procedure for an equivalent future site move.

Before deploying the cutover:

1. Initialize the public data repository's `main` branch using
   `templates/data-repository/`, including `.github/workflows/publish.yml` and
   `.gitattributes`. The caller is pinned to the reviewed tooling implementation;
   update both SHA pins together on future upgrades. Configure Pages to deploy
   through GitHub Actions, and permit the pinned reusable workflow/actions.
2. Restore the existing online role keys into the data repository's
   `TUF_ONLINE_KEYS` repository secret from private custody. Preserve the exact
   committed public root chain. Rebuild deployed Go candidate readers for
   signed format 2. Do not generate a new root or expose root private keys.
3. Stop the old tooling-site publisher and wait for any active run to finish.
   Complete the move before its 48-hour timestamp expires. The reusable
   workflow replaces the tooling repository's former standalone scheduler.
4. Dispatch **the data repository's** `publish.yml` with `bootstrap=false` and
   `migration_site=https://moonlarkstudios.github.io/romd-dat-catalogs`.
   Migration is accepted only when the new site's timestamp returns HTTP 404.
   It verifies the old site's root chain, freshness, index, and latest data,
   then publishes a metadata version greater than that verified old version.
   HTTP errors, an existing destination, or expired old metadata fail closed.
5. Verify the new site's deployed signatures, exact version, and DAT hashes.
   Update ROMD's configured Site and RSS subscriptions to the canonical data
   site. Carry forward the existing pinned root and verified TUF cache when
   moving clients; ROMD's cache directory is site-specific, so plan that cache
   transfer before enabling its new location. Do not reset trust to bypass
   expiry or rollback checks. The old RSS URL does not redirect automatically.
6. On subsequent manual runs leave `migration_site` empty and `bootstrap=false`.
   Daily runs always use these normal settings and restore only from the data
   site's verified index. Set `DAT_PUBLISH_ENABLED=true` in the data repository
   for the daily schedule. `SYNTHETIC_PUBLISH_ENABLED` is no longer consulted.
7. Keep `REDUMP_PSX_PUBLISH_ENABLED` absent/false until public redistribution
   is qualified. Synthetic catalogs can validate the site migration first.
   No new upstream DATs belong in the data repository before qualification.

`bootstrap=true` is only for a new independent deployment with an absent site;
it is mutually exclusive with migration. The existing deployment must migrate.
The initial migration can read the old signed format-1 recovery archive;
subsequent runs fetch only latest DATs from signed format-2 commit URLs. Routine
runs create no Releases or recovery archives. Existing Releases remain intact
for old references. Legacy `stage --release-base` supports local demo bundles
only and is not called by this workflow.

Every run checks in the same Go/toolchain environment, restores authenticated
latest state, acquires/validates into staging, exports DATs, conditionally commits
and pushes data, signs the index/feed, verifies public commit URLs, and deploys
Pages. A final check verifies the deployed exact metadata version and restores
latest state. A competing data push fails normally; retry the workflow. Never
force-push or introduce remote checkpoints to resolve it.

The data checkout includes history so `git count-objects -vH` in the run summary
can expose stored-object growth. This is a measurement, not a storage budget or
a precise GitHub quota reading. Also inspect current document sizes and actual
change frequency before broader coverage. No cleanup engine, bounded history,
or incremental patch format is introduced. Issue #12 tracks growth measurement.

## Publication failure and pause

Provider failures are signed as failed health while keeping the last document
and its successful/change timestamps. Installed ROMD catalogs remain usable;
the reader can return a retained artifact for new subscription or update review while current signed metadata and artifact integrity still verify. Refresh health remains visible separately. Missing artifacts and trust, expiry, or integrity failures still block acquisition. A paused existing PSX catalog
is reported as `publication_paused`, without deleting its DAT. A never-enabled
catalog is not invented. Keep ordinary metadata publication enabled while
pausing acquisition. Previously public DATs remain in Git history.

Failures during restoration, signing, data push, public hash verification, or
Pages deployment are visible Actions failures. A data commit made before a
later failure is harmless: retry exports the same bytes and reuses that commit.
Pages never advances before public documents are verified. Metadata versions
advance by one from the authenticated restored version. A new independent
publication starts at 1. Moving repositories or rerunning a failed job therefore
cannot reset the sequence. All data-repository publications use one concurrency
group. There must be no simultaneous publisher at the old site during cutover.

If metadata expires during an outage, restoration deliberately fails closed.
Do not bootstrap or disable expiry checks. An operator must use a previously
verified local state, refresh/rotate the root if necessary, and stage a higher
version. Automatic outage recovery is deferred until demonstrated need.

### Private key custody

Back up private signing keys in an access-controlled secrets manager and verify
that restored values exactly match the original files before publication. Keep
account names, vaults, item identifiers, recovery access, and operator-specific
backup locations in private operational records, outside this repository.

Preserve the exact JSON representation of each versioned public root and its
root and online key sets. Record root versions, expiry dates, and public-root
SHA-256 fingerprints with the private backup. These are TUF Ed25519 keys, not
SSH keys. Retain previous versions during rotation.

The root private key is for offline rotation only and must never be available
to CI. Only `online.json` supplies the `TUF_ONLINE_KEYS` Actions secret. Transfer
secret values through stdin or restricted private files, never command-line
arguments, chat, logs, or committed templates. CI must not have access to the
root-key backup.


## Key rotation

Prepare a new key set locally using the old root private key:

```sh
./bin/distribution rotate --previous-root trust/1.root.json \
  --offline-key .keys/synthetic-initial/offline-root.json \
  --out .keys/synthetic-rotation-2
cp .keys/synthetic-rotation-2/public/2.root.json trust/2.root.json
```

The new root is signed by both the old and new root keys and replaces all online
role keys. Keep every versioned public root so an old client can traverse the
chain. Pause scheduled publishing, merge the new public root, update
`TUF_ONLINE_KEYS` from the new `online.json`, run and verify a manual publication,
then resume scheduling. Substitute the latest root/version for later rotations.
Never place the offline root key in CI. Old online keys are rejected by staging
after rotation; this is exercised by tests with an old pinned client root.


## Validation scope

Local tests exercise format-1 to format-2 migration with the same TUF cache,
destination-absence migration gates (including HTTP/network failures),
real local Git commits, unchanged ZIP repackaging, failed acquisition retention,
changed documents, older immutable URLs after branch advancement, latest-only
restoration, and immutable URL restrictions. Existing signature, wrong-root,
expiry, rollback, tampering, and candidate-binding tests remain.
These local tests complement the live synthetic runs recorded above. Neither
establishes public Redump redistribution or a new ROMD activation/hardware test.

Qualification status (2026-09-06): the current official redump.info site explicitly
permits metadata reuse. The standard PlayStation DAT is qualified for centralized
mirroring; see [scope, provenance, and evidence](redump-psx-qualification.md).
Enable only `REDUMP_PSX_PUBLISH_ENABLED` after deploying the reviewed tooling pin.
Do not infer qualification for No-Intro, restricted exports, BIOS, or other platforms.


## Gated No-Intro SNES slice

The reusable workflow also selects `no-intro/snes/standard` when
`NOINTRO_SNES_PUBLISH_ENABLED=true`. This scoped mirror was enabled after operator
approval on 2026-09-07 UTC; see the production evidence in the qualification record.
[Public-catalog reconciliation and operator redistribution approval](nointro-snes-qualification.md)
are recorded for this slice.
Definition availability is not publication authorization. A never-enabled
catalog is not invented; a paused existing catalog retains its DAT.

After qualification, update both data caller pins to the reviewed tooling commit
and use the existing once-daily/manual workflow. The anonymous official form
runs centrally, preserving all selected public entries. No credentials or new
secrets are required. Existing online role keys sign the catalog index and RSS;
root keys stay offline. Export writes `no-intro/snes/standard.dat` before signing
its exact immutable commit URL. ZIP changes alone do not commit DATs or notify
RSS; routine metadata renewal still occurs. Failures retain the last document.

Deploy compatible reader binaries before publishing the extended registry:
older strict readers reject the newly supported provider. No ROMD runtime code
is changed by the publisher PR; review the signed candidate contract described
in the qualification document.


## No-Intro rollout walkthrough

The reader is the `distribution` executable built by this repository. It verifies
TUF metadata, reads shared reference data, and downloads a system-bound candidate
for ROMD. Older builds reject the newly supported `no-intro` provider or its
hyphenated Git paths. Even a disabled catalog definition enters the shared signed
registry, so reader deployment precedes the publisher pin update. No schema,
root-key, or ROMD database migration is required by this change.

1. **Review and merge the tooling change when authorized.** Record its actual
   reviewed commit on `main`; a squash merge may have a different SHA from the PR
   head. Build that revision with `mise install` and `mise run check`.
2. **Hand the reader contract to the ROMD application owner.** Package/deploy the
   resulting `bin/distribution` wherever ROMD runs `catalogs` or `candidate`.
   If a consumer implements verification itself, it must accept the new registry
   provider and immutable path. Preserve its independently pinned public root,
   existing TUF cache, publisher identity, and installed DATs. Never reset cached
   metadata versions to make an update pass. The candidate output contract stays
   `candidate.json` plus byte-exact `candidate.dat`, with shared `systemId: snes`.
   The application owner confirms its configured reader was replaced and current
   PSX/reference-data checks still work; publisher CI cannot establish deployment
   inside ROMD. No application checkout is edited by this publisher task.
3. **Review and deploy the data caller pin.** Update both its reusable-workflow
   reference and `tooling_ref` to the same reviewed tooling commit. Merge the data
   PR only after compatible readers are in place. Keep SNES disabled for this
   first run if a separate metadata-only compatibility check is useful.
4. **Enable the approved SNES slice and run one publication.** Set only
   `NOINTRO_SNES_PUBLISH_ENABLED=true` in the data repository and invoke the
   existing workflow, with bootstrap false and migration site empty. Existing
   online keys and root metadata are reused. The workflow restores authenticated
   state, downloads centrally, commits the uncompressed DAT before referencing it,
   signs the new index/RSS, verifies public hashes, deploys Pages, and restores the
   deployed result to verify its signature chain. It must advance existing metadata
   versions. The normal daily schedule then uses the same path.
5. **Verify the client-facing result before calling the slice live.** Use the
   updated reader's `catalogs` and `candidate --catalog no-intro/snes/standard
   --name 'Nintendo - Super Nintendo Entertainment System'` against the existing
   Pages URL and pinned root, using its normal cache and a fresh output directory.
   Check system `snes`, the published document hash/size, and a full immutable Git
   commit URL. Confirm RSS's corresponding content event. A second unchanged run
   must keep the DAT commit/content event; metadata signatures may renew. ROMD's
   preview/approval/activation acceptance remains owned by the application work.

The operator approved this rollout, and both publisher PRs are merged. SNES is
enabled and the signed public catalog independently verifies with the rebuilt
reader; see [production evidence](nointro-snes-qualification.md#evidence-and-acceptance-boundaries).
Deployment of that reader inside ROMD and application activation acceptance remain
with the application owner. No new signing keys or credentials were needed.
