package handler

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

func GetCF() (cookie, ua string, err error) {
	// Launch headless Chrome
		l := launcher.New().
		Headless(false).
		Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) TrackTasteBot/1.0").
		MustLaunch()
	browser := rod.New().ControlURL(l).MustConnect()
	defer browser.MustClose()

	// make stealth page
	page, err := stealth.Page(browser)
	if err != nil {
		return "", "", err
	}

	page.MustEval(`() => {
		Object.defineProperty(navigator, 'webdriver', {get: () => false});
		return true;
	}`)

	time.Sleep(time.Duration(rand.Intn(700)+300) * time.Millisecond)

	page.MustNavigate("https://musicstax.com/track/29vY6gIKRje259YNZ7FyDb").MustWaitLoad()

	// wait turnslite iframe
	iframe := page.MustElementR("iframe", "cloudflare")
	frame := iframe.MustFrame().MustWaitLoad()

	//  click checkbox
	box := frame.MustElement(`input[type="checkbox"]`)
	box.MustClick()

    // wait get cf_clearance cookie
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	for {
		select {
		case <-time.After(500 * time.Millisecond):
			for _, c := range browser.MustGetCookies() {
				if c.Name == "cf_clearance" {
					ver, _ := proto.BrowserGetVersion{}.Call(browser)
					ua = ver.UserAgent
					return c.Value, ua, nil
				}
			}
		case <-ctx.Done():
			return "", "", ctx.Err()
		}
	}
}
