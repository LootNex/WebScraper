package driver

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	// "github.com/tebeka/selenium/firefox"
)

func InitDriver() (selenium.WebDriver, func(), error) {

	driverPath, err := exec.LookPath("chromedriver")
	if err != nil {
		return nil, nil, err
	}

	service, err := selenium.NewChromeDriverService(driverPath, 9175)
	if err != nil {
		return nil, nil, err
	}

	caps := selenium.Capabilities{"browserName": "chrome"}
	chromecaps := chrome.Capabilities{Args: []string{
		"--headless=new",
		"--disable-extensions",
		"--disable-gpu",
		"--no-sandbox",
		"--enable-unsafe-swiftshader",
		"--disable-blink-features=AutomationControlled",
		"user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"accept-language=en-US,en;q=0.9",
		"--disable-web-security",
		"--disable-dev-shm-usage",
		"--hide-scrollbars",
		"--disable-infobars",
		"--disable-popup-blocking",
		"--disable-notifications",
		fmt.Sprintf("--user-data-dir=%s", fmt.Sprintf("/tmp/chrome-user-dir-%d", os.Getpid())),

	}}

	caps.AddChrome(chromecaps)

	fmt.Println("Почти")

	driver, err := selenium.NewRemote(caps, "http://localhost:9175/wd/hub")
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("запустился")

	finish := func() {
		driver.Quit()
		service.Stop()
	}

	return driver, finish, nil
}