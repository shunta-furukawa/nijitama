# 開発引き継ぎ
最優先の動線は学習者選択 → ５つの島 → 青いドラゴンの小数パズル → 報酬 → 保護者履歴。学習に直結する操作と、自分でやり直せる遊びを増やす。学校学年をアクセス制限に使わない。

- 起動・配備は README.md。Go 1.21、cgo、SQLite3、React 18、TypeScript、Vite。標準運用はGoプロセス1つで、外部CDNを使わない。
- 正解・報酬をクライアントで確定しない。所有者と問題番号を確認し、完了と履歴を同じトランザクションで保存する。
- 新単元は curriculum と docs/curriculum-design.md の対応を更新する。固定コンテンツと学習者記録を分ける。
- TSXを編集したら npm run build と、ルートで node scripts/standalone.cjs を実行する。生成済みJSの直接修正は避ける。
- gofmt、go vet ./...、go test ./...、NIJI_TEST_NODE="$(command -v node)" go test -count=1 ./server、scripts/build.sh を確認する。実機未検証項目は区別する。
- 名前・PIN・DB・バックアップはコミットしない。リモートへのpush、PR、Issue、クラウド配備、重要なファイル削除、破壊的DB変更には事前確認が必要。
