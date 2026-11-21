# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 2.7 account service：驗證與 bcrypt（見 ROADMAP）
- 最新 tag：step-02-06

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`（`dbtest.go` 也寫死），需學習者自行更新 ROADMAP 的決定
- `origin` 是 `git@github.com:brian-huanggg/esg.git`（SSH），CI 已在這個 repo 跑通

## 最近的注意事項（最多 3 條）
- `account.Handler.CreateAccount` 目前是 501 stub（借用 `default` response），Step 2.7 之後換成真實作並手寫 422 驗證；密碼規則暫定 8–72 字元
- `Repository.Create` 靠 MariaDB error 1062 加上 key 名稱 `uq_accounts_employee_id` 判斷工號重複，migration 裡的 UNIQUE 名稱不能改
- 整合測試用 `dbtest.Start`（每個 package 的 `TestMain` 呼叫一次）加每個測試開頭 `dbtest.Reset`；需要 Docker
