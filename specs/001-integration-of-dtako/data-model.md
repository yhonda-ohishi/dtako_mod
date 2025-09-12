# データモデル定義: dtako_mod

**作成日**: 2025-09-12  
**フェーズ**: 1 - 設計

## エンティティ定義

### 1. DtakoRow（運行データ）
車両の運行記録を表すエンティティ

**フィールド**:
| フィールド名 | 型 | 必須 | 説明 |
|------------|---|------|------|
| id | string | ✓ | 一意識別子（主キー） |
| date | date | ✓ | 運行日 |
| vehicle_no | string | ✓ | 車両番号 |
| driver_code | string | ✓ | 運転手コード |
| route_code | string | ✓ | ルートコード |
| distance | decimal(10,2) | ✓ | 走行距離（km） |
| fuel_amount | decimal(10,2) | ✓ | 燃料消費量（L） |
| created_at | timestamp | ✓ | 作成日時 |
| updated_at | timestamp | ✓ | 更新日時 |

**バリデーション**:
- date: 未来日不可
- distance: 0以上
- fuel_amount: 0以上

**インデックス**:
- PRIMARY: id
- INDEX: date, vehicle_no, driver_code, route_code

### 2. DtakoEvent（イベントデータ）
システムまたは運用イベントを記録

**フィールド**:
| フィールド名 | 型 | 必須 | 説明 |
|------------|---|------|------|
| id | string | ✓ | 一意識別子（主キー） |
| event_date | datetime | ✓ | イベント発生日時 |
| event_type | string | ✓ | イベントタイプ |
| vehicle_no | string | ✓ | 車両番号 |
| driver_code | string | ✓ | 運転手コード |
| description | text | ✓ | イベント詳細 |
| latitude | decimal(10,8) | | 緯度（オプション） |
| longitude | decimal(11,8) | | 経度（オプション） |
| created_at | timestamp | ✓ | 作成日時 |
| updated_at | timestamp | ✓ | 更新日時 |

**バリデーション**:
- event_date: 未来日時不可
- event_type: 定義済みタイプのみ
- latitude: -90 ～ 90
- longitude: -180 ～ 180

**イベントタイプ**:
- START: 運行開始
- STOP: 運行終了
- BREAK: 休憩
- ACCIDENT: 事故
- MAINTENANCE: メンテナンス
- OTHER: その他

**インデックス**:
- PRIMARY: id
- INDEX: event_date, event_type, vehicle_no, driver_code

### 3. DtakoFerry（フェリー運航データ）
フェリー運航記録

**フィールド**:
| フィールド名 | 型 | 必須 | 説明 |
|------------|---|------|------|
| id | string | ✓ | 一意識別子（主キー） |
| date | date | ✓ | 運航日 |
| route | string | ✓ | 航路名 |
| vehicle_no | string | ✓ | 船舶番号 |
| driver_code | string | ✓ | 船長コード |
| departure_time | datetime | ✓ | 出発時刻 |
| arrival_time | datetime | ✓ | 到着時刻 |
| passengers | int | ✓ | 乗客数 |
| vehicles | int | ✓ | 車両数 |
| created_at | timestamp | ✓ | 作成日時 |
| updated_at | timestamp | ✓ | 更新日時 |

**バリデーション**:
- arrival_time > departure_time
- passengers >= 0
- vehicles >= 0

**航路**:
- ROUTE_A: A航路
- ROUTE_B: B航路
- ROUTE_C: C航路

**インデックス**:
- PRIMARY: id
- INDEX: date, route, vehicle_no, driver_code, departure_time, arrival_time

### 4. ImportResult（インポート結果）
データインポート操作の結果

**フィールド**:
| フィールド名 | 型 | 必須 | 説明 |
|------------|---|------|------|
| success | boolean | ✓ | 成功フラグ |
| imported_rows | int | ✓ | インポート件数 |
| message | string | ✓ | 結果メッセージ |
| imported_at | timestamp | ✓ | インポート実行日時 |
| errors | []string | | エラーメッセージリスト |

## リレーションシップ

### エンティティ間の関係
- DtakoRow, DtakoEvent, DtakoFerry は独立したエンティティ
- vehicle_no, driver_codeは共通だが、外部キー制約なし（柔軟性のため）

### データ整合性
- 各テーブルは独立して管理
- インポート時の重複はUPSERTで処理
- 削除は論理削除ではなく物理削除

## 状態遷移

### インポートプロセス
```
[未インポート] → [インポート中] → [インポート完了]
                       ↓
                  [エラー発生] → [部分的完了]
```

### データライフサイクル
1. **作成**: 本番DBから取得
2. **保存**: ローカルDBに挿入/更新
3. **参照**: API経由で取得
4. **更新**: 再インポートで上書き
5. **削除**: 手動削除のみ（自動削除なし）

## パフォーマンス考慮事項

### インデックス戦略
- 日付範囲検索用: date, event_date
- 車両/運転手検索用: vehicle_no, driver_code
- 複合インデックスで検索最適化

### バッチ処理
- インポート: 1000件単位
- 検索結果: ページネーション（100件/ページ）

### キャッシュ
- 現時点では実装なし
- 将来的にRedis検討可能

## セキュリティ

### データ保護
- 個人情報（運転手コード）はハッシュ化を検討
- SQLインジェクション対策: パラメータバインディング

### アクセス制御
- 現時点では認証なし
- 将来的にAPIキー認証を検討

## 拡張性

### 将来の拡張ポイント
1. 新しいイベントタイプの追加
2. 追加フィールド（コスト、評価等）
3. リレーションの追加（運転手マスタ等）
4. 集計ビューの作成