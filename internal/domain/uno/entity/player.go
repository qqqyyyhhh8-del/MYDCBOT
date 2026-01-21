package entity

type Player struct {
	ID       string
	Username string
	Hand     []*Card
	HasUno   bool
}

func NewPlayer(id, username string) *Player {
	return &Player{
		ID:       id,
		Username: username,
		Hand:     make([]*Card, 0),
	}
}

func (p *Player) AddCard(card *Card) {
	p.Hand = append(p.Hand, card)
	p.HasUno = false
}

func (p *Player) AddCards(cards []*Card) {
	p.Hand = append(p.Hand, cards...)
	p.HasUno = false
}

func (p *Player) RemoveCard(index int) *Card {
	if index < 0 || index >= len(p.Hand) {
		return nil
	}
	card := p.Hand[index]
	p.Hand = append(p.Hand[:index], p.Hand[index+1:]...)
	return card
}

func (p *Player) GetCard(index int) *Card {
	if index < 0 || index >= len(p.Hand) {
		return nil
	}
	return p.Hand[index]
}

func (p *Player) HandSize() int {
	return len(p.Hand)
}

func (p *Player) CallUno() {
	p.HasUno = true
}

func (p *Player) HasWon() bool {
	return len(p.Hand) == 0
}
