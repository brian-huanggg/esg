# ROADMAP

依照 `docs/PRD.md`，從零做出 ESG Idea Contest 平台。

## 原則

- **每一步都能跑**：做完任何一步，專案都能 build、test，也能啟動
- **每個里程碑前後端對應**：步驟只動前端或後端其中一邊，但每個里程碑結束時，
  一定能在瀏覽器完整操作一個功能
- **行數預算**：前期每步 30–80 行，後期每步 100–200 行
- **越早做成本越低**：規格先行、JSON log、request-id、測試、覆蓋率、CI 都在 M1 就導入
- **基礎設施先簡化**：本機用 compose 的單機 MariaDB。M4 第一次部署到 k3s，使用平台的共用 MariaDB，
  並用 External Secrets Operator（ESO）把機密從 Vault 同步成 K8s Secret；之後每個里程碑都走一次 GitOps 更新。
  平台的環境與慣例見 [`docs/Platform.md`](docs/Platform.md)
- **細節逐步展開**：M0–M1 已經細到每一步。其餘里程碑先列大綱，做到前一個里程碑時再展開

## 進度怎麼看

- 完成的步驟：`git tag --sort=creatordate`（tag 名稱：`step-MM-SS`）
- 目前狀態：`PROGRESS.md`
- 這份文件**不放勾選框**，避免跟 git tag 不一致

## 目標資料夾結構

```
.
├── CLAUDE.md / ROADMAP.md / PROGRESS.md
├── Makefile                 # 所有常用指令的入口
├── api/
│   └── openapi.yaml         # API 規格，前後端共用
├── backend/
│   ├── go.mod
│   ├── cmd/api/main.go      # 程式進入點，只負責組裝
│   ├── migrations/          # golang-migrate 的 SQL 檔
│   └── internal/
│       ├── shared/          # log、request-id、DB、HTTP 錯誤等共用工具
│       ├── account/         # 註冊、帳號
│       ├── auth/            # 登入、Session、CSRF、鎖定
│       ├── proposal/        # 提案、團隊成員
│       ├── slides/          # 簡報驗證與上傳流程
│       ├── storages/        # MinIO 存取、pending_file_deletion
│       ├── event/           # 截止時間與狀態
│       ├── audit/           # 稽核紀錄
│       └── admin/           # 管理員功能
├── frontend/                # Vite + React + TS
├── deploy/
│   └── compose.yaml         # 本機開發用
├── docs/
│   ├── PRD.md
│   └── Platform.md          # 部署平台的上下文（M4 起）
└── .github/workflows/       # GitHub Actions
```

Helm chart 和 ArgoCD 設定**不在這個 repo**，在 platform-gitops 的 `apps/esg/`（見 `docs/Platform.md`）。

## 技術版本

| 項目 | 版本 / 選擇 |
|---|---|
| Go | 1.26 |
| Node | 24 LTS |
| React + TS | Vite 範本 |
| API 規格 | OpenAPI 3.0、oapi-codegen v2（後端：std-http + strict server）、openapi-typescript + openapi-fetch（前端） |
| 資料庫 | MariaDB（版本見待決事項） |
| 物件儲存 | MinIO |
| 測試 | `go test`、testcontainers-go、Vitest + Testing Library |
| CI | GitHub Actions |
| 部署 | Docker（distroless）、Helm、ArgoCD、k3s、Harbor、ESO + Vault JWT auth（平台既有，見 `docs/Platform.md`） |

## 每一步的格式

```
#### M.S 標題
- 目標：做完後能看到什麼
- 新概念：
- 檔案：新增 / 修改哪些檔案
- 行數：約多少行
- 驗收：用什麼指令確認，預期看到什麼
- 對應：PRD 的需求 
```

---

## M0 環境準備

> 完成時：工具都裝好，有一個空的 git repo，知道每個資料夾要放什麼。

#### 0.1 安裝並確認工具
- 目標：`go version`、`node -v`、`docker version`、`git --version`、`make -v` 都能正常輸出
- 新概念：什麼是 toolchain、PATH、為什麼要固定版本
- 檔案：無
- 行數：0（只有指令）
- 驗收：五個指令都印出版本號，Go 是 1.26.x，Node 是 24.x
- 對應：PRD「技術」

