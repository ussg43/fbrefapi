package scraper

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	// "strings"
	"github.com/gocolly/colly"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Player struct {
	Name     string
	Position string
	Club     string
	URL      string
}

func NewPlayer(name, club string) (*Player, error) {
	url, err := GetUrl(name, club)
	if err != nil {
		return nil, err
	}
	position, err := getPosition(url)
	if err != nil {
		return nil, err
	}
	p := Player{Name: name, Position: position, Club: club, URL: url}
	fmt.Printf("new player created: %s", name)
	return &p, nil
}

func GetUrl(name, club string) (string, error) {
	var url string
	var err error
	scraped := false
	name = strings.ReplaceAll(name, " ", "+")
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.AllowedDomains("https://fbref.com", "fbref.com"),
	)

	target := fmt.Sprintf("https://fbref.com/en/search/search.fcgi?hint=%s&search=%s&pid=&idx=", name, name)

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println(r.Request.URL)
		if RUrl := r.Request.URL.String(); RUrl != target {
			log.Println(RUrl)
			url = RUrl
			scraped = true
		} else if club != "" {
			c.OnHTML("div#players.current", func(h *colly.HTMLElement) {
				h.ForEach("div.search-item", func(i int, e *colly.HTMLElement) {
					if fuzzy.MatchFold(club, e.ChildText("div.search-item-team")) {
						url = "https://fbref.com" + e.ChildAttr("div.search-item-url > a", "href")
						scraped = true
					}
				})
			})
		}
	})
	c.Visit(target)
	log.Println("new url: " + url)
	if !scraped{
		err = ErrFetchURL
		return "", err
	}
	return url, nil
}



// func Compare(p []*Player) string {
// 	map1 := p[0].GetP90()
// 	map2 := p[1].GetP90()
// 	msg := ""

// 	for key := range map1 {
// 		if map1[key] > map2[key] {
// 			msg += fmt.Sprintf("%s: **%.2f** - %.2f\n", key, map1[key], map2[key])
// 		} else {
// 			msg += fmt.Sprintf("%s: %.2f - **%.2f**\n", key, map1[key], map2[key])
// 		}
// 	}
// 	return msg
// }

func getPosition(Url string) (string, error) {
	var post string
	var err error
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("retrieving position at: ", r.URL.String())
	})

	c.OnHTML("div#all_similar", func(h *colly.HTMLElement) {
		post = h.ChildText("div.current > a.sr_preset")
		log.Println(post)
		switch post {
		case "Goalkeepers":
			post = "GK"
		case "Att Mid / Wingers":
			post = "AM"
		case "Center Backs":
			post = "CB"
		case "Fullbacks":
			post = "FB"
		case "Forwards":
			post = "FW"
		case "Midfielders":
			post = "MF"
		default:
			post = ""
		}
	})
	log.Println(post)

	if post == ""{
		err = ErrPositionNotFound
		return "", err
	}
	c.Visit(Url)
	return post, nil
}

func (p *Player) GetP90() (map[string]float64, error) {
	var err error
	tableID := fmt.Sprintf(`table[id=scout_summary_%s]`, p.Position)
	m := make(map[string]float64)

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("retrieving p90 stats at:", r.URL.String())
	})

	c.OnHTML(tableID, func(h *colly.HTMLElement) {

		h.ForEach("tr", func(i int, h *colly.HTMLElement) {
			if key := h.ChildText("th"); !reflect.DeepEqual(key, "StatisticPer 90Percentile") {
				// s := strings.ReplaceAll(h.ChildText("td.right"),"%","")
				s := h.ChildAttr("td.right", "csk")
				val, _ := strconv.ParseFloat(s, 64)
				if key != "" {
					m[key] = val
				}
			}
		})
		if len(m) == 0{
			err = ErrScrapeFailed
		}
	})

	if err != nil{
		return nil, err
	}
	c.Visit(p.URL)
	return m, nil
}

func (p *Player) GetSeasonal(season string) (map[string]float64, error) {
	url := p.URL[0:38] + "all_comps" + p.URL[37:] + "-Stats---All-Competitions"
	log.Println(url)
	tableID := `table[id=stats_standard_collapsed]`
	m := make(map[string]float64)
	var err error

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.AllowedDomains("https://fbref.com", "http://fbref.com", "fbref.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("retrieving seasonal stats at: ", r.URL.String())
	})

	c.OnHTML(tableID, func(h *colly.HTMLElement) {
		h.ForEach("tr", func(i int, h *colly.HTMLElement) {
			if h.ChildText("th.left") == season {
				h.ForEach("td", func(i int, e *colly.HTMLElement) {
					key := e.Attr("data-stat")
					val, _ := strconv.ParseFloat(e.Text, 64)
					if val != 0 {
						m[key] = val
					}
				})
			}
		})
		if len(m) == 0{
			err = ErrScrapeFailed
		}
	})
	if err != nil{
		return nil, err
	}
	c.Visit(url)
	return m, nil
}
