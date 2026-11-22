# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：4.6 平台端（學習者以 Vault admin 執行）：`k3s-jwt-auth.sh <pubkeys.pem> app esg` 把 `eso-esg` 改成 `app`；`vault kv put secret/k3s/esg/db username=esg password=-`；在共用 MariaDB 建 database `esg` 與帳號（platform-gitops 新增 `scripts/create-mariadb-esg.sh`）
- 最新 tag：step-04-05（`make push` 推三個 amd64 image 到 Harbor，tag `3944807`）

## 卡住的地方
- （無）

## 下次開始前要先處理
- ROADMAP 和 PRD 仍寫 MariaDB 10.11 / 10.5；實際用 `mariadb:11.4`，需學習者自行更新
- ROADMAP 還沒有 2.8a / 2.8b / 3.4a 到 3.9c 的拆分記錄，需學習者自行更新
- `origin` 是 SSH；step-02-08a 到 04-05 都還沒 push（push 會觸發 notes 網站部署，先問學習者）
- `make db-exec` 免密碼的 Makefile 修改（`MYSQL_PWD`）學習者尚未套用，見 notes/3.1
- 4.6 起要改 platform-gitops（`~/Workspace/tools/gcp-vm-playground/gitops`，SSH clone）；只限 `apps/esg/**`、`registry/esg-*`，改之前先說明

## 最近的注意事項（最多 3 條）
- `make push` 要在 commit 後跑（工作區不乾淨會拒絕）；先 `docker login harbor.pf.internal`（見 docs/Platform.md §6.1）
- 4.9 提醒：nginx 啟動會寫 `/etc/nginx/conf.d/`，chart 開 `readOnlyRootFilesystem` 時要替它和 `/tmp` 掛 `emptyDir`
- cookie 是 `Secure: true`，本機 http 驗收要用 Chrome / Firefox（Safari 不收）
