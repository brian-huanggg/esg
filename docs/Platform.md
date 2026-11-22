# Platform

這份文件說明 ESG 專案部署到的實驗室平台：有哪些服務、怎麼連、誰負責哪一段、要照什麼慣例寫。
給 Agent 和學習者當作 M4 之後的上下文。這裡只放公開資訊，**沒有任何機密**；
機密一律放在 Vault，用 `vault kv put <path> field=-`（從 stdin 讀）寫入。

> 平台的權威文件在另外兩個 repo（見「三個 repo」）。這份文件和它們衝突時，以它們為準，並回報差異。

## 1. 三個 repo，各管一段

```
ESG repo（這個 repo）          platform-gitops                       devtools
程式、Dockerfile、compose  ->  Helm chart、values-dev.yaml 的 tag ->  VM、網路、Vault policy/role、
make push 推 image 到 Harbor    ExternalSecret、NetworkPolicy、        Harbor project 與 robot
                                registry/ 的 ArgoCD Application
```

| repo | 位置 | 負責 | ESG 的 Agent 可以改嗎 |
|---|---|---|---|
| ESG | 這個 repo（GitHub） | 程式碼、三個 image、本機 compose | 可以（依 `CLAUDE.md`） |
| platform-gitops | GitLab `simon-dev1/platform-gitops`，本機 `~/Workspace/tools/gcp-vm-playground/gitops`，只能用 SSH clone | 叢集裡的一切：chart、Secret 的來源、網路規則、ArgoCD 登記 | 只限 `apps/esg/**`、`registry/esg-*.yaml`，以及共用檔案裡 **esg 那一條**（例如 MariaDB 的 NetworkPolicy）。改之前先說明 |
| devtools | `~/Workspace/tools/gcp-vm-playground/devtools` | Terraform、Ansible、Vault 伺服器、Harbor 安裝 | 不可以。需要時由學習者在 devtools 的 session 處理 |

**交付邊界**：ESG repo 的工作在「image 推進 Harbor」結束。部署由 platform-gitops 決定。
ESG repo 裡**沒有** Helm chart。

## 2. 平台上有什麼

全部在 GCP，沒有外部 IP。Mac 經 Tailscale（gw）連進來，devbox 經 VPC peering 連。

| 服務 | 位址 | 說明 |
|---|---|---|
| k3s（單節點，amd64） | API `https://k3s.pf.internal:6443`；app `https://<name>.apps.pf.internal`（`*.apps.pf.internal -> 10.30.0.10`） | Ingress 是 Traefik。API 憑證由 k3s 自己的 CA 簽 |
| Harbor | `https://harbor.pf.internal`，project `lab`（私有） | 憑證由 lab CA 簽 |
| Vault | `https://vault.asia-east1-b.c.brhuangj-playground.internal:8200`（叢集內用 `https://10.10.0.10:8200`） | 獨立 VM，不在 k3s 裡；已 unseal（KMS 自動 unseal） |
| ArgoCD | 叢集內 `argocd` namespace，沒有 Ingress | app-of-apps：`root` 同步 platform-gitops 的 `registry/` |
| External Secrets Operator（ESO） | 叢集內 `external-secrets` namespace | 把 Vault 的值同步成 K8s Secret。**不是 VSO** |
| 共用 MariaDB | `mariadb.databases.svc.cluster.local:3306` | 所有專案共用一台，每個專案自己的 database 和帳號 |

- lab CA：devtools 的 `ansible/files/vault-ca.crt`（公開憑證）。Mac 已信任，並放在
  `~/.docker/certs.d/harbor.pf.internal/ca.crt`
- VM 每天 23:00 自動關機。gw 關機時 Mac 連不到平台，要先開機，或走 IAP 臨時通道（見 devtools `docs/architecture.md` §4.4）

## 3. 為什麼是 ESO + JWT，不是 VSO + Kubernetes auth

Vault 在另一個 GCP 專案，**Vault 連不到 k3s API**。Vault 的 `kubernetes` auth 需要回呼 k3s 做
TokenReview，所以不能用。平台改用：

- Vault `jwt` auth，掛在 `auth/k3s`。Vault 用 k3s 的公鑰在本機驗簽，不連叢集
- 每個 namespace 有 ServiceAccount `eso-vault`，用 audience `vault-k3s` 的短效 token 登入
- 每個 namespace 有自己的 `SecretStore`（不用 `ClusterSecretStore`，避免別的 namespace 讀到）
- ArgoCD 只部署 `ExternalSecret`，從不碰機密的值。Git 裡不會出現任何機密

