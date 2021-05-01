package main

import (
	"train_delay/notify"
	"train_delay/train"
)

func main() {
	n := notify.NewNotify()
	ml := train.NewTrain(n)
	ml.DoScrape()
}
