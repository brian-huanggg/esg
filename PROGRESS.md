# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：1.10 第一條 GitHub Actions CI
- 最新 tag：step-01-09

## 卡住的地方
- （無）

## 下次開始前要先處理
- `origin` 還指向 `playground` repo，push 之前要先確認；1.10 前要設定 GitHub 的 HTTPS credential

## 最近的注意事項（最多 3 條）
- 寫完測試後故意把預期值改錯，確認它**會失敗**（1.9 抓到兩次假的綠燈）
- `ContextHandler.WithAttrs` / `WithGroup` 沒有測試，第一次用 `logger.With(...)` 時要補
- 洋蔥順序寫錯（Logging 在 RequestID 外面），目前沒有測試抓得到，要靠看伺服器 log 驗收
