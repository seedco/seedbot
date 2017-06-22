# Seedbot
<p align="center">
<img src="https://s3-us-west-2.amazonaws.com/seed-assets/seedbot_small.png">
</p>

A Slackbot to interact with Seed.

## Libraries used
### Seed Go Client
To access Seed data, Seedbot uses [seed-go](https://www.github.com/seedco/seedgo), a go client that talks to the [Seed API](http://docs.seed.co/v1.0/docs)

### Nlopes Slack Client
https://github.com/nlopes/slack

# Installation

```
go install github.com/seedco/seedbot

```

# Usage
Start the bot via command line
```
seedbot --seed-token <seed api access token> --slack-token <slack bot token>
```
To get a seed api access token, email support@seed.co
To get a slack bot token, follow instructions to create a slack bot

1. visit https://[domain].slack.com/apps/manage/custom-integrations
2. Goto "Bots"
3. Click "Add Configuration"
4. Select "@seedbot" as the name for the bot
5. The "API Token" is the slack token

## Transactions
A set of commands to retrieve transactions for different date ranges

```
@seedbot transactions last week
```

Example output:
```
06/13/2017    The UPS Store #1603         -27.98
06/13/2017    The UPS Store #1603         -27.98
06/13/2017    Tang 190                    -32.32
06/08/2017    Mobile Check Deposit #0139    1.00
06/07/2017    Ralphs                      -17.24
06/06/2017    Ralphs                      -14.02
06/05/2017    McDonald's F12146            -8.45
```

The format is

```
@seedbot transactions <Date Directive>
```

Where the `Date Directive` as described below, translate into a transactions query with `[Start Date, End Date)`

### Date Directives

| Date Directive | Example | Start Date | End Date |
| -------------- | --------| ---------- | -------- |
| MM/dd/YYYY | `6/12/2017` | 6/12/2017 | 6/13/2017 |
| Month dd, YYYY | `June 12, 2017` | 6/12/2017 | 6/13/2017 |
| MM/yy | `6/12` | 6/12/2017 | 6/13/2017 |
| Month dd | `June 12` |  6/12/2017 | 6/13/2017 |
| YYYY | `2017` | 1/1/2017 | 1/1/2018 |
| today | `today` | today | tomorrow |
| yesterday | `yesterday` | yesterday | today |
| this week&#124;month&#124;year | `this week`<br>`this year` | monday of this week<br>1/1 of this year | monday of next week<br>1/1 of next year |
| last week&#124;month&#124;year | `last week`<br>`last year` | monday of last week<br>1/1 of last year | monday of this week<br>1/1 of this year |


## Balance
Retrieve the balance for your account


```
@seedbot balance
```

Example output:
```
Posted:          $378.24
Pending Debits:  $60.30
Pending Credits: $0.00
Lock Box:        $0.00
Available        $317.94
```
