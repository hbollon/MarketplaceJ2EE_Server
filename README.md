<h1 align="center">MarketplaceJ2EE_Server</h1>

> Disclaimer: This is a project related to academic work. It may be not complete.

---

## Table of Contents

- [Presentation](#presentation)
- [Architecture](#architecture)
  - [SOAP](#soap)
  - [REST](#rest)
  - [GraphQL](#graphql)
- [Base de donnée](#base-de-donnée)
- [Deploiement](#deploiement)
- [Paiement](#paiement)
- [Flags](#flags)
- [Author](#author)

## Presentation

Bienvenue sur la partie backend de mon projet marketplace !

Ce projet repose sur une architecture composé de deux serveurs:

- Un serveur J2EE supportant une api SOAP ainsi qu'une REST
- Un serveur Go avec une api GraphQL (situé dans le sous-répertoire ```graphql/```)

Le tout est hébergé sur un VPS sous Ubuntu server 20.04.1 LTS et avec un certificat SSL auto-signé. Toutes les requêtes ce font donc en https uniquement.
Ce choix d'hébergement à été choisi pour sa flexibilité et car Heroku tourne avec Tomcat qui causais des problèmes avec la partie J2EE, j'utilise donc Glassfish 5 sur le VPS.
En ce qui concerne le stockage des articles, tout ce passe dans une base de données **PostgreSQL**.

Vous pouvez retrouver le client web [ici](https://github.com/hbollon/MarketplaceJ2EE_Client)

## Architecture

Comme vu précédemment nous avons 3 api différentes:

- SOAP : Utilisé pour calculé et obtenir les frais d'envoi d'un article.
- REST : Utilisé pour effectué les actions de paiement (voir partie paiement pour plus de détails)
- GraphQL : Utilisé pour communiquer avec la base de donnée.

La partie GraphQL nécessitant des *secrets*, notamment pour la base de donnée, j'utilise donc un fichier **.env** ignoré sur git et créer localement à partir de ```.env.exemple```.

### SOAP

La partie SOAP sert uniquement pour le calcul des frais de livraison.
Le wsdl est généré automatiquement par le serveur J2EE et ressemble à ça en l'état:

```xml
<!--  Published by JAX-WS RI (http://jax-ws.java.net). RI's version is Metro/2.4.2 (UNKNOWN-17d59c9; 2018-09-26T15:16:30-0700) JAXWS-RI/2.3.1 JAXWS-API/2.3.1 JAXB-RI/2.3.1 JAXB-API/2.3.1 svn-revision#unknown.  -->
<!--  Generated by JAX-WS RI (http://javaee.github.io/metro-jax-ws). RI's version is Metro/2.4.2 (UNKNOWN-17d59c9; 2018-09-26T15:16:30-0700) JAXWS-RI/2.3.1 JAXWS-API/2.3.1 JAXB-RI/2.3.1 JAXB-API/2.3.1 svn-revision#unknown.  -->
<definitions xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd" xmlns:wsp="http://www.w3.org/ns/ws-policy" xmlns:wsp1_2="http://schemas.xmlsoap.org/ws/2004/09/policy" xmlns:wsam="http://www.w3.org/2007/05/addressing/metadata" xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/" xmlns:tns="http://MarketplaceServer.bitsplease.com/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns="http://schemas.xmlsoap.org/wsdl/" targetNamespace="http://MarketplaceServer.bitsplease.com/" name="DeliveryFeeService">
    <types>
        <xsd:schema>
            <xsd:import namespace="http://MarketplaceServer.bitsplease.com/" schemaLocation="http://localhost:8080/MarketplaceServer-1.0-SNAPSHOT/services/DeliveryFee?xsd=1" />
        </xsd:schema>
    </types>
    <message name="CalculateDelivery">
        <part name="parameters" element="tns:CalculateDelivery" />
    </message>
    <message name="CalculateDeliveryResponse">
        <part name="parameters" element="tns:CalculateDeliveryResponse" />
    </message>
    <portType name="DeliveryFee">
        <operation name="CalculateDelivery">
            <input wsam:Action="http://MarketplaceServer.bitsplease.com/DeliveryFee/CalculateDeliveryRequest" message="tns:CalculateDelivery" />
            <output wsam:Action="http://MarketplaceServer.bitsplease.com/DeliveryFee/CalculateDeliveryResponse" message="tns:CalculateDeliveryResponse" />
        </operation>
    </portType>
    <binding name="DeliveryFeePortBinding" type="tns:DeliveryFee">
        <soap:binding transport="http://schemas.xmlsoap.org/soap/http" style="document" />
        <operation name="CalculateDelivery">
            <soap:operation soapAction="" />
            <input>
                <soap:body use="literal" />
            </input>
            <output>
                <soap:body use="literal" />
            </output>
        </operation>
    </binding>
    <service name="DeliveryFeeService">
        <port name="DeliveryFeePort" binding="tns:DeliveryFeePortBinding">
            <soap:address location="http://localhost:8080/MarketplaceServer-1.0-SNAPSHOT/services/DeliveryFee" />
        </port>
    </service>
</definitions>
```

### REST

L'Api REST, quant à elle, remplit le rôle de passerelle avec l'api **Mangopay**.
Il y a deux endpoints:

#### - ```/mangopay/pay```:

Cet endpoint consomme du json représentant un client (l'acheteur) et un produit. 
Avec ceci, le serveur va effectué la routine suivante:

- Vérifié sur le client a déjà un compte client _MangoPay_
    - Si oui, il récupère son _clientid_
    - Si non, il l'inscrit et récupère son _clientid_
- Effectue une requête pour obtenir une _Web PayIn_ auprès de _MangoPay_ au nom de l'acheteur et pour l'article désiré
- Retourne l'url de la _Web PayIn_ reçue au client web

Exemple de json:
```json
{
    "client": {
        "firstName": "Hugo",
        "lastName": "Bollon",
        "email": "hugo.bollon@gmail.com",
    },
    "product": {
        "seller": {
            "firstName":"Laurent",
            "lastName":"Cutting",
            "email":"laurent.cutting@gmail.com",
            "walletId":104312216 // MangoPay wallet id, pour savoir quel seller créditer
        },
        "name":"Valheim",
        "description":"Valheim est un jeu vidéo [...]",
        "quantity":1,
        "weight":10,
        "price":16.99,
        "fees":2, // Frais de livraison calculés au préalable via SOAP
        "assetUrl":"https://www.actugaming.net/wp-content/uploads/2021/02/valheim-cover.jpg"
    }
}
```

#### - ```/seller/register```:

Cet endpoint consomme du json représentant un nouveau seller.
Avec ceci, le serveur va effectué la routine suivante:

- Vérifié sur le client a déjà un compte client _MangoPay_
  - Si oui, il récupère son _clientid_
  - Si non, il l'inscrit et récupère son _clientid_
- Effectue une requête pour créer un nouveau _Wallet_ pour la marketplace
- Retourne le _walletId_

Cet endpoint est utilisé par le serveur GraphQL quand un nouveau seller est inscrit sur la Marketplace afin d'obtenir son walletId pour la base de donnée.

Exemple de json:
```json
{
    "seller": {
        "firstName": "Hugo",
        "lastName": "Bollon",
        "email": "hugo.bollon@gmail.com",
    },
}
```

### GraphQL
En ce qui concerne l'API GraphQL, elle est géré par un serveur Golang et permet de faire le lien entre le client et la base de donnée.

J'ai défini deux modèles GraphQL:

- **seller**: 
```go
var sellerType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Seller",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"firstName": &graphql.Field{
			Type: graphql.String,
		},
		"lastName": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"walletId": &graphql.Field{
			Type: graphql.Int,
		},
	},
})
```

- **product**:
```go
var productType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"quantity": &graphql.Field{
			Type: graphql.Int,
		},
		"weight": &graphql.Field{
			Type: graphql.Float,
		},
		"price": &graphql.Field{
			Type: graphql.Float,
		},
		"asset_url": &graphql.Field{
			Type: graphql.String,
		},
		"seller": &graphql.Field{
			Type: sellerType,
		},
	},
})
```

Ce dernier comporte un champ _seller_ qui est donc un GraphQL Object et permet de sélectionner (lors d'un "GET" par exemple) les champs à retourner pour celui-ci en plus de ceux de product (voir plus bas).

#### Queries

##### - ```products```:

Retourne tout les produits du catalogue
Exemple:
```graphql
products {
    id
    name
    description
    quantity
    weight
    price
    asset_url
    seller {
        id
        firstName
        lastName
        email
        walletId
    }
}
```

##### - ```product```:

Cherche et retourne un produit spécifique via le nom ou l'id
Exemple:
```graphql
product (name: "Valheim") {
    id
    name
    description
    quantity
    weight
    price
    asset_url
    seller {
        id
        firstName
        lastName
        email
        walletId
    }
}
```

##### - ```sellProduct```:

Ajoute un article au catalogue (base de donnée)
Exemple:
```graphql
sellProduct (
    name: ""
    description: ""
    quantity: 10 // Optional
    weight: 10.50
    price: 59.99
    asset_url: "" // Optional
    seller: <seller email>
) 
```

##### - ```sellers```:

Retourne tout les sellers inscrit à la marketplace
Exemple:
```graphql
sellers {
    id
    firstName
    lastName
    email
    walletId
}
```

##### - ```seller```:

Cherche et retourne un seller spécifique via l'email
Exemple:
```graphql
seller (email: "hugo.bollon@gmail.com") {
    id
    firstName
    lastName
    email
    walletId
}
```

##### - ```registerSeller```:

Inscrit un seller à la marketplace (appelle l'api rest pour lui crée un compte client et un wallet _MangoPay_)
Exemple:
```graphql
registerSeller (
    firstName: ""
    lastName: ""
    email: ""
) 
```

**Tout les endpoints nécessitant des paramètres sont validés et retourne une erreur si la validation des inputs ne passe pas**

## Base de donnée

J'ai choisis d'utiliser une base de donnée **PostgreSQL** car je n'avais pas besoin des avantages d'une base NoSQL dans le cadre de ce TP et que je voulais employé la puissance de Postgre comparé à MySQL.
Toutes les interactions avec cette dernière se font via l'API GraphQL (et donc par **Golang**).

Cette BDD est composé de deux tables:

- **Seller**: id, firstName, lastName, email, walletId
- **Product**: id, name, description, quantity, weight, price, asset_url, seller_id (foreign key)

Structure des tables détaillée (obtenue via _DESCRIBE TABLE_):

```bash
                                      Table "public.seller"
   Column   |          Type          | Collation | Nullable |              Default               
------------+------------------------+-----------+----------+------------------------------------
 id         | integer                |           | not null | nextval('seller_id_seq'::regclass)
 first_name | character varying(40)  |           | not null | 
 last_name  | character varying(40)  |           | not null | 
 email      | character varying(255) |           | not null | 
 wallet_id  | integer                |           | not null | 
Indexes:
    "seller_pkey" PRIMARY KEY, btree (id)
    "seller_email_key" UNIQUE CONSTRAINT, btree (email)
    "seller_wallet_id_key" UNIQUE CONSTRAINT, btree (wallet_id)
Referenced by:
    TABLE "product" CONSTRAINT "fk_seller" FOREIGN KEY (seller_id) REFERENCES seller(id) ON DELETE SET NULL

                                     Table "public.product"
   Column    |         Type          | Collation | Nullable |               Default               
-------------+-----------------------+-----------+----------+-------------------------------------
 id          | integer               |           | not null | nextval('product_id_seq'::regclass)
 name        | character varying(40) |           | not null | 
 description | text                  |           | not null | 
 quantity    | integer               |           | not null | 
 weight      | real                  |           | not null | 
 price       | real                  |           | not null | 
 asset_url   | text                  |           |          | 
 seller_id   | integer               |           | not null | 
Indexes:
    "product_pkey" PRIMARY KEY, btree (id)
    "product_name_key" UNIQUE CONSTRAINT, btree (name)
Foreign-key constraints:
    "fk_seller" FOREIGN KEY (seller_id) REFERENCES seller(id) ON DELETE SET NULL
```

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