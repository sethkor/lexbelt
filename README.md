# lexbelt
A tool to provision Amazon Lex using YAML or JSON

Amazon Lex is amazing, except it has very low automation.  There's no CloudFormation for it CFN for it.

Also, when trying to create or update anything, the CLI wants you to pass in a checksum so it can figure out if it needs to update or create, this can be annoying.

`lexbelt` fixes all this.

```
usage: lexbelt [<flags>] <command> [<args> ...]

AWS Lex CLI utilities

Flags:
  -h, --help             Show context-sensitive help (also try --help-long and --help-man).
      --profile=PROFILE  AWS credentials/config file profile to use
      --region=REGION    AWS region
  -v, --version          Show application version.

Commands:
  help [<command>...]
    Show help.

  put-slot-type --name=NAME <file>
    Adds or updates a slot type

  put-intent --name=NAME <file>
    Adds or updates an intent

  put-bot --name=NAME [<flags>] <file>
    Adds or updates a bo
```
## Getting lexbelt
Easiest way to install if you're on a Mac or Linux (amd64 or arm64)  is to use [Homebrew](https://brew.sh/)

Type:

```
brew tap sethkor/tap
brew install lexbelt
```

For other platforms take a look at the releases in Github.  I build binaries for:

|OS            | Architecture                           |
|:------------ |:-------------------------------------- |
|Mac (Darwin)  | amd64 (aka x86_64)                     |
|Linux         | amd64, arm64, 386 (32 bit) |
|Windows       | amd64, 386 (32 bit)                   |

Let me know if you would like a particular os/arch binary regularly built.

```
your-lex-workspace
   ├──slots
   ├──intents
   └──bots

```

The yaml/json syntax for slots, intents and bots are all based directly of the AWS API Put API calls, so any attribute 
supported by the AWS API will be supported in the API now or in the future can be included in a yaml file

The tool won't handle mixing json and yaml within a bot (e.g. yaml bot, json intent or slot, etc), pick one and stick with it.

## Lex Bot Odd Behaviour
AWS Lex does do some weird stuff.

### Bot versioning
For instance if your attempt to create a new version of a slot or an intent and
nothing has actually changed, you'll get the last version number returned.  Smart.  However, if you try to create a new
version of abot and nothing has actualy changed, you'll get a new version number.  This inconsistency is a bit annoying. 

### 409 ConflictException
Once a new bot is published, there seems to be some asynchronous AWS magik still going on in the background.  Any 
subsequent request to publish the bot again can trigger a HTTP response 409 ConflictException.  Wait a minute and try it
again, it will work. I'm guessing this is related to the build process.

### AWS UI doesn't refresh or is slow
This happens quite often.  Some times it requires waiting for the page to load or in the case of an alias, clicking on
other settings before clicking on Aliases again.

### TODO
* Any other feature requested.
* Windows Testing