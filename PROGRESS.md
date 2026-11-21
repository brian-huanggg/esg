# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 2.8b handler 接上 service（`handler.go`、`main.go`、`handler_test.go`）
- 最新 tag：step-02-08a

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`（`dbtest.go` 也寫死），需學習者自行更新 ROADMAP 的決定
- ROADMAP 還沒有 2.8a / 2.8b 的拆分記錄，需學習者自行更新
- `origin` 是 `git@github.com:brian-huanggg/esg.git`（SSH），CI 已在這個 repo 跑通

## 最近的注意事項（最多 3 條）
- `account.Service.Register` 回 `*ValidationError`（`Fields []FieldError`，欄位名 `employeeId` / `name` / `password`）；handler 對應 422，`ErrEmployeeIDTaken` 對應 409
- `Account.CreatedAt` 已由 `Repository.Create` 填好；`api.Error.RequestId` 必填，用 `middleware.RequestIDFrom(ctx)` 取
- `account.Handler` 仍是 501 stub；`main.go` 要手動 `NewRepository` → `NewService` → Handler，`newHandler` 需要能接 `*sql.DB`
