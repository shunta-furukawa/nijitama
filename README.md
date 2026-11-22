# にじのたまご
Niji no Tamago — A playful learning adventure across five subjects.

家庭のiPadやChromebookで、小数の石をつないで青いドラゴンに橋をかける学習ゲーム。５つの世界のうち、現在遊べるのは小学4年相当の算数「小数のかけ橋」です。学校の学年に関係なく挑戦できます。

## 最初の起動
Macを家庭内Wi-Fiにつなぎ、Finderで **Start.command** をダブルクリックします。初回は保護者用の番号（数字6桁以上）を2回入力します。次回からはダブルクリックだけで起動します。画面に表示された `.local` のURLをiPadのSafariで開き、「共有」→「ホーム画面に追加」で登録してください。名前で開けない端末では一緒に表示されるIPv4のURLを使います。

ソースから起動するときはGo 1.21以上、Cコンパイラ、SQLite 3の開発環境が必要です。配布用フォルダではビルド済みの実行ファイルを起動します。Node.jsは家庭での起動には不要です。

「保護者の部屋」でニックネーム・学年・アイコン・子ども用の4桁の番号を登録。「人をえらぶ」で戻り、自分のアイコンと番号で島へ入ります。橋を5本かけると3エレメントが届きます。途中で閉じたら、同じ子を選び「小数のかけ橋へ」で再開できます。保護者の部屋では学習者ごとの完了回数、時間、直近の正誤、単元の進み具合、復習候補が見られます。

## 開発とテスト
通常起動は同梱Reactと生成済みJSを使います。フロントエンドを編集する場合は別途Node.js 20とnpmを用意し、web/ で npm install、npm run dev を実行します。Goサーバーを起動した状態でViteが表示するローカルURLを開きます。依存取得にはネット接続が必要です。

```sh
go test ./...
go vet ./...
(cd web && npm run build)
node scripts/standalone.cjs
NIJI_TEST_NODE="$(command -v node)" go test -count=1 -v ./server
./scripts/build.sh
```

Goテストは一時DBを使い、生成、採点、理解度、API権限、古い回答の拒否、永続化を検証します。NIJI_TEST_NODE 指定時はReactの画面をDOM環境で操作し、実際のAPIとSQLiteで登録から保護者履歴まで検証します。実ブラウザ検証は別途必要です。

## MacBookへの配備
対象のOS・CPUで `./scripts/build.sh` を実行し、生成された `build/` 一式を配備します。その中の `Start.command` からも起動できます。初期設定はWi-FiのプライベートIPv4で8765番を待ち受けます。有線接続などでは `NIJI_LAN_IP`、ポート変更には `NIJI_PORT` を指定します。起動中はターミナルを開き、停止はCtrl+C。既に検証用サーバーが同じポートを使っている場合は、そのターミナルで先に停止してください。ルーターのポート転送は不要です。

家庭用ランチャーの保存先は `~/Library/Application Support/nijitama/` です。保護者番号は初回に同フォルダの `parent.pin` へ本人だけが読める権限で保存します。`NIJI_DATA` で保存先を変更でき、`NIJI_PARENT_PIN` を明示した場合はそちらが優先されます。既存DBには元と同じ番号を使ってください。一時DBの自動移行や再起動後の自動起動は行いません。

ローカル開発は保護者番号を環境変数に設定して `./scripts/dev.sh` を起動します。既定はループバックの8080番です。Linuxなどで配備する場合は `build/` に移動し、保存先・保護者番号・家庭内のプライベートIPv4を環境変数で指定して `./nijitama` を直接起動します。

## バックアップと復元
sqlite3 コマンドのオンラインバックアップ例です。DBには家庭の履歴が入るので、コピーはアクセスを制限した場所に保存します。保護者番号は別に保管してください。

```sh
umask 077
sqlite3 "${NIJI_DATA:-$HOME/Library/Application Support/nijitama}/nijitama.db" ".backup 'nijitama-backup.db'"
sqlite3 nijitama-backup.db 'PRAGMA integrity_check;'
```

復元はサービスを停止してから行います。既存DBを保持するため、新しい空ディレクトリにバックアップを nijitama.db としてコピーし、NIJI_DATA をそのディレクトリの絶対パスに変更して再起動してください。元の保護者番号で履歴を確認します。

## 実機での確認
iPadの縦横とChromebookで学習者を2人登録し、片方だけでヒント・誤答・中断・再開・5橋完了を試してください。もう一人の残高が0であること、再起動後も元の子に3エレメントと履歴が残ることを確認します。キーボードのみでの操作、動きを減らす設定、3〜5分程度で遊べるかも確認対象です。

[製品理念](docs/product-vision.md) · [構成](docs/architecture.md) · [教材設計](docs/curriculum-design.md) · [次の開発](docs/roadmap.md)
