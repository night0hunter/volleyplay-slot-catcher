package handler

import (
	"context"
	"fmt"
	"os"
	"strconv"

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

	inputFields, err := h.driver.FindElements(selenium.ByClassName, "ng-empty")
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

	buttons, err := h.driver.FindElements(selenium.ByClassName, "button")
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
	elems, err := h.driver.FindElements(selenium.ByClassName, "ng-binding")
	if err != nil {
		return errors.Wrap(err, "driver.FindElements")
	}

	buttonBook, err := h.driver.FindElement(selenium.ByClassName, "button-block")
	if err != nil {
		return errors.Wrap(err, "driver.FindElements")
	}

	var val int

	for _, elem := range elems {
		text, err := elem.Text()
		if err != nil {
			return errors.Wrap(err, "elem.Text")
		}

		val, err = strconv.Atoi(text)
		if err != nil {
			continue
		}

		break
	}

	if val != 0 {
		err = buttonBook.Click()
		if err != nil {
			return errors.Wrap(err, "buttonBook.Click")
		}

		fmt.Println("Вы успешно записаны на занятие!")
	}

	h.driver.Refresh()

	return nil
}
