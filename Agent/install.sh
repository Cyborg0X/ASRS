#!/bin/bash
sleep 1s
echo -e  "\033[1;32mInstalling ASRS dependencies......"
mkdir ~/golang  2>/dev/null 
cd ~/golang || exit

distro=$(cat /etc/*-release)
 
if [[ "$distro" =~ "ID_LIKE="debian"" ]]; then
    echo -e "\033[1;33mThe system is debian-based [OK]\033[0m"
    
    if ping -c 4 google.com > /dev/null 2>&1 ; then 
        echo -e  "\033[1;32mInternet conectivity [OK]\033[0m"
    else
     
        echo -e  "\033[1;31mERROR: connection lost\033[0m"
        exit
    fi
    echo -e  "\033[1;32mDownloading Golang please wait .....\033[0m"
    sudo apt install wget -y > /dev/null 2>&1
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
    echo -e  "\033[1;32mAll Dependencies installed [OK]\033[0m"
    go1version=$(go version)
    printf "\033[1;32m%s installed\033[0m\n" "$go1version"
    echo -e  "\033[1;32mFinshied...\033[0m"
    exit


elif [[ "$distro" =~ 'ID_LIKE="rhel fedora"' ]]; then
    #dffgdgdg
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
    echo -e  "\033[1;32mAll Dependencies installed [OK]\033[0m"
    go1version=$(go version)
    printf "\033[1;32m%s installed\033[0m\n" "$go1version"
    echo -e  "\033[1;32mFinshied...\033[0m"
    exit

else
    echo -e  "\033[1;31mERROR: Operating system is not supported\033[0m"
fi