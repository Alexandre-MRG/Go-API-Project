package dbPostgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type postgresDb struct { // Declare
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	Db       *sql.DB
}

func InitializeDB() (*postgresDb, error) {

	DbStruct := postgresDb{ // Initialize & Set
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "passpassEstiam",
		Dbname:   "postgres",
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DbStruct.Host, DbStruct.Port, DbStruct.User, DbStruct.Password, DbStruct.Dbname)

	db, err := sql.Open("postgres", psqlInfo) // vérification des paramètres de connexion à la db
	DbStruct.Db = db
	if err != nil {
		LogError("Erreur lors de la vérification des paramètres de connexion à la db", err)
	}
	//defer DbStruct.Db.Close()

	err = DbStruct.Db.Ping() // ping de la db pour vérifier si la connexion est bien établie
	if err != nil {
		LogError("Erreur lors de l'ouverture de la connexion à la db", err)
	}

	DbStruct.CreateTablePlayers()
	DbStruct.CreateTableWallets()

	return &DbStruct, err
}

func (dbStruct postgresDb) CreateTablePlayers() {
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS players (
		id SERIAL PRIMARY KEY,
		username TEXT,
		password TEXT,
		pin_code TEXT
	  );`

	_, err := dbStruct.Db.Exec(sqlStatement)
	if err != nil {
		LogError("Erreur lors de la création de la table 'players'", err)
	}
}

func (dbStruct postgresDb) CreateTableWallets() {
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS wallets (
		id SERIAL PRIMARY KEY,
		id_player INT,
		wallet_address TEXT,
		currency_code TEXT
	  );`

	_, err := dbStruct.Db.Exec(sqlStatement)
	if err != nil {
		LogError("Erreur lors de la création de la table 'wallets'", err)
	}
}

func (dbStruct postgresDb) AddPlayer(_username string, _password string, _pin_code string) int {

	sqlStatement := `
	INSERT INTO players (username, password, pin_code)
	VALUES ($1, $2, $3)
	RETURNING id`
	id := 0
	err := dbStruct.Db.QueryRow(sqlStatement, _username, _password, _pin_code).Scan(&id)
	if err != nil {
		LogError("Erreur lors de l'insertion d'un player dans la table 'players'", err)
	}
	fmt.Println("New player ID is:", id)

	return id
}

func (dbStruct postgresDb) AddWallet(_id_player int, _wallet_address string, _currency_code string) string {

	sqlStatement := `
	INSERT INTO wallets (id_player, wallet_address, currency_code)
	VALUES ($1, $2, $3)
	RETURNING wallet_address`
	wallet_address := ""
	err := dbStruct.Db.QueryRow(sqlStatement, _id_player, _wallet_address, _currency_code).Scan(&wallet_address)
	if err != nil {
		LogError("Erreur lors de l'insertion d'un wallet dans la table 'wallets'", err)
	}
	fmt.Println("New wallet address is:", wallet_address)

	return wallet_address
}

func (dbStruct postgresDb) CloseDB() {
	dbStruct.Db.Close()
}

func LogError(errorMessage string, err error) {
	log.Println(errorMessage, err)
}
