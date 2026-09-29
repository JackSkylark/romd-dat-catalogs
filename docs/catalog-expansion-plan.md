# SNES closeout and catalog expansion

Recorded 2026-09-18 UTC. Scope: close out the existing SNES catalog and plan
additional catalogs. No new catalog is enabled by this document.

Batch 1 implementation and local qualification are now recorded in
[Game Boy family qualification](nointro-handheld-qualification.md). All three
complete-public acquisitions and database membership comparisons passed.
All three are now publicly deployed and passed isolated ROMD catalog acceptance.
Batch 1 is complete within that scope. The operator authorized immediate GBC/GBA
release after Game Boy acceptance, superseding the scheduled-observation gate;
no scheduled handheld check is claimed.

[Batch 2 qualification and rollout](nointro-home-console-qualification.md) covers
NES (Headered), Genesis and N64 (BigEndian). All three are now publicly deployed
and passed isolated ROMD import, unchanged-check, failure-retention and recovery
acceptance. Batch 2 is complete within catalog scope.
[Batch 3 qualification and rollout](nointro-batch-three-qualification.md) covers
Master System, Game Gear, PC Engine/TurboGrafx-16 and 32X. All four are publicly
deployed and passed the same isolated ROMD acceptance gates. Batch 3 is complete
within catalog scope. [Batch 4 qualification and rollout](redump-sega-qualification.md)
covers Saturn, Sega CD and Dreamcast, now publicly deployed and accepted in ROMD.
All four planned expansion batches are complete within catalog scope.

## SNES closeout

SNES acquisition and public distribution were already approved and deployed;
the README's pending qualification status was stale. The original
[qualification](nointro-snes-qualification.md) remains the inclusion-policy and
redistribution record. ROMD separately records successful system setup,
review, activation, unchanged detection, failure retention and retry in
`docs/reports/snes-subscription-acceptance-2026-09.md` in the ROMD repository.
That acceptance concerns catalog ingestion, not ROM import or hardware launch.

A fresh `distribution candidate` and `distribution catalogs` check against
https://moonlarkstudios.github.io/romd-dat-data passed using `trust/1.root.json`
and a new disposable cache. This verifies the current signature chain and
artifact bytes; it is not a retained-cache rollback test.

| Verified property | Value |
| --- | --- |
| Signed publication | 3021 |
| Catalog | `no-intro/snes/standard` |
| Health | `healthy` |
| Last successful acquisition | 2026-09-17 09:38:51 UTC |
| Last document change | 2026-09-14 19:59:22 UTC |
| DAT version | `20260913-204915` |
| Entries / files | 4,359 / 4,359 |
| Nodump / baddump | 28 / 18 |
| Document bytes | 1,911,777 |
| SHA-256 | `b4978fc6328e9f901c7ed62b32baf41485491d4e5642f13f8f5f4cb257c3b7ad` |
| Immutable data commit | `acfb8f8878be29d845526d687136bda42ae5a726` |

The authenticated snapshot includes the SNES change event. Compared with the
recorded initial version (`20260818-050713`, 4,358 entries), this is evidence of
a real document revision passing through publication. It does not establish
which individual records changed or prove activation of this revision in ROMD.
The original public database ID reconciliation was not repeated for this version.

