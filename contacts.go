package main

type Contact struct {
	ID                 string        `json:"id"`
	PhoneLabel         *string       `json:"phoneLabel"`
	Country            string        `json:"country"`
	Address            *string       `json:"address"`
	Source             *string       `json:"source"`
	Type               string        `json:"type"`
	LocationID         string        `json:"locationId"`
	Website            *string       `json:"website"`
	DND                bool          `json:"dnd"`
	State              *string       `json:"state"`
	BusinessName       *string       `json:"businessName"`
	CustomFields       []string      `json:"customFields"`
	Tags               []string      `json:"tags"`
	DateAdded          string        `json:"dateAdded"`
	AdditionalEmails   []string      `json:"additionalEmails"`
	Phone              *string       `json:"phone"`
	CompanyName        *string       `json:"companyName"`
	AdditionalPhones   []string      `json:"additionalPhones"`
	DateUpdated        string        `json:"dateUpdated"`
	City               *string       `json:"city"`
	DateOfBirth        *string       `json:"dateOfBirth"`
	FirstNameLowerCase string        `json:"firstNameLowerCase"`
	LastNameLowerCase  string        `json:"lastNameLowerCase"`
	Email              string        `json:"email"`
	AssignedTo         *string       `json:"assignedTo"`
	Followers          []string      `json:"followers"`
	ValidEmail         *string       `json:"validEmail"`
	PostalCode         *string       `json:"postalCode"`
	BusinessID         *string       `json:"businessId"`
	SearchAfter        []interface{} `json:"searchAfter"`
}

type ContactsResponse struct {
	Contacts []Contact `json:"contacts"`
	Total    int       `json:"total"`
	TraceID  string    `json:"traceId"`
}
