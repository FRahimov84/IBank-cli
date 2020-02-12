package main

import (
	"database/sql"
	"fmt"
	core "github.com/FRahimov84/IBank-core"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("can't open data base %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("can't close data base %v", err)
		}
	}()
	err = db.Ping()
	if err != nil {
		log.Fatalf("can't ping data base %v", err)
	}
	err = core.Init(db)
	if err != nil {
		log.Fatalf("can't init database %v", err)
	}
	fmt.Println("Добро пожаловать")
	loop(db)
}

func loop(db *sql.DB) {
	for ; ; {
		var cmd string
		fmt.Println("\nДоступные вам команды..", unauthorizedСommands)
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Errorf("can't scan command %w", err)
		}
		switch cmd {
		case "1":
			fmt.Println("Для входа запоните следуюшие поля...")
			handleLogin(db)
		case "2":
			fmt.Println("Список банкоматов...")
			handleGetATMsList(db)
		case "q":
			return
		default:
			fmt.Printf("Неправильная команда %s\n", cmd)
		}
	}
}

func handleLogin(db *sql.DB) {
	var login string
	var pass string
	fmt.Println("Логин: ")
	_, err := fmt.Scan(&login)
	if err != nil {
		fmt.Errorf("can't scan login %w", err)
	}
	fmt.Println("Пароль: ")
	_, err = fmt.Scan(&pass)
	if err != nil {
		fmt.Errorf("can't scan password %w", err)
	}
	user_id, name, err := core.Login(db, login, pass)

	if err != nil {
		fmt.Println("Не удалось войти в систему! неправильный логин или пароль")
		return
	}
	fmt.Printf("Вы вошли в систему как %s!\n", name)
	operationLoopAuthorized(user_id, db)
	fmt.Printf("Выход из системы. Пока %s.\n", name)
}

func operationLoopAuthorized(user_id int, db *sql.DB) {
	for ; ; {
		var cmd string
		fmt.Println("\nДоступные вам команды...", authorizedСommands)
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Errorf("can't scan command %w", err)
			continue
		}
		switch cmd {
		case "1":
			fmt.Println("Список вашых счетов...")
			handleGetBillsOfUser(db, user_id)
		case "2":
			fmt.Println("Перевод денег другому клиенту...")
			handleTransfer(db, user_id)
		case "3":
			fmt.Println("Оплатить услугу...")
			handlePayService(db, user_id)
		case "4":
			fmt.Println("Список Банкоматов...")
			handleGetATMsList(db)
		case "5":
			fmt.Println("Список Услуг...")
			handleGetServicesList(db)
		case "q":
			return
		default:
			fmt.Printf("Неправильная команда %s\n", cmd)
		}
	}
}

func handleGetBillsOfUser(db *sql.DB, user_id int) {
	list, err := core.UserBills(db, user_id)
	if err != nil {
		fmt.Errorf("can't get Bills list! %w", err)
		return
	}
	fmt.Printf("%s\t%s\t%s\n", "id", "Баланс", "Заблокирован")
	for _, value := range list {
		fmt.Printf("%v\t%v\t%v\n", value.Id, value.Balance, value.Locked)
	}
}

func handleTransfer(db *sql.DB, user_id int) {
	for ; ; {
		var cmd string
		fmt.Println("\nДоступные вам команды...", transferCommands)
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Errorf("can't scan command %w", err)
			continue
		}
		switch cmd {
		case "1":
			fmt.Println("по номеру счёта...")
			handleTransferByBill(db, user_id)
		case "2":
			fmt.Println("по номеру телефона...")
			handleTransferByPhone(db, user_id)
		case "q":
			return
		default:
			fmt.Printf("Неправильная команда %s\n", cmd)
		}
	}
}