ESG 的 role 是 `eso-esg`。目前是 `pull-only`（只能讀 Harbor pull token），4.6 會改成 `app` 類型。

## 4. ESG 在平台上的樣子（M4 完成後）

```
瀏覽器 --HTTPS--> Traefik（esg.apps.pf.internal）
                    |
                    v
              esg-frontend（Nginx :8080）--/api/--> esg-backend（:8080）--> mariadb.databases:3306
                                                                              ^
                    ArgoCD PreSync：esg-migrate Job（golang-migrate up）------+
```

- **同源**：跟 compose 一樣，Ingress 只接到 frontend，Nginx 再把 `/api/` 轉給 backend。
  所以 frontend 要設定 `BACKEND_URL=http://esg-backend:8080`
- **HTTPS**：Cookie 有 `Secure`，所以 Ingress 一定要走 `websecure`，並設定 `traefik.ingress.kubernetes.io/router.tls: "true"`。
  `*.apps.pf.internal` 的正式憑證還沒做，Traefik 會給它的預設憑證，瀏覽器第一次會出現警告，要手動接受

### 4.1 命名

| 項目 | 值 |
|---|---|
| namespace | `esg`（已存在） |
| image | `harbor.pf.internal/lab/esg-backend`、`lab/esg-frontend`、`lab/esg-migrate` |
| image tag | git short sha（`git rev-parse --short HEAD`），不用 `latest` |
| 平台架構 | `linux/amd64`。Apple Silicon 的 Mac 要用 `docker buildx build --platform linux/amd64` |
| Ingress host | `esg.apps.pf.internal` |
| database / 帳號 | `esg` / `esg` |
| Vault 路徑 | `secret/k3s/esg/db`，欄位 `username`、`password` |
| K8s Secret | `esg-db`（ESO 產生）、`harbor-pull`（ESO 產生，已存在） |

Vault 路徑選擇說明：其他專案用 `secret/k3s/mariadb/apps/<project>`，但那個類型（`project`）的 policy
讀不到 Harbor pull token。ESG 既要 DB 帳密也要拉 image，所以用 `app` 類型：可以讀 `secret/k3s/esg/*`
和 `secret/k3s/harbor-pull`。

### 4.2 backend 拿到的環境變數

`esg-db` Secret 由 ExternalSecret 的 template 組成，backend 用 `envFrom` 注入：

| key | 值 |
|---|---|
| `DB_HOST` | `mariadb.databases.svc.cluster.local` |
| `DB_PORT` | `3306` |
| `DB_NAME` | `esg` |
| `DB_USER` / `DB_PASSWORD` | 來自 Vault |
| `DATABASE_URL` | `mysql://<user>:<password>@tcp(mariadb.databases.svc.cluster.local:3306)/esg`，給 migrate Job 用 |

密碼只用英數字，因為它會直接放進 `DATABASE_URL`，特殊字元要另外 URL encode。

## 5. platform-gitops 裡 ESG 的檔案

```
apps/esg/
├── secrets/                     # registry/esg-secrets.yaml，sync-wave 3
│   ├── serviceaccount.yaml      # eso-vault（已存在）
│   ├── secretstore.yaml         # role eso-esg（已存在）
│   ├── vault-ca-configmap.yaml  # （已存在）
│   ├── externalsecret-harbor-pull.yaml  # （已存在）
│   └── externalsecret-db.yaml   # 4.7 新增，寫法參考 apps/mall/secrets/externalsecret-db.yaml
├── backend/                     # 4.8 新增；registry/esg-backend.yaml，sync-wave 4
│   ├── Chart.yaml / values.yaml / values-dev.yaml
│   └── templates/ deployment.yaml、service.yaml、migrate-job.yaml
└── frontend/                    # 已存在；registry/esg-frontend.yaml，sync-wave 4
    └── 4.9 加 BACKEND_URL、開 Ingress（TLS）
```

共用檔案裡 ESG 要加的部分：
- `apps/mariadb/policies/networkpolicy.yaml`：在第一條規則的 `from` 加 `esg` 的 namespaceSelector，否則連不到 3306
- `apps/network-policies/deny-metadata-esg.yaml`：已存在（擋住 Pod 取得 VM 的 GCP token）

### 5.1 慣例（照抄既有寫法）

- chart 參考 `apps/authkit/app/`，它是同類型的 app（API + Web + migrate + Secure cookie）
- 安全設定：`runAsNonRoot`、`readOnlyRootFilesystem`、`capabilities.drop: [ALL]`、`seccompProfile: RuntimeDefault`、
  `automountServiceAccountToken: false`；每個容器都設定 `resources`（單節點，資源有限）
