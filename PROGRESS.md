# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：1.9 前端測試與覆蓋率門檻
- 最新 tag：step-01-08

## 卡住的地方
- （無）

## 下次開始前要先處理
- `origin` 還指向 `playground` repo，push 之前要先確認

## 最近的注意事項（最多 3 條）
- `npm run dev` 不做型別檢查，驗收前端時要跑 `npm run build`（`tsc -b`）
- 看到 `Unexpected token '<'` ＝ 拿到 HTML 而不是 JSON：多半是路徑打錯或 proxy 沒生效
- 洋蔥順序寫錯（Logging 在 RequestID 外面），目前沒有測試抓得到，要靠看伺服器 log 驗收
