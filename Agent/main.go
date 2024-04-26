package main

/*

echo "Checking packages installed....."
packages=("rsync" "snapper" "ssh" "snort" "openssh-server" "openssh-client")






Echo "Installing packages and dependecies started......."
if [[ "$linuxdistro" =~ "ID_LIKE=debian" ]]; then
    install="sudo apt install"
    checkpkg="dpkg -s"

elif [[ "$linuxdistro" =~ "ID_LIKE=rhel fedora" ]]; then
    install="sudo yum install"
    checkpkg="yum list installed"
fi



echo "Installing ASRS agent packages and dependencies.........."

linuxdistro=$(cat /etc/*-release)


echo "installing Rsync for syncing files......."
sleep 3s
dpkg -s rsync
eval "$install rsync"

*/