func handleTransferByBill(db *sql.DB, user_id int) {
	var addressee_id int
	fmt.Print("Введите счет пользователя которому хотите осуществить перевод:\nid: ")
	_, err := fmt.Scan(&addressee_id)
	if err != nil {
		fmt.Printf("can't scan addressee id %v\n", err)
		return
	}
	ok, err, addressee_balance := core.CheckBill(db, addressee_id)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	if !ok {
		fmt.Printf("счет id: %v не существует!", addressee_id)
		return
	}
	var amount int
	fmt.Print("Введите сумму перевода:\nсумма: ")
	_, err = fmt.Scan(&amount)
	if err != nil {
		fmt.Printf("can't scan amount %v\n", err)
		return
	}
	if amount < 1 {
		fmt.Printf("Сумма перевода долбжна быть больше 0 %v\n", err)
		return
	}
	fmt.Println("доступные вам счета с которых можно осуществить перевод:")
	bills, err := core.AvailableBills(db, user_id, amount)
	if err != nil {
		fmt.Printf("Произошла ошибка!!! %v", err)
		return
	}
	fmt.Printf("%s\t%s\n", "id", "Баланс")
	for _, bill := range bills {
		fmt.Printf("%v\t%v\n", bill.Id, bill.Balance)
	}
	var chosed_id int
	fmt.Println("Введите id счета с которого перевести:")
	_, err = fmt.Scan(&chosed_id)
	if err != nil {
		fmt.Printf("can't scan bill id %v", err)
	}
	for _, value := range bills {
		if value.Id == chosed_id {
			err = core.TransferBillToBill(db, value.Id, value.Balance, addressee_id, addressee_balance, amount)
			if err != nil {
				fmt.Printf("не удалось осуществить перевод %v", err)
				return
			}
			fmt.Println("Перевод выполнен.")
			return
		}
	}
	fmt.Printf("Нет такого счета в списке! вы ввели id: %v", chosed_id)
}

func handleTransferByPhone(db *sql.DB, user_id int) {
	var addressee_phone string
	fmt.Print("Введите номер телефона пользователя которому хотите осуществить перевод:\nPhone Number: ")
	_, err := fmt.Scan(&addressee_phone)
	if err != nil {
		fmt.Printf("can't scan phone %v\n", err)
		return
	}
	var amount int
	fmt.Print("Введите сумму перевода:\nсумма: ")
	_, err = fmt.Scan(&amount)
	if err != nil {
		fmt.Printf("can't scan amount %v\n", err)
		return
	}
	if amount < 1 {
		fmt.Printf("Сумма перевода долбжна быть больше 0 %v\n", err)
		return
	}
	addressee_id, addressee_balance, err := core.GetAnyBill(db, addressee_phone, amount)

	fmt.Println("доступные вам счета с которых можно осуществить перевод:")
	bills, err := core.AvailableBills(db, user_id, amount)
	if err != nil {
		fmt.Printf("Произошла ошибка!!! %v", err)
		return
	}
	fmt.Printf("%s\t%s\n", "id", "Баланс")
	for _, bill := range bills {
		fmt.Printf("%v\t%v\n", bill.Id, bill.Balance)
	}
	var chosed_id int
	fmt.Println("Введите id счета с которого перевести:")
	_, err = fmt.Scan(&chosed_id)
	if err != nil {
		fmt.Printf("can't scan bill id %v", err)
	}
	for _, value := range bills {
		if value.Id == chosed_id {
			err = core.TransferBillToBill(db, value.Id, value.Balance, addressee_id, addressee_balance, amount)
			if err != nil {
				fmt.Printf("не удалось осуществить перевод %v", err)
				return
			}
			fmt.Println("Перевод выполнен.")
			return
		}
	}
	fmt.Printf("Нет такого счета в списке! вы ввели id: %v", chosed_id)
}

func handlePayService(db *sql.DB, user_id int) {
	var service_id int
	fmt.Print("Введите ID услуги: ")
	_, err := fmt.Scan(&service_id)
	if err != nil {
		fmt.Errorf("can't scan service_id %w", err)
		return
	}
	err = core.PayService(db, service_id, user_id)
	if err != nil {
		fmt.Printf("Не удалось оплатить услугу! %v", err)
		return
	}
	fmt.Printf("Услуга номер: %v была оплачена.\n", service_id)
}

func handleGetATMsList(db *sql.DB) {
	list, err := core.ATMsList(db)
	if err != nil {
		fmt.Errorf("can't get ATMs list! %w", err)
		return
	}
	fmt.Printf("%s\t%s\t%s\n", "id", "Адрес", "Заблокирован")
	for _, value := range list {
		fmt.Printf("%v\t%s\t%v\n", value.Id, value.Address, value.Locked)
	}
}

func handleGetServicesList(db *sql.DB) {
	list, err := core.ServicesList(db)
	if err != nil {
		fmt.Errorf("can't get Services list! %w", err)
		return
	}
	fmt.Printf("%s\t%s\t%s\n", "id", "Наименование", "Цена")
	for _, value := range list {
		fmt.Printf("%v\t%s\t%v\n", value.Id, value.Name, value.Price)
	}
}
