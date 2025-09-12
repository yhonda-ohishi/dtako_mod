# dtako_mod

Production data import module for ryohi_sub_cal2 router system.

## 概要

`dtako_mod`は、本番環境から`dtako_rows`、`dtako_events`、`dtako_ferry`の3つのテーブルのデータを取り込むためのGoモジュールです。ryohi_sub_cal2ルーターのサブモジュールとして動作します。

## 機能

- **dtako_rows**: 車両運行データの管理
- **dtako_events**: イベントデータの管理  
- **dtako_ferry**: フェリー運航データの管理

各モジュールは以下の機能を提供します：
- データの一覧取得（日付範囲フィルター）
- 個別データの取得
- 本番環境からのデータインポート（UPSERT対応）

## インストール

```bash
go get github.com/yhonda-ohishi/dtako_mod
```

## API エンドポイント

### dtako_rows
- `GET /dtako/rows` - データ一覧取得
- `GET /dtako/rows/{id}` - 個別データ取得
- `POST /dtako/rows/import` - データインポート

### dtako_events
- `GET /dtako/events` - イベント一覧取得
- `GET /dtako/events/{id}` - 個別イベント取得
- `POST /dtako/events/import` - イベントインポート

### dtako_ferry
- `GET /dtako/ferry` - フェリーデータ一覧取得
- `GET /dtako/ferry/{id}` - 個別フェリーデータ取得
- `POST /dtako/ferry/import` - フェリーデータインポート

## テスト

```bash
make test           # 全テスト実行
make test-contract  # 契約テストのみ
make test-integration # 統合テストのみ
```

## ライセンス

MIT License

## 作者

yhonda-ohishi
