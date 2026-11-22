# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：5.4 service + handler：422 / 409 / 201 + `ETag` + `Location`；`main.go` 接線
- 最新 tag：step-05-03（叢集 DB 的 000004 要等 5.13 由 PreSync Job 套用）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- step-04-11 到 step-05-03 的 commit 尚未 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1
- platform-gitops（`~/Workspace/tools/gcp-vm-playground/gitops`）只限 `apps/esg/**`、`registry/esg-*`、`apps/mariadb/policies/networkpolicy.yaml`；commit / push 由學習者自己做

## 最近的注意事項（最多 3 條）
- 5.4 換掉 `ProposalHandler` stub：201 要填 `ETag`（帶引號，如 `"1"`）和 `Location`；auth map 已登記兩個 operation
- 5.4 handler 用 `proposal.Repository`：`ErrDuplicateTitle` 轉 409，`ErrNotFound` 轉 404；leader id 從 session context 取
- `updated_at` 沒有 `ON UPDATE`，5.10 的 UPDATE 要自己 `SET updated_at = UTC_TIMESTAMP(3)`
