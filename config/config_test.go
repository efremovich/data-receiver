package config_test

import (
	"fmt"
	"testing"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/pkg/aconf/v3"
)

func TestParse(t *testing.T) {
	t.Setenv("BROKER_CONSUMER_URL", "test")
	t.Setenv("BROKER_PUBLISHER_URL", "test")

	t.Setenv("MARKETPLACE_0_NAME", "test")
	t.Setenv("MARKETPLACE_0_ID", "test")
	t.Setenv("MARKETPLACE_0_TOKEN", "test")
	t.Setenv("MARKETPLACE_0_TYPE", "test")

	t.Setenv("MARKETPLACE_1_NAME", "test1")
	t.Setenv("MARKETPLACE_1_ID", "test1")
	t.Setenv("MARKETPLACE_1_TOKEN", "test1")
	t.Setenv("MARKETPLACE_1_TYPE", "test1")

	t.Setenv("MARKETPLACE_2_NAME", "test2")
	t.Setenv("MARKETPLACE_2_ID", "test2")
	t.Setenv("MARKETPLACE_2_TOKEN", "test2")
	t.Setenv("MARKETPLACE_2_TYPE", "test2")

	c := new(config.Config)

	if err := aconf.Load(c); err != nil {
		t.Fatal(err)
	}

	c.FillMarketPlaceMap()
	fmt.Println(c)
}
