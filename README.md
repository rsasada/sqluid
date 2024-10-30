# sqluid

最小のrdbmsをつくってみました

## プロジェクトの説明
sqluidは、シンプルなRDBMSの軽量な実装を目指したプロジェクトです。このプロジェクトは、学習目的や小規模なデータベースアプリケーションに適しています。


## 参考資料

### postgresqlのアーキテクチャー
https://www.fujitsu.com/jp/products/software/resources/feature-stories/postgres/article-index/architecture-overview/
https://edbjapan.com/webinar/PostgreSQ_Basics_Architecture_220202.pdf

### Lexer
posgreを参考にしました
https://www.postgresql.org/docs/current/sql-syntax-lexical.html

### Parser
私の過去プロジェクトである自作シェルのParserを参考にしました。Thanks you, JohnSan

https://github.com/rsasada/MiniShell

### BTree
GoによるBtreeの実装
https://qiita.com/oko1977/items/822c0b3168716ebfbf0c

BTeee入門
https://qiita.com/kiyodori/items/f66a545a47dc59dd8839

### Backend(executer)

disk上でのデータの永続化についてはこちらを参考にしました
https://cstack.github.io/db_tutorial/
