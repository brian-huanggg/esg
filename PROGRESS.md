# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 3.4b auth service（`auth.Service.Login`，見 ROADMAP M3 的 3.4）
- 最新 tag：step-03-04a（`account.Repository.FindByEmployeeID` + 整合測試）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a / 3.4b 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 03-04a 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1

## 最近的注意事項（最多 3 條）
- 前端有改動時，commit 前先跑 `make lint-frontend` 和 `make format-frontend`
- 3.4b：`session.Repository` 的 `Create(ctx, tokenHash []byte, accountID int64)`；`account.Repository.FindByEmployeeID` 回 `*Account`（含 `ID int64`、`PasswordHash`），查不到回 `account.ErrNotFound`
- 3.4b：service 用自己定義的小 interface 接 repo，測試手寫 fake；帳號不存在也要跑 dummy bcrypt
