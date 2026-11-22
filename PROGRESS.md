# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：5.2 規格 `POST /api/proposals`、`GET /api/proposals/{id}`（回 `ETag`）、501 stub
- 最新 tag：step-05-01（本機 DB 已套用 000004；叢集 DB 要等 5.13 由 PreSync Job 套用）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- step-04-11、step-05-01 的 commit 尚未 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1
- platform-gitops（`~/Workspace/tools/gcp-vm-playground/gitops`）只限 `apps/esg/**`、`registry/esg-*`、`apps/mariadb/policies/networkpolicy.yaml`；commit / push 由學習者自己做

## 最近的注意事項（最多 3 條）
- collation 是 `utf8mb4_uca1400_ai_ci`：title 唯一性不分大小寫和重音，5.3 寫 1062 測試時可以利用
- `updated_at` 沒有 `ON UPDATE`，5.10 的 UPDATE 要自己 `SET updated_at = UTC_TIMESTAMP(3)`
- cookie 是 `Secure: true`，本機 http 驗收要用 Chrome / Firefox（Safari 不收）
