# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 3.5 登入、登出 handler 與 Cookie（見 ROADMAP M3）
- 最新 tag：step-03-04b（`auth.Service.Login` + fake 單元測試）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a / 3.4b 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 03-04b 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1

## 最近的注意事項（最多 3 條）
- 前端有改動時，commit 前先跑 `make lint-frontend` 和 `make format-frontend`
- 3.5：`auth.NewService(accounts, sessions)` 還沒在 `main.go` 接線；`Login` 回明文 token、失敗回 `auth.ErrInvalidCredentials`；DB 錯誤是其他 error（回 500）
- 3.5：201 要填 `Set-Cookie`（cookie 名 `session`）；錯誤 body 的 `requestId` 要填；token 存的是「編碼後字串」的 SHA-256，3.6 middleware 直接對 Cookie 值雜湊
