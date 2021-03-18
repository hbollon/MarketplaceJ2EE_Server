#!/bin/bash

# Stop Glassfish server
~/Softwares/glassfish/glassfish5/glassfish/bin/asadmin stop-domain domain2

# Check if the Go server is already running and ask for killing if yes
cd ~/Softwares/MarketplaceJ2EE_Server
PID_FILE="graphql/pid.txt"
RUNNING = false
if test -f "$PID_FILE"; then
    echo $PID_FILE
    if kill -0 $(head -n 1 $PID_FILE) > /dev/null 2>&1; then
        read -p "The Go server is already running. Do you want to kill it and continue? (y/n)" -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]
        then
            kill $(head -n 1 $PID_FILE)
        else
            RUNNING = true
        fi
    fi
fi

# Update repo
git reset --hard
git pull

# Install maven dependencies and build project
mvn clean install

# Start glassfish domain and redeploy project
~/Softwares/glassfish/glassfish5/glassfish/bin/asadmin start-domain domain2
~/Softwares/glassfish/glassfish5/glassfish/bin/asadmin undeploy MarketplaceServer-1.0-SNAPSHOT
~/Softwares/glassfish/glassfish5/glassfish/bin/asadmin --port 4848 --host localhost deploy target/MarketplaceServer-1.0-SNAPSHOT.war
echo "J2EE server launched !"

# Launch graphql go server
if !RUNNING then 
    cd graphql
    mkdir ../outputs
    touch ../outputs/graphql_server.log
    go run ./ &>> ../outputs/graphql_server.log & echo $! > graphql/pid.txt
    echo "Graphql go server launched !"
else 
    echo "Graphql go server starting skipped !"
fi
