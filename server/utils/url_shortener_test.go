package utils

import (
	"testing"
)

const key = "e0dba740-fc4b-4977-872c-d360239e6b1a"

func BenchmarkShortLinkGenerator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateShortLink("https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html", key)
	}
}
