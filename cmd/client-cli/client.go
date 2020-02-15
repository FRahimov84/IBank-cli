package main

import (
	"database/sql"
	"fmt"
	core "github.com/FRahimov84/IBank-core"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("can't start application client %e", err)
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
			fmt.Println("ЧТо-то пошло не так")
			log.Printf("can't scan command %e", err)
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
		fmt.Println("Что-то аошло не так")
		log.Printf("can't scan login in handleLogin() err: %e", err)
	}
	fmt.Println("Пароль: ")
	_, err = fmt.Scan(&pass)
	if err != nil {
		fmt.Println("Что-то аошло не так")
		log.Printf("can't scan password in handleLogin() err: %e", err)
	}
	user_id, name, err := core.Login(db, login, pass)
	if err != nil {
		fmt.Println("Не удалось войти в систему! неправильный логин или пароль")
		log.Printf("can't login, in handleLogin() err: %e",err)
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
			fmt.Println("Не удалось прочесть команду")
			log.Printf("can't scan command in operationLoopAuthorized() error: %e", err)
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
		fmt.Println("Не удалось получить счета пользователя")
		log.Printf("can't get Bills list! in handleGetBillsOfUser() error: %e", err)
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
			fmt.Println("Что-то пошло не так")
			log.Printf("can't scan command in handleTransfer error: %e", err)
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
		fmt.Println("Что-то пошло не так")
		log.Printf("can't scan addressee id in handleTransferByBill() err: %e\n", err)
		return
	}
	ok, err, addressee_balance := core.CheckBill(db, addressee_id)
	if err != nil {
		fmt.Println("Этот счет заблокирован/несуществует")
		log.Printf("can't checkbill in handleTransferByBill err: %e", err)
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
		fmt.Println("Что-то пошло не так")
		log.Printf("can't scan amount in handleTransferByBill() err: %e\n", err)
		return
	}
	if amount < 1 {
		fmt.Println("Сумма перевода долбжна быть больше 0")
		return
	}
	fmt.Println("доступные вам счета с которых можно осуществить перевод:")
	bills, err := core.AvailableBills(db, user_id, amount)
	if err != nil {
		fmt.Println("Что-то пошло не так при получении доступных счетов")
		log.Printf("can't get bills in handleTransferByBill() err: %e", err)
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
		fmt.Println("Что-то пошло не так")
		log.Printf("can't scan bill id in handleTransferByBill() err: %e", err)
		return
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
		log.Printf("can't scan phone in handleTransferByPhone() err: %e\n", err)
		return
	}
	var amount int
	fmt.Print("Введите сумму перевода:\nсумма: ")
	_, err = fmt.Scan(&amount)
	if err != nil {
		log.Printf("can't scan amount in handleTransferByPhone() err: %e\n", err)
		fmt.Println("Что-то пошло не так")
		return
	}
	if amount < 1 {
		fmt.Println("Сумма перевода долбжна быть больше 0 ")
		return
	}
	addressee_id, addressee_balance, err := core.GetAnyBill(db, addressee_phone, amount)

	fmt.Println("доступные вам счета с которых можно осуществить перевод:")
	bills, err := core.AvailableBills(db, user_id, amount)
	if err != nil {
		fmt.Println("Произошла ошибка!!!")
		log.Printf("can't get available bills in  handleTransferByPhone() err: %e\n",err)
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
		fmt.Println("Что-то пошло не так")
		log.Printf("can't scan bill id in  handleTransferByPhone() err: %e", err)
	}
	for _, value := range bills {
		if value.Id == chosed_id {
			err = core.TransferBillToBill(db, value.Id, value.Balance, addressee_id, addressee_balance, amount)
			if err != nil {
				fmt.Println("не удалось осуществить перевод")
				log.Printf("can't transfer money in  handleTransferByPhone() err: %e\n", err)
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
		fmt.Println("Что-то пошло не так")
		log.Printf("can't scan service_id in handlePayService() err: %e\n", err)
		return
	}
	err = core.PayService(db, service_id, user_id)
	if err != nil {
		fmt.Println("Не удалось оплатить услугу!")
		log.Printf("can't pay for service in handlePayService() err: %e\n", err)
		return
	}
	fmt.Printf("Услуга номер: %v была оплачена.\n", service_id)
}

func handleGetATMsList(db *sql.DB) {
	list, err := core.ATMsList(db)
	if err != nil {
		fmt.Println("Не удалось получить список банкоматов")
		log.Printf("can't get ATMs list! in handleGetATMsList() err: %e\n", err)
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
		fmt.Println("Что-то пошло не так")
		log.Printf("can't get Services list! in handleGetServicesList() err: %e", err)
		return
	}
	fmt.Printf("%s\t%s\t%s\n", "id", "Наименование", "Цена")
	for _, value := range list {
		fmt.Printf("%v\t%s\t%v\n", value.Id, value.Name, value.Price)
	}
}
