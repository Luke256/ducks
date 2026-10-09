# Ducks フロントエンド

学園祭のイベント・ポスター・物販・売上・来場者数を管理する Vue 3 + TypeScript + Vue Router + Vite のSPAです。
既存のGoバックエンドとデータ形式、画面URLを引き継いでいます。

## 開発

Node.js 24系を使用してください。buildpackが別のメジャーバージョンを選ばないよう、`package.json` の `engines.node` を `24.x` に指定しています。

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
| `/visitors`                                                    | イベント別の来場者数カウント、まとめて追加、現在の10分間の訂正、累計・履歴 |
| `/sales/items`、`/sales/items/new`、`/sales/items/:itemId`     | 共通の商品マスターの登録・編集・削除、画像の差し替え・カテゴリ・検索 |
| `/sales/stocks`、`/sales/stocks/new`、`/sales/stocks/:stockId` | イベントごとの販売商品・価格の登録、説明の編集、削除             |

`/` は `/event`、`/sales` は `/sales/cashier` に移動します。
対象イベントと販売カテゴリの選択は同じタブ内のセッションに保存されます。
PCでは左メニュー、スマートフォンでは下部メニューを使用します。

画像は送信前に最大800×800pxへ縮小・圧縮します。回収状況の更新はその場で保存します。
商品画像は商品マスターの詳細から「編集する」で差し替えられます。新画像を選ばなければ現在の画像を保持します。
商品画像は全イベントで共通です。商品情報と画像は別APIで保存するため、画像更新に失敗した場合は商品情報だけ保存されることがあります。選択画像と入力は保持するので、確認後に再度保存してください。

ポスターには写真を1〜10枚登録できます。複数選択や追加選択ができ、選択した写真は送信前に削除できます。編集では既存の写真と追加した写真を同じ一覧に表示します。既存の写真は「削除」ボタンを押すと赤枠が付き、ボタンが「元に戻す」に変わります。追加した写真は「削除」ボタンを押し、確認ダイアログで承認すると一覧から取り除きます。既存の写真の削除は保存時に反映されます。全写真を差し替える場合は、既存の写真をすべて削除対象にして新しい写真を追加してください。保存後に0枚または11枚以上になる変更は送信しません。
ポスター名・設置場所と写真も別APIで保存します。写真の更新に失敗した場合は選択を保持し、名前・設置場所が保存済みであれば画面に表示します。
送信中は操作を無効化し、失敗したフォーム・会計は入力を保持します。
レジで対象イベントを変更すると、会計中の商品をリセットします。未確定の商品があるまま画面移動する場合は確認を表示します。
売上金額・販売数は現在の検索条件に一致する記録を集計し、履歴は新しい順に表示します。

来場者数はサーバーの現在時刻で10分ごとに記録され、履歴は日本時間の新しい順に表示します。
「−1人（訂正）」は現在の10分間だけを減らします。過去の時間帯の編集はAPIにありません。
来場者数の送信にも重複防止キーがないため、通信が途切れた場合は更新して人数を確認してから再送信してください。

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

`npm run build` の出力は `dist/` です。buildpackでNode.jsアプリとして配信する場合は、`frontend` をアプリのルートにして次の順序で実行してください。

```sh
npm ci --include=dev
npm run build
npm prune --omit=dev
npm start
```

`npm start` は本番依存の `serve` で `dist/` を配信します。ViteやTypeScriptなどの開発用依存を削除した状態でも起動できます。環境変数 `PORT` を使用し、未設定の場合は3000番で待ち受けます。Vue RouterのURLを直接開いた場合も `index.html` を返します。
`npm run preview` は開発環境でビルドを確認するためのVite previewです。

buildpackのインストール段階では開発用依存も必要です。`vite: not found` がビルド中に出る場合は `npm ci --omit=dev` や `NPM_CONFIG_PRODUCTION=true` で省かれていないか確認し、`npm ci --include=dev` 相当でインストールしてください。
Paketo Node.js buildpackの場合は `BP_NODE_RUN_SCRIPTS=build` を設定します。リポジトリのルートからビルドするなら `BP_NODE_PROJECT_PATH=frontend` も設定してください（[Paketo公式ドキュメント](https://paketo.io/docs/howto/nodejs/)）。

`node: error while loading shared libraries: libatomic.so.1` は、Node.jsを起動するLinux環境に共有ライブラリが足りないエラーです。[Node.js公式資料](https://github.com/nodejs/node/blob/main/BUILDING.md#official-binary-platforms-and-toolchains)では、Node 25以降の公式Linuxバイナリに `libatomic` が必要とされています。バージョン指定を24系に絞ったうえで、イメージを再ビルドしてください。
Paketoの `BP_NODE_VERSION` は `package.json` より優先されます。設定済みなら `24.*` に変更するか、上書き設定を外してください。ほかのbuildpackでも、管理画面のNodeバージョン指定があれば24系に合わせ、ログで実際の選択結果を確認します。
24系でも同じエラーが出る場合は、選ばれたバイナリやベースイメージを確認し、配信環境側で `libatomic` を含むイメージへ更新する必要があります。npmパッケージの追加ではOSの共有ライブラリは補えません。

APIは静的配信サーバーから転送しません。`VITE_API_URL` に公開APIのURL（例：`https://api.example.com/api/v1`）をビルド時に設定してください。同一オリジンの `/api/v1` を使う場合は、配信基盤側でバックエンドへの転送を設定します。`API_PROXY_TARGET` は開発時のみの設定です。

別のWebサーバーで `dist/` を配信する場合は、次の設定が必要です。

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
