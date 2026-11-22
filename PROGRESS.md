# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：4.8 platform-gitops：`apps/esg/backend/` chart（Deployment、Service、PreSync migrate Job + `wait-for-db`），先用 `helm template` 檢查
- 最新 tag：step-04-07（Secret `esg-db` 由 ESO 同步，6 個 key；MariaDB NetworkPolicy 已放行 `esg`）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a 到 3.9c 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 04-07 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1
- platform-gitops（`~/Workspace/tools/gcp-vm-playground/gitops`，SSH clone）只限 `apps/esg/**`、`registry/esg-*`、`apps/mariadb/policies/networkpolicy.yaml`，改之前先說明；commit / push 由學習者自己做
- platform-gitops 的 `.gitignore` 已改成只擋根目錄 `/secrets/*`；新檔若在 `git status` 消失，用 `git check-ignore -v` 查

## 最近的注意事項（最多 3 條）
- 4.8：chart 用 `envFrom: secretRef: esg-db`；migrate Job 讀 `DATABASE_URL`，backend 讀 `DB_*`
- 4.9 提醒：nginx 啟動會寫 `/etc/nginx/conf.d/`，chart 開 `readOnlyRootFilesystem` 時要替它和 `/tmp` 掛 `emptyDir`
- cookie 是 `Secure: true`，本機 http 驗收要用 Chrome / Firefox（Safari 不收）
