# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 2.5 規格先行：POST /api/accounts（見 ROADMAP）
- 最新 tag：step-02-04

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`（`dbtest.go` 也寫死），需學習者自行更新 ROADMAP 的決定
- `origin` 是 `git@github.com:brian-huanggg/esg.git`（SSH），CI 已在這個 repo 跑通

## 最近的注意事項（最多 3 條）
- 整合測試用 `dbtest.Start`（每個 package 的 `TestMain` 呼叫一次）加每個測試開頭 `dbtest.Reset`；需要 Docker
- `dbtest/` 已排除在覆蓋率計算之外（Makefile 的 sed）
- 部署到 k3s 後不執行 migration down；部署時的 migration 方式（建議 PreSync hook Job）留到 M11
