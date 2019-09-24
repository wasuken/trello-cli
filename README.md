# trello-cli

## install

### build

```bash
$ git clone https://github.com/wasuken/trello-cli.git
$ cd trello-cli
$ go build
```

### setting config file

```bash
mkdir -p ~/.config/trello-cli
touch ~/.config/trello-cli/config.toml
```

#### write config toml

```toml
[API]
Apikey = <your-key>
Token = <your-token>
Member = <your-id>
```

## Basic Usage


### list boards

```
$ trello-cli boards
```

output

```
name:Pythonセキュリティプログラミング, id:<id>
name:世界で戦うプログラミング力を鍛える本, id:<id>
name:日常, id:<id>
```

### list boards

```
$ trello-cli lists <board id>
```

output

```
Pythonセキュリティプログラミング
name:未着手, id:<list-id>
    name:６章, id:<card-id>
name:着手, id:<list-id>
    name:５章, id:<card-id>
name:完了, id:<list-id>
    name:環境構築, id:<card-id>
    name:１章から４章までザッと, id:<id>
```

## others

* addCard

add the <card-id> to the <list-id>

```bash
$ trello-cli addCard <list-id> <card-name> <card-description>
```

* removeCard

remove the <card-id>

```bash
$ trello-cli removeCard <card-id>
```

* moveCard

move the <card-id> to <list-id>

```bash
$ trello-cli moveCard <card-id> <list-id>
```


## TODO

* if festival -> using [urfave/cli](https://github.com/urfave/cli)
