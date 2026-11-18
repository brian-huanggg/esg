# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：1.5 用產生的介面取代手寫 handler
- 最新 tag：step-01-04

## 卡住的地方
- （無）

## 下次開始前要先處理
- （無）

## 最近的注意事項（最多 3 條）
- 改了 `api/openapi.yaml` 要跑 `make gen`；`api.gen.go` 要 commit，但不要手改
- oapi-codegen v2.8.2：enum 常數叫 `Ok`（不是 `HealthStatusOk`），strict server 只用標準庫
- Go 指令要在 backend/ 裡執行；`make` 系列在專案根目錄執行（Makefile 用 Tab）
