# 古いPullRequestをクローズするスクリプト

作成してから3ヶ月以上経過しているOpen中のPullRequestをCloseするスクリプト

## 準備

```
$ go mod download
```

## 設定

以下のコマンドで設定ファイルをコピー

```
$ mv config.template.toml config.toml 
```

config.tomlを書き換え

```
[GitHub]
token = ""
owner = ""
repositoryName = ""
ignoreLabel = ""
```



| 名前 | 内容 |
----|---- 
| token  |  GitHub APIのトークンを設定  |
| owner  |  オーナー名  |
| repositoryName  |  リポジトリ名  |
| ignoreLabel  |  削除除外のラベル名  |

## 実行


```
$ go run main.go 
```

