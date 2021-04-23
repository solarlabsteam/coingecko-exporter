# coingecko-exporter

![Latest release](https://img.shields.io/github/v/release/solarlabsteam/coingecko-exporter)
[![Actions Status](https://github.com/solarlabsteam/coingecko-exporter/workflows/test/badge.svg)](https://github.com/solarlabsteam/coingecko-exporter/actions)

coingecko-exporter is a Prometheus scraper that fetches the exchange rates from Coingecko.

## How can I set it up?

First of all, you need to download the latest release from [the releases page](https://github.com/solarlabsteam/coingecko-exporter/releases/). After that, you should unzip it and you are ready to go:

```sh
wget <the link from the releases page>
tar xvfz coingecko-exporter-*.*-amd64.tar.gz
cd coingecko-exporter-*.*-amd64.tar.gz
./coingecko-exporter
```

That's not really interesting, what you probably want to do is to have it running in the background. For that, first of all, we have to copy the file to the system apps folder:

```sh
sudo cp ./coingecko-exporter /usr/bin
```

Then we need to create a systemd service for our app:

```sh
sudo nano /etc/systemd/system/coingecko-exporter.service
```

You can use this template (change the user to whatever user you want this to be executed from. It's advised to create a separate user for that instead of running it from root):

```
[Unit]
Description=Coingecko Exporter
After=network-online.target

[Service]
User=<username>
TimeoutStartSec=0
CPUWeight=95
IOWeight=95
ExecStart=coingecko-exporter
Restart=always
RestartSec=2
LimitNOFILE=800000
KillSignal=SIGTERM

[Install]
WantedBy=multi-user.target
```

Then we'll add this service to the autostart and run it:

```sh
sudo systemctl enable coingecko-exporter
sudo systemctl start coingecko-exporter
sudo systemctl status coingecko-exporter # validate it's running
```

If you need to, you can also see the logs of the process:

```sh
sudo journalctl -u coingecko-exporter -f
```

## How can I scrape data from it?

Here's the example of the Prometheus config you can use for scraping data:

```yaml
scrape-configs:
  - job_name:       'coingecko'
    scrape_interval: 15s
    metrics_path: /metrics/rates/usd # replace USD with other base currency if you like
    static_configs:
      - targets:
        - <list of tokens you want to monitor>
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_currency
      - source_labels: [__param_address]
        target_label: instance
      - target_label: __address__
        replacement: <node IP or hostname>:9400
```

Then restart Prometheus and you're good to go!

## How can I configure it?

You can pass the artuments to the executable file to configure it. Here is the parameters list:

- `--listen-address` - the address with port the node would listen to. For example, you can use it to redefine port or to make the exporter accessible from the outside by listening on `127.0.0.1`. Defaults to `:9400` (so it's accessible from the outside on port 9400)

## How can I contribute?

Bug reports and feature requests are always welcome! If you want to contribute, feel free to open issues or PRs.
