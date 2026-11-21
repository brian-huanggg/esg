# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：M3 已完成。進入 M4（第一次部署），先讀 ROADMAP 的 M4 大綱，向學習者確認 4.1 的定義後再開始
- 最新 tag：step-03-09c（`Login.test.tsx`；M3 瀏覽器走查通過）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a 到 3.9c 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 03-09c 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1

## 最近的注意事項（最多 3 條）
- 前端測試慣例：`vi.stubGlobal('fetch', ...)` 依網址分流，每次回新的 `Response`；模組層級狀態要在 `afterEach` 清掉
- 首頁目前沒有受保護 API 呼叫，等有受保護頁面時，要在瀏覽器補驗「session 過期跳回登入」
- `auth` package 測試名稱不能重複：已有 `fakeAccounts`、`fakeSessions`、`fakeAccountByID`、`fakeLogin`、`fakeSessionDeleter`
