# Oasis Mission Control

**Oasis Mission Control** provides a comprehensive set of monitoring metrics and alerts for Oasis validator node operators. We utilized the power of Grafana + Telegraf and extended the monitoring & alerting with a custom built go server. The system also sends configurable emergency alerts and calls to you based on your Pager Duty account configuration.

Ideally, Oasis Mission Control should be installed on a validator node directly. These instructions assume the user will install Oasis Mission Control on the validator. [Contact us](https://chainflow.io/contact) if you'd like to discuss how to run Oasis Mission Control on a remote server instead.

## Install Prerequisites
- **Go 15.x+**
- **Grafana 6.7+**
- **InfluxDB 1.7+**
- **Telegraf 1.14+**
- **Oasis node**

## 1.You can run this installation script to setup the monitoring tool

- One click installation to download grafana, telegraf and influxdb.
- The script will also create the required databases (oasis and telegraf). If you configure different database names make sure to create those in influxDB.
- If you are running telegraf on a remote node run the script using the --remote-hosted flag so the script don't start telegraf on the current node (ex: ./install_script.sh --remote-hosted)
- Here you can find the
[script file](https://github.com/Chainflow/oasis-mission-control/script/install_script.sh).
- To run script 
```bash
chmod +x install_script.sh
./install_script.sh

```
- After installation, configure the config.toml and start the monitoring tool server.
- Follow the additional steps listed below to import and set-up the grafana dashboards.

## 2. Install Manually

### Install Grafana for Ubuntu
Download the latest .deb file and extract it by using the following commands.

```sh
$ cd $HOME
$ sudo -S apt-get install -y adduser libfontconfig1
$ wget https://dl.grafana.com/oss/release/grafana_6.7.2_amd64.deb
$ sudo -S dpkg -i grafana_6.7.2_amd64.deb
```

Start the grafana server.
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

**Note :** If you would like to cusomize the configuration, edit `influxdb.conf` at `/etc/influxdb/influxdb.conf` and restart the server after saving the changes. You can find a sample 'influxdb.conf' [file here](https://github.com/jheyman/influxdb/blob/master/influxdb.conf).

Start telegraf

```sh
$ sudo -S apt-get update && sudo apt-get install telegraf
$ sudo -S service telegraf start
```

## Install and configure Oasis Mission Control

### Get the code

```bash
$ git clone https://github.com/Chainflow/oasis-mission-control.git
$ cd oasis-mission-control
$ cp example.config.toml config.toml
```

## After the installation completes, follow the below steps to configure the monitoring tool and grafana dashboards

### Configure the following variables in `config.toml`

- *tg_chat_id*

    Telegram chat ID to receive Telegram alerts, required for Telegram alerting.
    
- *tg_bot_token*

    Telegram bot token, required for Telegram alerting. The bot should be added to the chat you'd like alerts sent to and should have send message permission in that chat.

    - To create a telegram bot and to find `tg_bot_token` and `tg_chat_id` [follow these steps](##Steps-to-creat-telegram-bot) .

- *email_address*

    E-mail address to receive e-mail notifications, required for e-mail alerting.

- *sendgrid_token*

    Sendgrid mail service api token, required for e-mail alerting.
    
- *missed_blocks_threshold*

    Configure the threshold to receive  **Missed Block Alerting**, e.g. a value of 10 would alert you every time you've missed 10 consecutive blocks.

- *alert_time1* and *alert_time2*

    These are for regular status updates. To receive validator status daily (twice), configure these parameters in the form of "02:25PM". The time here refers to UTC time.

- *voting_power_threshold*

    Configure the threshold to receive alert when the voting power reaches or drops below the indicated threshold.

- *num_peers_threshold*

    Configure this threshold to receive an alert if the number of connected peers falls below the threshold.

- *enable_telegram_alerts*

    Configure **yes** if you wish to get telegram alerts otherwise make it **no** .

- *enable_email_alerts*

    Configure **yes** if you wish to get email alerts otherwise make it **no** .

- *validator_addr*

    Address of the validator you'd like to monitor.

- *validator_hex_addr*

    Validator hex address, which is useful to determine last proposed block, missed blocks and voting power.
    
    - To find this hex address, run this command on your oasis node.
    (`oasis-node identity tendermint show-consensus-address --datadir /root/dir/node`)

- *pagerduty_email*

   Email address of your configred pager duty service to receive alerts of emergency missed blocks.
    Note : Have to give email address which was generated by the creation of a service in pager duty.

    Learn more about Pager Duty [here](https://www.pagerduty.com/). And if you know of an effective and open source Pager Duty alternative pleae [contact us](https://chainflow.io/contact).

- *emergency_missed_blocks_threshold*

    Give threshold to notify you about continuous missed blocks. This alert will be sent to your configured pager duty email address and is used to trigger pager duty workflows.

- *block_diff_threshold*

    An integer value to receive block difference alerts, e.g. a value of 2 would alert you if your validator falls 2 or more blocks behind the chain's current block height.

- *epoch_diff_threshold*

     An integer value to receive woker epoch difference alerts, e.g. a value of 2 would alert you if your validator falls 2 or more behind the chain's worker epoch number.

- *network_url*

    The network url this is used to provide network information like network height and epoch number. 

- *network_node_name*

    The node name of the network which will be used to query params in network url.

    **Note:**  So the default network_url and network_node_name uses the SimplyVC api server to get the network block height and epoch number. If you want to use another network url and node name, change the default values. Otherwise keep the default values to use the SimplyVC api server.

After populating config.toml, check and see if you have connected to influxdb and created a database.

If your connection throws error "database not found", create a database

```bash
$   cd $HOME
$   influx
>   CREATE DATABASE oasis
> exit
```

### Then build and run the monitoring binary

```bash
$ go build -o oasis-chain-monit && ./oasis-chain-monit
```

We have finished the installation and started the server. Now let's configure the Grafana dashboard.

## Grafana Dashboards

Oasis Mission Control provides three dashboards

1. Summary (Which gives quick look at validator and system health)
2. Validator Monitoring Metrics (These are the more detailed validator metrics which are calculated and stored in influxdb)
3. System Metrics (These are the usual system metrics that telegraf provides)

### 1. Summary Dashboard
This dashboard displays a quick quick look at validator and system health. It includes following details.

- Validator identity (Hex Address and Address)
- Validator summary (Oasis Node Status, Validator Status, No.Of peers, Current validator Block Height, Worker Epoch Number, Block Height Difference and Epoch Difference)
- CPU usage, RAM Usage, Memory usage and disk usage are displayed from System details.


### 2. Validator monitoring metrics
The following list of metrics are displayed in this dashboard.

- Hex Address : Displays hex address of the validator.
- Address : Displays public key of the validator.
- Oasis Node Status :  Displays whether the node is running or not in the form of UP and DOWN.
- Validator Status :  Displays the validator health. Shows Voting if the validator is in active state or else Jailed.
- Oasis Node Version : Displays the version of oasis_node currently running.
- Block Time Difference : Displays the time difference between previous block and current block on the validator.
- Current Block Height - Validator : Displays the current block height committed by the validator.
- Current Block Height - Network : Displays the current block height of a network.
- Block Diffrence : Displays the difference between validator and network block heights.
- Last Missed Block Range : Displays the last continuous missed blocks range based on the threshold given in the config.toml
- Blocks Missed In last 48h : Displays the count of blocks missed by the validator in last 48 hours.
- No.of Peers : Displays the total number of peers connected to the validator.
- Voting Power : Displays the voting power of the validator.
- Escrow : Displays the delegation balance of the validator.
- Address Balance : Displays the account balance of the validator.
- Worker Epoch Number - Validator: Displays the current worker epoch number of the validator.
- Worker Epoch Number - Network : Displays the current worker epoch number of a network.
- Worker Epoch Difference : Displays the difference between validator and network epoch numbers.
- No Of Proposals : Displays number of consenus signed proposals (Retreived from prometheus metrics)
- Abci Db Size : Displays the db size (MiB) (Accroding to prometheus metrics)

   **Note:** 
    - To get `Abci Db Size` and `No Of Proposals` you have to enable the prometheus flag before starting oasis node. (ex: `--metrics.mode pull` `--metrics.address 0.0.0.0:3001` you have to add these falgs while running oasis node) .
   - To get `Voting Power` you should have enabled the consensusrpc flag. (ex: `--worker.consensusrpc.enabled` should be added when oasis-node is started.)


**Note:** The above mentioned metrics will be calculated and displayed according to the validator address which was configured in config.toml.

For alerts regarding system metrics, a Telegram bot can be set up on the dashboard itself. A new notification channel can be added for the Telegram bot by clicking on the bell icon on the left hand sidebar of the dashboard. 

This will let the user configure the Telegram bot ID and chat ID. **A custom alert** can be set for each graph in a Grafana dashboard by clicking on the edit button and adding alert rules.

### 3. System Monitoring Metrics
These are provided by telegraf.

-  For the list of system monitoring metrics, you can refer `telgraf.conf`. You can replace the file with your original telegraf.conf file which will be located at /telegraf/etc/telegraf (installation directory).

## How to import these dashboards in your Grafana

### 1. Login to your Grafana dashboard
- Open your web browser and go to http://<your_ip>:3000/. `3000` is the default HTTP port that Grafana listens to if you haven’t configured a different port.
- If you are a first time user type `admin` for the username and password in the login page.
- You can change the password after login.

### 2. Create Datasources

- Before importing the dashboards you have to create datasources for InfluxDBTelegraf, InfluxDBVCF and Prometheus Datasources.
- To create datasoruces go to configuration and select Data Sources.
- After that you can find Add data source, select InfluxDB from Time series databases section.
- Then to create `InfluxDBVCF` Datasource, follow these configurations. In place of name give InfluxDBVCF, in place of URL give url of influxdb where it is running (ex: http://ip_address:8086). Finaly in InfluxDB Details section give Database name as `oasis` (If you haven't created a database with different name). You can give User and Password of influx if you have set anthing, otherwise you can leave it empty.
- After this configuration click on Save & Test. Now you have a working InfluxDBVCF datasource.

- Repeat same steps to create an `InfluxDBTelegraf` Datasource. In place of name give InfluxDBTelegraf, give URL of telegraf where it is running (ex: http://ip_address:8086). Give Database name as telegraf, user and password (If you have configured any). 

- After this configuration click on Save & Test. Now you have a working InfluxDBTelegraf datasource.

- You have to repeat the same steps to create Prometheus Datasource, but you need to select `Prometheus` Data source from the list and configure accrodingly. In place of name give `Prometheus` and the URL of prometheus which was running on your validator (ex: http://ip_address:9090). 

- After this configuration click on Save & Test. Now you have a working Prometheus datasource.

### 3. Import the dashboards

- To import **summary**, click the *plus* button present on left hand side of the dashboard. Click on import and load the summary.json present in the grafana_template folder.

- To import the json file of the **validator monitoring metrics** click the *plus* button present on left hand side of the dashboard. Click on import and load the validator_monitoring_metrics.json present in the grafana_template folder. 

- Select the datasources and click on import.

- To import **system monitoring metrics** click the *plus* button present on left hand side of the dashboard. Click on import and load the system_monitoring_metrics.json present in the grafana_template folder.

- While creating this dashboard if you face any issues at valueset, change it to empty and then click on import by selecting the datasources.

- *For more info about grafana dashboard imports you can refer [here](https://grafana.com/docs/grafana/latest/reference/export_import/).*

## Alerting (Telegram and Email)
 A custom alerting module has been developed to alert on key validator health events. The module uses data from influxdb and trigger alerts based on user-configured thresholds.

- Alert when the missed blocks count reaches or exceeds **missed_blocks_threshold** which is user configured in *config.toml*.
- Alert when no.of peers count falls below of **num_peers_threshold** which is user configured in *config.toml*
- Alert when the oasis node status is not running on validator instance.
- Alert when the block difference between validator and network reaches or exceeds the **block_diff_threshold** which is user configured in config.toml.
- Alert about validator health, i.e. whether it's voting or jailed. You can get alerts twice a day based on the time you have configured **alert_time1** and **alert_time2** in *config.toml*. This is a useful sanity check, to confirm the validator is voting (or alerting you if it's jailed).
- Alert when the voting power of your validator drops below **voting_power_threshold** which is user configured in *config.toml*
- Alert when the worker epoch difference between validator and network reaches or exceeds the **epoch_diff_threshold** which is user configured in config.toml.

## Alerting - emergency notification

There is also an "emergency" alert that can be configured. When triggered, it sends an email to your pagerduty account to trigger a pagerduty workflow. This alert is configured as described above and repeated here for convenience.

- *pagerduty_email*

    Give mail address of pager duty service to send alerts of emergency missed blocks.
    Note : Have to give mail address which was generated after creation of a service in pager duty.

    Learn more about pagerduty [here](https://www.pagerduty.com/).

- *emergency_missed_blocks_threshold*

    Give threshold to notify about continuous missed blocks to your pager duty account. so that it will send mails, messages and makes you a call about alerts.

## Steps to create telegram bot
To create telegram bot and to configure tg_bot_token and tg_chat_id, one can follow the below steps.

- In the first step search for `@BotFather` , go to that chat and click on start.
- Then enter `/newbot` and the bot will respond, follow those instructions to create a bot.
- From `example 3` you can observe tg_bot_token which is marked, and also to go to the bot chat click on the `Testing_intro_bot`.

| example 1     | example 2      | example 3      |
|------------|-------------|-------------|
| <img src="https://github.com/Chainflow/oasis-mission-control/blob/implementation/images/start.jpg" width="230"> | <img src="https://github.com/Chainflow/oasis-mission-control/blob/implementation/images/new_bot.jpg" width="250"> | <img src="https://github.com/Chainflow/oasis-mission-control/blob/implementation/images/bot_created.jpg" width="250"> |

- To know your Telegram chat id, search for `@my_id_bot` then you can find bot with name `What's my ID`.
- Use this bot to get your personal ID or add it to any group to see its ID. Then that will become your tg_chat_id .

## Feedback and Questions

Please create issues [in this repo](https://github.com/Chainflow/oasis-mission-control) to ask questions and/or suggest additional metrics, alerts or features.
