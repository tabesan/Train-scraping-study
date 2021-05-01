package train

import (
	"fmt"
	"time"
	ln "train_delay/notify"

	"github.com/PuerkitoBio/goquery"
)

type trainStatus struct {
	//Target page's url
	url       string
	preStatus string
	newStatus string
	notify    *ln.Notify
	normally  string
}

// Create new trainStatus
func NewTrain(n *ln.Notify) *trainStatus {
	train := new(trainStatus)

	train.url = "https://subway.osakametro.co.jp/guide/subway_information.php"
	train.preStatus = ""
	train.newStatus = ""
	train.notify = n
	train.normally = "現在、10分以上の列車の遅れはございません。"

	return train
}

// Whether or not the status has changed
func (m *trainStatus) checkAlt() bool {
	return m.preStatus != m.newStatus
}

// Check normally operation
func (m *trainStatus) checkNormally(status string) bool {
	return status == m.normally
}

// Filter out lines isn't contained int the MyLine.myLine
//func (m *MyLine) contains(l string) bool {
//	for _, v := range m.myLine {
//		if l == v {
//			return true
//		}
//	}
//	return false
//}

// Create new status message
func (m *trainStatus) createMessage(l string, s string) {
	if m.checkNormally(s) {
		m.newStatus += "\n" + "【" + l + "】" + "\n" + "通常運行"
	} else {
		m.newStatus += "\n" + "【" + l + "】" + "\n" + "遅延または運休"
	}
}

// Get operational status
func (t *trainStatus) getStatus() {
	doc, err := goquery.NewDocument(t.url)
	if err != nil {
		fmt.Println(err.Error())
	}

	t.newStatus = "\n"
	mainContent := doc.Find("div.mainContents > table > tbody").Children()
	mainContent.Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			s.Each(func(j int, ss *goquery.Selection) {
				line := s.Find("td.cs-tdLine").Text()
				status := s.Find("td.cs-tdTxt").Text()
				t.createMessage(line, status)
			})
		}
	})
}

// Scraping at 1-minute intervals
func (t *trainStatus) DoScrape() {
	for range time.Tick(60 * time.Second) {
		t.getStatus()
		if t.checkAlt() {
			t.preStatus = t.newStatus
			t.notify.SendNotify(t.newStatus)
		}
	}
}
