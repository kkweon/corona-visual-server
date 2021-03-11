package model

import (
	"corona-visual-server/internal/config"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMarshalResponse(t *testing.T) {
	input := `<response>
    <header>
        <resultCode>00</resultCode>
        <resultMsg>NORMAL SERVICE.</resultMsg>
    </header>
    <body>
        <items>
            <item>
                <accDefRate>3.9193080566</accDefRate>
                <accExamCnt>210144</accExamCnt>
                <accExamCompCnt>191692</accExamCompCnt>
                <careCnt>7165</careCnt>
                <clearCnt>247</clearCnt>
                <createDt>2020-03-10 10:20:27.00</createDt>
                <deathCnt>54</deathCnt>
                <decideCnt>7513</decideCnt>
                <examCnt>18452</examCnt>
                <resutlNegCnt>184179</resutlNegCnt>
                <seq>69</seq>
                <stateDt>20200310</stateDt>
                <stateTime>00:00</stateTime>
                <updateDt>2020-03-10 10:20:27.27</updateDt>
            </item>
		</items>
	</body>
</response>`

	var response Response
	err := xml.Unmarshal([]byte(input), &response)
	assert.NoError(t, err, "it should not return an error")

	assert.InDeltaf(t, response.Body.Items.Item[0].AccDefRate, 3.9193080566, 1e-7, "accDefRate did not match")

	expected := time.Date(2020, 3, 10, 10, 20, 27, 0, config.SeoulTZ)
	assert.Equal(t, expected, response.Body.Items.Item[0].CreateDt.Time)
}
