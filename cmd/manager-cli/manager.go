package main

import (
	"bufio"
	"database/sql"
	"fmt"
	core "github.com/FRahimov84/IBank-core"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("can't open database %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("can't close database %v", err)
		}
	}()
	err = db.Ping()
	if err != nil {
		log.Fatalf("can't ping database %v", err)
	}
	err = core.Init(db)
	if err != nil {
		log.Fatalf("can't init database %v", err)
	}
	log.Println("start operation loop")
	loop(db)

}
func loop(db *sql.DB) {
	for ; ; {
		var cmd string
		fmt.Println(commands)
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Errorf("can't scan commad %w", err)
		}
		switch cmd {
		case "1":
			fmt.Println("Adding user...\nFill the fields->")
			handleAddUser(db)
		case "2":
			fmt.Println("Adding bill to user...\nFill the fields->")
			handleAddBillToUser(db)
		case "3":
			fmt.Println("Adding service...\nFill the fields->")
			handleAddService(db)
		case "4":
			fmt.Println("Adding ATM...\nFill the fields->")
			handleAddATM(db)
		case "5":
			fmt.Println("Export")
			handleExport(db)
		case "6":
			fmt.Println("Import")
			handleImport(db)
		case "q":
			return
		default:
			fmt.Printf("Неправильная команда %s\n", cmd)
		}
	}

}
func handleAddUser(db *sql.DB) {
	var name string
	var surname string
	var phone string
	var login string
	var pass string
	fmt.Print("Имя: ")
	_, err := fmt.Scan(&name)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		fmt.Errorf("can't scan field name %w", err)
		return
	}
	fmt.Print("Фамилия: ")
	_, err = fmt.Scan(&surname)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		fmt.Errorf("can't scan field surname %w", err)
		return
	}
	fmt.Print("Телефон: ")
	_, err = fmt.Scan(&phone)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		fmt.Errorf("can't scan field PhoneNumber %w", err)
		return
	}
	fmt.Print("login: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		fmt.Errorf("can't scan field login %w", err)
		return
	}
	fmt.Print("password: ")
	_, err = fmt.Scan(&pass)
	if err != nil {
		fmt.Errorf("can't scan field password %w", err)
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		return
	}
	err = core.AddUser(db, login, pass, name, surname, phone, false)
	if err != nil {
		fmt.Errorf("can't add new user: %w", err)
		fmt.Println("Не удалось добавить пользователя!")
	} else {
		fmt.Println("Новый пользователь добавлен!")
	}
}

func handleAddBillToUser(db *sql.DB) {
	var id int
	fmt.Print("id пользователя: ")
	_, err := fmt.Scan(&id)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		fmt.Errorf("can't scan field user_id %w", err)
		return
	}
	var balance int
	fmt.Print("Баланс: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		fmt.Errorf("can't scan field balance %w", err)
		return
	}
	err = core.AddBillToUser(db, id, balance, false)
	if err != nil {
		fmt.Errorf("can't add bill to user: %w", err)
		fmt.Println("Не удалось добавить счет пользователю!")
	} else {
		fmt.Printf("счет пользовалелю id: %v добавлен!\n", id)
	}
}

func handleAddService(db *sql.DB)  {
	fmt.Print("Название Услуги: ")
	reader := bufio.NewReader(os.Stdin)
	service, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		fmt.Errorf("can't read field service %w", err)
		return
	}
	var price int
	fmt.Print("Цена: ")
	_, err = fmt.Scan(&price)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		fmt.Errorf("can't scan field price %w", err)
		return
	}
	err = core.AddService(db, service, price)
	if err != nil {
		fmt.Errorf("can't add service: %w", err)
		fmt.Println("Не удалось добавить услугу!")
	} else {
		fmt.Println("Услуга добавлена!")
	}
}

func handleAddATM(db *sql.DB)  {
	fmt.Print("Аддресс банкомата: ")
	reader := bufio.NewReader(os.Stdin)
	address, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		fmt.Errorf("can't read field address %w", err)
		return
	}
	err = core.AddATM(db, address, false)
	if err != nil {
		fmt.Errorf("can't add ATM: %w", err)
		fmt.Println("Не удалось добавить Банкомат!")
	} else {
		fmt.Println("Банкомат добавлен!")
	}
}

func handleExport(db *sql.DB) {
	for ; ; {
		var cmd string
		fmt.Println(exportCommands)
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Errorf("can't scan command %w", err)
		}
		switch cmd {
		case "1":
			fmt.Println("Exporting users...")
			list, err := core.UsersList(db)
			if err != nil {
				fmt.Errorf("can't get users list %w", err)
			}
			handleExportUsers(list)
		case "2":
			fmt.Println("Exporting bills...")
			handleExportBills(db)
		case "3":
			fmt.Println("Exporting ATMs...")
			handleExportATMs(db)
		case "q":
			return
		default:
			fmt.Printf("Неправильная команда %s\n", cmd)
		}
	}
}

func handleExportUsers(users []core.UserList) {
	for ; ; {
		var cmd string
		fmt.Println(jsonXmlCommands)
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Errorf("can't scan command %w", err)
		}
		switch cmd {
		case "1":

		case "2":

		case "q":
			return
		default:
			fmt.Printf("Неправильная команда %s\n", cmd)
		}
	}
}

func handleExportUsersToJson(){

}

func handleExportXml(db *sql.DB)  {

}

func handleImport(db *sql.DB){
	for ; ; {
		var cmd string
		fmt.Println(jsonXML)
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Errorf("can't scan commad %w", err)
		}
		switch cmd {
		case "1":
			fmt.Println("Exporting to JSON format...")
			handleExportJson(db)
		case "2":
			handleExportXml(db)
		case "q":
			return
		default:
			fmt.Printf("Неправильная команда %s\n", cmd)
		}
	}
}