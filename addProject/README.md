# GitHub の issue 作成時の webhook 経由で GitHub Project に自動で追加

## デプロイ

### 初回

```
$ export GITHUB_APP_PRIVATE_KEY=$(cat private-key.pem)

$ gcloud functions deploy GitHubEvent --set-env-vars GITHUB_APP_PRIVATE_KEY=$GITHUB_APP_PRIVATE_KEY,GITHUB_APP_ID=$GITHUB_APP_ID,GITHUB_OWNER=$GITHUB_OWNER,INSTALLATION_ID=$INSTALLATION_ID,GITHUB_PROJECT_ID=$GITHUB_PROJECT_ID --runtime go119 --entry-point GitHubEvent --region asia-northeast1
--trigger-http
```

### 2 回目以降

```
$ gcloud functions deploy GitHubEvent --region asia-northeast1
```
