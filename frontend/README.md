# Ducks フロントエンド

学園祭のイベント・ポスター・物販・売上を管理する Vue 3 + TypeScript + Vue Router + Vite のSPAです。
既存のGoバックエンドとデータ形式、画面URLを引き継いでいます。

## 開発

Node.js 20.19以上（22系は22.12以上）を使用してください。

```sh
npm ci
cp .env.example .env
npm run dev
```

Windows PowerShellでは `Copy-Item .env.example .env` で環境設定をコピーできます。
http://localhost:3000 を開いてください。Goバックエンドは別途起動します。

| 環境変数           | 用途                                                          | 既定値                  |
| ------------------ | ------------------------------------------------------------- | ----------------------- |
| `VITE_API_URL`     | ブラウザからアクセスするAPIのベースURL（`/api/v1`まで含める） | `/api/v1`               |
| `API_PROXY_TARGET` | 開発サーバーが `/api` を転送するバックエンド                  | `http://localhost:8080` |

既存の `.env` にある `NEXT_PUBLIC_API_URL` も、`VITE_API_URL` が未設定なら使用します。
新しい設定では `VITE_API_URL` に移行してください。APIのURLはビルド時に組み込まれるため、変更後は再ビルドが必要です。
旧設定に `NODE_ENV=development` が残っていても、`npm run build` は本番用にビルドします。
相対URLを使う開発環境ではViteがAPIを転送します。別オリジンのAPIに直接接続する場合はCORS設定が必要です。

## 画面と機能

| URL                                                            | 機能                                                             |
| -------------------------------------------------------------- | ---------------------------------------------------------------- |
| `/event`、`/event/:eventId`                                    | イベントの一覧・作成・編集・削除、運営画面へのリンク             |
| `/poster`、`/poster/new`、`/poster/detail/:posterId`           | ポスターの登録・編集・削除、画像・設置場所、回収状況の更新・検索 |
| `/sales/cashier`                                               | カテゴリ別商品選択、明細・数量調整、会計、お釣り表示             |
| `/sales/orders`                                                | 対象イベントの売上履歴、カテゴリ・商品名での検索、合計、記録削除 |
| `/sales/items`、`/sales/items/new`、`/sales/items/:itemId`     | 共通の商品マスターの登録・編集・削除、写真・カテゴリ・検索       |
| `/sales/stocks`、`/sales/stocks/new`、`/sales/stocks/:stockId` | イベントごとの販売商品・価格の登録、説明の編集、削除             |

`/` は `/event`、`/sales` は `/sales/cashier` に移動します。
対象イベントと販売カテゴリの選択は同じタブ内のセッションに保存されます。
PCでは左メニュー、スマートフォンでは下部メニューを使用します。

画像は登録前に最大800×800pxへ縮小・圧縮します。回収状況の更新はその場で保存します。
送信中は操作を無効化し、失敗したフォーム・会計は入力を保持します。
レジで対象イベントを変更すると、会計中の商品をリセットします。未確定の商品があるまま画面移動する場合は確認を表示します。
売上金額・販売数は現在の検索条件に一致する記録を集計し、履歴は新しい順に表示します。

価格の変更・在庫数・決済・注文単位での取消は既存APIにないため追加していません。
お預かり金額とお釣りは計算補助のみで、サーバーには保存しません。
既存APIに会計の重複防止キーがないため、通信が途切れた場合は再送信前に売上履歴を確認してください。

## 検証

```sh
npm run lint
npm test
npm run build
```

テストはテスト用API応答を使用し、既存ルート、登録・更新・削除、画像圧縮、会計・お釣り、送信失敗・二重送信・イベント切り替えを検証します。
実データを使ったAPI・データベース・画像ストレージとの接続確認は別途必要です。

## 本番配信

`npm run build` の出力は `dist/` です。静的ファイルをWebサーバーで配信してください。
以前のNext.js standaloneサーバーは使用しません。`npm start` / `npm run preview` はビルドをローカル確認するためのVite previewです。
本番環境には次の設定が必要です。

1. Vue RouterのURLを直接開いた場合も `index.html` を返す。
2. `VITE_API_URL=/api/v1` を使うなら、`/api` をGoバックエンドへ転送する。
3. HTTPSのフロントエンドから接続するAPI・画像もHTTPSで公開する。

Nginxの例（`root` とバックエンドのアドレスを配信環境に合わせて変更）:

```nginx
server {
    listen 80;
    root /srv/ducks/dist;
    index index.html;

    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
    }

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```
