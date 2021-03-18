package main

import (
	"bytes"
	"errors"
	"net"
	"os"
	"os/exec"
	"regexp"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"

	winapi "github.com/iamacarpet/go-win64api"
	ps "github.com/mitchellh/go-ps"
)

// TODO: пинговалка сервера, шлюза, Л44, для более быстрого определения неисправностей
// TODO: Возможно - логгирование запусков Утилиты
// TODO: точка входа для какого нибудь powershell скрипта
// TODO: прокоментировать код

// создание глобального экзепляра Fyne app
var a = app.New()

// Точка входа
func main() {

	ip, err := externalIP()
	if err != nil {
		dlgWindow("Error!", err.Error(), "Ok")
	}
	ip = "Адрес: " + ip

	hostname, err := os.Hostname()
	if err != nil {
		dlgWindow("Error!", err.Error(), "Ok")
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

// очистка кеша 1с
func clearCash() {

	ScriptMayWork := true

	processList, err := ps.Processes()
	if err != nil {
		dlgWindow("Error!", err.Error(), "Ok")
	}

	for x := range processList {
		var process ps.Process
		process = processList[x]
		if process.Executable() == "1cv8c.exe" {
			dlgWindow("Ошибка!", "Обнаружен запущенный процесс 1С:Предприятия! Закройте 1С и повторите попытку.", "Ок")
			ScriptMayWork = false
		}
	}

	if ScriptMayWork == true {

		UserHomeDir, err := os.UserHomeDir()
		if err != nil {
			dlgWindow("Error!", err.Error(), "Ok")
		}

		RoamingOneCDirPath := UserHomeDir + "/AppData/Roaming/1C/1cv8/"
		LocalOneCDirPath := UserHomeDir + "/AppData/Local/1C/1cv8/"

		filesInRoaming, err := os.ReadDir(RoamingOneCDirPath)
		if err != nil {
			dlgWindow("Error!", err.Error(), "Ok")
		}

		filesInLocal, err := os.ReadDir(LocalOneCDirPath)
		if err != nil {
			dlgWindow("Error!", err.Error(), "Ok")
		}

		for _, fileInRoaming := range filesInRoaming {
			matched, err := regexp.MatchString(`[\-]`, fileInRoaming.Name())
			if err != nil {
				dlgWindow("Error!", err.Error(), "Ok")
			}
			if matched == true {
				os.RemoveAll(RoamingOneCDirPath + fileInRoaming.Name())
			}
		}

		for _, fileInLocal := range filesInLocal {
			matched, err := regexp.MatchString(`[\-]`, fileInLocal.Name())
			if err != nil {
				dlgWindow("Error!", err.Error(), "Ok")
			}
			if matched == true {
				os.RemoveAll(LocalOneCDirPath + fileInLocal.Name())
			}
		}

		dlgWindow("Очистка кеша", "Кэш 1С очищен успешно!", "Ок")
	}
}

// Настройка Удаленки
func RemoteConnection() {

	winapi.FirewallDisable(1)
	winapi.FirewallDisable(2)
	winapi.FirewallDisable(3)
	winapi.FirewallDisable(4)

	dlgWindow("Remote connection", "Настройка прошла успешно!", "Ок")
}

// Отрисовка интерфейса
func UI(ip string, hostname string) {

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

// Вывод информационного диалогового окна
func dlgWindow(title string, content string, btnText string) {

	windw := a.NewWindow(title)

	label := widget.NewLabel(content)
	windw.SetContent(widget.NewVBox(
		label,
		widget.NewButton(btnText, func() {
			windw.Close()
		}),
	))

	windw.Show()
}

// Запуск скрипта powershell по кнопке, с указанием пути к скрипту. Проверить, работает или нет вообще.
func runPowershellFileScript(path string) string {
	ps, _ := exec.LookPath("powershell.exe")
	psArg := "-File" + path

	command := exec.Command(ps, psArg)

	var stdout bytes.Buffer
	command.Stdout = &stdout

	command.Run()

	return stdout.String()
}
