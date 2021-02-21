#!/bin/bash

# Stop Glassfish server
~/Softwares/glassfish/glassfish5/glassfish/bin/asadmin stop-domain domain2

# Update repo
cd ~/Softwares/MarketplaceJ2EE_Server
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
cd graphql
mkdir ../outputs
touch ../outputs/graphql_server.log
go run ./ &>> ../outputs/graphql_server.log &
echo "Graphql go server launched !"