The five most recent publication runs returned success, including
[run 35206102774](https://github.com/MoonlarkStudios/romd-dat-data/actions/runs/35206102774)
on September 17. Workflow success alone can include retained artifacts after
acquisition failure; the signed SNES health and timestamps above establish the
current acquisition result. Current PSX health is also healthy.

`mise run check` passed: race tests (some cached), vet, formatting, workflow
validation and all three CLI builds. The first sandboxed attempt could not
write the Go cache; the authorized rerun passed. No ROMD application suite or
browser acceptance was rerun for this documentation closeout.

The live registry remains schema 2, while this checkout authors schema 3.
Before the next tooling deployment, verify ROMD's pinned reader against schema
3 and restore the existing live schema-2 state. Preserve trust-root identity,
metadata version progression, catalog identities and artifact history.

## Recommended order

These are proposed qualification batches, not declarations of support.
The [official No-Intro listing](https://datomatic.no-intro.org/index.php?page=download)
was checked on September 18 and lists the cartridge systems below. Listing
presence does not qualify acquisition, completeness or redistribution.
All proposed ROMD keys already exist in `definitions/system-keys.json`.

| Batch | Systems and ROMD keys | Reason and scope |
| --- | --- | --- |
| 1 | Game Boy (`gb`), Game Boy Color (`gbc`), Game Boy Advance (`gba`) | Reuse the Standard DAT path and cartridge ingestion. Keep GBA e-Reader, Video and Multiboot as separately scoped upstream catalogs. |
| 2 | NES (`nes`), Genesis (`genesis`), Nintendo 64 (`n64`) | Expand home-console coverage. Review NES header representation and N64 byte order explicitly; FDS and 64DD are separate scopes. |
| 3 | Master System (`sms`), Game Gear (`gg`), PC Engine/TurboGrafx-16 (`tg16`), 32X (`32x`) | Reuse qualified cartridge publication with explicit regional naming and provider bindings. |
| 4 | Saturn (`saturn`), Sega CD (`segacd`), Dreamcast (`dc`) | Extend the Redump disc path after cartridge batches; qualify track counts, multi-disc identity, archive size and complete public coverage individually. |

Start with Game Boy as a single canary, then complete batch 1. Keep PS2,
GameCube, Wii, DS/3DS, arcade and computer collections out of the first batches
until catalog sizes, representation distinctions and ingestion behavior are
measured. This ordering is an implementation recommendation based on reuse,
not a claim that deferred systems cannot work.

## Work required per batch

1. Record the exact upstream system ID, header name, representation and ROMD
   key from official sources. Do not guess IDs or create production definitions
   from display names. Preserve separate upstream variants.
2. Acquire one complete official document with the reviewed inclusion policy.
   Record version, bytes, hash, counts, nodump/baddump handling and exclusions.
   Reconcile eligible public membership where available, as for SNES; overview
   counters alone are not exact acceptance counts. Standing operator approval
   covers all Redump and No-Intro catalogs; preserve bytes and notices as
   recorded in `AGENTS.md`.
3. Set measured anomaly floors and prove the existing 16 MiB input / 32 MiB
   expanded limits suffice. If not, design and test bounded capacity changes
   before enrollment. Retain exact bytes and notices.
4. Add reviewed definitions and synthetic fixtures for the new identity and
   representation. Exercise wrong-header/system rejection, unchanged bytes,
   changed candidates, rate limiting and last-working-artifact retention. Run
   `mise run check` and the synthetic smoke publication.
5. Extend the existing daily workflow with explicit per-catalog enablement and
   pause handling. Reuse one provider adapter across a batch: preserve the
   five-second No-Intro request-start gap across catalog boundaries, and stop
   provider requests after 429/503 rather than starting a fresh process that
   forgets cooldown. Persist retry deadlines through the signed state. Avoid
   parallel provider jobs and automatic enrollment of all listed systems.
6. Prove compatibility of the schema-3 authoring revision with the deployed
   ROMD reader and schema-2 restoration before updating immutable workflow pins.
   Add focused multi-system fixtures to catch selection/routing leakage.
7. Publish the canary through the existing signed Git/Pages workflow. Verify
   the authenticated candidate, immutable asset and change event, then exercise
   ROMD discovery, review, approval, exact-byte active download, unchanged check,
   failure retention and retry. Recheck PSX and SNES discovery as regressions.
8. Record real acquisition evidence separately from synthetic changes. Observe
   the next scheduled healthy check before enabling the rest of the batch.
   Pause a failing catalog and retain its working artifact; a genuine upstream
   version change is recorded when available, not fabricated to close a gate.

All four expansion batches are now complete in acquisition, signed publication
and isolated ROMD catalog-ingestion scope. Further systems require a new scoped
qualification plan, including provider identity, public coverage, representation,
capacity and application acceptance. No further systems are enabled by this plan.
The operator authorized immediate rollout of batches 2–4 after qualification;
no scheduled observation is claimed for those batches.

## Batch 5: GameCube, PlayStation 2 and Wii

After the original four batches, the operator requested GameCube and selected PS2
and Wii as companion systems. [Qualification and rollout evidence](redump-gamecube-ps2-wii-qualification.md)
records exact public membership, alternate-export binary checksum comparison,
PS2 cue reconciliation, measured floors and signed publication. All three passed
isolated ROMD discovery, reviewed import, exact field/byte preservation (including
file sizes above 4 GiB), unchanged checks, outage retention and recovery. This
completes batch 5 in catalog scope; compressed game-image conversion and emulator
launch are separate work. No additional systems or rollout automation were enabled.


## Batch 6: DS, 3DS, PSP and PC Engine CD

The requested Dreamcast, Wii, Saturn and PC Engine cartridge subscriptions were
already live and passed fresh signed discovery and ROMD unchanged checks.
[Nintendo DS, Nintendo 3DS, PSP and PC Engine CD qualification](nds-3ds-psp-pcecd-qualification.md)
adds the four missing subscriptions. DS and 3DS explicitly select decrypted
representations. PSP and PC Engine CD use standard Redump disc catalogs; CD is a
separate `tgcd` subscription from cartridge `tg16`. All four passed signed
publication, reviewed import, full field/byte reconciliation, unchanged checks,
outage retention and recovery in isolated ROMD. The enabled real catalog count
is now 22. No other system variants or automation were enabled.


## Batch 7: Xbox, Xbox 360, PS3 and Wii U

The operator approved these four systems after batch 6. [Qualification and rollout
evidence](redump-xbox-ps3-wiiu-qualification.md) records complete standard exports,
Xbox 360's four observed upstream exclusions, cue and binary checks, and signed
publication 3037. All 26 real catalogs are healthy, including NES after upstream
export generation completed. The CLI selection cap was raised to 32 with
full-registry regression coverage. Wii U passed reviewed import and unchanged
checking in isolated ROMD; Xbox, Xbox 360 and PS3 acceptance is pending explicit
browser approval, so this batch is not yet closed in application acceptance scope.
