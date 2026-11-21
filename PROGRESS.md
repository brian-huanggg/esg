# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：M3 登入與 Session 的第一步。ROADMAP 的 M3 只有大綱，尚未編號，先與學習者拆出 3.1
- 最新 tag：step-02-10（M2 完成：驗證、409/422 欄位錯誤、Register 測試）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 02-10 都還沒 push（push 會觸發 notes 網站部署，先問學習者）

## 最近的注意事項（最多 3 條）
- 前端有改動時，commit 前先跑 `make lint-frontend` 和 `make format-frontend`
- 前端測試 mock fetch：`client.ts` 的 `api` 每次呼叫才查 `globalThis.fetch`，且有 `baseUrl`
- M3 前端要處理 401 導回登入頁；錯誤檢查順序見 CLAUDE.md 第 10 節
