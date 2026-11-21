# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 2.9 前端 router 與註冊表單
- 最新 tag：step-02-08b（後端 `POST /api/accounts` 已完整可用）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`（`dbtest.go` 也寫死），需學習者自行更新 ROADMAP 的決定
- ROADMAP 還沒有 2.8a / 2.8b 的拆分記錄，需學習者自行更新
- `origin` 是 `git@github.com:brian-huanggg/esg.git`（SSH）；step-02-08a、step-02-08b 都還沒 push

## 最近的注意事項（最多 3 條）
- API 前綴維持 `/api`（不加 `/v1`）：前後端同源、同 repo 一起部署，沒有外部 client
- 錯誤回應格式：`{code, message, requestId, fields?}`；碼為 `VALIDATION_FAILED`（422）、`EMPLOYEE_ID_TAKEN`（409）、`INTERNAL`（500）；欄位名 `employeeId` / `name` / `password`
- 201 回應是 `{id, employeeId, name, createdAt}`，沒有密碼；前端 2.9 的表單欄位要對應這份 spec
