# クイックスタートガイド: dtako_mod

## 概要
dtako_modを使用して本番環境からデータをインポートし、ローカルで利用する手順です。

## 前提条件
- Go 1.21以上
- MySQL 5.7以上（ローカルと本番）
- Git

## セットアップ

### 1. リポジトリのクローン
```bash
git clone https://github.com/yhonda-ohishi/dtako_mod.git
cd dtako_mod
```

### 2. 依存関係のインストール
```bash
go mod download
```

### 3. データベースのセットアップ
```bash
# ローカルデータベースの作成
mysql -u root -p < schema.sql
```

### 4. 環境設定
```bash
# .envファイルの作成
cp .env.example .env

# .envファイルを編集して接続情報を設定
# 本番DB
PROD_DB_HOST=production.example.com
PROD_DB_PORT=3306
PROD_DB_USER=prod_user
PROD_DB_PASSWORD=prod_password
PROD_DB_NAME=production_db

# ローカルDB
LOCAL_DB_HOST=localhost
LOCAL_DB_PORT=3306
LOCAL_DB_USER=root
LOCAL_DB_PASSWORD=
LOCAL_DB_NAME=dtako_local
```

## 使用方法

### サーバーの起動（スタンドアロンモード）
```bash
go run cmd/server/main.go
```

### ryohi_sub_cal2への統合
```go
// ryohi_sub_cal2のmain.goで
import (
    "github.com/go-chi/chi/v5"
    dtako "github.com/yhonda-ohishi/dtako_mod"
)

func main() {
    r := chi.NewRouter()
    
    // dtako_modのルートを登録
    dtako.RegisterRoutes(r)
    
    http.ListenAndServe(":8080", r)
}
```

## 基本的な操作例

### 1. 運行データのインポート（過去1ヶ月）
```bash
curl -X POST http://localhost:8080/dtako/rows/import \
  -H "Content-Type: application/json" \
  -d '{
    "from_date": "2025-08-01",
    "to_date": "2025-08-31"
  }'
```

期待される応答:
```json
{
  "success": true,
  "imported_rows": 1234,
  "message": "Imported 1234 rows from 2025-08-01 to 2025-08-31",
  "imported_at": "2025-09-12T10:00:00Z"
}
```

### 2. インポートしたデータの確認
```bash
# 運行データの取得
curl "http://localhost:8080/dtako/rows?from=2025-08-01&to=2025-08-31"
```

### 3. イベントデータのインポート（特定タイプのみ）
```bash
curl -X POST http://localhost:8080/dtako/events/import \
  -H "Content-Type: application/json" \
  -d '{
    "from_date": "2025-08-01",
    "to_date": "2025-08-31",
    "event_type": "ACCIDENT"
  }'
```

### 4. フェリーデータのインポート（特定航路）
```bash
curl -X POST http://localhost:8080/dtako/ferry/import \
  -H "Content-Type: application/json" \
  -d '{
    "from_date": "2025-08-01",
    "to_date": "2025-08-31",
    "route": "ROUTE_A"
  }'
```

## テストの実行

### ユニットテスト
```bash
go test ./...
```

### 統合テスト
```bash
go test -tags=integration ./tests/integration
```

### コントラクトテスト
```bash
go test ./tests/contract
```

## トラブルシューティング

### 本番DBに接続できない
1. ファイアウォール設定を確認
2. 認証情報が正しいか確認
3. ネットワーク接続を確認

### インポートが遅い
1. 日付範囲を狭める
2. バッチサイズを調整（環境変数: BATCH_SIZE）
3. データベースインデックスを確認

### 重複データエラー
- 正常動作です。UPSERTにより既存データは更新されます

## 検証シナリオ

### シナリオ1: 初回データインポート
1. ローカルDBが空の状態を確認
2. 1週間分のデータをインポート
3. インポート結果を確認（件数が正しいか）
4. データを取得して内容を確認

### シナリオ2: データ更新
1. 既存データがある状態で同じ期間を再インポート
2. エラーが発生しないことを確認
3. データが更新されていることを確認

### シナリオ3: エラーハンドリング
1. 本番DBを停止（またはネットワーク切断）
2. インポートを実行
3. 適切なエラーメッセージが返ることを確認

### シナリオ4: 大量データ処理
1. 3ヶ月分のデータインポートを実行
2. タイムアウトしないことを確認
3. メモリ使用量が適切か確認
4. 全データが正しくインポートされたか確認

## パフォーマンスベンチマーク

### 目標値
- 1ヶ月分のデータ（約30万レコード）: 5分以内
- API応答時間: 200ms以内（95パーセンタイル）
- メモリ使用量: 500MB以下

### 測定方法
```bash
# ベンチマークテストの実行
go test -bench=. ./tests/benchmark
```

## 次のステップ

1. 認証機能の追加（APIキー）
2. スケジュール実行機能
3. データ集計API
4. 監視・アラート機能

## サポート

問題が発生した場合は、GitHubのIssueを作成してください：
https://github.com/yhonda-ohishi/dtako_mod/issues