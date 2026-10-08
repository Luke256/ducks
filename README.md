# Ducks
traP 工大祭用 ツール群

- `frontend`: Vue 3 + TypeScript + Vite。イベント・ポスター・物販・売上管理。
- `backend`: Go / Echo。既存の `/api/v1` APIを使用。

backend の日時は、DB 保存・内部処理を UTC（+00:00）に統一し、API レスポンスを返すときに JST（+09:00）へ変換します。対象は売上の `created_at` と来客数の `bucket_start` です。同じ瞬間を保ったまま、タイムゾーンの表現を変換します。

フロントエンドの起動・環境変数・本番配信については [frontend/README.md](frontend/README.md) を参照してください。