#### 0.2 初始化 repo 與資料夾骨架
- 目標：有 git repo、`.gitignore`，以及上面列出的頂層資料夾
- 新概念：git 的 working tree / staging / commit、`.gitignore` 為什麼重要（尤其是 secret）
- 檔案：`.gitignore`、各資料夾的 `.gitkeep`
- 行數：約 20
- 驗收：`git status` 乾淨，`git log` 有第一個 commit，`git tag` 有 `step-00-02`
- 對應：PRD「架構 / 模組化單體」

---

## M1 Walking Skeleton：前後端打通（已完成）

> 完成時：瀏覽器顯示後端 `/api/health` 的回應。後端 log 是帶 request-id 的 JSON，
> `make test` 檢查覆蓋率門檻，push 時 CI 自動跑測試。
> 每一步的細節看 `notes/1.X-*.md` 和 tag `step-01-01`～`step-01-10`。

| 步驟 | 內容 |
|---|---|
| 1.1 | 第一個 Go 程式（module、`main`） |
| 1.2 | 手寫 `/api/health` handler（`net/http`、JSON） |
| 1.3 | 第一個測試（`httptest`、table-driven）與 Makefile |
| 1.4 | 規格先行：`api/openapi.yaml`、oapi-codegen、`make gen` |
| 1.5 | 改用 generated strict server 介面實作 health |
| 1.6 | JSON log（`slog`）與 request-id middleware |
| 1.7 | Vite + React 專案骨架 |
| 1.8 | 前端 proxy 與 openapi-typescript / openapi-fetch 型別化 client |
| 1.9 | Vitest 前端測試；後端覆蓋率門檻 80%（排除 generated code） |
| 1.10 | GitHub Actions CI：gen 一致性、測試、覆蓋率、前端 build |

---

## M2 註冊（已完成）

> 完成時：在瀏覽器用工號、姓名、密碼註冊。資料庫裡的密碼是 bcrypt 雜湊。重複工號會看到錯誤。
> 已決定：MariaDB 11.4 LTS；重複工號回 409；`go-sql-driver/mysql`、`golang-migrate`、
> `testcontainers-go`、`x/crypto/bcrypt`、`react-router` 已核准。
> 每一步的細節看 `notes/2.X-*.md` 和 tag `step-02-01`～`step-02-10`（2.8 拆成 08a / 08b）。

| 步驟 | 內容 |
|---|---|
| 2.1 | compose 啟動 MariaDB（volume、healthcheck、`.env`） |
| 2.2 | Go 連線與 connection pool（`database/sql`、環境變數設定） |
| 2.3 | migration：`accounts` 表（雙 ID、UNIQUE 工號） |
| 2.4 | testcontainers 整合測試底座（`dbtest`） |
| 2.5 | 規格先行：`POST /api/accounts`、統一錯誤格式、501 stub |
| 2.6 | account repository（參數化 SQL、辨識 error 1062） |
| 2.7 | account service：輸入驗證與 bcrypt cost 10 |
| 2.8a | repository 回讀 `created_at`、修正密碼長度上限規格 |
| 2.8b | handler 接上 service：201 / 409 / 422 |
| 2.9 | 前端：router 與註冊表單 |
| 2.10 | 前端：欄位驗證、伺服器錯誤顯示、Testing Library 測試 |

---

## M3 登入與 Session（已完成）

> 完成時：在瀏覽器登入、登出，重新整理後仍維持登入，閒置 30 分鐘後自動登出。
> 鎖定與 IP 限速留到 M10；3.4 的 service 只留清楚的失敗出口。
> 每一步的細節看 `notes/3.X-*.md` 和 tag `step-03-01`～`step-03-09c`（3.8 拆成 08a / 08b，3.9 拆成 09a / 09b / 09c）。

