# lexbelt
A tool to provision Amazon Lex using YAML or JSON

Amazon Lex is amazing, except it has very low automation.  There's no CloudFormation for it CFN for it.

Also, when trying to create or update anything, the CLI wants you to pass in a checksum so it can figure out if it needs to update or create, this can be annoying.

Lexbelt fixes all this.

```
usage: lexbelt [<flags>] <command> [<args> ...]

AWS Lex CLI utilities

Flags:
  -h, --help             Show context-sensitive help (also try --help-long and --help-man).
      --profile=PROFILE  AWS credentials/config file profile to use
      --region=REGION    AWS region
      --verbose          Verbose Logging - not implemented yet
  -v, --version          Show application version.

Commands:
  help [<command>...]
    Show help.

  put-slot-type --name=NAME <file>
    Adds or updates a slot type

  put-intent --name=NAME <file>
    Adds or updates an intent

  put-bot --name=NAME [<flags>] <file>
    Adds or updates a bot
```