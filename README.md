# Go-API-Project
Une API Rest en Golang qui simule la création d'un utilisateur et son wallet. L'API contacte un service de blockchain mock.

## Introduction

Avant toute chose il faut télécharger Go pour votre système : https://go.dev/dl/

Le projet est décomposé en deux parties distinctes.

- API_Go : cette api peut être contactée par l'utilisateur pour créer son compte et le wallet associé
- Mock-Blockchain_GO : ce serveur peut être contacté par l'API_Go pour créer un wallet et en recevoir les informations

## Démarrage

- En se plaçant à la racine du répertoire *API_Go*, ouvrir un terminal et lancer la commande :
```go
go run .
```

- En se plaçant à la racine du répertoire *Mock_Blockchain_GO*, ouvrir un terminal et lancer la commande :
```go
go run .
```

## Utilisation

Pour utiliser le service il faut contacter l'API_Go en utilisant un outil comme *Postman* ou *Rester*. Voici le type de requête à faire :

![image](https://user-images.githubusercontent.com/62239442/190095932-9aebf4c4-729b-4aa4-a7c1-579b595ba1ce.png)

## Tests

Pour lancer les tests, ouvrir un terminal dans un répertoire contenant des fichiers au format *fichier_test.go*. 

Exécuter l'une des commandes suivantes : 
- Lancer les tests unitaires : 
```go
go test
```
- Lancer le benchmark d'une fonction :
```go
go test -bench=HandleNewPlayer
```
