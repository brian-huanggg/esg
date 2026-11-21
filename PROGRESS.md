# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：4.4 compose 全服務（migrate 跑完才啟動 backend）；驗收（檢查點 1）：瀏覽器註冊、登入、登出
- 最新 tag：step-04-03（`frontend/Dockerfile`、`nginx.conf`；經 Nginx 打 `/api/health` 回 200）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a 到 3.9c 的拆分記錄，需學習者自行更新；`ROADMAP.md` 目前有未 commit 的修改，是學習者自己的
- `origin` 是 SSH；step-02-08a 到 04-03 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1

## 最近的注意事項（最多 3 條）
- 前端 image 的後端位址用 `BACKEND_URL`（envsubst template，預設 `host.docker.internal:8080`）；compose 要覆寫成 `http://backend:8080`，且 backend 要先起來，否則 Nginx 報 `host not found in upstream`
- 容器連主機 DB 用 `host.docker.internal`（OrbStack 可用）；distroless 沒有 shell/curl，health 要從主機打
- migrate image 的 ENTRYPOINT 固定 `-path`，DSN 由呼叫端用 args 傳入（exec 形式不展開 env）；基底 `migrate/migrate:v4.18.3`
