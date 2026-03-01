---
name: sync-main
description: mainブランチに切り替えてgit pullで最新化
argument-hint:
allowed-tools: Bash
---

# mainブランチの同期

以下を実行してください:

1. `git checkout main` でmainブランチに切り替える
2. `git pull` でリモートの最新を取得する
3. 結果を報告する（取得したコミット数、現在の状態）
