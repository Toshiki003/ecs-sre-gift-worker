# アプリ実装タスクリスト

> `app-requirements.md` から導出した実装タスクを管理する。

## 凡例
- [ ] 未着手
- [x] 完了
- 🚧 進行中

---

## フェーズ 1: 要件確定・設計

- [x] `app-requirements.md` のプロジェクト概要を確定
- [x] MVPスコープを確定（セクション5.1〜5.3）
- [x] 技術スタック選定（Go, RDS, SQS, ECS Fargate, Datadog Agent）
- [x] データモデル設計（gifts, send_requests テーブル）
- [x] API仕様設計（POST /gifts, POST /gifts/{id}/send, GET /gifts/{id}, /healthz, /readyz）
- [x] 状態遷移設計（CREATED→QUEUED→SENT/FAILED）
- [x] 冪等性・重複排除方式の設計（Idempotency-Key）
- [x] 監視・アラート設計（Datadog ダッシュボード1枚 + アラート6〜8個）
- [x] インフラ構成設計（VPC, ALB, ECS, RDS, SQS, ECR, IAM, CloudWatch）

---

## フェーズ 2: Go プロジェクト基盤

- [x] T-2.1: Go プロジェクト初期化（go.mod, ディレクトリ構成）
  - `cmd/api/`, `cmd/worker/`, `internal/` 構成
  - Go modules 初期化
- [ ] T-2.2: 設定管理（環境変数ベース）
  - DB接続情報、SQS URL、ポート番号等を環境変数から読み込み
- [ ] T-2.3: 構造化ログ基盤
  - JSON構造化ログ（gift_id, request_id, status, service フィールド対応）
  - 要件8.3準拠
- [ ] T-2.4: DB接続・マイグレーション
  - PostgreSQL接続プール設定
  - gifts テーブル / send_requests テーブルのマイグレーションSQL
  - 要件4.3準拠
- [ ] T-2.5: ヘルスチェックエンドポイント
  - GET /healthz（Liveness: プロセス生存確認）
  - GET /readyz（Readiness: DB接続確認）
  - 要件5.1.1準拠

---

## フェーズ 3: API エンドポイント実装

- [ ] T-3.1: エラーレスポンス共通基盤
  - 共通エラーフォーマット `{ error: { code, message }, request_id }`
  - 要件5.1.0準拠
- [ ] T-3.2: POST /gifts（ギフト作成）
  - 入力バリデーション（amount: 必須・正の整数、message: 最大500文字）
  - ギフトコード自動生成（UUID）
  - DB保存（status=CREATED）
  - 要件5.1.1 / 6.1準拠
- [ ] T-3.3: GET /gifts/{gift_id}（ギフト状態照会）
  - gift_id 存在チェック（404）
  - 全フィールド返却
  - 要件5.1.1準拠
- [ ] T-3.4: POST /gifts/{gift_id}/send（ギフト送信依頼）
  - Idempotency-Key ヘッダー必須チェック（400）
  - gift_id 存在チェック（404）
  - 状態チェック（SENT→409再送不可、QUEUED+異なるrequest_id→409）
  - SELECT FOR UPDATE による排他ロック
  - send_requests INSERT + gifts.status=QUEUED UPDATE（トランザクション内）
  - 冪等性: 同一(gift_id, request_id)は既存結果返却（200）
  - 要件5.1.1 / 5.1.3 / 6.2準拠
- [ ] T-3.5: SQSメッセージ投入（送信依頼連携）
  - トランザクション確定後にSQSへ { gift_id, request_id } を投入
  - SQS投入失敗時はログ記録（DB側はQUEUED確定済み）
  - 要件5.1.1 / 7.2準拠
- [ ] T-3.6: API ユニットテスト
  - 各エンドポイントの正常系・異常系テスト
  - 冪等性の動作確認テスト

---

## フェーズ 4: Worker 実装

- [ ] T-4.1: SQSポーリング基盤
  - SQSメッセージ受信ループ
  - graceful shutdown 対応
- [ ] T-4.2: メッセージ処理ロジック
  - gift_id, request_id の取得・バリデーション
  - 重複処理チェック（status が SENT/FAILED なら正常消化）
  - 擬似送信処理（20%ランダム失敗）
  - 要件5.1.2準拠
- [ ] T-4.3: 状態遷移更新（条件付きUPDATE）
  - 成功時: `UPDATE gifts SET status='SENT' WHERE id=? AND status='QUEUED'`
  - 失敗時（maxReceiveCount到達）: `UPDATE gifts SET status='FAILED', last_error=? WHERE id=? AND status='QUEUED'`
  - ApproximateReceiveCount による最終失敗判定
  - 要件5.1.2 / 6.2 / 6.3準拠
