package quotes

import "math/rand"

type QuoteRepo interface {
	GetQuote() Quote
}

type Repo struct {
	quotes []Quote
}

func NewRepository() *Repo {
	var quotes = []Quote{
		{
			QuoteText: "You create your own opportunities. Success doesn’t just come and find you–you have to go out and get it.",
		},
		{
			QuoteText: "Laughter meets encouragement with our hilarious list of funny motivational quotes about life. ",
		},
		{
			QuoteText: "To live is the rarest thing in the world. Most people exist, that is all.",
		},
		{
			QuoteText: "The worst part of being okay is that okay is far from happy.",
		},
		{
			QuoteText: "Pain is inevitable. Suffering is optional.",
		},
		{
			QuoteText: "Be kind, for everyone you meet is fighting a hard battle.",
		},
		{
			QuoteText: "The Six Golden Rules of Writing: Read, read, read, and write, write, write.",
		},
		{
			QuoteText: "To produce a mighty book, you must choose a mighty theme.",
		},
	}

	return &Repo{
		quotes: quotes,
	}
}

func (r *Repo) GetQuote() Quote {
	randomIndex := rand.Int() % len(r.quotes)

	return r.quotes[randomIndex]
}
