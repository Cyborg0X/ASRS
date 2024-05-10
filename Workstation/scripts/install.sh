#!/bin/bash
sleep 1s
if [ -z "$1" ];then
    printf "\033[1;32mUsage: ./install.sh [options] command\nplease choose your OS:\n -debian\n -ubuntu\n -fedora\n\033[0m"
    exit
fi

cd ~/golang || exit
if [[ $1 == "-debian" || $1 == "-ubuntu" ]]; then
    # INSTALL FOR DEBIAN
    echo -e  "\033[1;32mInstalling ASRS Workstation dependencies......"
    sudo mkdir -p /etc/ASRS_WS/.config  2>/dev/null
    cd /etc/ASRS_WS/.config && sudo mkdir SSH_config rsync_config snapshot_config connection_config env_var 2>/dev/null
    mkdir /etc/ASRS_WS/.database 2>/dev/null
    cd /etc/ASRS_WS/.database && sudo mkdir status_logs snapshot_logs detection_marker rsync_logs 2>/dev/null
    logs=("status_logs" "snapshot_logs" "detection_marker" "rsync_logs")
    confingu=("SSH_config" "rsync_config" "snapshot_config" "connection_config" "env_var")
    for log in "${logs[@]}"; do
        cd /etc/ASRS_WS/.database/"$log" && touch "$log".json
    done
    for confe in "${confingu[@]}"; do
        cd /etc/ASRS_WS/.config/"$confe" && touch "$confe".json
    done

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
    packages=("rsync" "snapper" "ssh" "openssh-server" "openssh-client")
    
    for pkg in "${packages[@]}"; do
        echo -e  "\033[1;32minstalling $pkg ......\033[0m"
        sudo apt install "$pkg" -y > /dev/null 2>&1
        inst_pid=$!
        wait $inst_pid
        echo -e  "\033[1;32m$pkg installed [OK]\033[0m"
    done
    
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
    cd /etc/ASRS_WS/.config && sudo mkdir SSH_config rsync_config snapshot_config connection_config env_var 2>/dev/null
    mkdir /etc/ASRS_WS/.database 2>/dev/null
    cd /etc/ASRS_WS/.database && sudo mkdir status_logs snapshot_logs detection_marker rsync_logs 2>/dev/null
    logs=("status_logs" "snapshot_logs" "detection_marker" "rsync_logs")
    confingu=("SSH_config" "rsync_config" "snapshot_config" "connection_config" "env_var")
    for log in "${logs[@]}"; do
        cd /etc/ASRS_WS/.database/"$log" && touch "$log".json
    done
    for confe in "${confingu[@]}"; do
        cd /etc/ASRS_WS/.config/"$confe" && touch "$confe".json
    done

    cd ~/golang || exit
    echo -e  "\033[1;32mDownloading Golang please wait .....\033[0m"
    sudo sudo dnf update && sudo dnf install wget tar -y > /dev/null 2>&1
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
    packages=("rsync" "snapper" "ssh" "snort" "openssh-server" "openssh-client")
    
    for pkg in "${packages[@]}"; do
        echo -e  "\033[1;32minstalling $pkg ......\033[0m"
        sudo apt install "$pkg" -y > /dev/null 2>&1
        inst_pid=$!
        wait $inst_pid
        echo -e  "\033[1;32m$pkg installed [OK]\033[0m"
    done

    echo -e  "\033[1;32mAll Dependencies installed [OK]\033[0m"
    go1version=$(go version)
    printf "\033[1;32m%s installed\033[0m\n" "$go1version"
    echo -e  "\033[1;32mFinshied...\033[0m"
    exit

else
    echo -e  "\033[1;31mERROR: Operating system is not supported\033[0m"
fi