# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：Step 3.9b 前端測試：Login、RequireAuth、401 全域攔截（Testing Library + mock fetch）；M3 在瀏覽器完整走一遍
- 最新 tag：step-03-09a（`api/client.ts` 的 `onResponse` 攔 401，`AuthProvider` 註冊 handler）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a / 3.4b / 3.5a / 3.5b / 3.6a / 3.6b / 3.7a / 3.7b / 3.8a / 3.8b / 3.9a / 3.9b 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 03-09a 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1

## 最近的注意事項（最多 3 條）
- 3.9b：`App.test.tsx` 的 `mockFetch` 已依網址分流（`/me` 固定回 200）；Login、RequireAuth 的測試可沿用同樣寫法
- 3.9b：測 401 攔截要用 `/me` 以外的受保護網址，因為 `/api/sessions` 和 `/api/accounts/me` 在例外清單裡
- `auth` package 測試名稱不能重複：已有 `fakeAccounts`、`fakeSessions`、`fakeAccountByID`、`fakeLogin`、`fakeSessionDeleter`
