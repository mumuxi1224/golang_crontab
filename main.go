package main

func main() {
	crontab := getCrontab()
	crontab.Run()
}
