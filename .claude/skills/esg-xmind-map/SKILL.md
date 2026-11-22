---
name: esg-xmind-map
description: Build or update the ESG project's study mind maps in Xmind (module maps, cross-cutting maps, index map, interview map) from a fixed template. Use when asked to create, extend, or review an ESG Xmind map, or to sync a map after a roadmap step is finished.
---

# ESG Xmind maps

Purpose: help the learner understand the ESG codebase and design, and prepare for interviews
(architecture design, code understanding, software concepts). One Xmind file covers one functional module,
split into sheets by purpose: study, code, interview.

## Ground rules

- Read-only on the repo. Never edit `backend/`, `frontend/`, `api/`, `ROADMAP.md`. Writes go to Xmind only.
- Before any map call, run `xmind_whoami`, then `xmind_get_mindmap_guide` for the map kind. The guide owns the
  markdown syntax, `design` fields and marks. Do not restate it here.
- One Xmind file per map in the table below. A file holds several sheets (see "Sheet layout"). Create the file with
  `xmind_create_mindmap` (first sheet) and add the other sheets with `xmind_add_sheet`, all in the same turn.
  Write boundaries, summaries, labels and notes inline in each call; these marks work only on create and add_sheet.
  Name each added sheet with the `title` parameter of `xmind_add_sheet` (use the fixed names in "Sheet layout").
  No tool renames an existing sheet, and `xmind_create_mindmap` has no sheet name parameter: the first sheet is
  created as `Map 1`. Tell the learner to rename it to `1-Overview` by hand, every time.
  Run dependent Xmind write calls one at a time. Parallel writes to one file failed with
  "Duplicated ID found during deserialization process".