| 步驟 | 內容 |
|---|---|
| 3.1 | `sessions` 表（只存 `token_hash`） |
| 3.2 | 規格：`POST /api/sessions`、`DELETE /api/sessions/current`、`GET /api/accounts/me`、`cookieAuth` |
| 3.3 | session repository（閒置 30 分鐘、絕對 8 小時，單一 SQL） |
| 3.4 | auth service（`crypto/rand`、SHA-256、dummy bcrypt） |
| 3.5 | 登入、登出 handler 與 Cookie |
| 3.6 | 認證 middleware 與 `/accounts/me` |
| 3.7 | CSRF（`X-CSRF-Token`） |
| 3.8a / 3.8b | 前端登入頁；auth context、受保護路由、登出按鈕 |
| 3.9a / 3.9b / 3.9c | 全域 401 處理；401 與 `RequireAuth` 測試；Login 測試與瀏覽器走查 |

---

## M4 第一次部署（已完成）

> 完成時：同一份程式跑在 k3s 上（`https://esg.apps.pf.internal`），改一行程式碼後走完
> build → push Harbor → 更新 platform-gitops 的 tag → ArgoCD sync。平台上下文見 [`docs/Platform.md`](docs/Platform.md)。
> 已決定：獨立 migrate image（ArgoCD PreSync Job + `wait-for-db`）；chart 和 Application 放在 platform-gitops 的
> `apps/esg/`（ESG repo 不放 chart）；用平台共用 MariaDB；機密由 ESO 經 Vault JWT auth（role `eso-esg`，`app` 類型 policy）同步成 K8s Secret。
> 每一步的細節看 `notes/4.X-*.md` 和 tag `step-04-01`～`step-04-11`。

| 步驟 | 內容 |
|---|---|
| 4.1 | 後端 Dockerfile（multi-stage、distroless nonroot） |
| 4.2 | migrate image（`migrate/migrate` + `backend/migrations`） |
| 4.3 | 前端 Dockerfile + `nginx.conf`（SPA fallback、`/api` 轉後端） |
| 4.4 | compose 全服務（migrate 完成才啟動 backend） |
| 4.5 | `make push`：git short sha tag、`linux/amd64`、推到 Harbor `lab/esg-*` |
| 4.6 | 平台端：Vault policy 改 `app`、`secret/k3s/esg/db`、共用 MariaDB 建 `esg` 帳號 |
| 4.7 | platform-gitops：ExternalSecret `esg-db`、MariaDB NetworkPolicy 放行 `esg` |
| 4.8 | platform-gitops：`apps/esg/backend/` chart（PreSync migrate Job） |
| 4.9 | platform-gitops：frontend `BACKEND_URL`、HTTPS Ingress |
| 4.10 | platform-gitops：`registry/esg-backend.yaml`，ArgoCD Synced / Healthy |
| 4.11 | GitOps 迴圈：改程式 → `make push` → 改 `values-dev.yaml` tag → sync |

- 之後每個里程碑的部署：`make push` → 改 platform-gitops `apps/esg/*/values-dev.yaml` 的 tag → push。
  新 migration 由 PreSync Job 自動套用；ArgoCD 不監看 Harbor，只有 Git 裡的 tag 改變才會部署

## M5 提案 CRUD（13 步，3 個垂直切片）

> 完成時：在瀏覽器建立、列表、編輯、刪除提案；兩個分頁同時編輯時，會看到衝突提示。
> 已決定：欄位是 `title`（1–100 字，同一 Leader 不可重複）和 `description`（最多 2000 字）；
> M5 的 POST 用 JSON，M7 加簡報時改成 multipart；分頁用 offset（`page` / `page_size` + `total`）；
> ETag 用 `version` 欄位。M5 只有 Leader 看得到自己的提案，擁有者條件直接寫進 SQL，查不到一律 404（M6 再擴充成 Leader 或 Member）。
> 每個切片結束都有一個瀏覽器檢查點。

