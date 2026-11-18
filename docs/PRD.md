# ESG Idea Contest：產品需求文件（PRD）

## 1. 概述

公司內部 IT 員工的 ESG 點子競賽平台。員工組隊提交提案，每份提案附一份簡報，
管理員在截止時間前後管理所有提案。

### 1.1 角色

| 角色 | 說明 |
| --- | --- |
| Leader | 建立提案的人，一份提案只有一位 Leader |
| Member | 被 Leader 加入團隊的組員，每隊最多 4 位 |
| Admin | 管理員，帳號由 Vault 注入；Admin 同時也可以是 Leader 或 Member |

### 1.2 名詞

| 名詞 | 說明 |
| --- | --- |
| 團隊 | 1 位 Leader + 最多 4 位 Member，對應 1 份提案（Proposal） |
| 提案（Proposal） | 團隊提交的點子，一定附有 1 份簡報 |
| 簡報 | 提案的附件，格式 PDF 或 PPTX |
| 競賽狀態 | `NOT_OPEN`、`OPEN`、`CLOSED`，由資料庫時間和截止時間即時推算 |

## 2. 功能需求

### 2.1 帳號

- 用 6 碼工號、姓名和密碼自行註冊，不需要 Email
- Admin 帳號由 Vault 注入並同步，不能自行註冊

### 2.2 登入

- Session 存在 server-side 的 MariaDB，只存 Session ID 的 SHA-256 雜湊
- 閒置 30 分鐘逾時，最長 8 小時（絕對逾時）
- 同一帳號連續登入失敗 5 次，鎖定 15 分鐘
- 同一 IP 在 15 分鐘內失敗 30 次，會被限速（throttle）

### 2.3 提案

- 建立時用一個 multipart 請求同時送出資料和簡報，兩者具備原子性（全部成功或全部失敗）
- 編輯是整筆替換
- 寫入一定要帶 `ETag` / `If-Match` 做樂觀鎖：缺少 `If-Match` 回 428，版本不符回 412
- 同一位 Leader 的提案名稱不可以重複
- 列表支援分頁（每頁最多 100 筆）、篩選和排序

### 2.4 團隊成員

- Leader 用工號加入組員
- Member 只能檢視提案和下載簡報
- Member 不能自行退出團隊

### 2.5 簡報

- 只支援 PDF 和 PPTX
- 檔案上限 5 MiB
- 有預期的檔案名稱（規則見第 6 節待確認事項）
- 重新上傳會覆蓋舊檔
- 刪除提案時一併刪除檔案
- 下載由後端串流，不使用預簽名 URL（presigned URL），也不做分塊傳輸

### 2.6 截止時間

- 競賽狀態（`NOT_OPEN` / `OPEN` / `CLOSED`）依資料庫時間即時推算，不存成欄位
- 不在 `OPEN` 時，一般使用者的寫入一律擋下（state 403），檢視與下載不受影響
- Admin 不受截止時間限制：Admin 的請求略過 state 403 檢查

### 2.7 Admin

- 設定截止時間
- 替任何提案建立和編輯簡報
- 檢視和管理所有提案與提案者
- 查詢 Audit Log

### 2.8 Audit Log

- 記錄：註冊、登入、提案寫入、下載，以及被拒絕的嘗試
- 稽核寫入失敗時，整筆操作一起失敗
- 不允許出現沒有稽核紀錄的成功寫入

### 2.9 加分功能

- 統計儀表板
- Excel 匯出
- Audit hash chain（防竄改）
- 孤兒檔案清理排程

## 3. 基本需求（非功能需求）

### 3.1 品質

- 後端測試覆蓋率 >= 80%
- 建立 API-Spec（OpenAPI），規格先行，再用 oapi-codegen 產生程式碼

### 3.2 效能

- 取得提案列表的 API：在 10,000 筆資料下，RPS >= 100，p95 < 300 ms

### 3.3 可觀測性

