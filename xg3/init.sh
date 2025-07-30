#!/bin/sh

# shellcheck disable=SC2112
function installsftpservice() {
    cd /tmp
    wget http://uuxia.cn:8087/temp/ssl/ca-certificates_20211016-1_all.ipk
    wget http://uuxia.cn:8087/temp/ssl/ca-bundle_20211016-1_all.ipk
    wget http://uuxia.cn:8087/temp/ssl/libustream-openssl20201210_2022-01-16-868fd881-2_aarch64_cortex-a53.ipk
    opkg install *.ipk
}

#sed -e 's,https://downloads.openwrt.org,https://mirror.bjtu.edu.cn/openwrt,g' \
#-e 's,https://downloads.openwrt.org,https://mirror.bjtu.edu.cn/openwrt,g' \
#-i.bak /etc/opkg/customfeeds.conf
function installsftp() {
    opkg update
    opkg install vsftpd openssh-sftp-server
    /etc/init.d/vsftpd enable
    /etc/init.d/vsftpd start
    pwd
}

function openssh() {
    curl http://192.168.10.1/cgi-bin/luci/pti/ssh_open
}

function opkgUpdate() {
    # 清除 opkg 缓存
    rm -rf /var/opkg-lists/*
    # 重新更新
    opkg update
}

#cp /etc/opkg/distfeeds.conf /etc/opkg/distfeeds.conf.bak
function changeOpkg() {
    cp /etc/opkg/distfeeds.conf /etc/opkg/distfeeds.conf.bak

#    https://downloads.openwrt.org/releases/21.02-SNAPSHOT/targets/mediatek/mt7981/packages
#    https://mirror.nju.edu.cn/immortalwrt/releases/21.02-SNAPSHOT/packages/aarch64_cortex-a53/luci/
#    https://mirrors.tuna.tsinghua.edu.cn/openwrt/releases/21.02-SNAPSHOT/packages/aarch64_cortex-a53/base
    sed -i 's_downloads.openwrt.org_mirror.nju.edu.cn/immortalwrt_' /etc/opkg/distfeeds.conf
}
# shellcheck disable=SC2120
function main() {
    echo "1. 解决wget ssl问题"
    echo "2. 安装sftp"
    echo "3. 开启SSH"
    echo "4. 更新源"
    echo "请输入编号:"
    read index
    clear
    case "$index" in
    [1]) (installsftpservice);;
    [2]) (installsftp);;
    [3]) (openssh);;
    [4]) (opkgUpdate);;
    *) echo "exit" ;;
  esac
}
main