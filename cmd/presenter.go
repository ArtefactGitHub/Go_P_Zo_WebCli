package main

import (
	"fmt"
	"time"
)

type Presenter struct{}

var gPresenter *Presenter = &Presenter{}

func (p *Presenter) TimeToStr(time time.Time) string {
	return fmt.Sprintf("%d年%d月%d日", time.Year(), time.Month(), time.Day())
}
