# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 2.8 handler 接上 service，後端完整可用（見 ROADMAP）
- 最新 tag：step-02-07

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`（`dbtest.go` 也寫死），需學習者自行更新 ROADMAP 的決定
- `origin` 是 `git@github.com:brian-huanggg/esg.git`（SSH），CI 已在這個 repo 跑通

## 最近的注意事項（最多 3 條）
- `account.Service.Register` 回 `*ValidationError`（`Fields []FieldError`，欄位名 `employeeId` / `name` / `password`）；2.8 handler 要對應成 422，`ErrEmployeeIDTaken` 對應 409
- `account.Handler.CreateAccount` 仍是 501 stub（借用 `default` response），2.8 換成真實作；`main.go` 要手動 `NewRepository` → `NewService` → Handler
- 密碼規則 8–72 bytes、姓名 trim 後 1–100 字元；`Repository.Create` 靠 key 名稱 `uq_accounts_employee_id` 判斷重複，migration 的 UNIQUE 名稱不能改
