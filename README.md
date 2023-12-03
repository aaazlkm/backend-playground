# backend-playground

## コマンド

cfn lint

```shell
cfn-lint aws/template/network.yaml
```

install pache

```shell
[ec2-user@ip-10-0-0-29 wordpress]$ sudo systemctl start httpd.service
[ec2-user@ip-10-0-0-29 wordpress]$ sudo systemctl enable httpd.service
```

## CFn についての注意

現状、ACM で Certificate を CFn を使用して作成しドメインに紐付けいてる。
この時、その証明書の検証のための C レコードを作成しないといけないが、これは CFn 上では作成することができない。
ACM のページに行き、GUI 上で C レコードを作成しないといけないことに注意。
詳しくはこちらの「[Setting up DNS validation](https://docs.aws.amazon.com/acm/latest/userguide/dns-validation.html#setting-up-dns-validation)」

## Secrets について

REVIEWDOG_TOKEN: review dog 用の api token
https://github.com/reviewdog/reviewdog
