# 構成とデータフロー
web/main.tsx は学習者選択、島、橋、保護者画面を持つ。UIは React 18、操作はボタンとフォーム。ドラゴンの独自SVGは dragon.ts に分離した。Vite版と同じTSXから standalone/app.js を生成し、React・ReactDOMのMITライセンス版を同梱する。通常の家庭利用ではパッケージ取得工程を省ける。

Goは静的ファイル、GET /catalog、POST /api を同一オリジンで配信する。APIはJSONの action で処理を選ぶ。子どもの操作は profiles → select → home → play → hint / answer、保護者操作は parent → create / report。logout はサーバー側のチケットを失効させる。回答時に送るのは run、step、id、3個のパーツ数。報酬・採点・学習者の照合はすべてサーバー側で行う。

server/learning.go は生成・判定・理解度、server/sqlite.go はOSのlibsqlite3との小さなcgo境界、server/schema.sql はDBバージョン1。配備先でのCコンパイルまたは同一OS・CPU向けのビルドが必要になる。ネットワーク経由のGo依存取得は不要。SQLに渡す外部文字列は16進リテラルからTEXTに変換し、数値は検証済みの整数として組み立てる。未知のスキーマバージョンでの起動を拒否する。

家庭の数人を対象に、1プロセス内でAPI処理を直列化する。1学習者・単元につき進行中セッションは最大1件。問題列と現在位置を保存し、ページを閉じても同じ橋から再開できる。未提出の石の配置は再開時にゼロになる。回答履歴・進行・報酬は1トランザクションで更新し、古いステップの再送は409になる。学習者別・単元別に履歴を集計し、完了セッションの報酬合計から成長表示を復元する。

保護者PINは環境変数から読み、子どものPINは学習者IDと保護者PINを用いたHMACで照合する。保護者PINはDBと対にして別途保管する。変更すると現行の学習者PINの再発行が必要になる。チケットはランダムな値のHttpOnly・SameSite=Strict Cookieで8時間有効。失敗回数はメモリ内で制限する。平文HTTPは信頼できる家庭内LANを前提とする。
