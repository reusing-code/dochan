# Deployment

Using Ubuntu 16.04
Server + Client, autostart using Systemd

dochan will be installed at `/mnt/nas/Development/dochan`
feel free to use any other directory

## Build

### Server

- Install go

  ```bash
  wget https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go1.11.2.linux-amd64.tar.gz
  echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
  ```

- Checkout source

  ```bash
  cd ~
  mkdir -p build && cd build
  git clone git@github.com:reusing-code/dochan.git
  ```

- Build

  TODO: requires GOPATH or modules

  ```bash
  cd dochan/api
  go build
  ```

- Copy
  ```bash
  mkdir -p /mnt/nas/Development/dochan
  cp api /mnt/nas/Development/dochan/server
  ```
- Cleanup (optional)
  ```bash
  cd ~/build
  rm -r dochan
  ```

### Client

- Install nodejs + npm

  ```bash
  curl -sL https://deb.nodesource.com/setup_11.x | sudo -E bash -
  sudo apt-get install nodejs
  ```

- Checkout source
  ```bash
  cd ~
  mkdir -p build && cd build
  git clone git@github.com:reusing-code/dochan-client.git
  ```
- Build
  ```bash
  cd dochan-client
  npm install
  npm run build
  ```
- Copy
  ```bash
  mkdir -p /mnt/nas/Development/dochan
  cp -r dist /mnt/nas/Development/dochan/
  ```
- Cleanup (optional)
  ```bash
  cd ~/build
  rm -r dochan-client
  ```

## Install

- Create user
  ```bash
  sudo useradd -r dochan
  ```
- Change file rights

  ```bash
  sudo chown -cR dochan:dochan /mnt/nas/Developtment/dochan
  ```

- Create systemd unit file

  ```bash
  vi /lib/systemd/system/dochan.service
  ```

  ```
  [Unit]
  Description=Dochan Service
  ConditionPathExists=/mnt/nas/Development/dochan
  After=network.target

  [Service]
  Type=simple
  User=dochan
  Group=dochan
  LimitNOFILE=8096

  Restart=on-failure
  RestartSec=10

  WorkingDirectory=/mnt/nas/Development/dochan
  ExecStart=/mnt/nas/Development/dochan/server --path="output" --assetPath="dist/" --dbFile="dochan.db"

  StandardOutput=syslog
  StandardError=syslog
  SyslogIdentifier=dochan

  [Install]
  WantedBy=multi-user.target
  ```

- Adapt access rights

  ```bash
  sudo chmod 755 /lib/systemd/system/dochan.service
  ```

- Enable autostart and manually start

  ```bash
  sudo systemctl enable dochan
  sudo systemctl start dochan
  ```

- Enable some group (e.g. jenkins) to start/stop/restart the service
  ```bash
  sudo visudo
  ```
  ```
  Cmnd_Alias DOCHAN_CMNDS = /bin/systemctl start dochan, /bin/systemctl stop dochan, /bin/systemctl restart dochan
  %jenkins ALL=(ALL) NOPASSWD: DOCHAN_CMNDS
  ```
  ```bash
  sudo systemctl daemon-reload
  ```
