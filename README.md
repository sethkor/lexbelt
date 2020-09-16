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

  provision [<flags>] <file>
    Provisions and builds an entire Lex bot including slots, intents and the actual bot
```
###Getting lexbelt
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

###Monobots
A monobot will provision everything you need for a lex bot.  This includes the slots, intents and the bot plus it will build it
lexbelt expects your lex yaml files in the following directory structure for a mono bot to be provisioned.
```
your-lex-workspace
   ├──slots
   ├──intents
   ├──bots
   └──monobot
```

You can take a look at examples/yaml/monobot/OrderFlowersMono.yaml to see an example monobot yaml file like so:
```
LexBotProvisioner:
  lexBotName: OrderFlowersOnceMore
  lexBot: OrderFlowersBot.yaml
  lexIntent: OrderFlowers.yaml
  lexSlot:
    - FlowerTypes.yaml
```

The yaml/json syntax for slots, intents and bots are all based directly of the AWS API Put API calls, so any attribute 
supported by the AWS API will be supported in the API now or in the future can be included in a yaml file

###TODO
* Publishing bots
* Any other feature requested.
* Windows Testing