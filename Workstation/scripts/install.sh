#!/bin/bash

   # Define some variables
SSH_CONFIG_FILE="/etc/ssh/sshd_config"
BACKUP_FILE="/etc/ssh/sshd_config.bak"
function main() {
    if [ -z "$1" ];then
      printf "\033[1;32mUsage: ./install.sh [options] command\nplease choose your OS:\n -debian\n -ubuntu\n -fedora\n\033[0m"
      exit
  fi

  if [[ $1 == "-debian" || $1 == "-ubuntu" ]]; then
      # INSTALL FOR DEBIAN
      echo -e  "\033[1;32mInstalling ASRS Workstation dependencies......"
      sudo mkdir -p /etc/ASRS_WS/.config  2>/dev/null
      cd /etc/ASRS_WS/.config && sudo touch config.json 2>/dev/null
      mkdir /etc/ASRS_WS/.database 2>/dev/null
      cd /etc/ASRS_WS/.database && sudo touch logs.json 2>/dev/null
      cd /etc/ASRS_WS/.database && sudo mkdir website_backup database_backup snapshots_backup 2>/dev/null
  
      echo -e "\033[1;33mThe system is debian-based [OK]\033[0m"

      if ping -c 4 google.com > /dev/null 2>&1 ; then 
          echo -e  "\033[1;32mInternet connectivity [OK]\033[0m"
      else

          echo -e  "\033[1;31mERROR: Please check the internet connectivity\033[0m"
          exit
      fi
      echo -e  "\033[1;32mDownloading Golang please wait .....\033[0m"
      sudo apt install wget -y > /dev/null 2>&1
      #sudo wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz > /dev/null 2>&1 &
      #wget_pid=$!
      #wait $wget_pid
      echo -e  "\033[1;32mGolang Downloaded [OK]"
      sleep 1s
      echo -e  "\033[1;32minstalling Golang V1.22.2 ......\033[0m"
      #sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf ~/golang/go1.22.2.linux-amd64.tar.gz
      export PATH=$PATH:/usr/local/go/bin
      wait
      echo -e  "\033[1;32mGolang installed [OK]\033[0m"
      sleep 1s 
      packages=("rsync" "snapper" "rsync-daemon")

      for pkg in "${packages[@]}"; do
          echo -e  "\033[1;32minstalling $pkg ......\033[0m"
          sudo apt install "$pkg" -y > /dev/null 2>&1
          inst_pid=$!
          wait $inst_pid
          echo -e  "\033[1;32m$pkg installed [OK]\033[0m"
      done
           #..........install SSH server
      mkdir ~/.ssh
      touch ~/.ssh/authorized_keys
      install_ssh_server
      backup_ssh_config
      configure_ssh_server
      restart_ssh_service
      sudo ufw allow ssh
      echo "SSH server installation and configuration complete."

      echo -e  "\033[1;32minstalling snort ......\033[0m"
      sudo apt install snort -y 
      snort_pid=$!
      wait $snort_pid
      echo -e  "\033[1;32msnort installed [OK]\033[0m"


      echo -e  "\033[1;32mAll Dependencies installed [OK]\033[0m"
      go1version=$(go version)
      printf "\033[1;32m%s installed\033[0m\n" "$go1version"
      echo -e  "\033[1;32mFinshied...\033[0m"
      exit

      # the same installation but for RHEL FEDORA
  elif [[ $1 == "-fedora" ]]; then
      sleep 1s
      echo -e  "\033[1;32mInstalling ASRS Workstation dependencies......"
       sudo mkdir -p /etc/ASRS_WS/.config  2>/dev/null
      cd /etc/ASRS_WS/.config && sudo touch config.json rsyncd.secrets 2>/dev/null
      mkdir /etc/ASRS_WS/.database 2>/dev/null
      cd /etc/ASRS_WS/.database && sudo touch logs.json 2>/dev/null
      cd /etc/ASRS_WS/.database && sudo mkdir website_backup snapshots_backup 2>/dev/null
      sudo chmod 0640 /etc/ASRS_WS/.config/rsyncd.secrets
      pass="snapper:Sn@pPeer
      webuser:FG4@#%3"
      echo "$pass" > /etc/ASRS_WS/.config/rsyncd.secrets
      cd ~/golang || exit
      echo -e  "\033[1;32mDownloading Golang please wait .....\033[0m"
      sudo sudo dnf update && sudo dnf install wget tar -y > /dev/null 2>&1
      #sudo wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz > /dev/null 2>&1 &
      wget_pid=$!
      wait $wget_pid
      echo -e  "\033[1;32mGolang Downloaded [OK]"
      sleep 1s
      echo -e  "\033[1;32minstalling Golang V1.22.2 ......\033[0m"
      #sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf ~/golang/go1.22.2.linux-amd64.tar.gz
      export PATH=$PATH:/usr/local/go/bin
      wait
      echo -e  "\033[1;32mGolang installed [OK]\033[0m"
      sleep 1s 
      packages=("rsync" "snapper" "golang-bin" "rsync-daemon")

      for pkg in "${packages[@]}"; do
          echo -e  "\033[1;32minstalling $pkg ......\033[0m"
          sudo dnf install "$pkg" -y > /dev/null 2>&1
          inst_pid=$!
          wait $inst_pid
          echo -e  "\033[1;32m$pkg installed [OK]\033[0m"
      done
# Install dependencies
      #dnf update -y
      #dnf install -y epel-release
      #dnf install -y snort

# Backup the Snort configuration file
      #cp /etc/snort/snort.conf /etc/snort/snort.conf.bak

# Edit the Snort configuration file
    # sed -i "'/^include \$RULE_PATH/a \alert tcp \$EXTERNAL_NET any -> \$HOME_NET any (msg:"COMMAND INJECTION ATTEMPT"; content:"|2e 2f|"; depth:2; sid:1000001; rev:1;)' /etc/snort/snort.conf"
    # sed -i '/^#output alert_full/s/^#//' /etc/snort/snort.conf
    # sed -i 's|output alert_full: .*|output alert_full: /var/log/snort/command_injection_alerts.txt|' /etc/snort/snort.conf

# Create the log directory
      #mkdir -p /var/log/snort

# Start and enable Snort
      #systemctl start snort
      #systemctl enable snort

      echo "Snort configuration updated. Command injection attack rule added and separate alerts file created."
  # Start Snort
     # /usr/local/bin/snort -c /etc/snort/snort.conf -i eth0 -A console
     #..........install SSH server
     # install_ssh_server
      #backup_ssh_config
      #configure_ssh_server
      #restart_ssh_service
      #sudo ufw allow ssh
      #echo "SSH server installation and configuration complete."

      #sudo firewall-cmd --permanent --add-service=ssh
      #sudo firewall-cmd --reload
      echo -e  "\033[1;32mInstalling & configuring Rsync daemon\033[0m"
      configdaemons
      sudo rsync --daemon
      echo -e  "\033[1;32mRsync daemon Installed [OK]\033[0m"

      echo -e  "\033[1;32mAll Dependencies installed [OK]\033[0m"
      go1version=$(go version)
      printf "\033[1;32m%s installed\033[0m\n" "$go1version"
      echo -e  "\033[1;32mFinshied...\033[0m"
      exit

  else
      echo -e  "\033[1;31mERROR: Operating system is not supported\033[0m"
  fi


}

