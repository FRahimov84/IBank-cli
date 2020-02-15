package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	core "github.com/FRahimov84/IBank-core"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("can't start application manager%e", err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("can't close file %e", err)
		}
	}()
	log.SetOutput(file)
	log.Print("start application")
	log.Print("open db")
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("can't open database %e", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("can't close database %e", err)
		}
	}()
	err = db.Ping()
	if err != nil {
		log.Fatalf("can't ping database %e", err)
	}
	err = core.Init(db)
	if err != nil {
		log.Fatalf("can't init database %e", err)
	}
	log.Println("start operation loop")
	fmt.Println("вы вошли в систему как манеджер!")
	loop(db)

}
func loop(db *sql.DB) {
	for ; ; {
		var cmd string
		fmt.Println(commands)
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Println("Не удалось прочитать команду")
			log.Printf("can't scan cmd in loop() error: %e", err)
		}
		switch cmd {
		case "1":
			fmt.Println("Добавление пользователя...\nЗаполните следующие поля->")
			handleAddUser(db)
		case "2":
			fmt.Println("Добавление счета пользователю...\nЗаполните следующие поля->")
			handleAddBillToUser(db)
		case "3":
			fmt.Println("Добавление услуги...\nЗаполните следующие поля->")
			handleAddService(db)
		case "4":
			fmt.Println("Добавление Банкомата...\nЗаполните следующие поля->")
			handleAddATM(db)
		case "5":
			fmt.Println("Export - экспортировтаь")
			handleExport(db)
		case "6":
			fmt.Println("Import - Импортировать")
			handleImport()
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
	var password string
	fmt.Print("Имя: ")
	_, err := fmt.Scan(&name)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		log.Printf("can't scan name in handleAddUser() error: %e", err)
		return
	}
	fmt.Print("Фамилия: ")
	_, err = fmt.Scan(&surname)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		log.Printf("can't scan surname in handleAddUser() error: %e", err)
		return
	}
	fmt.Print("Телефон: ")
	_, err = fmt.Scan(&phone)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		log.Printf("can't scan phone in handleAddUser() error: %e", err)
		return
	}
	fmt.Print("login: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		log.Printf("can't scan login in handleAddUser() error: %e", err)
		return
	}
	fmt.Print("password: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		log.Printf("can't scan password in handleAddUser() error: %e", err)
		return
	}
	err = core.AddUser(db, login, password, name, surname, phone, false)
	if err != nil {
		fmt.Println("Не удалось добавить пользователя!")
		log.Printf("can't add user in handleAddUser() error: %e", err)
		return
	}
	fmt.Println("Новый пользователь добавлен!")
}

func handleAddBillToUser(db *sql.DB) {
	var id int
	fmt.Print("id пользователя: ")
	_, err := fmt.Scan(&id)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		log.Printf("can't scan field user_id in handleAddBillToUser() error: %e", err)
		return
	}
	var balance int
	fmt.Print("Баланс: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		log.Printf("can't scan field balance in handleAddBillToUser() error: %e", err)
		return
	}
	if balance < 0 {
		fmt.Println("Баланс не должен быть меньше 0!")
		log.Printf("balance should be more than zero in handleAddBillToUser()")
		return
	}
	err = core.AddBillToUser(db, id, balance, false)
	if err != nil {
		log.Printf("can't add bill to user in handleAddBillToUser() error: %e", err)
		fmt.Println("Не удалось добавить счет пользователю!")
	} else {
		fmt.Printf("счет пользовалелю id: %v добавлен!\n", id)
	}
}

func handleAddService(db *sql.DB) {
	fmt.Print("Название Услуги: ")
	reader := bufio.NewReader(os.Stdin)
	service, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		log.Printf("can't read field service in handleAddService() error: %e", err)
		return
	}
	var price int
	fmt.Print("Цена: ")
	_, err = fmt.Scan(&price)
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		log.Printf("can't scan field price in handleAddService() error: %e", err)
		return
	}
	err = core.AddService(db, service, price)
	if err != nil {
		log.Printf("can't add service in handleAddService() error: %e", err)
		fmt.Println("Не удалось добавить услугу!")
	} else {
		fmt.Println("Услуга добавлена!")
	}
}

