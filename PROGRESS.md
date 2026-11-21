# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 3.6 認證 middleware 與 `GET /api/accounts/me`（見 ROADMAP M3）
- 之後：3.5b 登出（排在 3.6 之後，用 middleware 放進 context 的 token）
- 最新 tag：step-03-05a（`POST /api/sessions` 登入 + Set-Cookie）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a / 3.4b / 3.5a / 3.5b 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 03-05a 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1

## 最近的注意事項（最多 3 條）
- 前端有改動時，commit 前先跑 `make lint-frontend` 和 `make format-frontend`
- 3.6：cookie 名 `session`（`auth.cookieName`）；token 存的是「Cookie 值字串」的 SHA-256，middleware 直接對 Cookie 值雜湊；`session.Repository` 有 `FindActive` / `Touch` / `Delete`
- 3.6：`GetCurrentAccount` 與 `DeleteCurrentSession` 目前仍是 501 stub（`TestStubsReturn501`）；錯誤 body 要填 `requestId`
