# ROADMAP

依照 `docs/PRD.md`，從零做出 ESG Idea Contest 平台。

## 原則

- **每一步都能跑**：做完任何一步，專案都能 build、test，也能啟動
- **每個里程碑前後端對應**：步驟只動前端或後端其中一邊，但每個里程碑結束時，
  一定能在瀏覽器完整操作一個功能
- **行數預算**：前期每步 30–80 行，後期每步 100–200 行
- **越早做成本越低**：規格先行、JSON log、request-id、測試、覆蓋率、CI 都在 M1 就導入
- **基礎設施先簡化**：單機 MariaDB、用環境變數代替 Vault。M4 第一次部署到 k3s，
  之後每個里程碑都走一次 GitOps 更新。Galera 和 Vault 放到 M11
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
│   ├── compose.yaml         # 本機開發用
│   └── helm/esg/            # Helm chart（M4）
└── .github/workflows/       # GitHub Actions
```

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
| 部署 | Docker（distroless）、Helm、ArgoCD、k3s、Harbor |

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

## M2 註冊

> 完成時：在瀏覽器用工號、姓名、密碼註冊。資料庫裡的密碼是 bcrypt 雜湊。重複工號會看到錯誤。
>
> 已決定：MariaDB 用 11.4 LTS；重複工號回 409；新增套件 `go-sql-driver/mysql`、
> `golang-migrate`、`testcontainers-go`、`golang.org/x/crypto/bcrypt`、`react-router` 都已核准。

#### 2.1 compose 啟動 MariaDB
- 目標：`make db-up` 起一個 MariaDB 11.4，能用 client 連進去執行 `SELECT 1`
- 新概念：docker compose、volume（資料持久化）、healthcheck、`.env` 與 `.env.example`（secret 不進 git）
- 檔案：`deploy/compose.yaml`、`.env.example`、`Makefile`（`db-up` / `db-down`）
- 行數：約 30
- 驗收：`docker compose ps` 顯示 healthy；`SELECT 1` 回 1
- 對應：PRD「資料庫」

#### 2.2 Go 連線與 connection pool
- 目標：後端啟動時連上 DB 並 ping；DB 連不上就啟動失敗，log 是 JSON
- 新概念：`database/sql`、driver 與 `sql.Open` 的差別、pool 參數（MaxOpenConns、ConnMaxLifetime）、
  從環境變數讀設定、`parseTime` / `loc=UTC` DSN 參數
- 檔案：`backend/internal/shared/db.go`、修改 `main.go`、`go.mod`
- 行數：約 50
- 驗收：DB 開著時啟動成功；`make db-down` 後再啟動，看到錯誤 log 並結束
- 對應：後端 1.3

#### 2.3 migration：accounts 表
- 目標：`make migrate-up` 建出 `accounts`，`migrate-down` 能還原
- 新概念：migration 為什麼要版本化、up / down、BIGINT 主鍵加 UUID（雙 ID）、
  `UNIQUE`（工號）、`created_at` 用 DB 時間、CHECK 約束
- 檔案：`backend/migrations/000001_create_accounts.{up,down}.sql`、`Makefile`
- 行數：約 40
- 驗收：`SHOW CREATE TABLE accounts` 看到 UNIQUE；重複 insert 同工號被 DB 擋下
- 對應：PRD「雙 ID 策略」「資料庫約束優先」

#### 2.4 testcontainers：整合測試底座
- 目標：測試自動啟動真的 MariaDB、跑 migration，並提供乾淨的 DB 給每個測試
- 新概念：testcontainers-go、`TestMain`、為什麼不 mock DB、測試間的資料隔離
- 檔案：`backend/internal/shared/dbtest/dbtest.go`、一個驗證 UNIQUE 的測試
- 行數：約 70
- 驗收：`make test` 通過（需要 Docker）；CI 的 `ubuntu-latest` 內建 Docker，覆蓋率門檻仍通過
- 對應：測試 3.4

#### 2.5 規格先行：POST /api/accounts
- 目標：`openapi.yaml` 有 `POST /api/accounts`、`Account` schema 和統一的錯誤格式；
  前後端型別重新產生；後端先回 501 stub
- 新概念：request / response schema、`201` / `409` / `422`、錯誤 body 設計
  （`code`、`message`、欄位錯誤清單、request-id）、strict server 要求補齊介面
- 檔案：`api/openapi.yaml`、`backend/internal/account/handler.go`（stub）、`main.go`
- 行數：約 60
- 驗收：`make gen` 後 build 成功；CI 的 gen 一致性檢查通過
- 對應：PRD 錯誤順序表第 7 項（422）

#### 2.6 account repository
- 目標：`Create(ctx, account)` 寫入 DB；工號重複回傳 `ErrEmployeeIDTaken`
- 新概念：repository 層的職責、參數化 SQL（防 SQL injection）、辨識 MySQL error 1062、UUID 在哪產生
- 檔案：`backend/internal/account/repository.go`、`repository_test.go`
- 行數：約 80
- 驗收：整合測試覆蓋成功與重複；覆蓋率門檻通過
- 對應：專案不變量「SQL 一律參數化」

#### 2.7 account service：驗證與 bcrypt
- 目標：驗證輸入（工號 6 碼數字、姓名非空且有長度上限、密碼最短長度），
  密碼用 bcrypt cost 10 雜湊後交給 repository
- 新概念：為什麼用 bcrypt 而不是 SHA-256、cost 的意義、bcrypt 72 byte 上限、
  service 與 handler 的分工、用 fake repository 做單元測試
- 檔案：`backend/internal/account/service.go`、`service_test.go`
- 行數：約 80
- 驗收：測試覆蓋各種無效輸入；存進去的字串以 `$2a$10$` 開頭
- 對應：PRD「密碼只存 Hash」；後端 2.2

#### 2.8 handler 接上 service，後端完整可用
- 目標：`curl -X POST /api/accounts` 成功回 201（不含密碼），重複工號回 409，
  無效輸入回 422，log 帶 request-id
- 新概念：service error 對應 HTTP status、回應不洩漏 hash、`main` 只負責組裝
- 檔案：修改 `handler.go`、`main.go`，新增 `handler_test.go`
- 行數：約 100
- 驗收：三種情況的 curl；DB 裡密碼欄是 bcrypt；`make test` 通過
- 對應：PRD 2.1 帳號

#### 2.9 前端：router 與註冊表單
- 目標：`/register` 頁有三個欄位，送出後呼叫 API，成功顯示成功畫面
- 新概念：react-router、受控表單（controlled input）、`openapi-fetch` 的型別化 POST、提交中狀態
- 檔案：`frontend/src/main.tsx`、`App.tsx`、`pages/Register.tsx`
- 行數：約 100
- 驗收：瀏覽器註冊成功；Network 面板看到 201
- 對應：前端領域

#### 2.10 前端：欄位驗證與伺服器錯誤，加測試
- 目標：前端先擋明顯錯誤；後端 409 / 422 的欄位錯誤顯示在對應欄位旁
- 新概念：前端驗證只是 UX，不能取代後端；解析錯誤格式；Testing Library 測表單
- 檔案：修改 `Register.tsx`，新增 `Register.test.tsx`
- 行數：約 100
- 驗收：重複工號顯示錯誤；`make test` 通過；M2 里程碑在瀏覽器完整走一遍
- 對應：PRD 2.1 帳號；前端 3.6

## M3 登入與 Session（大綱，約 9 步）

> 完成時：在瀏覽器登入、登出，重新整理後仍維持登入，閒置 30 分鐘後自動登出。

- `sessions` 表：只存 SHA-256；閒置 30 分鐘、最長 8 小時（用 DB 時間判斷）
- Cookie 屬性（HttpOnly、Secure、SameSite）
- 認證 middleware（401）、CSRF 防護
- 前端：登入頁、登入狀態、受保護的路由、401 時導回登入頁
- 對應：PRD「登入」；[Session 與 Cookie]、[CSRF]

## M4 第一次部署（大綱，約 9 步）

> 完成時：同一份程式跑在 k3s 上，改一行程式碼後，走一次完整的
> build → push Harbor → 更新 values → ArgoCD sync。

- 後端 Dockerfile：multi-stage + distroless
- 前端 Dockerfile：Nginx + SPA，Nginx 把 `/api` 轉給後端（同源）
- compose 一次啟動所有服務
- Helm chart：backend、frontend、MariaDB（單節點）、migration Job
- 本機 build，push 到 Harbor（`lab` 專案）
- ArgoCD Application 指向 repo 裡的 chart，練習一次 GitOps 更新
- Secret 先由學習者手動建立（不進 git），M11 再改成 Vault
- 對應：PRD「部署」；[distroless]、維運 3.2、3.3

## M5 提案 CRUD（大綱，約 12 步）

> 完成時：在瀏覽器建立、列表、編輯、刪除提案；兩個分頁同時編輯時，會看到衝突提示。

- `proposals` 表、同一 Leader 的提案名稱 `UNIQUE`
- 建立 / 讀取 / 整筆替換 / 刪除
- 分頁（上限 100）、篩選、排序（排序要穩定）
- ETag / If-Match：缺少 If-Match 回 428，版本不符回 412
- 前端：列表頁、表單、分頁元件、412 衝突時的處理
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
- Admin 從環境變數注入並同步（M11 改成 Vault）
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

> 完成時：CI 包含四種安全掃描；機密由 Vault 注入；資料庫換成 Galera。

- GitHub Actions：Semgrep、Trivy（image）、gitleaks（secret）、相依套件掃描（Mend 或替代方案）
- CI build image 並 push（方式見待決事項）；更新 Helm values，觸發 ArgoCD
- Vault：Admin 帳號與 DB 密碼的注入方式
- MariaDB Galera（注意 k3s 單節點的資源）
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
| ArgoCD 怎麼讀 repo | 確認 k3s 能連到 github.com；私有 repo 需要 deploy key 或 token | M4 |
| 相依套件掃描 | Mend App，或 dependency-review-action / osv-scanner / Trivy | M11 |