func handleAddATM(db *sql.DB) {
	fmt.Print("Аддресс банкомата: ")
	reader := bufio.NewReader(os.Stdin)
	address, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("что-то пошло не так, попробуйте еще раз!")
		log.Printf("can't read field address in handleAddATM() error: %e", err)
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
		fmt.Println("Что вы хотите экспортировать?")
		fmt.Println(exportImportCommands)
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Println("что-то пошло не так, попробуйте еще раз!")
			log.Printf("can't scan command in handleExport() error: %e", err)
			continue
		}
		switch cmd {
		case "1":
			fmt.Println("Exporting users...")
			list, err := core.UsersList(db)
			if err != nil {
				log.Printf("can't get users list in handleExport() error: %e", err)
				fmt.Println("Не удалось получить пользователей!")
				continue
			}
			users := core.List{
				UsersList:    list,
				ATMsList:     nil,
				BillUserList: nil,
			}
			handleExportToFile(users, 1)
		case "2":
			fmt.Println("Exporting bills...")
			list, err := core.BillsWithUserList(db)
			if err != nil {
				fmt.Println("Не удалось получить счета!")
				log.Printf("can't get Bills list in handleExport() error: %e", err)
				continue
			}
			bills := core.List{
				UsersList:    nil,
				ATMsList:     nil,
				BillUserList: list,
			}
			handleExportToFile(bills, 2)
		case "3":
			fmt.Println("Exporting ATMs...")
			list, err := core.ATMsList(db)
			if err != nil {
				fmt.Println("Не удалось получить список банкоматов")
				log.Printf("can't get ATMs list in handleExport() error: %e", err)
				continue
			}
			atms := core.List{
				UsersList:    nil,
				ATMsList:     list,
				BillUserList: nil,
			}
			handleExportToFile(atms, 0)
		case "q":
			return
		default:
			fmt.Printf("Неправильная команда %s\n", cmd)
		}
	}
}

func handleExportToFile(list core.List, index int) {
	fileName := []string{"ATMs", "Users", "Bills"}
	for ; ; {
		var cmd string
		fmt.Println(jsonXmlCommands)
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Println("Что-то пошло не так")
			log.Printf("can't scan command in handleExportToFile() error: %e", err)
			continue
		}
		switch cmd {
		case "1":
			bytes, err := json.Marshal(list)
			if err != nil {
				fmt.Println("Не удалось форматировать в json")
				log.Printf("can't make to Json format in handleExportToFile() error: %v", err)
				return
			}
			err = ioutil.WriteFile(Directory+fileName[index]+".json", bytes, 0666)
			if err != nil {
				fmt.Println("Не удалось записать в файл")
				log.Printf("can't write to file %s.json in handleExportToFile() error: %e\n", fileName[index], err)
				return
			}
			fmt.Println("Успешно...")
			return
		case "2":
			bytes, err := xml.Marshal(list)
			if err != nil {
				fmt.Println("Не удалось форматировать в xml")
				log.Printf("can't make to Xml format in handleExportToFile() error: %e", err)
				return
			}
			err = ioutil.WriteFile(Directory+fileName[index]+".xml", bytes, 0666)
			if err != nil {
				fmt.Println("Не удалось записать в файл")
				log.Printf("can't write to file %s.xml in handleExportToFile() error: %e\n", fileName[index], err)
				return
			}
			fmt.Println("Успешно...")
			return
		case "q":
			return
		default:
			fmt.Printf("Неправильная команда %s\n", cmd)
		}
	}
}

