# DTako Module - Router メンテナンス指示書（Claude 向け）

## 🎯 目的
この指示書は、Ryohi Router プロジェクトでDTakoモジュールを統合・メンテナンスする際のClaude向けガイダンスです。

## 📋 基本情報

### プロジェクト構成
- **DTako Module**: `github.com/yhonda-ohishi/dtako_mod`（このリポジトリ）
- **Ryohi Router**: 親プロジェクト（DTakoモジュールを統合する側）
- **統合方式**: Go moduleのimportによる自動統合

### 重要な設計原則
1. **DTakoモジュールは独立したライブラリ**として設計されている
2. **ルートプレフィックスは親プロジェクトが制御**する
3. **Swagger統合は自動化**されている

## 🔧 Ryohi Router 側でのメンテナンス手順

### 1. DTakoモジュールの統合

```go
// main.go または router設定ファイルで
import (
    "github.com/go-chi/chi/v5"
    "github.com/yhonda-ohishi/dtako_mod"
    _ "github.com/yhonda-ohishi/dtako_mod/docs"  // Swagger自動統合
)

func setupRouter() *chi.Mux {
    r := chi.NewRouter()

    // DTakoモジュールを /dtako パスにマウント
    r.Route("/dtako", func(r chi.Router) {
        dtako_mod.RegisterRoutes(r)
    })

    return r
}
```

### 2. 環境変数の設定

Ryohi Router の`.env`ファイルに追加：

```env
# DTako Module Database Configuration
DB_HOST=localhost
DB_PORT=3307
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=dtako_local

# Alternative (legacy support)
LOCAL_DB_HOST=localhost
LOCAL_DB_PORT=3307
LOCAL_DB_USER=root
LOCAL_DB_PASSWORD=your_password
LOCAL_DB_NAME=dtako_local

# Debug Mode
DEBUG=false
```

### 3. DTakoモジュールの更新

```bash
# 最新版を取得
go get -u github.com/yhonda-ohishi/dtako_mod@latest

# 依存関係を整理
go mod tidy

# Swaggerドキュメントを再生成（Ryohi Router側で）
swag init -g cmd/main.go -o docs

# テスト実行
go test ./...
```

## ⚠️ 重要な注意事項とトラブルシューティング

### 🚫 やってはいけないこと

1. **DTakoモジュールのルート定義を変更しない**
   ```go
   // ❌ 間違い
   r.Get("/dtako/rows", handler) // 二重プレフィックスになる

   // ✅ 正しい
   r.Route("/dtako", func(r chi.Router) {
       dtako_mod.RegisterRoutes(r) // モジュール内部では /rows のみ定義
   })
   ```

2. **環境変数ファイルをサブディレクトリにコピーしない**
   - 必ずプロジェクトルートの`.env`を使用
   - テストは `go test ./...` でルートから実行

3. **Swaggerインスタンス名を変更しない**
   - DTako側は`InfoInstanceName: "swagger"`で統合用に設定済み

### 🔍 トラブルシューティング

#### 問題1: 404 Not Found エラー

**症状**: `/dtako/rows`にアクセスしても404エラー

**原因と解決**:
```go
// ❌ 間違ったマウント方法
r.Mount("/dtako", dtako_mod.RegisterRoutes())

// ✅ 正しいマウント方法
r.Route("/dtako", func(r chi.Router) {
    dtako_mod.RegisterRoutes(r)
})
```

#### 問題2: データベース接続エラー

**症状**: `Error 1049: Unknown database 'dtako_local'`

**解決手順**:
1. 環境変数を確認: `echo $DB_HOST $DB_PORT $DB_NAME`
2. データベースの存在確認: `mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD -e "SHOW DATABASES;"`
3. DTakoの診断ツール実行: `go run github.com/yhonda-ohishi/dtako_mod/cmd/diagnose`

#### 問題3: Swagger統合されない

**症状**: Swagger UIにDTakoのエンドポイントが表示されない

**解決**:
```go
import (
    _ "github.com/yhonda-ohishi/dtako_mod/docs"  // このimportを追加
    _ "your-project/docs"                        // 既存のdocs
)
```

## 🧪 テスト手順

### 単体テスト
```bash
# DTakoモジュールのテスト（モジュール内で）
go test ./tests/contract/... -v
go test ./tests/integration/... -v

# Ryohi Router全体のテスト
go test ./... -v
```

### 統合テスト（Ryohi Router側で）
```bash
# サーバー起動
go run cmd/main.go

# DTakoエンドポイントの確認
curl http://localhost:8080/dtako/rows
curl http://localhost:8080/dtako/events
curl http://localhost:8080/dtako/ferry_rows

# Swagger UI確認
# http://localhost:8080/swagger/index.html
```

## 📝 メンテナンス作業時のチェックリスト

### DTakoモジュール更新時
- [ ] `go get -u github.com/yhonda-ohishi/dtako_mod@latest` 実行
- [ ] `go mod tidy` で依存関係を整理
- [ ] 既存テストが通ることを確認
- [ ] Swagger UIでエンドポイントが正しく表示されることを確認
- [ ] 新機能の動作テスト実施

### 新環境セットアップ時
- [ ] データベースセットアップ実行
- [ ] 環境変数ファイル（`.env`）設定
- [ ] DTako診断ツールで接続確認
- [ ] 全エンドポイントの疎通テスト
- [ ] Swagger UIでドキュメント確認

### 問題発生時
- [ ] ログレベルをDEBUG=trueに設定
- [ ] DTako診断ツール実行
- [ ] 環境変数の値確認
- [ ] データベース接続確認
- [ ] ルーティング設定確認

## 🔄 定期メンテナンス

### 月次作業
- DTakoモジュールの最新版チェック
- セキュリティアップデートの適用
- パフォーマンス監視データの確認

### 四半期作業
- DTakoモジュールのAPIコンパチビリティ確認
- 統合テストシナリオの見直し
- ドキュメント更新の確認

## 📚 関連リソース

- [DTako Module Constitution](/.specify/memory/constitution.md)
- [Important Rules](IMPORTANT_RULES.md)
- [Database Connection Instructions](../ryohi_sub_cal2/docs/DTAKO_DATABASE_CONNECTION_FIX.md)
- [Route Fix Documentation](../ryohi_sub_cal2/docs/DTAKO_ROUTE_FIX.md)

## 🆘 緊急時の連絡先とエスカレーション

1. **即座にできること**: DTako診断ツールの実行
2. **ログ確認**: `DEBUG=true` でデバッグログ有効化
3. **ロールバック**: 前回正常動作していたバージョンに戻す

---

**Last Updated**: 2025-09-13
**Version**: 1.0.0
**Author**: Claude Code Assistant