# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：1.6B JSON log（`log/slog`、帶 request_id 的 handler、access log middleware）
- 最新 tag：step-01-05（1.6 拆成 A/B；1.6A 已 commit 未打 tag，1.6B 完成後打 `step-01-06`）

## 卡住的地方
- （無）

## 下次開始前要先處理
- （無）

## 最近的注意事項（最多 3 條）
- header 給使用者看、context 給內層程式用，是兩條獨立的路；context 一定要 `r.WithContext` 傳下去
- 1.6B 的洋蔥順序：RequestID 要包在 Logging 外面，log 才拿得到 id
- 改了 `api/openapi.yaml` 一定要跑 `make gen`：編譯器只讀 `.go`
