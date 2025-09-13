# DTako Module Constitution

## Core Principles

### I. Environment Variables Management
- **絶対にテストフォルダやサブディレクトリに`.env`ファイルをコピーしない**
- 環境変数の設定は必ずプロジェクトルートの`.env`ファイルから読み込む
- テストはプロジェクトルートから実行するか、適切な環境変数を設定して実行する
- セキュリティと保守性のため、設定は一元管理する

### II. Database Connection Pattern
- godotenv/autoloadを使用した自動的な環境変数読み込み
- DB_* 環境変数を優先し、LOCAL_DB_* にフォールバック
- シングルトン接続パターンで接続を管理
- 接続エラーは適切にログに記録

### III. Route Prefix Management
- DTako moduleのルートは親ルーターがマウントポイントを制御
- RegisterRoutes関数内では`/dtako`プレフィックスを付けない
- テストではヘルパー関数で`/dtako`にマウント
- 二重プレフィックス問題を防ぐ

### IV. Japanese Database Compatibility
- dtako_ferry_rowsテーブルの日本語カラム名に対応
- TIME型フィールドは文字列として取得後パース
- NULLフィールドは適切にハンドリング
- 文字エンコーディングはutf8mb4を使用

### V. Testing Strategy
- 契約テストで外部APIとの互換性を保証
- テストデータベースは本番と同じスキーマを使用
- REPLACE INTOでべき等なテストデータ挿入
- 各テストは独立して実行可能

## Security Requirements

### Database Credentials
- パスワードを環境変数で管理
- .envファイルは.gitignoreに追加
- .env.exampleで設定例を提供
- コマンドラインでパスワードを直接指定しない

### API Security
- 入力値の検証とサニタイゼーション
- SQLインジェクション対策としてプリペアドステートメント使用
- エラーメッセージに機密情報を含めない

## Development Workflow

### Commit Messages
- 日本語でのコミットメッセージを使用可能
- 変更内容を明確に記述
- プレフィックス（fix:, feat:, refactor:など）を使用

### Code Style
- Go標準のフォーマッティングを遵守
- エラーハンドリングは明示的に
- 不要なコメントは追加しない
- 既存のコードスタイルに従う

## Governance

- この憲法はすべての開発慣行に優先する
- 変更には文書化と承認が必要
- IMPORTANT_RULES.mdで追加ルールを管理
- 違反は即座に修正する

**Version**: 1.0.0 | **Ratified**: 2025-09-13 | **Last Amended**: 2025-09-13