package datamining

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	parser, err := NewParser()

	if err != nil {
		panic(err)
	}

	body, err := parser.RequestURL("https://www.bileter.ru/afisha/building/bolshoy_kontsertnyiy_zal_oktyabrskiy.html")

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("div.afishe-item").Each(func(i int, s *goquery.Selection) {
		infoBlock := s.Find("div.info-block")
		divName := infoBlock.Find("div.name").First()
		title := divName.Find("a").Text()
		fmt.Println("Название:", title)
		link, exist := divName.Find("a").Attr("href")
		if exist {
			fmt.Println("Ссылка:", fmt.Sprintf("https://www.bileter.ru/%s", link))
		}
		date := infoBlock.Find("div.date").Text()
		newDate := strings.ReplaceAll(date, " ", "")
		fmt.Println("Дата:", newDate)

		result := Event{
			Name:  title,
			URL:   link,
			Date:  newDate,
			Venue: "БКЗ \"Октябрьский\"",
		}

		if err := parser.SendToKafka(result); err != nil {
			log.Fatal(err)
		}
	})
}
