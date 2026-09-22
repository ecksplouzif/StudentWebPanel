#!/bin/bash
set -eEuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)"
project_dir=$(cd $SCRIPT_DIR; cd ..; pwd)
set -a
source $project_dir/.env
set +a
GOOS=linux GOARCH=amd64 go build -o "$project_dir/StudentWebPanel.new" $project_dir
scp -P $DEPLOY_PORT $project_dir/StudentWebPanel.new $SERVER_USER@$SERVER_IP:/opt/StudentWebPanel/
ssh -p $DEPLOY_PORT $SERVER_USER@$SERVER_IP "sudo systemctl stop StudentWebPanel && mv /opt/StudentWebPanel/StudentWebPanel /opt/StudentWebPanel/StudentWebPanel.old && mv /opt/StudentWebPanel/StudentWebPanel.new /opt/StudentWebPanel/StudentWebPanel && sudo systemctl start StudentWebPanel"
sleep 5 && curl -f $HEALTHZ_URL