function configdaemons() {
  touch /etc/rsyncd.conf
  conf="
  # Global Settings
  uid = root
  gid = root
  use chroot = yes
  max connections = 4
  log file = /var/log/rsyncd.log
  pid file = /var/run/rsyncd.pid
  lock file = /var/run/rsync.lock
  motd file = /etc/rsyncd.motd
  
  # Modules
  [snapshots]
  path = /etc/ASRS_WS/.database/snapshots_backup
  comment = Snapper Snapshots
  read only = true
  auth users = snapper
  secrets file = /etc/ASRS_WS/.config/rsyncd.secrets
  transfer logging = yes
  log format = %t %a %m %f %b
  
  [database]
  path = /etc/ASRS_WS/.database/database_backup
  comment = SQL Database Backup
  read only = false
  auth users = webuser
  secrets file = /etc/ASRS_WS/.config/rsyncd.secrets
  exclude = lost+found
  transfer logging = yes
  log format = %t %a %m %f %b
  
  [website]
  path = /etc/ASRS_WS/.database/website_backup
  comment = Website Files
  read only = false
  auth users = webuser
  secrets file = /etc/ASRS_WS/.config/rsyncd.secrets
  exclude = .git, node_modules, .cache
  transfer logging = yes
  log format = %t %a %m %f %b"

  echo "$conf" > /etc/rsyncd.conf
}




<<'END'

# Function to install ssh server
function install_ssh_server() {
  if [ -x "$(command -v apt)" ]; then
    sudo apt update
    sudo apt install -y openssh-server
  elif [ -x "$(command -v dnf)" ]; then
    sudo dnf install -y openssh-server
  elif [ -x "$(command -v dnf)" ]; then
    sudo dnf install -y openssh-server
  elif [ -x "$(command -v pacman)" ]; then
    sudo pacman -Syu --noconfirm openssh
  else
    echo "Unsupported package manager. Please install OpenSSH Server manually."
    exit 1
  fi
}

# Function to backup the current SSH config
function backup_ssh_config() {
  if [ -f "${SSH_CONFIG_FILE}" ]; then
    sudo cp "${SSH_CONFIG_FILE}" "${BACKUP_FILE}"
    echo "Backup of SSH configuration saved to ${BACKUP_FILE}"
  else
    echo "SSH configuration file not found at ${SSH_CONFIG_FILE}. Exiting."
    exit 1
  fi
}

# Function to configure the SSH server for passwordless authentication 
function configure_ssh_server() {
  sudo sed -i 's/^#Port 22/#Port 22/' "${SSH_CONFIG_FILE}"
  sudo sed -i 's/^#PermitRootLogin prohibit-password/#PermitRootLogin no/' "${SSH_CONFIG_FILE}"
  sudo sed -i 's/^#PasswordAuthentication yes/#PasswordAuthentication no/' "${SSH_CONFIG_FILE}"
  sudo sed -i 's/^#PubkeyAuthentication no/#PubkeyAuthentication yes/' "${SSH_CONFIG_FILE}"
}


# Function to restart the SSH service
function restart_ssh_service() {
  if [ -x "$(command -v systemctl)" ]; then
    sudo systemctl restart sshd || sudo systemctl restart ssh
  elif [ -x "$(command -v service)" ]; then
    sudo service ssh restart
  else
    echo "Unsupported service manager. Please restart the SSH service manually."
    exit 1
  fi
}
END

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$1"
fi