- Structured Logging：每一行 log 都是 JSON，都帶 request-id

### 3.4 容錯（Fault Tolerance）

- MinIO 故障時：刪除提案仍然成功，上傳與下載回 503

### 3.5 安全

- 密碼與 Session 只存 Hash：密碼用 bcrypt cost 10，Session 用 SHA-256
- 所有 SQL 查詢一律參數化，所有列表都有分頁上限
- 無權限一律回 404，並用 dummy bcrypt 讓回應時間一致，避免帳號被枚舉
- 加上 CSP 和其他 security headers

### 3.6 CI

- 使用 Azure Pipelines，包含以下安全掃描：
  - WhiteSource（相依套件）
  - Semgrep（靜態分析）
  - Base image scan
  - Secret detection

### 3.7 部署

- 透過 Azure Release + ArgoCD（GitOps）+ Helm 部署到 Kubernetes
- 容器使用 distroless base image

## 4. 技術

| 類別 | 選擇 |
| --- | --- |
| 後端 | Go 1.26 |
| 前端 | React + TypeScript |
| 資料庫 | MariaDB Galera 10.5 |
| 物件儲存 | MinIO |
| 機密 | Vault |
| 執行環境 | Kubernetes |
| API 程式碼產生 | oapi-codegen |
| Migration | golang-migrate |
| 容器映像 | distroless |
| 整合測試 | testcontainers-go |

## 5. 架構

```text
Browser ──> Ingress ──> Frontend (Nginx + SPA) ──/api──> Backend (Go) ──> MariaDB（資料、Session、計數器、時間）
                                                                      ├─> MinIO（簡報檔案）
                                                                      └─> Vault（Admin 帳號、機密）
```

### 5.1 模組化單體

採用模組化單體（modular monolith），模組如下：

| 模組 | 負責 |
| --- | --- |
| `auth` | 登入、登出、Session、鎖定與限速 |
| `account` | 註冊、使用者資料 |
| `admin` | Admin 同步、截止時間設定、管理功能 |
| `proposal` | 提案與團隊成員 |
| `slides` | 簡報的驗證、上傳與下載 |
| `storage` | 物件儲存（MinIO）的存取 |
| `event` | 競賽狀態與截止時間推算 |
| `audit` | 稽核紀錄 |
| `shared` | 共用元件（錯誤、log、資料庫連線等） |

### 5.2 設計原則

- **資料庫約束優先**：能交給資料庫約束（`UNIQUE`、`FOREIGN KEY`、`CHECK`）的規則，就交給資料庫
- **雙 ID 策略**：對外一律用 UUID，對內一律用 BIGINT 主鍵加速檢索
- **唯一的時鐘是資料庫時間（UTC）**：應用程式不用自己的時間判斷業務邏輯
- **前後端同源**：Nginx 把 `/api` 轉給後端，不需要 CORS
- **錯誤檢查順序固定**：

  | 順序 | 檢查 | 狀態碼 |
  | --- | --- | --- |
  | 1 | 未登入 | 401 |
  | 2 | CSRF 驗證失敗 | 403 |
  | 3 | 資源不存在或無權限 | 404 |
  | 4 | 角色不符（role） | 403 |
  | 5 | 缺少 `If-Match` / 版本不符 | 428 / 412 |
  | 6 | 競賽狀態不允許（state，Admin 略過） | 403 |
  | 7 | 輸入驗證失敗 | 422 |

- **儲存一致性**：
  - 每次上傳都寫入新的 object key，不覆寫既有物件
  - 被取代或刪除的檔案記到 `pending_file_deletion`，和資料異動在同一個 transaction 裡
  - 之後再非同步清除實際檔案

## 6. 待確認事項

- 「預期的檔案名稱」的具體規則（下載時的檔名格式？上傳時是否檢查檔名？）
- MariaDB 10.5 已 EOL，是否改用 10.11 LTS（ROADMAP 已列為 M2 開始前決定）
