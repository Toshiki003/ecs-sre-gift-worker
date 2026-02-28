# テンプレート要件定義書（dev-template）

> このドキュメントは、`dev-template` 自体の仕様を定義する。  
> テンプレート保守の実装・レビュー・変更判断は本書を基準に行う。

---

## 1. プロジェクト概要

**プロジェクト名**: dev-template  
**目的**: PR駆動開発の最低運用品質（変更安全性・レビュー可能性・セキュリティ注意）を、個人開発/小規模開発で再利用可能なテンプレートとして提供する。  
**対象ユーザー**: GitHubで開発する個人開発者・小規模チーム

### 1.1 成果物

- GitHub運用テンプレート（Issue/PR/CI/依存レビュー/Dependabot）
- ローカル補助スクリプト（`scripts/pr.sh`, `scripts/pr-finish.sh`）
- OptionalなAI補助（PRサマリ、Codexレビューコメント）
- 運用ドキュメント（README、SECURITY、manual-workflow、本要件、tasklist、decision-log）

### 1.2 スコープ外

- 特定アプリケーションの業務ロジック実装
- CodeQLの常時有効化（必要時に利用者が任意で導入）
- CODEOWNERSの強制設定（利用者判断）

---

## 2. 技術スタック

| カテゴリ | 技術 | 備考 |
|---------|------|------|
| ドキュメント | Markdown | README/運用手順/要件管理 |
| 自動化 | GitHub Actions | CI、Dependency Review、Optional AI |
| スクリプト | Bash | PR作業補助 |
| 外部CLI | `git`, `gh` | `scripts/` 実行時に必要 |

---

## 3. 機能要件

### 3.1 Always-on（必須）

| ID | 要件 | 受け入れ条件 |
|----|------|-------------|
| FR-001 | Issueテンプレートを提供する | `.github/ISSUE_TEMPLATE/` に `bug_report.yml` / `feature_request.yml` / `chore_task.yml` が存在する |
| FR-002 | PRテンプレートを提供する | `.github/PULL_REQUEST_TEMPLATE.md` が存在し、変更要約・影響・テスト・セキュリティ観点を入力できる |
| FR-003 | 依存関係レビューをPRで実行する | `.github/workflows/dependency-review.yml` が `pull_request` で動作する |
| FR-004 | Dependabotで依存更新PRを定期作成する | `.github/dependabot.yml` が `pip` / `composer` / `gomod` / `github-actions` を週次設定している |
| FR-005 | PR作業補助スクリプトを提供する | `scripts/pr.sh` と `scripts/pr-finish.sh` が存在する |
| FR-006 | CIを自動判定で実行する | `.github/workflows/ci.yml` が `composer.json` / `go.mod` / `pyproject.toml` or `requirements.txt` を検知して該当ジョブを実行する |

### 3.2 PR補助スクリプト要件

| ID | 要件 | 受け入れ条件 |
|----|------|-------------|
| FR-010 | `pr.sh` は作業開始ブランチを安全に作成する | ワーキングツリーがdirtyなら失敗し、`feat/<slug>-<timestamp>` 形式の新規ブランチを作成する |
| FR-011 | `pr-finish.sh` は保護ブランチ上で実行拒否する | デフォルトブランチ/`main`/`master` で実行時にエラー終了する |
| FR-012 | `pr-finish.sh` は不正ブランチ名を拒否する | 許可プレフィックス外は失敗し、`ALLOW_ANY_BRANCH=true` 時のみ回避可能 |
| FR-013 | `pr-finish.sh` はPR作成失敗を成功扱いしない | `gh pr create` 非0終了時に異常終了し、URLが妥当形式でない場合も失敗扱いにする |
| FR-014 | ラベル運用を自動化する | タイトルプレフィックスから機能ラベルを付与し、`ai-review` も付与する |

### 3.3 Optional（AI機能）

| ID | 要件 | 受け入れ条件 |
|----|------|-------------|
| FR-020 | AI機能は明示有効化時のみ動作する | `AI_ENABLED=true` のときのみ Optional ワークフローが実行される |
| FR-021 | PRサマリを外部LLMで生成できる | `.github/workflows/pr-summary.yml` がPRタイトル+差分（最大10,000 bytes、lockfile除外）を送信し、コメントを作成/更新する |
| FR-022 | Codexレビュー依頼コメントを自動化する | `.github/workflows/codex-review-comment.yml` が `ai-review` ラベルで `@codex review` コメントを投稿する |
| FR-023 | 外部送信ポリシーを文書化する | READMEとSECURITYにLLM外部送信ポリシーが明記されている |

### 3.4 Claude運用ドキュメント

| ID | 要件 | 受け入れ条件 |
|----|------|-------------|
| FR-030 | 仕様・タスク・判断履歴を保持する | `claude-ext/docs/template-requirements.md` / `template-tasklist.md` / `decision-log.md` が存在し、実体内容を持つ |
| FR-031 | AI非依存の手動手順を提供する | `claude-ext/docs/manual-workflow.md` で `scripts/` と `gh` によるPR運用手順が示される |

---

## 4. 非機能要件

| 項目 | 要件 |
|------|------|
| 安全性 | `pr-finish.sh` は保護ブランチを拒否し、PR作成失敗を隠蔽しない |
| 再現性 | テンプレートから新規作成したリポジトリで、追加セットアップなしに基本フローが成立する |
| 可搬性 | Bash + GitHub Actions + `gh` CLI が利用できる環境で動作する |
| セキュリティ | APIキーはSecrets管理を前提とし、LLM外部送信ポリシーをREADME/SECURITYで明示する |
| 保守性 | 仕様変更時は `decision-log.md` と関連ドキュメントを同一PRで更新する |

---

## 5. 運用ルール

- Optional AIはデフォルトOFF（`AI_ENABLED` 未設定）
- CodeQLは任意導入
- CODEOWNERSは利用者方針で設定
- 仕様変更PRでは最低限以下を同時更新する
  - `README.md`（利用者向け説明）
  - `SECURITY.md`（セキュリティ方針）
  - `claude-ext/docs/decision-log.md`（判断記録）

---

## 6. 検証チェックリスト

- [ ] `bash -n scripts/pr.sh scripts/pr-finish.sh` が成功する
- [ ] `.github/workflows/` の主要ワークフローが存在する（CI / dependency-review / optional 2種）
- [ ] READMEの手順と実ファイル構成に齟齬がない
- [ ] LLM外部送信ポリシーが README と SECURITY の両方で確認できる

---

**作成日**: 2026-02-24  
**最終更新**: 2026-02-24
