<h1 align="center">MarketplaceJ2EE_Server</h1>

> Disclaimer: This is a project related to academic work. It may be not complete.

---

## Table of Contents

- [Presentation](#presentation)
- [Architecture](#architecture)
- [Deploiement](#deploiement)
- [Paiement](#paiement)
- [Author](#author)

## Presentation

Bienvenue sur la partie backend de mon projet marketplace !

Ce projet repose sur une architecture composé de deux serveurs:

- Un serveur J2EE supportant une api SOAP ainsi qu'une REST
- Un serveur Go avec une api GraphQL (situé dans le sous-répertoire ```graphql/```)

Le tout est hébergé sur un VPS sous Ubuntu server 20.04.1 LTS et avec un certificat SSL auto-signé. Toutes les requêtes ce font donc en https uniquement.
Ce choix d'hébergement à été choisi pour sa flexibilité et car Heroku tourne avec Tomcat qui causais des problèmes avec la partie J2EE, j'utilise donc Glassfish 5 sur le VPS.
En ce qui concerne le stockage des articles, tout ce passe dans une base de données **PostgreSQL**, des données par défault y sont insérées lors de la création de la table pour l'instant car la soumission d'article n'a pas encore été implémenté.

Vous pouvez retrouver le client web [ici](https://github.com/hbollon/MarketplaceJ2EE_Client)

## Architecture

Comme vu précédemment nous avons 3 api différentes:

- SOAP : Utilisé pour calculé et obtenir les frais d'envoi d'un article.
- REST : Utilisé pour effectué les actions de paiement (voir partie paiement pour plus de détails)
- GraphQL : Utilisé pour communiquer avec la base de donnée (et récupérer tout les articles ou seulement certains)

La partie GraphQL nécessitant des *secrets*, notamment pour la base de donnée, j'utilise donc un fichier **.env** ignoré sur git et créer localement à partir de ```.env.exemple```.

## Deploiement

Tout le processus de déploiement se fait facilement et automatiquement grâce au script ```deploy.sh```. Il s'occupe de mettre à jour le repo, build le projet maven et déployer les différents serveurs:

```sh
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
    ./GraphQL-Server &>> ../outputs/graphql_server.log & echo $! > graphql.pid
    echo "Graphql go server launched !"
fi
```

## Paiement

Pour la gestion des paiement avec MangoPay j'ai opté pour un ***Web PayIn*** (page de paiement générée par MangoPay), cela nécessite d'enregistrer le client sur MangoPay.

Voici la routine effectué lors que l'ont fait une requête à l'api rest:

1. Nous regardons si l'utilisateur est "Client" sur MangoPay, si oui alors nous récupérons ses infos et passons directement à l'étape **3**
2. Dans le cas contraire, nous l'enregistrons
3. Ensuite, nous effectuons une requête à MangoPay pour obtenir un Web PayIn en lui fournissant les informations du client ainsi que celles de l'article acheté.
4. Enfin, nous redirigeons le client web vers l'url du Web PayIn reçu

Documentation: https://docs.mangopay.com/endpoints/v2.01/payins#e269_create-a-card-web-payin

## Author

👤 **Hugo Bollon**

* Github: [@hbollon](https://github.com/hbollon)
* LinkedIn: [@Hugo Bollon](https://www.linkedin.com/in/hugobollon/)
* Portfolio: [hugobollon.me](https://www.hugobollon.me)

## Show your support

Give a ⭐️ if this project helped you!