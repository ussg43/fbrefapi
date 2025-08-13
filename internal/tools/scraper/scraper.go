package scraper

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	// "strings"
	"github.com/gocolly/colly"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Player struct {
	Name     string
	Position string
	Club     string
	URL   string
	Stats    map[string]float64
}

func GetUrl(name, club string) string {
	var url string
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.AllowedDomains("https://fbref.com", "fbref.com"),
	)

	target := fmt.Sprintf("https://fbref.com/en/search/search.fcgi?hint=%s&search=%s&pid=&idx=", name, name)

	c.OnResponse(func(r *colly.Response) {
		log.Println(r.Request.URL)
		if RUrl := r.Request.URL.String(); RUrl != target {
			log.Println(RUrl)
			url = RUrl
		} else if club != "" {
			c.OnHTML("div#players.current", func(h *colly.HTMLElement) {
				h.ForEach("div.search-item", func(i int, e *colly.HTMLElement) {
					if fuzzy.MatchFold(club, e.ChildText("div.search-item-team")) {
						url = "https://fbref.com" + e.ChildAttr("div.search-item-url > a", "href")
					}
				})
			})
		}
	})

	

	c.Visit(target)
	log.Println(url)
	return url
}

func (P *Player) GetName() string {
	return P.Name
}

func NewPlayer(name, club string) *Player {
	url := GetUrl(name, club)
	position := getPosition(url)
	p := Player{Name: name, Position: position, Club: club, URL: url}

	return &p
}

func Compare(p []*Player) string {
	map1 := p[0].GetP90()
	map2 := p[1].GetP90()
	msg := ""

	for key := range map1 {
		if map1[key] > map2[key] {
			msg += fmt.Sprintf("%s: **%.2f** - %.2f\n", key, map1[key], map2[key])
		} else {
			msg += fmt.Sprintf("%s: %.2f - **%.2f**\n", key, map1[key], map2[key])
		}
	}
	return msg
}

func getPosition(Url string) string {
	var post string
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.AllowedDomains("https://fbref.com", "http://fbref.com", "fbref.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
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

	c.Visit(Url)
	return post
}

func (p *Player) GetP90() map[string]float64 {
	tableID := fmt.Sprintf(`table[id=scout_summary_%s]`, p.Position)
	m := make(map[string]float64)

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
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
	})

	c.Visit(p.URL)
	return m
}

func (p *Player) GetSeasonal(season string) map[string]float64 {
	url := p.URL[0:38] + "all_comps"+p.URL[37:]+"-Stats---All-Competitions"
	log.Println(url)
	tableID := `table[id=stats_standard_collapsed]`
	m := make(map[string]float64)

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.AllowedDomains("https://fbref.com", "http://fbref.com", "fbref.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting lol", r.URL.String())
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
	})

	c.Visit(url)
	return m
}
