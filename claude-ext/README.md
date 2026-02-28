# Claude Extension Kit

Claude Codeとの協働を効率化するためのプロジェクト拡張キットです。

## クイックスタート

```bash
# このテンプレートでは .claude/ と claude-ext/ は初期配置済み

# まずアプリ要件を記入
# claude-ext/docs/app-requirements.md
```

## 使い方（Claude Codeのスラッシュコマンド）

| コマンド | 説明 |
|---------|------|
| `/analyze` | リポジトリの実装状況を分析してレポート生成 |
| `/update-tasks` | 分析結果に基づいてタスクリスト更新（既定: `app-tasklist.md`） |
| `/implement-next` | 次の未着手タスクを実装してPR作成 |
| `/fix-review` | PRのレビュー指摘を取得して自動修正 |
| `/create-requirements` | ソースコードから要件定義書を生成（既定: `app-requirements.md`） |

## ディレクトリ構成

```
.claude/
├── CLAUDE.md           # プロジェクトルール（起動時に読む）
├── settings.json       # 権限・環境設定
├── skills/             # カスタムスラッシュコマンド
│   ├── analyze/
│   ├── update-tasks/
│   ├── implement-next/
│   ├── fix-review/
│   └── create-requirements/
└── rules/              # モジュール化されたルール
    ├── commit.md
    └── security.md

claude-ext/
├── docs/
│   ├── app-requirements.md        # アプリ実装要件（既定）
│   ├── app-tasklist.md            # アプリ実装タスク（既定）
│   ├── template-requirements.md   # テンプレート保守要件
│   ├── template-tasklist.md       # テンプレート保守タスク
│   ├── requirements.md            # 互換ポインタ
│   ├── tasklist.md                # 互換ポインタ
│   ├── decision-log.md            # 意思決定ログ
│   └── analysis-repo-template.md  # レポートテンプレート
└── prompts/
    └── outputs/                   # 生成レポート（Git除外）
```

## 開発フロー

1. **要件定義の作成**: `app-requirements.md` を編集（または `/create-requirements`）
2. **タスク更新**: `/analyze` → `/update-tasks` で `app-tasklist.md` を更新
3. **実装**: `/implement-next` で実装・PR作成
4. **レビュー対応**: `/fix-review` で指摘修正

テンプレート自体を保守する場合は `template-requirements.md` / `template-tasklist.md` を対象にします。

## カスタマイズ

### プロジェクト設定の変更

- 応答言語や運用方針は `.claude/CLAUDE.md` を編集
- 権限設定は `.claude/settings.json` を編集

### スキルの追加

`.claude/skills/<name>/SKILL.md` を作成すると、`/<name>` コマンドとして使用可能になります。
