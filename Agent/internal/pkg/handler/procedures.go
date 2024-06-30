package handler

import "fmt"

func Heal_now() {
	fmt.Println("Heal now done")
}

func Get_Status() {
	fmt.Println("PROCEDURE MESSAGE: STATUS HAS BEEN SENT")

}

func Restore_Backup(done chan bool) {
	<-done
	fmt.Println("restore backup done")
	done <- true
}
