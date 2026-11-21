# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 3.8 前端：登入頁、登入狀態、受保護路由（`openapi-fetch` middleware 自動帶 `X-CSRF-Token`，值從 `csrf_token` cookie 讀）。可能超過 80 行，超過時先提議拆步
- 最新 tag：step-03-07b（CSRF middleware 與登出清兩個 cookie）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a / 3.4b / 3.5a / 3.5b / 3.6a / 3.6b / 3.7a / 3.7b 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 03-07b 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1

## 最近的注意事項（最多 3 條）
- 3.8：`csrf_token` cookie 沒有 HttpOnly，前端 JS 讀它並放進寫入類請求的 `X-CSRF-Token` header；GET 不需要
- 3.8：後端是 `Secure` cookie，本機用 HTTP 開發時瀏覽器對 `localhost` 例外，但 curl 手動測要自己帶值
- `auth` package 測試名稱不能重複：已有 `fakeAccounts`、`fakeSessions`、`fakeAccountByID`、`fakeLogin`、`fakeSessionDeleter`
