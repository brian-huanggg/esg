---
name: esg-anki-export
description: Export selected interview questions from the ESG Xmind maps into Anki through AnkiConnect, one-way and only after the learner approves each batch. Use when asked to push, sync, or review Anki cards for the ESG project.
---

# ESG Xmind to Anki

Source of truth is the Xmind map. Anki is a copy for spaced repetition. Sync goes one way: map to Anki.
Never edit the maps from Anki data.

## Constants

- Endpoint: `http://127.0.0.1:8765` (AnkiConnect add-on, code `2055492159`). Anki desktop must be running.
- Deck: `🍎 Software:Interview`. Put cards directly in it. Do not create subdecks.
  The single colon is part of the name; `::` would be a subdeck separator, so never rewrite the name.
- Tags on every card: `esg`, `map-<map name without the ESG- prefix>`, `qid-<slug>`.
- Note types: `Basic` (fields `Front`, `Back`) for questions, `Cloze` (field `Text`, `Back Extra`) for fill-in concepts.

## Candidate selection

- Candidates are topics marked `{label: anki}` in an `ESG-*` file. Read only the interview sheets: `3-Interview` in module files,
  `2-Interview` in cross-cutting files, all three sheets of `ESG-Q-Interview`. Find them with `xmind_list_sheets`,
  then `xmind_read_mindmap` per sheet.
- Question text = topic title. Answer = the topic's note or its child topics, flattened to short HTML-free text.
- Skip a candidate with no answer and report it.
- Cards come only from questions the learner marked or approved. Do not mark new ones for export on your own.

## Card format

- Front: the question.
- Back: the answer, then a last line `Source: <map title> / <topic path>`.
- `qid` slug: `<map>-<short-kebab-of-question>`, stable across runs so a re-run finds the same card.

## Procedure (follow in order)

1. **Check the connection.** `curl -s http://127.0.0.1:8765 -d '{"action":"version","version":6}'`.
   If it fails, tell the learner to open Anki and install AnkiConnect. Stop.
2. **Check the deck.** Call `deckNames` and confirm `🍎 Software:Interview` exists. If missing, ask before creating it.
3. **Build the batch** from the map. Look up existing cards with `findNotes` query `tag:qid-<slug>`.
   Classify each as new, unchanged, or changed (compare fields via `notesInfo`).
4. **Show the batch to the learner** as a list: new / changed / skipped, with question text.
   This is an outward-facing write to their Anki collection, so wait for an explicit yes for each batch.
5. **Push after approval.**
   - New: `addNotes` with `options: {"allowDuplicate": false, "duplicateScope": "deck"}`.
   - Changed: `updateNoteFields`. Do not delete cards or reset scheduling.
6. **Report** counts of added, updated, skipped, failed, with the error text for any failure.

## Rules

- Never delete notes or decks. Never change scheduling. If the learner wants removal, tell them to do it in Anki.
- Run `curl` with `-s` and a JSON body built from the map text. Escape quotes properly; pass untrusted text as data, not as shell syntax.
- Keep the collection private: send only card content to `127.0.0.1`. No external services.
- Do not include secrets, personal data, or real credentials in any card.
