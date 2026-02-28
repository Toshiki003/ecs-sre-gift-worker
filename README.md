# dev-template

PR駆動開発を「毎回同じ品質」で回すための **テンプレートリポジトリ**です。  
（10個以上の個人開発・PoCを量産しても、Issue/PR/CI/セキュリティの最低ラインを揃える）

## 目的

- Private リポジトリでも運用できる、再現性の高い開発フローを用意する
- AI（Claude/Codex）を使う場合も、**使わない場合も**同じフローで回るようにする（ロックイン回避）
- 「変更の安全性」「セキュリティ」「レビュー」「サマリ」が残る形で開発する

---

## このテンプレでできること

### Always-on（AIなしでも成立）

- Issue テンプレ（Bug / Feature / Chore）
- PR テンプレ（変更点サマリ・影響範囲・テスト・セキュリティ/運用メモ）
- CI（PHP / Go / Python を自動判定して best-effort 実行）
- Dependency Review（PRに入る依存変更のチェック）
- Dependabot（依存更新の自動PR）
- `scripts/` によるローカル補助（PR作成を自動化）

### Optional（必要な時だけ、DEFAULT OFF）

- PR変更点サマリ生成（Codex Action など）
- Codex レビュー依頼コメント（ラベル駆動：`ai-review`）

> **重要**: Optional機能は、設定しなければ「自動でスキップ」されます。  
> つまり Claude/Codex の契約を終了しても、テンプレのフロー自体は壊れません。

---

## ディレクトリ構成

```text
dev-template/
  .github/
    CODEOWNERS                   # コードオーナー設定
    ISSUE_TEMPLATE/
      bug_report.yml
      feature_request.yml
      chore_task.yml
    PULL_REQUEST_TEMPLATE.md
    dependabot.yml
    workflows/
      ci.yml
      dependency-review.yml
      pr-summary.yml           # Optional
      codex-review-comment.yml # Optional
  .claude/                     # Claude Code 設定・スキル（Optional）
    CLAUDE.md
    rules/
    skills/
  claude-ext/                  # 要件・タスク管理ドキュメント（Optional）
    docs/
      app-requirements.md
      app-tasklist.md
      template-requirements.md
      template-tasklist.md
      requirements.md          # 互換ポインタ
      tasklist.md              # 互換ポインタ
      decision-log.md
      manual-workflow.md
  scripts/
    pr.sh                      # ブランチ作成
    pr-finish.sh               # コミット→Push→PR作成
  .gitignore
  README.md
  SECURITY.md
```

---

## 使い方（テンプレ repo → 量産）

### 1) このリポジトリを Template repository にする

GitHubのUI:

- Settings → General → **Template repository** をON

### 2) テンプレから新規リポジトリ作成（gh CLI）

```bash
OWNER="あなたのユーザー名 or 組織名"
TEMPLATE="dev-template"

for n in $(seq -w 1 10); do
  name="proj${n}"
  gh repo create "$OWNER/$name" --private --template "$OWNER/$TEMPLATE"
done
```

### 3) 初期設定（初回のみ）

テンプレートから作成した直後は、以下の初期設定を行ってください。

```bash
# Dependency Graph / Vulnerability Alerts の有効化
gh api "repos/OWNER/REPO/vulnerability-alerts" -X PUT

# または Claude Code で一括設定
# /setup-repo
```

> 未設定でもPRは作成可能ですが、Dependency Reviewが警告を出します。

---

## 日々の開発フロー（基本・AI不問）

### 1) Issue作成（任意）

Issueテンプレで作成して、タスクを明確化。

### 2) 実装

ローカルで実装し、変更を作る。（Claude Code 等のAIツールとの併用も可）

### 3) PR作成（スクリプトで自動）

- ブランチ作成だけ（作業開始用）

```bash
./scripts/pr.sh "feat: add x"
```

- コミット → Push → PR作成（作業完了用）

```bash
./scripts/pr-finish.sh "feat: add x"
```

