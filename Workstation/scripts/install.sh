#!/bin/bash

   # Define some variables
#SSH_CONFIG_FILE="/etc/ssh/sshd_config"
#BACKUP_FILE="/etc/ssh/sshd_config.bak"
function main() {
    if [ -z "$1" ];then
      printf "\033[1;32mUsage: ./install.sh [options] command\nplease choose your OS:\n -debian\n -ubuntu\n -fedora\n\033[0m"
      exit
  fi

  if  [[ $1 == "-fedora" ]]; then
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
  elif [[ $1 == "-debian" || $1 == "-ubuntu" ]]; then
      sleep 1s
      echo -e  "\033[1;32mInstalling ASRS Workstation dependencies......"
      sudo mkdir -p /etc/ASRS_WS/.config  2>/dev/null
      cd /etc/ASRS_WS/.config && sudo touch config.json pass.txt 2>/dev/null
      pas="12345"
      echo "$pas" > pass.txt
      sudo chmod 600 pass.txt
      sudo chown root:root pass.txt
      mkdir /etc/ASRS_WS/.database 2>/dev/null
      cd /etc/ASRS_WS/.database && sudo touch logs.json && sudo mkdir website_backup snapshots_backup 2>/dev/null

      cd ~/golang || exit
      echo -e  "\033[1;32mDownloading Golang please wait .....\033[0m"
      sudo sudo apt update && sudo apt install wget tar -y > /dev/null 2>&1
      sudo wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz > /dev/null 2>&1 &
      wget_pid=$!
      wait $wget_pid
      echo -e  "\033[1;32mGolang Downloaded [OK]"
      sleep 1s
      echo -e  "\033[1;32minstalling Golang V1.22.2 ......\033[0m"
      sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf ~/golang/go1.22.2.linux-amd64.tar.gz
      export PATH=$PATH:/usr/local/go/bin
      wait
      echo -e  "\033[1;32mGolang installed [OK]\033[0m"
      sleep 1s 
      packages=("rsync" "snapper" "golang-bin")

      for pkg in "${packages[@]}"; do
          echo -e  "\033[1;32minstalling $pkg ......\033[0m"
          sudo apt install "$pkg" -y > /dev/null 2>&1
          inst_pid=$!
          wait $inst_pid
          echo -e  "\033[1;32m$pkg installed [OK]\033[0m"
      done

      echo "Snort configuration updated. Command injection attack rule added and separate alerts file created."

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
  
  sudo apt install xinetd && export RSYNC_ENABLE=inetd
  sudo touch /etc/rsyncd.conf
  sudo touch /etc/rsyncd.secrets
  users="asrs:12345" 
  echo "$users" > /etc/rsyncd.secrets
  sudo chmod 600 /etc/rsyncd.secrets
  xconf="
# Simple configuration file for xinetd
#
# Some defaults, and include /etc/xinetd.d/

service rsync
{

  disable = no
  socket_type = stream
  wait = no
  user = root
  server = /usr/bin/rsync
  server_args = --daemon
  log_on_failure += USERID
  flags = IPv6

# Please note that you need a log_type line to be able to use log_on_success
# and log_on_failure. The default is the following :
# log_type = SYSLOG daemon info

}

includedir /etc/xinetd.d"
  Rconf="
# Global Settings
max connections = 4
log file = /var/log/rsync.log
timeout = 300
use chroot = false
[backup]
comment = ASRS BACKUP
path = /home/backup
read only = no
list = yes
uid = ws
gid = ws
auth users = asrs
secrets file = /etc/rsyncd.secrets
transfer logging = yes
log format = %t %a %m %f %b"

  echo "$Rconf" > /etc/rsyncd.conf
  echo "$xconf" > /etc/xinetd.conf
  sudo /etc/init.d/xinetd restart
  systemctl restart rsync.service
}




# shellcheck disable=SC2188
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