- `image.tag` 用 `required`，tag 空著就讓 render 失敗，不要默默拉 `latest`
- migrate Job：ArgoCD `PreSync` hook、`hook-delete-policy: BeforeHookCreation,HookSucceeded`。
  要加一個 `wait-for-db` initContainer：kube-router 要等 NetworkPolicy 把 Pod IP 加進允許清單，
  太早連線會被拒絕（見 authkit 的 `migrate-job.yaml` 註解）
- 檔名不要有 `secret` 字樣（`.gitignore` 會擋 `*secret*.yaml`），`externalsecret-*`、`secretstore` 除外
- commit 用 conventional commits，例如 `feat(esg): ...`

## 6. 誰做什麼

依 ESG `CLAUDE.md`，會改變狀態的指令（docker、kubectl、helm、vault、git push）由學習者執行。
學習者同時是平台管理員（Vault admin）。

| 動作 | 在哪裡做 | 誰 |
|---|---|---|
| build、tag、push image | Mac，這個 repo | 學習者 |
| 改 `apps/esg/**`、`registry/esg-*` | platform-gitops | Agent 改檔，學習者 commit / push（或照 ESG `CLAUDE.md` 的收尾規則） |
| Vault：role、policy、寫入帳密 | Mac（admin token）或 devbox | 學習者 |
| 在 MariaDB 建 database 與帳號 | devbox（有 kubectl 和 `devbox-vault`） | 學習者 |
| VM、防火牆、DNS、Harbor robot | devtools | devtools 的 session，ESG 這邊不處理 |

### 6.1 常用指令

```bash
# Harbor 登入（Mac，已用 vault login 取得 admin token；VAULT_ADDR、VAULT_CACERT 已設定）
vault kv get -field=token secret/devbox/harbor-push \
  | docker login harbor.pf.internal -u "$(vault kv get -field=username secret/devbox/harbor-push)" --password-stdin
# devbox 上可以直接用：harbor-login

# kubeconfig（Mac）：存成 0600 的本機檔，絕不 commit
umask 077; vault kv get -field=value secret/k3s/kubeconfig > ~/.kube/lab-k3s
export KUBECONFIG=~/.kube/lab-k3s
# devbox 上 ~/.kube/config 已自動放好

# 看 ArgoCD 狀態
kubectl -n argocd get applications esg-secrets esg-backend esg-frontend
kubectl -n argocd get application esg-backend -o jsonpath='{.status.sync.revision}'   # 確認 sync 的是新 commit
kubectl -n argocd annotate application esg-backend argocd.argoproj.io/refresh=normal --overwrite   # 立刻重新比對

# ArgoCD UI（沒有 Ingress）
kubectl -n argocd port-forward svc/argocd-server 8080:80
```

## 7. 一次部署的流程（GitOps 迴圈）

1. ESG repo：commit 程式改動
2. Mac：`make push` 產生三個 `harbor.pf.internal/lab/esg-*:<short sha>`，在 Harbor UI 確認 image 存在
3. platform-gitops：把 `apps/esg/backend/values-dev.yaml`、`apps/esg/frontend/values-dev.yaml` 的 tag 改成同一個 sha，commit、push 到 GitLab
4. ArgoCD 自動 sync（`root` 是 automated + prune + selfHeal）：先跑 PreSync 的 migrate Job，再 rollout
5. 瀏覽器開 `https://esg.apps.pf.internal` 確認

ArgoCD 不會 build image，也不會監看 Harbor。**只有 Git 裡的 tag 改變才會觸發部署**。

## 8. 常見問題

| 現象 | 原因 |
|---|---|
| Pod `CreateContainerConfigError` | `esg-db` 還沒產生。看 `kubectl -n esg get externalsecret`；ESO 產生後 Pod 會自動重試 |
| ExternalSecret `SecretSyncedError`，訊息有 `permission denied` | `eso-esg` 的 policy 讀不到該路徑（例如還是 `pull-only`） |
| `ImagePullBackOff` | tag 打錯、沒推上去，或 image 是 arm64 |
| migrate Job 一直 `connection refused` | MariaDB NetworkPolicy 沒放行 `esg`，或沒有 `wait-for-db` |
| 登入成功但重新整理又變成未登入 | 走了 http，或沒接受憑證警告，瀏覽器沒有存 `Secure` cookie |
| ArgoCD 顯示 Synced 但畫面沒變 | 它比對的還是舊 commit。看 `.status.sync.revision`，再 refresh |
| 所有 ExternalSecret 同時失敗 | k3s 重建或簽發金鑰輪替後，Vault 裡的公鑰沒更新（平台端處理） |