- One subject per file. Sheets are views of the same subject, so content must not repeat across sheets.
  Link files with `href` (`xmind_set_topic_attrs` set to the child file's `fileUrl`). Never copy the same content into two files.
- Linking to one sheet of a file is unverified. The pilot map must test `xmind_add_topic_link` and `xmind_get_topic`
  for sheet-level targets and record the result in the "Open questions" section below. Until then link to the file.
- Put all files in an Xmind folder named `ESG`. File title format: `ESG-<name>` from the table below.
- Leaves hold one concrete fact each: a file path with line, a status code, a number, a rule. No vague restatements.
- Code is referenced as `path:line`, never pasted. Verify each reference exists before writing it.
- Source of truth, in order: `docs/PRD.md`, `ROADMAP.md`, the code, then `notes/X.Y-*.md`. If code and docs disagree,
  show the difference in the map instead of picking silently.

## Map list (13)

| Title | Type | Covers | Roadmap |
|---|---|---|---|
| ESG-Index | index | Entry point. `href` to every map, status of each | all |
| ESG-Account-Auth | module | Register, login, session, CSRF, bcrypt, session hash; lockout, IP throttle, CSP | M2, M3, M10 |
| ESG-Proposal | module | Create, list, edit, delete, ETag, optimistic lock, 428/412, pagination, sort; team members and resource authorization | M5, M6 |
| ESG-Slides-Storage | module | Validation, multipart, MinIO, `pending_file_deletion`, dependency failure | M7 |
| ESG-Event-Admin | module | Deadline state, state 403, admin backend; dashboard and Excel export (bonus) | M8, M12 |
| ESG-Audit | module | Same-transaction audit, rejected attempts, invariant test; hash chain (bonus) | M9, M12 |
| ESG-Deploy-GitOps | module | Dockerfile, migrate image, Harbor, ESO and Vault, ArgoCD | M4 |
| ESG-Testing-CI | module | Test strategy, coverage gate, CI pipeline, security scans, image push, load test | M1, M11, M12 |
| ESG-X-Error-Model | cross-cutting | Error format, status codes, check order, 404 for no permission | all |
| ESG-X-DB-Design | cross-cutting | Dual ID, migrations, constraints, UTC clock, parameterized SQL, indexes, EXPLAIN | all |
| ESG-X-Shared-Infra | cross-cutting | `shared/`, slog, request-id, config, `main.go` wiring, spec-first and codegen, frontend proxy | M1 onward |
| ESG-X-Frontend-Arch | cross-cutting | Router, auth context, openapi-fetch, URL as query state, global 401 | M1 onward |
| ESG-Q-Interview | interview | Cross-module questions in three categories; Go vs Spring comparison under concepts | all |

Do not add maps without the learner's approval. Keep the total small.

## Status label

The first sheet's root carries one label: `planned`, `in-progress` or `done`.
Do not write `{label: ...}` on the root `#` line: it is not parsed and becomes literal title text.
Create the map without it, then call `xmind_set_topic_attrs` with `labels` on the root topic id, in a separate call
after the create call has returned. If the title was polluted, fix it with `xmind_edit_mindmap`.
Derive it from `git tag` and the roadmap table. Planned maps are built from PRD and ROADMAP only:
leave the "code index" branch empty except for "expected files", and say in the root note of sheet 1 that the file is a design draft.

## Sheet layout

| File type | Sheets |
|---|---|
| Module (7 files) | 3 sheets: `1-Overview`, `2-Flow-Code`, `3-Interview` |
| Cross-cutting (4 files) | 2 sheets: `1-Rules`, `2-Interview` |
| ESG-Q-Interview | 3 sheets: `1-Architecture`, `2-Code`, `3-Concepts` |
| ESG-Index | 1 sheet |

Sheet names are fixed so other agents and the Anki skill can find them with `xmind_list_sheets`.
Keep the interview questions of a file only on its interview sheet. That keeps answers away from the study sheets
and gives the Anki skill one place to read.

## Module file template

**Sheet `1-Overview`**
1. **Goal and PRD mapping**: what the user can do at the end, which PRD clauses apply, roadmap steps and tags.
2. **Design decisions and trade-offs**: each decision as a topic, with "why" and "alternative rejected" as children.
   Tie each to a project invariant when one applies.

**Sheet `2-Flow-Code`**
1. **Request flow**: the layers a request passes through, one topic per layer with its responsibility
   (spec, handler, service, repository, DB, frontend). Include the frontend path.
2. **Code index**: `path:line` per key symbol, grouped by layer. Planned maps list expected files only.
3. **Pitfalls**: concrete mistakes (wrong error order, `ORDER BY ?`, missing `updated_at`) with the fix.

**Sheet `3-Interview`**
1. Three branches: architecture design, code understanding, software concepts.
   Question text is the topic, the answer goes in a child or note. Mark candidates for Anki with `{label: anki}`.

## Cross-cutting file template

**Sheet `1-Rules`**
1. **Rule or principle**: the rule stated precisely.
2. **Where it appears**: one topic per module with `href` to that module's file.
3. **Trade-offs and alternatives**.
4. **Pitfalls**.

**Sheet `2-Interview`**: same format as a module's `3-Interview`.

## Index and interview files

- Index (one sheet): one branch per type (module, cross-cutting, interview), each leaf `href`s its file and shows the status label.
- Interview (three sheets): only questions that span modules, one sheet per category (architecture design, code understanding,
  software concepts). Module-specific questions stay in their module file. The Go vs Spring wiring comparison lives on `3-Concepts`.

## Question quality

- Ask "why" and "what if", not "what is". Example: "Why must If-Match be checked before body validation?"
- One question has one answer that fits in a note of a few lines.
- Prefer questions tied to a real decision in this repo over generic trivia.

## After a roadmap step finishes

1. Read the step's tag, `notes/X.Y-*.md`, and the changed files.
2. Update the affected module file: fill the code index on `2-Flow-Code`, flip the status label on `1-Overview`,
   and record any deviation from the design draft as a decision topic (sheet 1) and an interview question (sheet 3)
   ("planned X, shipped Y, why").
3. Update cross-cutting files only where a new rule appeared.
4. Update the index status. Use `xmind_edit_mindmap` with the target sheet; marks like `{boundary}` and `{summary}`
   only work on create and add_sheet.

## Workflow for a new file

1. Pick the title from the table. Read sources for that module.
2. Draft the full outline in chat only if the learner asks to review it; otherwise build directly.
3. Create the file in the `ESG` folder with sheet 1, add the remaining sheets with `xmind_add_sheet`,
   then set `href` links from the index file, all in the same turn.
4. Report the file's `fileUrl` and any section left empty on purpose.

## Open questions

- Can a topic link to a specific sheet of another file, or only to the file? (test in the pilot, then write the answer here)
  Pilot `ESG-Account-Auth` (2026-11-22) was built without a second file, so this is still untested. Test it when
  `ESG-Index` or the second map exists.
