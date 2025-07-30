#!/bin/bash

function installsftpservice() {
    cd /tmp
    wget http://uuxia.cn:8087/temp/ssl/ca-certificates_20211016-1_all.ipk
    wget http://uuxia.cn:8087/temp/ssl/ca-bundle_20211016-1_all.ipk
    wget http://uuxia.cn:8087/temp/ssl/libustream-openssl20201210_2022-01-16-868fd881-2_aarch64_cortex-a53.ipk
    opkg install *.ipk
}

function installsftp() {
    opkg update
    opkg install vsftpd openssh-sftp-server
    /etc/init.d/vsftpd enable
    /etc/init.d/vsftpd start
    pwd
}

# shellcheck disable=SC2120
function main() {
    echo "1. 解决wget ssl问题"
    echo "2. 安装sftp"
    echo "请输入编号:"
    read index
    clear
    case "$index" in
    [1]) (installsftpservice);;
    [2]) (installsftp);;
    *) echo "exit" ;;
  esac
}
main