# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：1.6 JSON log 與 request-id middleware
- 最新 tag：step-01-05

## 卡住的地方
- （無）

## 下次開始前要先處理
- （無）

## 最近的注意事項（最多 3 條）
- 改了 `api/openapi.yaml` 一定要跑 `make gen`：編譯器只讀 `.go`，不跑就什麼都不會發現
- 規格每加一個 operation，`StrictServerInterface` 就多一個方法；之後要決定多個 handler 怎麼組合
- Go 指令要在 backend/ 裡執行；`make` 系列在專案根目錄執行（Makefile 用 Tab）
