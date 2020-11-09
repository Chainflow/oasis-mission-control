# Validator Mission Control

**Validator Mission Control** provides a comprehensive set of metrics and alerts for Oasis validator node operators. We utilized the power of Grafana + Telegraf and extended the monitoring & alerting with a custom built go server. It also sends emergency alerts and calls to you based on your pagerduty account configuration.

It should be installed on a validator node directly. These instructions assume the user will install Validator Mission Control on the validator.

## Install Prerequisites
- **Go 14.x+**
- **Docker 19+**
- **Grafana 6.7+**
- **InfluxDB 1.7+**
- **Telegraf 1.14+**
- **Oasis node**

## 1.You can run this installation script to setup the monitoring tool

- One click installation to download grafana, telegraf and influxdb.
- Also it will create the databases(oasis and telegraf). If you configure different database names make sure to create those in influxDB.
- If you are running telegraf on someother node run script by giving --remote-hosted flag so that the script don't start telegraf in current node (ex: ./install_script.sh --remote-hosted)
- Here you can find
[script file](https://github.com/Chainflow/oasis-mission-control/script/install_script.sh).
- To run script 
```bash
chmod +x install_script.sh
./install_script.sh

```
- After installation you just need to configure the config.toml and start the monitoring tool server.
- Follow further steps to setup grafana dashboards.

## 2. Install Manually

### Install Grafana for Ubuntu
Download the latest .deb file and extract it by using the following commands

```sh
$ cd $HOME
$ sudo -S apt-get install -y adduser libfontconfig1
$ wget https://dl.grafana.com/oss/release/grafana_6.7.2_amd64.deb
$ sudo -S dpkg -i grafana_6.7.2_amd64.deb
```

Start the grafana server
```sh
$ sudo -S systemctl daemon-reload

$ sudo -S systemctl start grafana-server

Grafana will be running on port :3000 (ex:: https://localhost:3000)
```

### Install InfluxDB and Telegraf

```sh
$ cd $HOME
$ wget -qO- https://repos.influxdata.com/influxdb.key | sudo apt-key add -
$ source /etc/lsb-release
$ echo "deb https://repos.influxdata.com/${DISTRIB_ID,,} ${DISTRIB_CODENAME} stable" | sudo tee /etc/apt/sources.list.d/influxdb.list
```

Start influxDB

```sh
$ sudo -S apt-get update && sudo apt-get install influxdb
$ sudo -S service influxdb start

The default port that runs the InfluxDB HTTP service is :8086
```

**Note :** If you want cusomize the configuration, edit `influxdb.conf` at `/etc/influxdb/influxdb.conf` and don't forget to restart the server after the changes. You can find a sample 'influxdb.conf' [file here](https://github.com/jheyman/influxdb/blob/master/influxdb.conf).


Start telegraf

```sh
$ sudo -S apt-get update && sudo apt-get install telegraf
$ sudo -S service telegraf start
```

## Install and configure the Validator Mission Control

### Get the code

```bash
$ git clone https://github.com/Chainflow/oasis-mission-control.git
$ cd oasis-mission-control
$ cp example.config.toml config.toml
```

## After all installations you have to follow below steps to configure the monitoring tool and grafana dashboards

### Configure the following variables in `config.toml`

- *tg_chat_id*

    Telegram chat ID to receive Telegram alerts, required for Telegram alerting.
    
- *tg_bot_token*

    Telegram bot token, required for Telegram alerting. The bot should be added to the chat and should have send message permission.

- *email_address*

    E-mail address to receive mail notifications, required for e-mail alerting.

- *sendgrid_token*

    Sendgrid mail service api token, required for e-mail alerting.
    
- *missed_blocks_threshold*

    Configure the threshold to receive  **Missed Block Alerting**, e.g. a value of 10 would alert you every time you've missed 10 consecutive blocks.

- *alert_time1* and *alert_time2*

    These are for regular status updates. To receive validator status daily (twice), configure these parameters in the form of "02:25PM". The time here refers to UTC time.

- *voting_power_threshold*

    Configure the threshold to receive alert when the voting power reaches or drops below of the threshold given.

- *num_peers_threshold*

    Configure the threshold to get an alert if the no.of connected peers falls below the threshold.

- *enable_telegram_alerts*

    Configure **yes** if you wish to get telegram alerts otherwise make it **no** .

- *enable_email_alerts*

    Configure **yes** if you wish to get email alerts otherwise make it **no** .

- *val_operator_addr*

    Operator address of your validator which will be used to get staking, delegation and distribution rewards.

- *validator_hex_addr*

    Validator hex address useful to know about last proposed block, missed blocks and voting power.

- *staking_denom*

    Give stakig denom to display along with self delegation balance (ex:uatom or umuon)

- *pagerduty_email*

    Give mail address of pager duty service to send alerts of emergency missed blocks.
    Note : Have to give mail address which was generated after creation of a service in pager duty.

    You can refer this to know about pagerduty (https://www.pagerduty.com/)

- *emergency_missed_blocks_threshold*

    Give threshold to notify about continuous missed blocks to your pager duty account. so that it will send mails, messages and makes you a call about alerts.

- *block_diff_threshold*

    An integer value to receive Block difference alerts, e.g. a value of 2 would alert you if your validator falls 2 or more blocks behind the chain's current block height.

- *epoch_diff_threshold*

     An integer value to receive woker epoch difference alerts, e.g. a value of 2 would alert you if your validator falls 2 or more behind the chain's worker epoch number.

- *network_url*

    Give the network url, whic is useful to gather information like network height and epoch number. 

- *network_node_name*

    Give the node name of network which will be used as a query params in network url.

    **Note:**  We have configured the SimplyVc api server and running it on our mainnet to get the network block height and epoch number. So the given network_url and network_node_name has taken based on the configuration we did. 
    If you want to give another network url and node name you can configure it otherwise you can just keep which we have provided.

After populating config.toml, check if you have connected to influxdb and created a database which you are going to use.

If your connection throws error "database not found", create a database

```bash
$   cd $HOME
$   influx
>   CREATE DATABASE db_name    (ex : CREATE DATABASE oasis)
> exit
```

### Then build and run the monitoring binary

```bash
$ go build -o oasis-chain-monit && ./oasis-chain-monit
```

### Run using docker
```bash
$ docker build -t cfv .
$ docker run -d --name oasis-chain-monit cfv
```

We have finished the installation and started the server. Now lets configure the Grafana dashboard.

## Grafana Dashboards

Validator Mission Control provides three dashboards

1. Validator Monitoring Metrics (These are the metrics which we have calculated and stored in influxdb)
2. System Metrics (These are the metrics related to the system configuration which come from telegraf)
3. Summary (Which gives quick overview of validator and system metrics)


### 1. Validator monitoring metrics
The following list of metrics are displayed in this dashboard.

- Hex Address : Displays hex address of the validator.
- Address : Displays public key of the validator.
- Oasis Node Status :  Displays whether the node is running or not in the form of UP and DOWN.
- Validator Status :  Displays the validator health. Shows Voting if the validator is in active state or else Jailed.
- Oasis Node Version : Displays the version of gaia currently running.
- Block Time Difference : Displays the time difference between previous block and current block.
- Current Block Height - Validator : Displays the current block height committed by the validator.
- Current Block Height - Network : Displays the current block height of a network.
- Block Diffrence : Displays the difference of validator and network block heights.
- Last Missed Block Range : Displays the continuous missed blocks range based on the threshold given in the config.toml
- Blocks Missed In last 48h : Displays the count of blocks missed by the validator in last 48 hours.
- No.of Peers : Displays the total number of peers connected to the validator.
- Voting Power : Displays the voting power of the validator.
- Escrow : Displays delegation balance of the validator.
- Address Balance : Displays the account balance of the validator.
- Worker Epoch Number - Validator: Displays the worker epoch number of the validator.
- Worker Epoch Number - Network : Displays the worker epoch number of a network.
- Worker Epoch Difference : Displays the difference between validator and network worker epoch numbers.
- No Of Proposals : Displays number of consenus signed proposals (Will get this from prometheus metrics)
- Abci Db Size : Displays the db size (MiB) (Accroding to prometheus metrics)


**Note:** The above mentioned metrics will be calculated and displayed according to the validator address which will be configured in config.toml.

For alerts regarding system metrics, a Telegram bot can be set up on the dashboard itself. A new notification channel can be added for the Telegram bot by clicking on the bell icon on the left hand sidebar of the dashboard. 

This will let the user configure the Telegram bot ID and chat ID. **A custom alert** can be set for each graph in a Grafana dashboard by clicking on the edit button and adding alert rules.

### 2. System Monitoring Metrics
These are provided by telegraf.

-  For the list of system monitoring metrics, you can refer `telgraf.conf`. You can replace the file with your original telegraf.conf file which will be located at /telegraf/etc/telegraf (installation directory).
 
 ### 3. Summary Dashboard
This dashboard displays a quick information summary of validator details and system metrics. It includes following details.

- Validator identity (Hex Address and Address)
- Validator summary (Oasis Node Status, Validator Status, No.Of peers, Current validator Block Height, Worker Epoch Number, Block Height Difference and Epoch Difference) are the metrics being displayed from Validator details.
- CPU usage, RAM Usage, Memory usage and information about disk usage are the metrics being displayed from System details.

## How to import these dashboards in your Grafana

### 1. Login to your Grafana dashboard
- Open your web browser and go to http://<your_ip>:3000/. `3000` is the default HTTP port that Grafana listens to if you haven’t configured a different port.
- If you are a first time user type `admin` for the username and password in the login page.
- You can change the password after login.

### 2. Create Datasources

- Before importing the dashboards you have to create datasources of InfluxDBTelegraf, InfluxDBVCF and Prometheus Datasources.
- To create datasoruces go to configuration and select Data Sources.
- After that you can find Add data source, select InfluxDB from Time series databases section.
- Then to create `InfluxDBVCF` Datasource, follow these configurations. In place of name give InfluxDBVCF, in place of URL give url of influxdb where it is running (ex : http://ip_address:8086). Finaly in InfluxDB Details section give Database name as `oasis` (If you haven't created a database with different name). You can give User and Password of influx if you have set anthing, otherwise you can leave it empty.
- After this configuration click on Save & Test. Now you have a working Datasource of InfluxDBVCF.

- Repeat same steps to create `InfluxDBTelegraf` Datasource. In place of name give InfluxDBTelegraf, give URL of telegraf where it is running (ex: http://ip_address:8086). Give Database name as telegraf, user and password (If you have configured any). 

- After this configuration click on Save & Test. Now you have a working Datasource of InfluxDBTelegraf.

- You have to repeat the same steps to create Prometheus Datasource, but you need to select `Prometheus` Data source from the list and configure accrodingly. In place of name give `Prometheus`, and URL of prometheus which was running on your validator (ex : http://ip_address:9090). 

- After this configuration click on Save & Test. Now you have a working Datasource of Prometheus. 

### 3. Import the dashboards
- To import the json file of the **validator monitoring metrics** click the *plus* button present on left hand side of the dashboard. Click on import and load the validator_monitoring_metrics.json present in the grafana_template folder. 

- Select the datasources and click on import.

- To import **system monitoring metrics** click the *plus* button present on left hand side of the dashboard. Click on import and load the system_monitoring_metrics.json present in the grafana_template folder.

- While creating this dashboard if you face any issues at valueset, change it to empty and then click on import by selecting the datasources.

- To import **summary**, click the *plus* button present on left hand side of the dashboard. Click on import and load the summary.json present in the grafana_template folder.

- *For more info about grafana dashboard imports you can refer https://grafana.com/docs/grafana/latest/reference/export_import/*


## Alerting (Telegram and Email)
 A custom alerting module has been developed to alert on key validator health events. The module uses data from influxdb and trigger alerts based on user-configured thresholds.

 - Alert when the missed blocks count reaches or exceeds **missed_blocks_threshold** which is user configured in *config.toml*.
 - Alert when no.of peers count falls below of **num_peers_threshold** which is user configured in *config.toml*
- Alert when the oasis node status is not running on validator instance.
- Alert when the block difference between validator and network reaches or exceeds the **block_diff_threshold** which is user configured in config.toml.
- Alert about validator health, i.e. whether it's voting or jailed. You can get alerts twice a day based on the time you have configured **alert_time1** and **alert_time2** in *config.toml*. This is a useful sanity check, to confirm the validator is voting (or alerting you if it's jailed).
- Alert when the voting power of your validator drops below **voting_power_threshold** which is user configured in *config.toml*
- Alert when the worker epoch difference between validator and network reaches or exceeds the **epoch_diff_threshold** which is user configured in config.toml.
 


## Steps to creat telegram bot
To create telegram bot and to configure tg_bot_token and tg_chat_id, one can follow the below steps.

- In the first step search for `@BotFather` , go to that chat and click on start.
- Then do `/newbot` and the bot will respond, just follow those instcructions to create a bot.
- From `example 3` you can observe tg_bot_token which is marked, and also to go to the bot chat click on the `Testing_intro_bot`.

| example 1     | example 2      | example 2      |
|------------|-------------|-------------|
| <img src="https://github.com/Chainflow/oasis-mission-control/blob/implementation/images/start.jpg" width="230"> | <img src="https://github.com/Chainflow/oasis-mission-control/blob/implementation/images/new_bot.jpg" width="250"> | <img src="https://github.com/Chainflow/oasis-mission-control/blob/implementation/images/bot_created.jpg" width="250"> |

- To know your telegram chat id, search for `@my_id_bot` then you can find bot with name `What's my ID`.
- Use this bot to get your personal ID or add it to any group to see its ID. Then that will become your tg_chat_id .

## Feedback and Questions

Please feel free to create issues in this repo to ask questions and/or suggest additional metrics, alerts or features.
