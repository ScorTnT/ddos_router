#!/bin/sh
# 이 스크립트에는 OpenWRT에 설치해야 하는 패키지를 설치하는 명령어가 포함되어 있음.
opkg update && opkg install \
    # TODO: 필요한 패키지 목록 채워넣기