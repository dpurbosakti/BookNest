package rent

import "time"

type Rent struct {
	Id        uint      `json:"id"`
	StartRent time.Time `json:"start_rent"`
	EndRent   time.Time `json:"end_rent"`
}
