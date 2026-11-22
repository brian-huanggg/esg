# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：5.5 前端建立表單 + 詳情頁（檢查點 1：瀏覽器建立提案並看到詳情）
- 最新 tag：step-05-04（叢集 DB 的 000004 要等 5.13 由 PreSync Job 套用）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- step-04-11 到 step-05-04 的 commit 尚未 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1
- 5.4 的 409 與別人 id 回 404 只靠測試驗證，沒親手 `curl`
- platform-gitops（`~/Workspace/tools/gcp-vm-playground/gitops`）只限 `apps/esg/**`、`registry/esg-*`、`apps/mariadb/policies/networkpolicy.yaml`；commit / push 由學習者自己做

## 最近的注意事項（最多 3 條）
- 5.5 前端：POST 要帶 `X-CSRF-Token`（讀 `csrf_token` cookie）；建立成功後用 `Location` 導向詳情頁
- 錯誤 body 的 `fields[]` 可直接對應表單欄位的 422 訊息；409 的 `code` 是 `DUPLICATE_TITLE`
- `updated_at` 沒有 `ON UPDATE`，5.10 的 UPDATE 要自己 `SET updated_at = UTC_TIMESTAMP(3)`
