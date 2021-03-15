package main

import (
	"bytes"
	"errors"
	"log"
	"net"
	"os"
	"regexp"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/gen2brain/dlgs"
	winapi "github.com/iamacarpet/go-win64api"
	ps "github.com/mitchellh/go-ps"
)

// TODO: Отказаться от Concatinate
// TODO: отказаться от dlgs
// TODO: Кнопка убийства файрвола

// Точка входа
func main() {

	ip, err := externalIP()
	if err != nil {
		err = nil
	}
	ip = "Адрес: " + ip

	hostname, err := os.Hostname()
	if err != nil {
		err = nil
	}
	hostname = "Имя комп: " + hostname

	UI(ip, hostname)
}

// Получение IP Адреса компьютера
func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("Нет сети!")
}

// Вспомогательная функция конкатинации строк
func Concatinate(oneword string, twoword string) string {
	strings := []string{oneword, twoword}
	buffer := bytes.Buffer{}
	for _, val := range strings {
		buffer.WriteString(val)
	}

	return buffer.String()
}

// очистка кеша 1с
func clearCash() {

	ScriptMayWork := true

	processList, err := ps.Processes()
	if err != nil {
		log.Println("ps.Processes() Failed, are you using windows?")
		return
	}

	for x := range processList {
		var process ps.Process
		process = processList[x]
		if process.Executable() == "1cv8c.exe" {
			dlgs.Warning("Ошибка!", "Обнаружен процесс 1С:Предприятия! Закройте 1С и повторите попытку.")
			ScriptMayWork = false
		}
	}

	if ScriptMayWork == true {

		UserHomeDir, err := os.UserHomeDir()
		if err != nil {
			err = nil
		}

		RoamingOneCDirPath := Concatinate(UserHomeDir, "/AppData/Roaming/1C/1cv8/")
		LocalOneCDirPath := Concatinate(UserHomeDir, "/AppData/Local/1C/1cv8/")

		filesInRoaming, err := os.ReadDir(RoamingOneCDirPath)
		if err != nil {
			err = nil
		}

		filesInLocal, err := os.ReadDir(LocalOneCDirPath)
		if err != nil {
			err = nil
		}

		for _, fileInRoaming := range filesInRoaming {
			matched, err := regexp.MatchString(`[\-]`, fileInRoaming.Name())
			if err != nil {
				err = nil
			}
			if matched == true {
				os.RemoveAll(Concatinate(RoamingOneCDirPath, fileInRoaming.Name()))
			}
		}

		for _, fileInLocal := range filesInLocal {
			matched, err := regexp.MatchString(`[\-]`, fileInLocal.Name())
			if err != nil {
				err = nil
			}
			if matched == true {
				os.RemoveAll(Concatinate(LocalOneCDirPath, fileInLocal.Name()))
			}
		}

		dlgs.Info("Очистка кеша", "Кэш 1С очищен успешно!")
	}
}

// Настройка Удаленки
func RemoteConnection() {

	winapi.FirewallDisable(1)
	winapi.FirewallDisable(2)
	winapi.FirewallDisable(3)
	winapi.FirewallDisable(4)

	dlgs.Info("Remote connection", "Настройка прошла успешно!")
}

// Отрисовка интерфейса
func UI(ip string, hostname string) {
	a := app.New()
	w := a.NewWindow("Admin Tool")
	w.Resize(fyne.NewSize(300, 100))
	w.CenterOnScreen()

	ipAddr := widget.NewLabel(ip)
	CompName := widget.NewLabel(hostname)
	w.SetContent(widget.NewVBox(
		ipAddr,
		CompName,
		widget.NewButton("Очистить кеш 1С", clearCash),
		widget.NewButton("Удаленный доступ", RemoteConnection),
	))

	w.ShowAndRun()
}
