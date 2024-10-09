# sqluid

最小のrdbmsをつくってみました

## プロジェクトの説明
sqluidは、シンプルなRDBMSの軽量な実装を目指したプロジェクトです。このプロジェクトは、学習目的や小規模なデータベースアプリケーションに適しています。

## インストール手順
以下の手順に従って、プロジェクトをインストールしてください：

1. リポジトリをクローンします：
    ```sh
    git clone https://github.com/rsasada/sqluid.git
    ```
2. プロジェクトディレクトリに移動します：
    ```sh
    cd sqluid
    ```
3. 必要な依存関係をインストールします：
    ```sh
    # 依存関係のインストールコマンド（例）
    make install
    ```

## 使用例
以下のコマンドを使用して、sqluidを実行します：

```sh
# サンプルコマンド
./sqluid
```

## 参考資料

### 有名rdbmsのアーキテクチャー
https://www.fujitsu.com/jp/products/software/resources/feature-stories/postgres/article-index/architecture-overview/

### Lexer
posgreを参考にしました
https://www.postgresql.org/docs/current/sql-syntax-lexical.html

これを `README.md` ファイルに追加してください。
