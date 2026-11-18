# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：1.3 第一個測試與 Makefile
- 最新 tag：step-01-02

## 卡住的地方
- （無）

## 下次開始前要先處理
- 回答 0.2 自我檢查：.gitignore 規則對調的結果、如何把誤 add 的檔案移出暫存區

## 最近的注意事項（最多 3 條）
- Go 指令要在 backend/ 裡執行（go.mod 在那裡）；module path 是 github.com/brian-huanggg/esg-award-registration/backend
- go run 啟動的伺服器要在另一個終端機跑；用 ! 跑 curl 時加 -s 關掉進度條
- Go 版本由專案根目錄的 mise.toml 固定（go = "1.26"）