func handleImport() {

	files, err := ioutil.ReadDir("./results/")
	if err != nil {
		fmt.Println("не получается прочесть директорию result")
		log.Printf("can't read files from directory in handleImport() error: %e", err)
	}
	fileNames := make([]string, 0)
	t := 0
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".json") || strings.HasSuffix(file.Name(), ".xml")) {
			t++
			fmt.Printf("%v - %s\n", t, file.Name())
			fileNames = append(fileNames, file.Name())
		}
	}
	if t == 0 {
		fmt.Println(`Директория "results" не содержит файлы формата *.json и *.xml!`)
		return
	}
	fmt.Println(importInfo)
	var cmd int
	_, err = fmt.Scan(&cmd)
	if err != nil && cmd > 0 && cmd <= t {
		fmt.Println("Не правельная команда!")
		log.Printf("can't read command in handleImport() files: %v, cmd: %v, err: %e", t, cmd, err)
		return
	}
	bytes, err := ioutil.ReadFile(Directory + fileNames[cmd-1])
	if err != nil {
		log.Printf("in handleImport() can't open file: %s error: %e\n", fileNames[cmd-1], err)
		fmt.Println("Не удается открыть файл")
		return
	}
	if strings.HasSuffix(fileNames[cmd-1], ".json") {
		var list core.List
		err = json.Unmarshal(bytes, &list)
		if err != nil {
			fmt.Println("Не получается прочесть json файл")
			log.Printf("in handleImport() can't unmarshal json file: %s error: %e\n", fileNames[cmd-1], err)
			return
		}
		if list.ATMsList != nil {
			fmt.Printf("%s\t%s\t%s\n", "id", "Адрес", "Статус")
			for _, atm := range list.ATMsList {
				fmt.Printf("%v\t%v\t%v\n", atm.Id, atm.Address, atm.Locked)
			}
		}
		if list.BillUserList != nil {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\n", "id", "Баланс", "СтатусСчета", "Имя", "Фамилия", "Телефон", "СтатусПольз.")
			for _, bill := range list.BillUserList {
				fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\n", bill.Id, bill.Balance, bill.LockedBill, bill.UserName, bill.UserSurname, bill.UserPhone, bill.LockedUser)
			}
		}
		if list.UsersList != nil {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", "id", "Имя", "Фамилия", "Телефон", "Статус")
			for _, user := range list.UsersList {
				fmt.Printf("%v\t%v\t%v\t%v\t%v\n", user.Id, user.Name, user.Surname, user.Phone, user.Locked)
			}
		}
	}
	if strings.HasSuffix(fileNames[cmd-1], ".xml") {
		var list core.List
		err = xml.Unmarshal(bytes, &list)
		if err != nil {
			fmt.Println("Не получается данные привести в формат xml")
			log.Printf("can't unmarshal xml file: %s error: %e\n", fileNames[cmd-1], err)
			return
		}
		if list.ATMsList != nil {
			fmt.Printf("Список банкоматов\n%s\t%s\t%s\n", "id", "Адрес", "Статус")
			for _, atm := range list.ATMsList {
				fmt.Printf("%v\t%v\t%v\n", atm.Id, atm.Address, atm.Locked)
			}
		}
		if list.BillUserList != nil {
			fmt.Printf("Список счетов\n%s\t%s\t%s\t%s\t%s\t%s\t%s\n", "id", "Баланс", "СтатусСчета", "Имя", "Фамилия", "Телефон", "СтатусПольз.")
			for _, bill := range list.BillUserList {
				fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\n", bill.Id, bill.Balance, bill.LockedBill, bill.UserName, bill.UserSurname, bill.UserPhone, bill.LockedUser)
			}
		}
		if list.UsersList != nil {
			fmt.Printf("Список пользователей\n%s\t%s\t%s\t%s\t%s\n", "id", "Имя", "Фамилия", "Телефон", "Статус")
			for _, user := range list.UsersList {
				fmt.Printf("%v\t%v\t%v\t%v\t%v\n", user.Id, user.Name, user.Surname, user.Phone, user.Locked)
			}
		}
	}
}
