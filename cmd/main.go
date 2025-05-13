package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"volleyplaySlotCatcher/internal/cronjob"
	"volleyplaySlotCatcher/internal/handler"
	"volleyplaySlotCatcher/internal/utils"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("godotenv.Load: %s", err)
	}

	authorizationURL := os.Getenv("LOGIN_URL")
	classURL := os.Getenv("CLASS_URL")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// chromeService, err := selenium.NewChromeDriverService("./chromedriver", 4444)
	// if err != nil {
	// 	fmt.Println(errors.Wrap(err, "selenium.NewChromeDriverService"))
	// }
	// defer chromeService.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"--headless-new", // comment out this line for testing
	}})

	// create a new remote client with the specified options
	driver, err := selenium.NewRemote(caps, "http://172.17.0.2:4444/wd/hub")
	if err != nil {
		fmt.Println(errors.Wrap(err, "selenium.NewRemote"))
	}

	// maximize the current window to avoid responsive rendering
	err = driver.MaximizeWindow("")
	if err != nil {
		fmt.Println(errors.Wrap(err, "MaximizeWindow"))
	}

	handler := handler.New(driver)

	err = driver.Get(authorizationURL)
	if err != nil {
		fmt.Println(errors.Wrap(err, "driver.Get"))
	}

	err = handler.Authorize(ctx)
	if err != nil {
		fmt.Println(errors.Wrap(err, "handler.Authorize"))
	}

	time.Sleep(time.Second * 5)
	// err = driver.Get("http://localhost:3333/")
	err = driver.Get(classURL)
	if err != nil {
		fmt.Println(errors.Wrap(err, "driver.Get"))
	}

	cr := cronjob.NewCatchSlotCron(handler)

	c := cron.New(
		cron.WithLocation(utils.MoscowLocation),
		cron.WithParser(cron.NewParser(cron.Second|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow)),
	)

	_, err = c.AddJob("*/15 * * * * *", cr)
	if err != nil {
		fmt.Println(errors.Wrap(err, "c.AddJob"))
	}

	c.Start()

	fmt.Println("Press Ctrl+C to exit...")

	<-ctx.Done()

	fmt.Println("\nShutdown signal received. Exiting...")
}
