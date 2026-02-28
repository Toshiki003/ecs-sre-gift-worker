---
name: analyze
description: リポジトリの実装状況を分析してレポートを生成
argument-hint:
allowed-tools: Read, Glob, Grep, Write
---

# 現状分析レポートの生成

以下の手順で実装状況を分析してください:

1. **テンプレート読み込み**: `claude-ext/docs/analysis-repo-template.md` を読む
2. **効率的な調査**:
   - `git ls-files` コマンドを使用して、Git管理下のファイルのみ構造を把握する（`node_modules`などを除外するため）
   - または `tree -I 'node_modules|dist|.git'` を使用する
   - 必要な箇所のみ `grep` で詳細確認
3. **レポート出力**:
   - パス: `claude-ext/prompts/outputs/analysis-{{YYYYMMDD-HHmmss}}.md`
   - テンプレートの構造を維持

**注意**: コンテキスト節約のため、無関係なファイルの読み込みは避けてください。
