# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 2.10 前端欄位驗證與伺服器錯誤，加測試
- 最新 tag：step-02-09（`/register` 表單可註冊，201 驗收通過）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a、02-08b、02-09 都還沒 push

## 最近的注意事項（最多 3 條）
- 前端有改動時，commit 前先跑 `make lint-frontend` 和 `make format-frontend`（已寫進 CLAUDE.md 第 6 節）
- `api.POST` 的 HTTP 錯誤回在 `error`，不是 exception；2.10 要解析 `{code, message, requestId, fields?}`
- 2.10 要處理：409 `EMPLOYEE_ID_TAKEN`、422 `VALIDATION_FAILED`，以及 fetch 本身失敗（會 throw）
