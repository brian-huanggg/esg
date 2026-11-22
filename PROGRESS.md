# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：4.10 platform-gitops：`registry/esg-backend.yaml`（sync-wave 4）；push 後由 `root` 自動 sync；檢查點 2：瀏覽器 HTTPS 註冊、登入
- 最新 tag：step-04-09（frontend chart 的 `BACKEND_URL`、Ingress、`conf.d` emptyDir；platform-gitops 的變更由學習者自行 commit / push）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a 到 3.9c 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 04-09 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1
- platform-gitops（`~/Workspace/tools/gcp-vm-playground/gitops`，SSH clone）只限 `apps/esg/**`、`registry/esg-*`、`apps/mariadb/policies/networkpolicy.yaml`，改之前先說明；commit / push 由學習者自己做
- platform-gitops 的 `.gitignore` 已改成只擋根目錄 `/secrets/*`；新檔若在 `git status` 消失，用 `git check-ignore -v` 查
- 4.8 還沒驗證：`lab/mariadb-base` image 裡有沒有 `mariadb-admin`（wait-for-db 用）；4.10 Job 若卡在 initContainer 先查這個
- 4.9 的 frontend chart 還沒 commit / push 到 platform-gitops；4.10 前要先做，ArgoCD 才看得到

## 最近的注意事項（最多 3 條）
- 4.10：frontend 和 backend 的 `values-dev.yaml` tag 仍是空字串，ArgoCD sync 前要填 `make push` 推上去的 SHA，否則 `required` 會擋下來
- `/healthz` 在 `nginx.conf` 沒有專屬 location，目前靠 SPA fallback 回 `index.html` 的 200 通過探針
- cookie 是 `Secure: true`，本機 http 驗收要用 Chrome / Firefox（Safari 不收）
