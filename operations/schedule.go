package operations

import (
	"encoding/json"
	"net/url"
	"time"

	omise "github.com/omise/omise-go"
	"github.com/omise/omise-go/internal"
	"github.com/omise/omise-go/schedule"
)

// CreateSchedule represent create schedule API payload
//
// Example:
//
//	schd, create := &omise.Schedule{}, &operations.CreateSchedule{
//              Every:  3,
//              Period: schedule.PeriodWeek,
//              Weekdays: []schedule.Weekday{
//              schedule.Monday,
//              	schedule.Saturday,
//              },
//              StartDate: "2017-05-15",
//              EndDate:   "2018-05-15",
//              Customer:  "customer_id",
//              Amount:    100000,
//	}
//	if e := client.Do(schd, create); e != nil {
//		panic(e)
//	}
//
//	fmt.Println("created schedule:", schd.ID)
//
type CreateSchedule struct {
	Every          int
	Period         schedule.Period
	StartDate      string
	EndDate        string
	Customer       string
	Amount         int
	Currency       string
	Card           string
	Weekdays       []schedule.Weekday
	DaysOfMonth    []int
	WeekdayOfMonth string
}

func (req *CreateSchedule) MarshalJSON() ([]byte, error) {
	type charge struct {
		Customer string `json:"customer"`
		Amount   int    `json:"amount"`
		Currency string `json:"currency,omitempty"`
		Card     string `json:"card,omitempty"`
	}

	type on struct {
		Weekdays       []schedule.Weekday `json:"weekdays,omitempty"`
		DaysOfMonth    []int              `json:"days_of_month,omitempty"`
		WeekdayOfMonth string             `json:"weekday_of_month,omitempty"`
	}

	type param struct {
		Every     int             `json:"every"`
		Period    schedule.Period `json:"period"`
		StartDate *omise.Date     `json:"start_date,omitempty"`
		EndDate   omise.Date      `json:"end_date"`
		Charge    charge          `json:"charge"`
		On        *on             `json:"on,omitempty"`
	}

	p := param{
		Every:  req.Every,
		Period: req.Period,
		Charge: charge{
			Customer: req.Customer,
			Amount:   req.Amount,
			Currency: req.Currency,
			Card:     req.Card,
		},
	}

	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
		p.StartDate = (*omise.Date)(&startDate)
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
		p.EndDate = omise.Date(endDate)
	}

	switch {
	case p.Period == "week":
		p.On = &on{
			Weekdays: req.Weekdays,
		}
	case p.Period == "month" && req.DaysOfMonth != nil:
		p.On = &on{
			DaysOfMonth: req.DaysOfMonth,
		}
	case p.Period == "month" && req.WeekdayOfMonth != "":
		p.On = &on{
			WeekdayOfMonth: req.WeekdayOfMonth,
		}
	}

	return json.Marshal(p)
}

func (req *CreateSchedule) Op() *internal.Op {
	return &internal.Op{
		Endpoint:    internal.API,
		Method:      "POST",
		Path:        "/schedules",
		Values:      url.Values{},
		ContentType: "application/json",
	}
}

// ListSchedules represent list schedule API payload
//
// Example:
//
//	schds, list := &omise.ScheduleList{}, &ListSchedules{
//		List{
//			Limit: 100,
//			From: time.Now().Add(-1 * time.Hour),
//		},
//	}
//	if e := client.Do(schds, list); e != nil {
//		panic(e)
//	}
//
//	fmt.Println("# of schedules made in the last hour:", len(schds.Data))
//
type ListSchedules struct {
	List
}

func (req *ListSchedules) MarshalJSON() ([]byte, error) {
	return json.Marshal(req.List)
}

func (req *ListSchedules) Op() *internal.Op {
	return &internal.Op{
		Endpoint:    internal.API,
		Method:      "GET",
		Path:        "/schedules",
		ContentType: "application/json",
	}
}

// RetrieveSchedule
//
// Example:
//
//	schd := &omise.Schedule{ID: "schd_57z9hj228pusa652nk1"}
//	if e := client.Do(schd, &RetrieveSchedule{schd.ID}); e != nil {
//		panic(e)
//	}
//
//	fmt.Printf("schedule #schd_57z9hj228pusa652nk1: %#v\n", schd)
//
type RetrieveSchedule struct {
	ScheduleID string `query:"-"`
}

func (req *RetrieveSchedule) Op() *internal.Op {
	return &internal.Op{
		Endpoint: internal.API,
		Method:   "GET",
		Path:     "/schedules/" + req.ScheduleID,
	}
}