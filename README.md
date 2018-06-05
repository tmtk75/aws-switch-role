# README

Generate links for page of AWS to switch-role with your ~/.aws/config

```
Usage:
  open $(aws-switch-role | fzf | awk '{print $2}')
```

If you have this profile,
```
[profile foo]
role_arn = arn:aws:iam::112233445566:role/bar
```

this command generates the below link.
```
foo  https://signin.aws.amazon.com/switchrole?roleName=foo&account=112233445566&bar
```
