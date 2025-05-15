package handler

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
)

type handler struct {
	driver selenium.WebDriver
}

func New(drv selenium.WebDriver) *handler {
	return &handler{
		driver: drv,
	}
}

func (h *handler) Authorize(ctx context.Context) error {
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	inputFields, err := h.driver.FindElements(selenium.ByXPATH, "/html/body/ion-nav-view/ion-view/ion-content/div[1]/form/ion-list/div/label")
	if err != nil {
		return errors.Wrap(err, "driver.FindElements")
	}

	err = inputFields[0].SendKeys(email)
	if err != nil {
		return errors.Wrap(err, "login: SendKeys")
	}

	err = inputFields[1].SendKeys(password)
	if err != nil {
		return errors.Wrap(err, "password: SendKeys")
	}

	buttons, err := h.driver.FindElements(selenium.ByXPATH, "/html/body/ion-nav-view/ion-view/ion-content/div[1]/form/button")
	if err != nil {
		return errors.Wrap(err, "driver.FindElements")
	}

	for _, btn := range buttons {
		btnText, err := btn.Text()
		if err != nil {
			return errors.Wrap(err, "btn.Text")
		}

		if btnText == "Войти" {
			err = btn.Click()
			if err != nil {
				return errors.Wrap(err, "btn.Click")
			}

			return nil
		}
	}

	return nil
}

func (h *handler) CatchCron(ctx context.Context) error {
	time.Sleep(time.Second * 2)
	buttons, err := h.driver.FindElements(selenium.ByXPATH, "/html/body//button")
	if err != nil {
		return errors.Wrap(err, "driver.FindElements")
	}

	var buttonBook selenium.WebElement

	for _, btn := range buttons {
		btnText, err := btn.Text()
		if err != nil {
			return errors.Wrap(err, "btn.Text")
		}

		if btnText == "Записаться на занятие" || btnText == "Записаться в очередь" {
			buttonBook = btn
			break
		}
	}

	classNames, err := buttonBook.GetAttribute("class")
	if err != nil {
		return errors.Wrap(err, "buttonBook.GetAttribute")
	}

	buttonBookText, err := buttonBook.Text()
	if err != nil {
		return errors.Wrap(err, "buttonBook.Text")
	}

	if hasClassName(classNames, "disabled") && buttonBookText != "Записаться в очередь" {
		fmt.Println("Вы уже записаны на это занятие")

		h.driver.Close()

		os.Exit(1)
	}

	if buttonBookText == "Записаться на занятие" {
		err = buttonBook.Click()
		if err != nil {
			return errors.Wrap(err, "buttonBook.Click")
		}

		fmt.Println("Вы успешно записаны на занятие!")

		h.driver.Close()

		os.Exit(1)
	}

	fmt.Println("Свободных мест пока нет, продолжаем...")
	h.driver.Refresh()

	return nil
}

func hasClassName(classname string, searched string) bool {
	namesSlice := strings.Split(classname, " ")

	for _, class := range namesSlice {
		if class == searched {
			return true
		}
	}

	return false
}
