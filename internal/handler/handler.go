package handler

import (
	"context"
	"fmt"
	"os"
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
	err := h.driver.Get(os.Getenv("CLASS_URL"))
	if err != nil {
		return errors.Wrap(err, "driver.Get")
	}

	time.Sleep(time.Second * 2)

	buttons, err := h.driver.FindElements(selenium.ByXPATH, "/html/body//button")
	if err != nil {
		return errors.Wrap(err, "driver.FindElements")
	}

	buttonBook, err := findBookButton(buttons)
	if err != nil {
		return errors.Wrap(err, "findBookButton")
	}

	buttonBookText, err := buttonBook.Text()
	if err != nil {
		return errors.Wrap(err, "buttonBook.Text")
	}

	switch buttonBookText {
	case "Перезаписаться на другое занятие":
		fmt.Println("Вы уже записаны на это занятие")

		err = h.driver.Close()
		if err != nil {
			return errors.Wrap(err, "driver.Close")
		}

		os.Exit(1)
	case "Записаться на занятие":
		err = buttonBook.Click()
		if err != nil {
			return errors.Wrap(err, "buttonBook.Click")
		}

		err = h.driver.Refresh()
		if err != nil {
			return errors.Wrap(err, "driver.Refresh")
		}

		time.Sleep(time.Second * 2)

		buttons, err = h.driver.FindElements(selenium.ByXPATH, "/html/body//button")
		if err != nil {
			return errors.Wrap(err, "driver.FindElements")
		}

		buttonBook, err = findBookButton(buttons)
		if err != nil {
			return errors.Wrap(err, "findBookButton")
		}

		buttonBookText, err = buttonBook.Text()
		if err != nil {
			return errors.Wrap(err, "buttonBook.Text")
		}

		if buttonBookText != "Перезаписаться на другое занятие" {
			fmt.Println("Свободных мест пока нет, продолжаем...")

			err = h.driver.Refresh()
			if err != nil {
				return errors.Wrap(err, "driver.Refresh")
			}

			return nil
		}

		fmt.Println("Вы успешно записаны на занятие!")

		err = h.driver.Close()
		if err != nil {
			return errors.Wrap(err, "driver.Close")
		}

		os.Exit(1)
	default:
		fmt.Println("Свободных мест пока нет, продолжаем...")
		err = h.driver.Refresh()
		if err != nil {
			return errors.Wrap(err, "driver.Refresh")
		}
	}

	return nil
}

func findBookButton(buttons []selenium.WebElement) (selenium.WebElement, error) {
	for _, btn := range buttons {
		btnText, err := btn.Text()
		if err != nil {
			return nil, errors.Wrap(err, "btn.Text")
		}

		if btnText == "Записаться на занятие" {
			return btn, nil
		}

		if btnText == "Записаться в очередь" {
			return btn, nil
		}

		if btnText == "Перезаписаться на другое занятие" {
			return btn, nil
		}
	}

	return nil, nil
}
