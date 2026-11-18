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

## M1 Walking Skeleton：前後端打通

> 完成時：瀏覽器打開前端頁面，畫面顯示後端 `/api/health` 的回應。
> 後端的每一行 log 都是 JSON 並帶 request-id，`make test` 會檢查覆蓋率門檻，
> push 到 GitHub 時 CI 會自動跑測試。

#### 1.1 第一個 Go 程式
- 目標：`go run ./cmd/api` 印出一行文字
- 新概念：module、package、`main` 函式、`go.mod`
- 檔案：`backend/go.mod`、`backend/cmd/api/main.go`
- 行數：約 10
- 驗收：終端機印出 `hello esg`
- 對應：無（暖身）

#### 1.2 手寫一個 HTTP handler（暖身用，下一步會換掉）
- 目標：`curl localhost:8080/api/health` 回傳 `{"status":"ok"}`
- 新概念：`net/http`、handler、`ServeMux` 的路由寫法、JSON encoding、status code
- 檔案：修改 `main.go`
- 行數：約 30
- 驗收：curl 拿到 JSON，`curl -i` 看到 `200` 和 `Content-Type: application/json`
- 對應：後端 2.5

#### 1.3 第一個測試與 Makefile
- 目標：`make test` 跑過 health handler 的測試
- 新概念：`_test.go`、`httptest`、table-driven test、把 `main` 縮到最小（只呼叫 `run()`）
  讓程式容易測試、Makefile 的 target
- 檔案：`backend/cmd/api/main_test.go`、`Makefile`
- 行數：約 40
- 驗收：`make test` 顯示 PASS，`make cover` 印出覆蓋率
- 對應：PRD「後端測試覆蓋率 >= 80%」；測試 3.4

#### 1.4 規格先行：寫出 health 的 OpenAPI 規格
- 目標：`api/openapi.yaml` 描述 `GET /api/health`，並且能產生 Go 程式碼
- 新概念：OpenAPI 的 paths / schemas / responses、規格先行的意義、
  `go get -tool` 管理 oapi-codegen、`go generate`
- 檔案：`api/openapi.yaml`、`backend/oapi-codegen.yaml`、`backend/internal/shared/api/gen.go`
- 行數：約 50
- 驗收：`make gen` 產生 `api.gen.go`，build 成功
- 對應：PRD「oapi-codegen、建立 API-Spec」；後端 2.5「規格先行與程式碼生成」

#### 1.5 用產生的介面取代手寫 handler
- 目標：health 改成實作 strict server 介面，行為和 1.2 一樣
- 新概念：Go 的 interface、為什麼 generated code 不要手改、編譯器怎麼幫你檢查規格和實作一致
- 檔案：新增 `backend/internal/shared/health.go`，修改 `main.go` 和測試
- 行數：約 50
- 驗收：`make test` 通過，curl 的結果和 1.2 相同
- 對應：同 1.4

#### 1.6 JSON log 與 request-id middleware
- 目標：每個請求都有 request-id，回應 header 帶 `X-Request-Id`，每一行 log 都是帶 request-id 的 JSON
- 新概念：middleware 的洋蔥模型、`context.Context`、`log/slog`、
  request-id 從哪來（沿用 header 或新產生）
- 檔案：`backend/internal/shared/middleware/requestid.go`、`.../logging.go` 與對應測試
- 行數：約 80
- 驗收：curl 兩次，兩次的 request-id 不同；log 是 JSON 且有 `request_id` 欄位；測試通過
- 對應：PRD「Structured Logging」；後端 2.8「結構化日誌與 correlation ID」

#### 1.7 建立 React 專案並看懂結構
- 目標：`npm run dev` 打開 Vite 的預設頁面，然後清成只剩一個標題
- 新概念：npm 與 `package.json`、Vite、JSX/TSX、component、`index.html` → `main.tsx` → `App.tsx`
  的載入順序
- 檔案：Vite 產生的檔案，清理 `App.tsx`
- 行數：約 20（自己寫的部分）
- 驗收：瀏覽器顯示「ESG Idea Contest」
- 對應：前端領域

#### 1.8 前端呼叫後端：同源 proxy 與型別化的 client
- 目標：頁面顯示 `Backend: ok`，資料來自後端
- 新概念：同源政策（same-origin）、Vite dev proxy、openapi-typescript 從同一份 spec 產生型別、
  `useEffect` 與 `useState`、loading / error 狀態
- 檔案：`frontend/vite.config.ts`、`frontend/src/api/client.ts`、修改 `App.tsx`、`Makefile`（gen 目標）
- 行數：約 60
- 驗收：前後端都啟動後，畫面顯示 ok；把後端關掉，畫面顯示錯誤狀態
- 對應：PRD「前後端同源 不需要 CORS」；[同源政策]、[CORS]

#### 1.9 前端測試與覆蓋率門檻
- 目標：`make test` 同時跑前後端測試，後端覆蓋率低於 80% 時失敗
- 新概念：Vitest、Testing Library 的「用使用者的角度測」、mock fetch、
  覆蓋率門檻的計算方式（排除 generated code）
- 檔案：`frontend/src/App.test.tsx`、修改 `Makefile`
- 行數：約 50
- 驗收：測試通過；故意刪掉一個測試，看到門檻擋下來
- 對應：測試 3.4「在 CI 強制覆蓋率門檻」

#### 1.10 第一條 GitHub Actions CI
- 目標：push 到 GitHub 後，CI 自動跑 gen 檢查、後端測試、覆蓋率門檻、前端測試與 build
- 新概念：workflow / job / step、`actions/setup-go` 與 cache、
  為什麼 CI 要檢查「generated code 是否和 spec 一致」
- 檔案：`.github/workflows/ci.yaml`
- 行數：約 60
- 驗收：GitHub 的 Actions 頁面顯示綠勾；故意讓一個測試失敗，看到紅叉
- 對應：PRD「CI Pipeline」（安全掃描放在 M11）

---

## M2 註冊（大綱，約 8 步）

> 完成時：在瀏覽器用工號、姓名、密碼註冊，資料庫裡的密碼是 bcrypt 雜湊。

- compose 啟動 MariaDB；Go 連線與 connection pool
- golang-migrate：第一張 `accounts` 表（BIGINT 主鍵 + UUID、`UNIQUE` 工號、DB 時間欄位）
- testcontainers-go：整合測試啟動真的 MariaDB
- spec：`POST /api/accounts` 與錯誤格式（422 驗證錯誤）
- account 模組：repository（參數化 SQL）→ service（bcrypt cost 10）→ handler
- 前端：react-router、註冊表單、欄位驗證、顯示伺服器錯誤
- 對應：PRD「帳號」「密碼與 Session 只存 Hash」；後端 1.3、2.2

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
| MariaDB 版本 | PRD 指定 10.5（已 EOL），或改用 10.11 LTS | M2 開始前 |
| GitHub 認證 | 在 devbox 設定 GitHub 的 HTTPS credential | M1.10 前 |
| CI 怎麼 push image | GHCR，或在 devbox 架 self-hosted runner push 到 Harbor | M11 |
| ArgoCD 怎麼讀 repo | 確認 k3s 能連到 github.com；私有 repo 需要 deploy key 或 token | M4 |
| 相依套件掃描 | Mend App，或 dependency-review-action / osv-scanner / Trivy | M11 |
