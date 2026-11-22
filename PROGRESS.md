# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：4.5 `make push`：tag 用 git short sha，推到 Harbor `lab`；`imagePullSecret`；驗收：Harbor 看得到三個 image
- 最新 tag：step-04-04（compose 全服務；檢查點 1 通過，Chrome 註冊、登入、登出成功）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a 到 3.9c 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 04-04 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1

## 最近的注意事項（最多 3 條）
- compose 用法：`make compose-up` / `make compose-down`（`db-down` 已改名）；`make db-up` 只起 `db`；前端對外 `127.0.0.1:8081`
- cookie 是 `Secure: true`，本機 http 驗收要用 Chrome / Firefox（Safari 不收）
- 三個 image 的本機 build 名稱是 `deploy-backend`、`deploy-frontend`、`deploy-migrate`；4.5 要重新 tag 成 Harbor 路徑；Harbor 位址與帳號要先問學習者
