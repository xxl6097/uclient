#!/bin/bash
module=$(grep "module" go.mod | cut -d ' ' -f 2)
options=("windows:amd64" "windows:arm64" "linux:amd64" "linux:arm64" "linux:arm:7" "linux:arm:5" "linux:mips64" "linux:mips64le" "linux:mips:softfloat" "linux:mipsle:softfloat" "linux:riscv64" "linux:loong64" "darwin:amd64" "darwin:arm64" "freebsd:amd64" "android:arm64")
#options=("linux:amd64")
#options=("linux:amd64" "windows:amd64")
version=0.0.0
versionDir="$module/pkg"
appname="ubus"


# shellcheck disable=SC2120
function buildgo() {
  builddir=$1
  appname=$2
  version=$3
  appdir=$4
  os=$5
  arch=$6
  extra=$7
  dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}
  flags='';
  if [ "${os}" = "linux" ] && [ "${arch}" = "arm" ] && [ "${extra}" != "" ] ; then
    if [ "${extra}" = "7" ]; then
      flags=GOARM=7;
      dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}hf
    elif [ "${extra}" = "5" ]; then
      flags=GOARM=5;
      dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}
    fi;
  elif [ "${os}" = "windows" ] ; then
    dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}.exe
    if [ "${arch}" = "amd64" ]; then
        go generate ${appdir}
    fi
  elif [ "${os}" = "linux" ] && ([ "${arch}" = "mips" ] || [ "${arch}" = "mipsle" ]) && [ "${extra}" != "" ] ; then
    flags=GOMIPS=${extra};
  fi;
  #echo "build：GOOS=${os} GOARCH=${arch} ${flags} ==> ${dstFilePath}"
  printf "build：GOOS=%-7s GOARCH=%-8s ==> %s\n" ${os} ${arch} ${dstFilePath}

  filename=$(basename "$dstFilePath")
  binName="-X '${versionDir}.BinName=${filename}'"
#  echo "--->env CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} ${flags} go build -trimpath -ldflags "$ldflags $binName -linkmode internal" -o ${dstFilePath} ${appdir}"
  env CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} ${flags} go build -trimpath -ldflags "$ldflags $binName -linkmode internal" -o ${dstFilePath} ${appdir}
  if [ "${os}" = "windows" ] ; then
    if [ "${arch}" = "amd64" ]; then
        rm -rf ${appdir}/resource.syso
    fi
  fi;
}

# builddir：输出目录
# appname：应用名称
# version：应用版本
# appdir：main.go目录
# disname：显示名
# describe：描述
function buildMenu() {
  builddir=$1
  appname=$2
  version=$3
  appdir=$4
  disname=$5
  describe=$6
  ldflags=$(buildLdflags $appname $disname $describe)
  PS3="请选择需要编译的平台："
  select arch in "${options[@]}"; do
      if [[ -n "$arch" ]]; then
        IFS=":" read -r os arch extra <<< "$arch"
        buildgo $builddir $appname $version $appdir $os $arch $extra
        return $?
      else
        echo "输入无效，请重新选择。"
      fi
  done
}

# builddir：输出目录
# appname：应用名称
# version：应用版本
# appdir：main.go目录
# disname：显示名
# describe：描述
function buildAll() {
  builddir=$1
  appname=$2
  version=$3
  appdir=$4
  disname=$5
  describe=$6
  ldflags=$(buildLdflags $appname $disname $describe)
  for arch in "${options[@]}"; do
      IFS=":" read -r os arch extra <<< "$arch"
      buildgo $builddir $appname $version $appdir $os $arch $extra
  done
  #wait
}

function buildLdflags() {
  local ldflags
  ldflags="-s -w"
  echo "${ldflags[*]-}"
}


function showBuildDir() {
  # 检查是否输入路径参数
  if [ -z "$1" ]; then
      echo "用法: $0 <路径>"
      exit 1
  fi

  # 验证路径是否存在且为目录
  if [ ! -d "$1" ]; then
      echo "错误: 路径 '$1' 不存在或不是目录！"
      exit 1
  fi

  # 获取指定路径下的所有直接子目录（非递归）
  dirs=()
  while IFS= read -r dir; do
      dirs+=("$dir")
  done < <(find "$1" -maxdepth 1 -type d ! -path "$1" | sort)

  # 检查是否有子目录
  if [ ${#dirs[@]} -eq 0 ]; then
      echo "路径 '$1' 下没有子目录！"
      exit 0
  fi

  # 生成交互式菜单
  echo "请选择要操作的目录："
  PS3="输入序号 (1-${#dirs[@]}): "
  select dir in "${dirs[@]}"; do
      if [[ -n "$dir" ]] && [[ $REPLY -ge 1 && $REPLY -le ${#dirs[@]} ]]; then
          echo "您选择的目录是: $dir"
          break
#          return $dir
      else
          echo "无效输入！请输入有效序号。"
      fi
  done
}

# shellcheck disable=SC2120
function buildInstaller() {
  showBuildDir ./cmd
  builddir="./release"
  #appname=$(basename "$dir")
#  appname="srvinstaller"
  appdir=${dir}
  disname="${appname}应用程序"
  describe="一款基于GO语言的服务安装程序"
  rm -rf ${builddir}
  buildMenu $builddir $appname "$version" $appdir $disname $describe
#  buildAll $builddir $appname "$version" $appdir $disname $describe
}

buildInstaller