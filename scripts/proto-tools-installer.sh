#!/bin/bash

set -ue

DESTDIR=${DESTDIR:-}
PREFIX=${PREFIX:-/usr/local}
UNAME_S="$(uname -s 2>/dev/null)"
UNAME_M="$(uname -m 2>/dev/null)"
PROTOC_VERSION=3.13.0

f_abort() {
  local l_rc=$1
  shift

  echo $@ >&2
  exit ${l_rc}
}

case "${UNAME_S}" in
Linux)
  PROTOC_ZIP="protoc-${PROTOC_VERSION}-linux-x86_64.zip"
  ;;
Darwin)
  PROTOC_ZIP="protoc-${PROTOC_VERSION}-osx-x86_64.zip"
  ;;
*)
  f_abort 1 "Unknown kernel name. Exiting."
esac

TEMPDIR="$(mktemp -d)"

trap "rm -rvf ${TEMPDIR}" EXIT

f_print_installing_with_padding() {
  printf "Installing %30s ..." "$1" >&2
}

f_print_done() {
  echo -e "\tDONE" >&2
}

f_ensure_tools() {
  ! which curl &>/dev/null && f_abort 2 "couldn't find curl, aborting" || true
}

f_ensure_dirs() {
  mkdir -p "${DESTDIR}/${PREFIX}/bin"
  mkdir -p "${DESTDIR}/${PREFIX}/include"
}

f_needs_install() {
  if [ -x $1 ]; then
    echo -e "\talready installed. Skipping." >&2
    return 1
  fi

  return 0
}

f_install_protoc() {
  f_print_installing_with_padding proto_c
  f_needs_install "${DESTDIR}/${PREFIX}/bin/protoc" || return 0

  pushd "${TEMPDIR}" >/dev/null
  curl -o "${PROTOC_ZIP}" -sSL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}"
  unzip -q -o ${PROTOC_ZIP} -d ${DESTDIR}/${PREFIX} bin/protoc; \
  unzip -q -o ${PROTOC_ZIP} -d ${DESTDIR}/${PREFIX} 'include/*'; \
  rm -f ${PROTOC_ZIP}
  popd >/dev/null
  f_print_done
}

f_install_protoc_gen_gocosmos() {
  f_print_installing_with_padding protoc-gen-gocosmos

  if ! grep "github.com/gogo/protobuf => github.com/regen-network/protobuf" go.mod &>/dev/null ; then
    echo -e "\tPlease run this command from somewhere inside the cosmos-sdk folder."
    return 1
  fi

  go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos 2>/dev/null
  f_print_done
}

f_install_protoc_gen_grpc_gateway() {
  f_print_installing_with_padding protoc-gen-grpc-gateway
  #f_needs_install "${DESTDIR}/${PREFIX}/bin/protoc-gen-grpc-gateway" || return 0

  pushd "${TEMPDIR}" >/dev/null
  go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
  popd >/dev/null
  f_print_done
}

f_install_protoc_gen_swagger() {
  f_print_installing_with_padding protoc-gen-swagger
  f_needs_install "${DESTDIR}/${PREFIX}/bin/protoc-gen-swagger" || return 0

  if ! which npm &>/dev/null ; then
    echo -e "\tNPM is not installed. Skipping."
    return 0
  fi

  pushd "${TEMPDIR}" >/dev/null
  go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
  npm install -g swagger-combine
  popd >/dev/null
  f_print_done
}

f_ensure_tools
f_ensure_dirs
f_install_protoc
#f_install_protoc_gen_gocosmos
f_install_protoc_gen_grpc_gateway
#f_install_protoc_gen_swagger