| 步驟 | 邊 | 內容 | 驗收 |
|---|---|---|---|
| **A** | | **建立與檢視** | |
| 5.1 | 後 | migration `proposals`：雙 ID、`leader_account_id` FK、`version`、`UNIQUE(leader_account_id, title)` | `make migrate-up` / `down` 都成功 |
| 5.2 | 規格 | `POST /api/proposals`、`GET /api/proposals/{id}`（回 `ETag`）、501 stub | `make gen` 一致，`make test` 通過 |
| 5.3 | 後 | repository：Create、依擁有者取得、1062 轉成重複名稱錯誤 | 整合測試：別人的提案拿到 not found |
| 5.4 | 後 | service + handler：422 / 409 / 201 + `ETag` + `Location`；`main.go` 接線 | `curl` 建立回 201、帶 `ETag`；別人的 id 回 404 |
| 5.5 | 前 | 建立表單 + 詳情頁 | **檢查點 1**：瀏覽器建立提案並看到詳情 |
| **B** | | **列表** | |
| 5.6 | 規格+後 | `GET /api/proposals?page&page_size&sort&order&q`；repository：排序欄位白名單、`id` 當 tie-breaker、`LIKE` 跳脫、`COUNT` | 整合測試：同值排序穩定、`page_size` 上限 100 |
| 5.7 | 後 | handler：預設值、參數驗證 422、回 `items` + `total` | `curl` 分頁、排序、搜尋結果正確 |
| 5.8 | 前 | 列表頁 + 分頁元件，查詢條件放在 URL（`useSearchParams`） | **檢查點 2**：翻頁、排序、搜尋，重新整理後條件保留 |
| **C** | | **編輯、刪除與樂觀鎖** | |
| 5.9 | 規格 | `PUT` / `DELETE /api/proposals/{id}`：`If-Match` 必填，加上 428 / 412 | `make gen` 一致 |
| 5.10 | 後 | repository：`UPDATE / DELETE … WHERE id=? AND version=?`，影響 0 列時分出 404 和 412 | 整合測試：舊版本更新失敗 |
| 5.11 | 後 | handler：檢查順序 404 → 428 / 412 → 422 | `curl` 不帶 `If-Match` 回 428，舊 ETag 回 412 |
| 5.12 | 前 | 編輯頁、刪除按鈕、412 衝突提示（重新載入最新版） | **檢查點 3**：兩個分頁同時編輯，後送出的看到衝突 |
| 5.13 | 部署 | GitOps 更新（新 migration 由 PreSync Job 套用） | k3s 上 HTTPS 完成檢查點 1–3 |

- 重用：repository、1062、`dbtest`、錯誤格式都跟 M2 一樣；service 很薄，跟 handler 合成一步
- 陷阱（5.6）：`ORDER BY ?` 綁的是值，不是欄位名，所以排序只能用白名單 map 轉成欄位名
- 陷阱（5.11）：428 / 412 排在 422 之前，所以要在驗證 body 之前就先檢查 `If-Match`
- 對應：PRD「提案」；後端 2.2、2.5；前端 3.6

## M6 團隊與授權（大綱，約 8 步）

> 完成時：Leader 用工號加成員；成員只能檢視；沒有權限的人看到的是 404。

- 成員表與人數上限（Leader + 最多 4 位）
- 資源層級授權；沒有權限一律 404；dummy bcrypt 讓回應時間一致
- 固定的錯誤檢查順序（401 → CSRF → 404 → 403 → 428/412 → 403 → 422）
- 前端：成員管理、依角色顯示可用的操作
- 對應：PRD「團隊成員」「無權限一律顯示404」；後端 4.2；資安 2.4

## M7 簡報上傳與下載（大綱，約 12 步）

> 完成時：建立提案時一起上傳簡報；下載由後端串流；MinIO 掛掉時，刪除仍然成功，上傳和下載回 503。

- compose 加入 MinIO；storages 模組
- 檔案驗證：PDF / PPTX（檢查 magic bytes，不只看副檔名）、5 MiB 上限
- multipart：資料和檔案同一個請求、原子性
- 每次都寫新的 object key；被取代的檔案寫進 `pending_file_deletion`，跟資料異動在同一個 transaction
- 後端串流下載、預期的檔案名稱（Content-Disposition）
- 依賴故障：MinIO 掛掉時的行為與測試
- 前端：上傳元件、進度、下載
- 對應：PRD「簡報」「Fault Tolerance」；後端 4.1、3.5

