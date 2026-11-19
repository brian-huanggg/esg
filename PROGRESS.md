# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：1.7 建立 React 專案並看懂結構
- 最新 tag：step-01-06

## 卡住的地方
- （無）

## 下次開始前要先處理
- （無）

## 最近的注意事項（最多 3 條）
- 洋蔥順序寫錯（Logging 在 RequestID 外面），目前沒有測試抓得到，要靠看伺服器 log 驗收
- 要讓 log 帶 request_id，一定要用 `logger.InfoContext(ctx, ...)`，不是 `Info`