`pr-finish.sh` はタイトルプレフィックス（`feat:` / `fix:` / `docs:` 等）からラベルを自動判定し、`ai-review` ラベルも付与します。

### 4) GitHub上でチェック

- CI が通ること
- dependency review がOKであること
- PRテンプレのサマリ/テスト/リスクが書けていること

### 5) OKなら merge

あなたは **PRサマリを確認して merge するだけ**に近づきます。

> Claude Code を使う場合は、上記フローを自動化するスキルが利用できます。
> 詳細は「[Claude Code 連携](#claude-code-連携optional)」セクションを参照してください。

---

## Optional（AI機能）を有効化する

> **Claude Code を使う場合**: `/setup-repo` を実行すると、以下の手順1〜2（Variables/Secrets設定）を対話的に一括設定できます。
> Dependency Graph・Branch Protection も同時に適用されます。

### 1) Repo Variables に `AI_ENABLED=true` を設定

GitHub: Settings → Secrets and variables → Actions → Variables

- `AI_ENABLED=true`

### 2) PRサマリ生成を使いたい場合（任意）

デフォルトでは **Groq（groq/compound）** を使用します。

**最小構成（Groq）**: Secret に `LLM_API_KEY` を設定するだけでOK。

GitHub: Settings → Secrets and variables → Actions → **Secrets**

- `LLM_API_KEY` — Groq の APIキー（https://console.groq.com/keys で取得）

> APIキーは **リポジトリにコミットしません**。Secretsに保存します。

**別のLLMプロバイダーに切り替えたい場合**:

Settings → Secrets and variables → Actions → **Variables** で以下を追加します。

| Variable名 | 説明 | デフォルト値（未設定時） |
|------------|------|------------------------|
| `LLM_API_BASE` | APIベースURL | `https://api.groq.com/openai/v1` |
| `LLM_MODEL` | モデル名 | `groq/compound` |

<details>
<summary>例: OpenAI に切り替える場合</summary>

| 種別 | 名前 | 値 |
|------|------|-----|
| Secret | `LLM_API_KEY` | OpenAI APIキー |
| Variable | `LLM_API_BASE` | `https://api.openai.com/v1` |
| Variable | `LLM_MODEL` | `gpt-4o-mini` |

</details>

<details>
<summary>例: Gemini に切り替える場合</summary>

| 種別 | 名前 | 値 |
|------|------|-----|
| Secret | `LLM_API_KEY` | Google AI Studio の APIキー |
| Variable | `LLM_API_BASE` | `https://generativelanguage.googleapis.com/v1beta/openai` |
| Variable | `LLM_MODEL` | `gemini-2.5-flash` |

Gemini は OpenAI互換エンドポイントを提供しているため、上記の設定だけで動作します。

</details>

<details>
<summary>旧バージョン（OPENAI_API_KEY）からの移行</summary>

1. Secret `OPENAI_API_KEY` を削除
2. Secret `LLM_API_KEY` を追加（OpenAI のキーをそのまま設定可）
3. Variable `LLM_API_BASE` に `https://api.openai.com/v1` を設定
4. Variable `LLM_MODEL` に `gpt-4o-mini` を設定

</details>

<details>
<summary>GitHub UIでの設定手順</summary>

**Secretの設定（APIキー等の機密情報）**:

1. GitHubリポジトリページ → **Settings** タブ
2. 左メニュー **Secrets and variables** → **Actions**
3. **Secrets** タブ → **New repository secret**
4. Name に `LLM_API_KEY`、Secret にAPIキーを入力 → **Add secret**

**Variableの設定（APIベースURL・モデル名等）**:

1. GitHubリポジトリページ → **Settings** タブ
2. 左メニュー **Secrets and variables** → **Actions**
3. **Variables** タブ → **New repository variable**
4. Name と Value を入力 → **Add variable**

> **SecretとVariableの違い**: Secretは暗号化されログに表示されません（APIキー向き）。Variableはワークフローログに表示されます（URL・モデル名向き）。

</details>

### 3) LLM外部送信ポリシー（必読）

PRサマリ機能（`.github/workflows/pr-summary.yml`）は、外部LLM APIに以下を送信します。

- PRタイトル
- `git diff origin/<base>...HEAD` の先頭10,000 bytes（`*.lock` と `package-lock.json` は除外）

以下を含む可能性があるPRでは、PRサマリ機能を有効化しないでください。

- APIキー、トークン、認証情報
- 個人情報、顧客データ、契約上の秘匿情報
- 未公開の脆弱性情報、インシデント情報

運用ルール:

- デフォルトは `AI_ENABLED` 未設定（OFF）
- 有効化は、リポジトリ管理者が情報分類と利用規約を確認した場合のみ
- 機密データを扱う期間・プロジェクトは `AI_ENABLED=false` で運用する

詳細は `SECURITY.md` を参照してください。

### 4) Codexレビュー依頼（ラベル駆動）

PRにラベル `ai-review` を付けると、`@codex review` コメントが自動で付きます。
（Codex GitHub連携が有効な場合、レビューが返ります）

---

## Claude Code 連携（Optional）

Claude Code を使ってローカル開発を加速するための設定・スキルが含まれています。
**設定しなくてもテンプレートの基本フローには影響しません。**

### セットアップ

このテンプレートから新規リポジトリを作成した場合、`.claude/` と `claude-ext/` は初期状態で含まれています。追加セットアップは不要です。

既存リポジトリへ後から導入する場合は、以下を手動で追加してください。

- `.claude/`
- `claude-ext/`

### 構成

| ディレクトリ | 役割 |
|-------------|------|
| `.claude/` | Claude Code のプロジェクト設定・ルール・スキル定義 |
| `claude-ext/docs/` | アプリ実装用 (`app-*`) とテンプレート保守用 (`template-*`) の要件/タスク・意思決定ログ等 |

### 利用可能なスキル（スラッシュコマンド）

| コマンド | 説明 |
|---------|------|
| `/setup-repo` | リポジトリ初期設定（Dependency Graph・Branch Protection・Variables）を一括適用 |
| `/analyze` | リポジトリの実装状況を分析してレポート生成 |
| `/update-tasks` | 分析結果に基づいてタスクリストを更新 |
| `/implement-next` | タスクリストから次の未着手タスクを実装しPR作成 |
| `/fix-review` | PRのレビュー指摘を取得して自動修正 |
| `/create-requirements` | ソースコードから要件定義書を生成 |

### Claude Code 利用時の開発フロー

詳細は [`.claude/CLAUDE.md`](.claude/CLAUDE.md) の「推奨開発フロー」を参照してください。

基本的な流れ:

1. `app-requirements.md` に要件を記載
2. `/analyze` → `/update-tasks` で `app-tasklist.md` を更新
3. `/implement-next` でタスクを実装・PR作成
4. レビュー後、`/fix-review` で指摘を修正
5. マージ → 次のタスクへ

テンプレート自体を保守する場合は `template-requirements.md` / `template-tasklist.md` を使用します。

---

## CodeQL を導入したい場合

静的セキュリティ解析が必要になったら、以下の手順で追加できます。
（Public リポジトリは無料。Private リポジトリは GitHub Advanced Security が必要）

1. `.github/workflows/codeql.yml` を作成:

```yaml
name: CodeQL

on:
  pull_request:
  push:
    branches: [main]
  schedule:
    - cron: '0 3 * * 1'

permissions:
  contents: read
  actions: read
  security-events: write

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - uses: github/codeql-action/init@v4
      - uses: github/codeql-action/analyze@v4
```

2. GitHub: Settings → Code security → Code scanning を有効化

---

## セキュリティ方針

詳細は `SECURITY.md` を参照してください。

---

## ライセンス

個人用途のテンプレとして自由に利用してください（必要なら後で追記）。
