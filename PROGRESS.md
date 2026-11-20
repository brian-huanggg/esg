# PROGRESS

> 每次整份覆寫，上限 20 行。哪些步驟做完了看 `git tag`，歷史看 `git log`。

## 目前
- 下一步：M2 第一步（ROADMAP 的 M2 目前只有大綱，開始前先和導師拆出第一步的定義）
- 最新 tag：step-01-10（M1 完成）

## 卡住的地方
- （無）

## 下次開始前要先處理
- `git tag` 在本機和 origin 都是空的，`step-01-01`～`step-01-09` 的 tag 不存在；要決定是否補打
- `origin` 是 `git@github.com:brian-huanggg/esg.git`（SSH），CI 已在這個 repo 跑通

## 最近的注意事項（最多 3 條）
- 寫完測試後故意把預期值改錯，確認它**會失敗**（1.9 抓到兩次假的綠燈；1.10 用在 CI 紅叉驗收）
- Makefile 要同時在 macOS 與 Ubuntu 跑：避免 `sed -i`，recipe 縮排用 Tab
- `ContextHandler.WithAttrs` / `WithGroup` 沒有測試，第一次用 `logger.With(...)` 時要補
