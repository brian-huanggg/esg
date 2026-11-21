# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 3.8b 前端：Auth context（以 `GET /api/accounts/me` 判斷登入）、`<RequireAuth>` 路由守衛、登出按鈕（`DELETE /api/sessions/current`）。預計約 70 行
- 最新 tag：step-03-08a（CSRF middleware 與登入頁）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a / 3.4b / 3.5a / 3.5b / 3.6a / 3.6b / 3.7a / 3.7b / 3.8a / 3.8b 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 03-08a 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1

## 最近的注意事項（最多 3 條）
- 3.8b：Login 成功後目前只 `navigate('/')`，Auth context 完成後要改成更新登入狀態再導頁
- 3.8b：登出請求是 DELETE，會自動帶 `X-CSRF-Token`；成功後兩個 cookie 被清掉
- `auth` package 測試名稱不能重複：已有 `fakeAccounts`、`fakeSessions`、`fakeAccountByID`、`fakeLogin`、`fakeSessionDeleter`
