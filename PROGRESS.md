# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：4.11 GitOps 迴圈：改一行 → `make push` → 改 `apps/esg/*/values-dev.yaml` 的 tag → push → sync
- 最新 tag：step-04-10（`registry/esg-backend.yaml`、三個 tag 填 `8d85f2c`；platform-gitops commit `fb45d84` 已由學習者 push；檢查點 2 通過）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a 到 3.9c 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 04-10 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1；注意它連的是本機 compose DB，不是叢集 DB
- platform-gitops（`~/Workspace/tools/gcp-vm-playground/gitops`，SSH clone）只限 `apps/esg/**`、`registry/esg-*`、`apps/mariadb/policies/networkpolicy.yaml`，改之前先說明；commit / push 由學習者自己做
- platform-gitops 的 `.gitignore` 已改成只擋根目錄 `/secrets/*`；新檔若在 `git status` 消失，用 `git check-ignore -v` 查

## 最近的注意事項（最多 3 條）
- 4.11：ArgoCD 不監看 Harbor，只有 Git 裡的 tag 改變才會部署；刷新可用 `argocd.argoproj.io/refresh=hard` annotation
- backend 兩個 Pod 各重啟 1 次（原因未查，已穩定）；可用 `kubectl -n esg logs deploy/esg-backend --previous` 看
- cookie 是 `Secure: true`，本機 http 驗收要用 Chrome / Firefox（Safari 不收）
