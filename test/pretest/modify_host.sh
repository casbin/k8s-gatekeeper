#!/bin/bash

#get ip
ip=$(ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $2}'|tr -d "addr:")
#add to /etc/host
echo "${ip} webhook.domain.local" >> /etc/hosts
