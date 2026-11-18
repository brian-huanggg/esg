# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：1.4 規格先行：寫出 health 的 OpenAPI 規格
- 最新 tag：step-01-03

## 卡住的地方
- （無）

## 下次開始前要先處理
- （無）

## 最近的注意事項（最多 3 條）
- `make test` / `make cover` 在專案根目錄執行；Makefile 指令行開頭必須是 Tab
- Go 指令要在 backend/ 裡執行（go.mod 在那裡）；module path 是 github.com/brian-huanggg/esg-award-registration/backend
- 覆蓋率檔案請命名為 coverage.out（.gitignore 只忽略這個名字）