- [ ] T-4.4: Worker ユニットテスト
  - メッセージ処理の正常系・異常系テスト
  - 条件付きUPDATEの動作確認

---

## フェーズ 5: コンテナ化・ローカル開発環境

- [ ] T-5.1: Dockerfile（API）
  - マルチステージビルド（Go バイナリ）
  - 最小イメージ（distroless or alpine）
- [ ] T-5.2: Dockerfile（Worker）
  - マルチステージビルド（Go バイナリ）
  - 最小イメージ（distroless or alpine）
- [ ] T-5.3: docker-compose.yml（ローカル開発用）
  - API + Worker + PostgreSQL + LocalStack(SQS) の構成
  - 動作確認用の一括起動環境

---

## フェーズ 6: Terraform インフラ定義

- [ ] T-6.1: Terraform 基盤（main.tf, variables.tf, outputs.tf）
  - provider設定、backend設定、共通変数定義
- [ ] T-6.2: VPC / Subnet / Security Group（vpc.tf）
  - public subnet（MVP: NAT回避）
  - SG: ALB用、ECS用、RDS用
  - 要件8.4 / 8.5準拠
- [ ] T-6.3: ECR リポジトリ（ecr.tf）
  - api用 / worker用の2リポジトリ
- [ ] T-6.4: SQS Queue + DLQ（sqs.tf）
  - 標準キュー + Dead Letter Queue
  - Redrive Policy（maxReceiveCount=5）
  - 要件6.3準拠
- [ ] T-6.5: RDS PostgreSQL（rds.tf）
  - 最小構成（db.t3.micro等）
  - Secrets Manager でパスワード管理
  - 要件8.4準拠
- [ ] T-6.6: ALB / Target Group / Listener（alb.tf）
  - ヘルスチェック: /healthz
  - 要件7.1準拠
- [ ] T-6.7: ECS Cluster / Task Definition / Service（ecs.tf）
  - API サービス（Fargate + Datadog Agent サイドカー）
  - Worker サービス（Fargate + Datadog Agent サイドカー）
  - 要件3.1 / 7.1準拠
- [ ] T-6.8: IAM Role / Policy（iam.tf）
  - task execution role（ECR pull, CloudWatch Logs）
  - task role（SQS, Secrets Manager）
  - 最小権限原則（要件8.4準拠）
- [ ] T-6.9: CloudWatch Log Group（cloudwatch.tf）
  - API / Worker 用ロググループ
- [ ] T-6.10: terraform fmt / validate 確認

---

## フェーズ 7: CI/CD パイプライン

- [ ] T-7.1: PRトリガー ワークフロー更新
  - Go テスト実行（go test ./...）
  - terraform fmt チェック
  - terraform validate チェック
  - 要件10.1準拠
- [ ] T-7.2: デプロイワークフロー（mainマージトリガー）
  - Docker Build → ECR Push → ECS Task Definition 更新 → ECS Deploy
  - 要件10.2準拠

---

## フェーズ 8: 監視（Datadog）

- [ ] T-8.1: Datadog ダッシュボード定義
  - API/ALB/ECS/DB/Queue を1枚で俯瞰
  - 要件8.2準拠
- [ ] T-8.2: Datadog アラート定義（必須6個）
  - ALB 5xx 増加
  - ALB HealthyHostCount 低下
  - ECS running < desired（API）
  - ECS running < desired（Worker）
  - SQS Queue depth 増加 or 最古遅延増加
  - Worker 失敗ログ急増
  - 要件8.2準拠
- [ ] T-8.3: Datadog アラート定義（任意2個）
  - RDS CPU 高止まり
  - RDS connections 急増

---

## フェーズ 9: ドキュメント・運用

- [ ] T-9.1: docs/architecture.md
  - 責務分割、設計判断（NAT有無、DB選定、冪等性方針）
  - 要件12準拠
- [ ] T-9.2: docs/observability.md
  - SLO仮置き、アラート一覧、ダッシュボード設計
  - 要件8.1 / 8.2 / 12準拠
- [ ] T-9.3: docs/runbook.md
  - 一次対応手順（切り分け→暫定→恒久）
  - 障害シナリオ: Worker停止、SQS滞留、DB接続失敗
  - 要件11.1 / 13準拠
- [ ] T-9.4: docs/cost.md
  - コスト概算、NAT有無の判断、public subnet運用の理由と対策
  - 要件8.5 / 12準拠
- [ ] T-9.5: README.md 更新
  - アーキテクチャ概要、構成図、起動手順、デプロイ手順
  - 要件12準拠

---

## 次に着手するタスク

- [ ] **T-2.2: 設定管理**（環境変数ベースの設定読み込み）

---

**最終更新**: 2026-03-01
