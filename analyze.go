package rbsa

import (
	. "github.com/badgerodon/lalg"
	"bufio"
	"fmt"
	"http"
	"os"
	"github.com/badgerodon/statistics"
	"strings"
	"strconv"
	"time"
)

var (
	cache = NewCache(100)
)

func getData(symbol string) (Vector, os.Error) {
	t := time.LocalTime()
	
	y := t.Year
	m := t.Month
	
	if m == 1 {
		m = 12
		y--
	} else {
		m--
	}

	client := new(http.Client)
	vec, err := cache.Get(fmt.Sprint(y, ":", m, ":", symbol), func() (interface{}, os.Error) {
		r, _, err := client.Get("http://ichart.finance.yahoo.com/table.csv?s=" + http.URLEscape(symbol) +
			fmt.Sprint("&a=", (m - 1), "&b=5&c=", (y - 4), 
				"&d=", (m - 1), "&b=5&c=", y,
				"&ignore=.csv"))
		if err != nil {
			return nil, err
		}	
		defer r.Body.Close()
		
		csv := bufio.NewReader(r.Body)
		
		vec := NewVector(37)
		
		for i := 0; i <= len(vec); i++ {
			line, err := csv.ReadString('\n')
			if err != nil {
				break
			}
			
			// Skip the headers
			if i == 0 {
				continue
			}
			
			// Read the data
			parts := strings.Split(line, ",", 7)
			if len(parts) < 6 {
				continue
			}
			
			v, err := strconv.Atof64(strings.Trim(parts[6], "\r\n"))
			
			if err != nil {
				v = 0
			}
			
			vec[i-1] = v
		}
		vec = statistics.Relativize(vec)
		return vec, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return vec.(Vector), nil
	
}

//http://ichart.finance.yahoo.com/table.csv?s=%5EGSPC&a=00&b=3&c=1950&d=05&e=2&f=2011&g=m&ignore=.csv
func Analyze(id string) (map[string]float64, os.Error) {
	indices := map[string]string{
		"IWB": "Large Cap",
		"IWD": "Large Cap Value",
		"IWF": "Large Cap Growth",
		"IWM": "Small Cap",
		"IWN": "Small Cap Value",
		"IWO": "Small Cap Growth",
		"IWR": "Mid Cap",
		"EEM": "Emerging Markets",
		"ICF": "Real Estate",
		"EFA": "International",
		"AGG": "Fixed Income",
	}

	alg := New()
	for k, _ := range indices {
		data, err := getData(k)
		if err != nil {
			return nil, err
		}
		alg.AddIndex(k, data)
	}
	
	data, err := getData(id)
	
	if err != nil {
		return nil, err
	}
	
	solution, err := alg.Run(data)
	
	if err != nil {
		return nil, err
	}
	
	newSolution := make(map[string]float64)
	for k, v := range solution {
		newSolution[indices[k]] = v
	}
	
	return newSolution, nil
}