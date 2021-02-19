#!/bin/bash

# Restart Glassfish domain and deploy java j2ee server
~/Softwares/glassfish/glassfish5/glassfish/bin/asadmin stop-domain domain2
cd ~/Softwares/MarketplaceJ2EE_Server
git reset --hard
git pull
mvn clean install
~/Softwares/glassfish/glassfish5/glassfish/bin/asadmin start-domain domain2
~/Softwares/glassfish/glassfish5/glassfish/bin/asadmin --port 4848 --host localhost deploy target/MarketplaceServer-1.0-SNAPSHOT.war
echo "J2EE server launched !"

# Launch graphql go server
cd graphql
mkdir ../outputs
touch ../outputs/graphql_server.log
go run ./ &>> ../outputs/graphql_server.log &
echo "Graphql go server launched !"