## M8 截止時間與 Admin（大綱，約 10 步）

> 完成時：Admin 設定截止時間；截止後一般使用者無法寫入，但可以檢視和下載；Admin 不受截止時間限制。

- 狀態依 DB 時間即時推算（NOT OPEN / OPEN / CLOSED）
- state 403；Admin 略過 state 檢查
- Admin 帳號從 Vault 經 ESO 同步的 Secret 以環境變數注入，啟動時同步（Vault 路徑 `secret/k3s/esg/*`）
- Admin 後台：所有提案、提案者，以及替他人建立和編輯
- 前端：截止狀態顯示、Admin 頁面
- 對應：PRD「截止時間」「Admin」；後端 2.2

## M9 稽核（大綱，約 7 步）

> 完成時：註冊、登入、寫入、下載和被拒絕的嘗試都有紀錄；稽核寫入失敗時，整筆操作一起失敗。

- `audit_logs` 表；稽核和業務寫入在同一個 transaction
- 被拒絕的嘗試（401 / 403 / 404）也要記錄
- 稽核不變量測試：不存在「成功了但沒有稽核紀錄」的寫入
- Admin 查詢頁（分頁、篩選）
- 對應：PRD「Audit Log」；後端 4.3；測試 2.3

## M10 安全強化（大綱，約 7 步）

> 完成時：連續登入失敗五次會鎖定，同一 IP 太多失敗會被限速，回應帶完整的 security headers。

- 帳號鎖定：連續失敗 5 次鎖 15 分鐘（計數器存 DB）
- IP throttle：15 分鐘內 30 次失敗
- CSP（含 nonce 或 hash）與其他 security headers，設定在 Nginx 或後端
- 對應：PRD「登入」「加上CSP」；後端 2.6；資安 2.3、2.4

## M11 CI 安全掃描與正式基礎設施（大綱，約 10 步）

> 完成時：CI 包含四種安全掃描；CI 自動 push image 並更新 platform-gitops 的 tag。（Vault 與 ESO 已在 M4 完成）

- GitHub Actions：Semgrep、Trivy（image）、gitleaks（secret）、相依套件掃描（Mend 或替代方案）
- CI build image 並 push（方式見待決事項）；更新 platform-gitops 的 `apps/esg/*/values-dev.yaml`，觸發 ArgoCD
- MariaDB Galera：資料庫是平台共用的，換 Galera 屬於平台（platform-gitops / devtools）的工作，不在 ESG repo 做；
  到 M11 時再決定要不要做（注意 k3s 單節點的資源）
- 對應：PRD「CI」「部署」；資安 2.6；維運 1.3、2.3

## M12 效能驗收與加分功能（大綱，約 10 步）

- 一萬筆資料的 seed、壓測工具、RPS >= 100、p95 < 300ms；EXPLAIN 與索引調整
- 加分：統計儀表板、Excel 匯出、Audit hash chain、孤兒檔案排程
- 對應：PRD「基本需求」「加分功能」；後端 3.6、4.3、4.8；測試 4.1

---

## 待決事項

| 項目 | 選項 | 何時決定 |
|---|---|---|
| MariaDB 版本 | 已決定：11.4 LTS（PRD 的 10.5 已 EOL） | 已決定 |
| GitHub 認證 | 在 devbox 設定 GitHub 的 HTTPS credential | M1.10 前 |
| CI 怎麼 push image | GHCR，或在 devbox 架 self-hosted runner push 到 Harbor | M11 |
| ArgoCD 怎麼讀 repo | 已決定：ArgoCD 只讀 platform-gitops（GitLab），不讀 ESG repo | 已決定 |
| Vault 是否已安裝並 unseal | 已決定：Vault 是獨立 VM，已 unseal（KMS 自動 unseal） | 已決定 |
| chart 放哪裡 | 已決定：platform-gitops `apps/esg/`，跟平台其他 app 統一管理 | 已決定 |
| 資料庫 | 已決定：平台共用 MariaDB（`databases` namespace），不自己跑 StatefulSet | 已決定 |
| 相依套件掃描 | Mend App，或 dependency-review-action / osv-scanner / Trivy | M11 |
