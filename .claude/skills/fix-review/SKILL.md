---
name: fix-review
description: 現在のPRのレビュー指摘を取得して自動修正する
argument-hint:
allowed-tools: Read, Write, Edit, Glob, Grep, Bash
---

# レビュー指摘の自動修正

現在のブランチに関連するPRのレビューコメントを取得し、指摘内容に基づいて修正を行います。

## 1. PR特定

1. `gh pr view --json number,title,url` で現在ブランチのPRを特定する
2. PRが存在しない場合は「このブランチにはPRがありません」と報告して終了

## 2. レビューコメント取得

1. `gh pr view --json reviews,comments` でレビューコメントを取得する
2. `gh api repos/{owner}/{repo}/pulls/{number}/comments` でインラインコメントも取得する
3. 未解決の指摘を一覧化する

## 3. 指摘分析

1. 各指摘の内容を分析し、以下に分類する:
   - **自動修正可能**: コード変更で対応できるもの
   - **要確認**: 仕様判断が必要なもの（ユーザーに確認）
2. 要確認の指摘がある場合はユーザーに報告して判断を仰ぐ

## 4. 修正実装

1. 自動修正可能な指摘に対して修正を実装する
2. 既存のコード規約・パターンに従う

## 5. コミット・プッシュ

1. 変更があれば以下を実行する:
   ```
   git add -A
   git commit -m "fix: レビュー指摘の修正

   Co-Authored-By: Claude <noreply@anthropic.com>"
   git push
   ```

## 6. 修正報告

以下を報告する:
- 修正した指摘の一覧
- 要確認としてスキップした指摘（あれば）
- PRのURL
