---
name: setup-repo
description: リポジトリの初期設定（Dependency Graph・Branch Protection・Variables）を一括適用
argument-hint:
allowed-tools: Bash
---

# リポジトリ初期設定

現在のリポジトリに対して、以下の初期設定を一括で適用してください。

## 手順

### 1. リポジトリ情報の取得

```bash
gh api repos/{owner}/{repo}
```

- `owner` / `repo` は `gh repo view --json owner,name` から取得する
- デフォルトブランチ名を確認する

### 2. Dependency Graph / Vulnerability Alerts の有効化

```bash
gh api repos/{owner}/{repo}/vulnerability-alerts -X PUT
```

- HTTP 204 が返れば成功

### 3. デフォルトブランチの保護設定

```bash
gh api repos/{owner}/{repo}/branches/{default_branch}/protection -X PUT \
  -H "Accept: application/vnd.github+json" \
  --input - <<'EOF'
{
  "required_pull_request_reviews": {
    "required_approving_review_count": 0
  },
  "required_status_checks": null,
  "enforce_admins": false,
  "restrictions": null,
  "allow_force_pushes": false,
  "allow_deletions": false
}
EOF
```

- PRマージ必須（承認0名＝セルフマージ可）
- Force push 禁止
- ブランチ削除禁止

### 4. GitHub Variables の設定

ワークフローが使用するVariablesを設定する。

```bash
# 現在の設定を確認
gh variable list
```

以下のVariablesについて、未設定のものをユーザーに確認してから設定する:

| Variable名 | 必須 | デフォルト値 | 用途 |
|------------|------|------------|------|
| `AI_ENABLED` | Yes | - | AI機能（PR要約・Codexレビュー）の有効化フラグ |
| `LLM_API_BASE` | No | `https://api.groq.com/openai/v1` | LLM APIのベースURL |
| `LLM_MODEL` | No | `groq/compound` | 使用するLLMモデル |

設定手順:
1. `gh variable list` で既存の設定を確認する
2. `AI_ENABLED` が未設定の場合、`true` を推奨としてユーザーに値を確認し、`gh variable set AI_ENABLED --body "<値>"` で設定する
3. `LLM_API_BASE` が未設定の場合、デフォルト値 `https://api.groq.com/openai/v1` を提示し、変更したい場合のみユーザーに値を確認して `gh variable set LLM_API_BASE --body "<値>"` で設定する
4. `LLM_MODEL` が未設定の場合、デフォルト値 `groq/compound` を提示し、変更したい場合のみユーザーに値を確認して `gh variable set LLM_MODEL --body "<値>"` で設定する

### 5. Secrets の手動設定ガイド

Secretsはセキュリティ上の理由からCLIでの設定を行わず、WebUIでの手動設定を案内する。

`{owner}` / `{repo}` は手順1で取得した値を使用する。

以下のメッセージをユーザーに表示する:

```
⚠️ 以下のSecretはGitHub WebUIから手動で設定してください:

  - LLM_API_KEY（必須）: LLM APIの認証キー

  設定ページ: https://github.com/{owner}/{repo}/settings/secrets/actions

  ※ GITHUB_TOKEN は自動提供されるため設定不要です。
```

### 6. 設定結果の確認・レポート

適用後、以下を実行して結果を確認する:

```bash
# Vulnerability Alerts の状態確認（204=有効）
gh api repos/{owner}/{repo}/vulnerability-alerts -X GET

# Branch Protection の状態確認
gh api repos/{owner}/{repo}/branches/{default_branch}/protection

# Variables の設定確認
gh variable list

# Secrets の設定確認（値は表示されない）
gh secret list
```

最終的に、適用した設定の一覧をユーザーに表示する:

| 設定項目 | 状態 |
|---------|------|
| Dependency Graph / Vulnerability Alerts | 有効 / 失敗 |
| Branch Protection: PRマージ必須 | 有効 / 失敗 |
| Branch Protection: Force push 禁止 | 有効 / 失敗 |
| Branch Protection: ブランチ削除禁止 | 有効 / 失敗 |
| Variable: AI_ENABLED | 設定済み(値) / 未設定 |
| Variable: LLM_API_BASE | 設定済み(値) / 未設定(デフォルト使用) |
| Variable: LLM_MODEL | 設定済み(値) / 未設定(デフォルト使用) |
| Secret: LLM_API_KEY | 設定済み / 未設定(要手動設定) |

**注意**: エラーが発生した場合は、エラー内容とともに手動での設定手順を案内してください。
