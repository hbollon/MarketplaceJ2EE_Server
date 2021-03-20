#!/bin/bash

# Check if the deploy script is already running
for pid in $(pidof -x deploy.sh); do
    echo $pid
    if [ $pid != $$ ]; then
        echo "Deploy script already running!"
        exit 2
    fi
done

# Stop Glassfish server
~/Softwares/glassfish/glassfish5/glassfish/bin/asadmin stop-domain domain2

# Check if the Go server is already running and ask for killing if yes
cd ~/Softwares/MarketplaceJ2EE_Server
PID_FILE="graphql/graphql.pid"
RUNNING=false
if test -f "$PID_FILE"; then
    echo $PID_FILE
    if kill -0 $(head -n 1 $PID_FILE) > /dev/null 2>&1; then
        read -p "The Go server is already running. Do you want to kill it and continue? (y/n)" -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]
        then
            kill $(head -n 1 $PID_FILE)
        else
            $RUNNING=true
        fi
    fi
fi

# Update git repository
git reset --hard
git pull

# Install maven dependencies and build project
mvn clean install

# Start glassfish domain, login and redeploy project
~/Softwares/glassfish/glassfish5/glassfish/bin/asadmin start-domain domain2
~/Softwares/glassfish/glassfish5/glassfish/bin/asadmin login
~/Softwares/glassfish/glassfish5/glassfish/bin/asadmin undeploy MarketplaceServer-1.0-SNAPSHOT
~/Softwares/glassfish/glassfish5/glassfish/bin/asadmin --port 4848 --host localhost deploy target/MarketplaceServer-1.0-SNAPSHOT.war
echo "J2EE server launched !"

# Launch graphql go server
if $RUNNING; then
    echo "Graphql go server restart skipped !"
else
    cd graphql
    mkdir ../outputs
    touch ../outputs/graphql_server.log
    go build -o GraphQL-Server ./

    read -p "Do you want to reset the database? (y/n)" -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]
    then
        ./GraphQL-Server --reset-db &>> ../outputs/graphql_server.log & echo $! > graphql.pid
    else
        ./GraphQL-Server &>> ../outputs/graphql_server.log & echo $! > graphql.pid
    fi

    echo "Graphql go server launched !"
fi