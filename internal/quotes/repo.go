package quotes

import "math/rand"

type QuoteRepository interface {
	GetAnyQuote() Quote
}

type Repo struct {
	quotes []Quote
}

func NewRepository() *Repo {
	var quotes = []Quote{
		{
			Text: "Never decide you are smart enough. Be wise enough to recognize that there is always more to learn",
		},
		{
			Text: "Intend to be as wise as nature, for she never gets pace or cadence wrong.",
		},
		{
			Text: "A loving heart is the truest wisdom.",
		},
		{
			Text: "The worst part of being okay is that okay is far from happy.",
		},
		{
			Text: "Pain is inevitable. Suffering is optional.",
		},
		{
			Text: "Wisdom is trusting the timing of the universe.",
		},
		{
			Text: "Wise is the one who walks against the grain.",
		},
		{
			Text: "To produce a mighty book, you must choose a mighty theme.",
		},
	}

	return &Repo{
		quotes: quotes,
	}
}

func (r *Repo) GetAnyQuote() Quote {
	randomIndex := rand.Int() % len(r.quotes)
	return r.quotes[randomIndex]
}
