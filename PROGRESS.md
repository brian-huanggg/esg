# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 2.4 testcontainers：整合測試底座（見 ROADMAP）
- 最新 tag：step-02-03

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新 ROADMAP 的決定
- `origin` 是 `git@github.com:brian-huanggg/esg.git`（SSH），CI 已在這個 repo 跑通

## 最近的注意事項（最多 3 條）
- `.env` 密碼只用英數字（`openssl rand -hex 16`）；`MARIADB_*` 只在第一次建資料目錄時生效，改密碼要 `down -v`
- `make migrate-up` / `migrate-down` 需要 `.env` 和 DB 在跑（`make db-up`）；down 要寫 `down 1`
- 能交給 DB 約束的規則就交給 DB：check-then-act 擋不住併發，UNIQUE 才擋得住
