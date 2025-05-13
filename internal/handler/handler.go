package handler

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
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
	seconds := rand.Intn(10)
	time.Sleep(time.Duration(seconds))
	elems, err := h.driver.FindElements(selenium.ByClassName, "ng-binding")
	if err != nil {
		return errors.Wrap(err, "driver.FindElements")
	}

	buttons, err := h.driver.FindElements(selenium.ByClassName, "button-block")
	if err != nil {
		return errors.Wrap(err, "driver.FindElements")
	}

	var buttonBook selenium.WebElement

	for _, btn := range buttons {
		btnText, err := btn.Text()
		if err != nil {
			return errors.Wrap(err, "btn.Text")
		}

		if btnText == "Записаться на занятие" {
			buttonBook = btn

			break
		}

		if btnText == "Записаться в очередь" {
			fmt.Println("Свободных мест пока нет, продолжаем...")

			return nil
		}
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

	classNames, err := buttonBook.GetAttribute("class")
	if err != nil {
		return errors.Wrap(err, "buttonBook.GetAttribute")
	}

	if hasClassName(classNames, "disabled") {
		fmt.Println("Вы уже записаны на это занятие")

		h.driver.Close()

		os.Exit(1)
	}

	if val != 0 {
		err = buttonBook.Click()
		if err != nil {
			return errors.Wrap(err, "buttonBook.Click")
		}

		fmt.Println("Вы успешно записаны на занятие!")

		h.driver.Close()

		os.Exit(1)
	}

	if val == 0 {
		fmt.Println("Свободных мест пока нет, продолжаем...")
	}

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
