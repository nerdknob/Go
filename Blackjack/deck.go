package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

var (
	suits = []string{"Clubs", "Spades", "Hearts", "Diamonds"}

	ranks = map[string]int{
		"Two":   2,
		"Three": 3,
		"Four":  4,
		"Five":  5,
		"Six":   6,
		"Seven": 7,
		"Eight": 8,
		"Nine":  9,
		"Ten":   10,
		"Jack":  10,
		"Queen": 10,
		"King":  10,
		"Ace":   11,
	}

	images = map[string]fyne.Resource{
		"AceClubs":   imgAceClubsSvg,
		"TwoClubs":   imgTwoClubsSvg,
		"ThreeClubs": imgThreeClubsSvg,
		"FourClubs":  imgFourClubsSvg,
		"FiveClubs":  imgFiveClubsSvg,
		"SixClubs":   imgSixClubsSvg,
		"SevenClubs": imgSevenClubsSvg,
		"EightClubs": imgEightClubsSvg,
		"NineClubs":  imgNineClubsSvg,
		"TenClubs":   imgTenClubsSvg,
		"JackClubs":  imgJackClubsSvg,
		"QueenClubs": imgQueenClubsSvg,
		"KingClubs":  imgKingClubsSvg,

		"AceSpades":   imgAceSpadesSvg,
		"TwoSpades":   imgTwoSpadesSvg,
		"ThreeSpades": imgThreeSpadesSvg,
		"FourSpades":  imgFourSpadesSvg,
		"FiveSpades":  imgFiveSpadesSvg,
		"SixSpades":   imgSixSpadesSvg,
		"SevenSpades": imgSevenSpadesSvg,
		"EightSpades": imgEightSpadesSvg,
		"NineSpades":  imgNineSpadesSvg,
		"TenSpades":   imgTenSpadesSvg,
		"JackSpades":  imgJackSpadesSvg,
		"QueenSpades": imgQueenSpadesSvg,
		"KingSpades":  imgKingSpadesSvg,

		"AceHearts":   imgAceHeartsSvg,
		"TwoHearts":   imgTwoHeartsSvg,
		"ThreeHearts": imgThreeHeartsSvg,
		"FourHearts":  imgFourHeartsSvg,
		"FiveHearts":  imgFiveHeartsSvg,
		"SixHearts":   imgSixHeartsSvg,
		"SevenHearts": imgSevenHeartsSvg,
		"EightHearts": imgEightHeartsSvg,
		"NineHearts":  imgNineHeartsSvg,
		"TenHearts":   imgTenHeartsSvg,
		"JackHearts":  imgJackHeartsSvg,
		"QueenHearts": imgQueenHeartsSvg,
		"KingHearts":  imgKingHeartsSvg,

		"AceDiamonds":   imgAceDiamondsSvg,
		"TwoDiamonds":   imgTwoDiamondsSvg,
		"ThreeDiamonds": imgThreeDiamondsSvg,
		"FourDiamonds":  imgFourDiamondsSvg,
		"FiveDiamonds":  imgFiveDiamondsSvg,
		"SixDiamonds":   imgSixDiamondsSvg,
		"SevenDiamonds": imgSevenDiamondsSvg,
		"EightDiamonds": imgEightDiamondsSvg,
		"NineDiamonds":  imgNineDiamondsSvg,
		"TenDiamonds":   imgTenDiamondsSvg,
		"JackDiamonds":  imgJackDiamondsSvg,
		"QueenDiamonds": imgQueenDiamondsSvg,
		"KingDiamonds":  imgKingDiamondsSvg,
	}
)

type Deck struct {
	Cards []*Card
}

func buildDeck() *Deck {
	var d Deck

	for _, s := range suits {
		for r, v := range ranks {
			c := &Card{
				Rank:  r,
				Suit:  s,
				Value: v,
				Image: canvas.Image{Resource: getCardImg(r, s)},
			}
			d.Cards = append(d.Cards, c)
		}
	}
	return &d
}

func getCardImg(r string, s string) fyne.Resource {
	card := fmt.Sprintf("%v%v", r, s)
	for k, v := range images {
		if k == card {
			return v
		}
	}
	mgs := fmt.Sprintf("Could not find card image for %v of %v", r, s)
	log.Fatal(mgs)
	return nil
}

func (d *Deck) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) { d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i] })
}

func (d *Deck) DealCard() *Card {
	card := d.Cards[0]
	d.Cards = d.Cards[1:]

	return card
}
