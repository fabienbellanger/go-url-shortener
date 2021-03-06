package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const key = "e0dba740-fc4b-4977-872c-d360239e6b1a"

func TestShortLinkGenerator(t *testing.T) {
	initialLink_1 := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortLink_1, _ := GenerateShortLink(initialLink_1, key)

	initialLink_2 := "https://www.eddywm.com/lets-build-a-url-shortener-in-go-with-redis-part-2-storage-layer/"
	shortLink_2, _ := GenerateShortLink(initialLink_2, key)

	initialLink_3 := "https://spectrum.ieee.org/automaton/robotics/home-robots/hello-robots-stretch-mobile-manipulator"
	shortLink_3, _ := GenerateShortLink(initialLink_3, key)

	assert.Equal(t, shortLink_1, "jTa4L57P")
	assert.Equal(t, shortLink_2, "d66yfx7N")
	assert.Equal(t, shortLink_3, "dhZTayYQ")
}

func BenchmarkShortLinkGenerator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateShortLink("https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html", key)
	}
}
