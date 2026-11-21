# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 3.7 CSRF（登入時發 CSRF token；寫入類請求帶 `X-CSRF-Token`；失敗回 403；順序 401 → CSRF；用登出 `DELETE` 驗證）
- 最新 tag：step-03-05b（登出完成）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a / 3.4b / 3.5a / 3.5b / 3.6a / 3.6b 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 03-05b 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1
- 3.5b 的 curl 手動驗證（204、Set-Cookie Max-Age=0、再打 /me 回 401）學習者尚未回報

## 最近的注意事項（最多 3 條）
- 前端有改動時，commit 前先跑 `make lint-frontend` 和 `make format-frontend`
- `auth` package 測試名稱不能重複：已有 `fakeAccounts`、`fakeSessions`（service 測試）、`fakeAccountByID`、`fakeLogin`、`fakeSessionDeleter`（handler 測試）
- 3.7 會碰到 `Authenticate` middleware 的錯誤順序，記得 CSRF 檢查放在 401 之後
