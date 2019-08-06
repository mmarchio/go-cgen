package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Rule struct {
	Trait string
	Value int
}

type Trait struct {
	Name string `json: "name"`
	Value int `json: "value"`
}

type Info struct {
	Name string `json: "name"`
	Player string `json: "player"`
	Chronicle string `json: "chronicle"`
	Nature string `json: "nature"`
	Demeanor string `json: "demeanor"`
	Concept string `json: "concept"`
	Clan string `json: "clan"`
	Generation int `json: "generation"`
	Sire string `json: "sire"`
}

type Attributes struct {
	Physical []Trait `json: "physical"`
	Social []Trait `json: "social"`
	Mental []Trait `json: "mental"`
}

type Abilities struct {
	Talents []Trait `json: "talents"`
	Skills []Trait `json: "skills"`
	Knowledges []Trait `json: "knowledges"`
}

type Advantages struct {
	Disciplines []Trait `json: "disciplines"`
	Backgrounds []Trait `json: "backgrounds"`
	Virtues []Trait `json: "virtues"`
}

type Character struct {
	Info Info `json: "info"`
	Attributes Attributes `json: "attributes"`
	Abilities Abilities `json: "abilities"`
	Advantages Advantages `json: "advantages"`
	Vitals []Trait `json: "vitals"`
}

type Params struct {
	Name string
	Player string
	Chronicle string
	Nature string
	Demeanor string
	Sire string
}

func main() {
	args := os.Args
	port := args[1:]

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","application/json")
		fmt.Fprintf(w, generateCharacter(r))
	})

	http.ListenAndServe(concat(":", string(port[0])), nil)
}

func concat(str1, str2 string) string {
	b := bytes.Buffer{}
	b.WriteString(str1)
	b.WriteString(str2)
	return b.String()
}

func getParams(path string) Params {
	a := strings.Split(path, "/")
	params := Params{}
	params.Name = "Name"
	params.Player = "Player"
	params.Chronicle = "Chronicle"
	params.Nature = "Nature"
	params.Demeanor = "Demeanor"
	params.Sire = "Sire"

	switch len(a) {
	case 0:
		break
	case 1:
		break
	case 2:
		params.Name = a[1]
		break
	case 3:
		params.Name = a[1]
		params.Player = a[2]
		break
	case 4:
		params.Name = a[1]
		params.Player = a[2]
		params.Chronicle = a[3]
		break
	case 5:
		params.Name = a[1]
		params.Player = a[2]
		params.Chronicle = a[3]
		params.Nature = a[4]
		break
	case 6:
		params.Name = a[1]
		params.Player = a[2]
		params.Chronicle = a[3]
		params.Nature = a[4]
		params.Demeanor = a[5]
		break
	case 7:
		params.Name = a[1]
		params.Player = a[2]
		params.Chronicle = a[3]
		params.Nature = a[4]
		params.Demeanor = a[5]
		params.Sire = a[6]
		break

	}
	return params
 }


func generateCharacter(r *http.Request) string {
	params := getParams(r.URL.Path)

	info := Info{params.Name, params.Player, params.Chronicle, params.Nature, params.Demeanor, "NPC", selectClan(), 13, params.Sire}

	rule := getRules(info)

	attributePriority := priority(7, 5, 3)
	abilityPriority := priority(13, 9, 5)

	physical := make([]Trait,3)
	social := make([]Trait, 3)
	mental := make([]Trait, 3)

	physical[0] = NewTrait("strength", 1)
	physical[1] = NewTrait("dexterity",1)
	physical[2] = NewTrait("stamina",1)

	social[0] = NewTrait("charisma",1)
	social[1] = NewTrait("manipulation",1)
	social[2] = NewTrait("appearance",1)

	mental[0] = NewTrait("perception",1)
	mental[1] = NewTrait("intelligence",1)
	mental[2] = NewTrait("wits",1)

	talents := make([]Trait,10)
	skills := make([]Trait, 10)
	knowledges := make([]Trait, 10)

	talents[0] = NewTrait("acting", 0)
	talents[1] = NewTrait("alertness", 0)
	talents[2] = NewTrait("athletics", 0)
	talents[3] = NewTrait("brawl", 0)
	talents[4] = NewTrait("dodge", 0)
	talents[5] = NewTrait("empathy", 0)
	talents[6] = NewTrait("intimidation", 0)
	talents[7] = NewTrait("leadership", 0)
	talents[8] = NewTrait("streetwise", 0)
	talents[9] = NewTrait("subterfuge", 0)

	skills[0] = NewTrait("animal ken", 0)
	skills[1] = NewTrait("drive", 0)
	skills[2] = NewTrait("etiquette", 0)
	skills[3] = NewTrait("firearms", 0)
	skills[4] = NewTrait("melee", 0)
	skills[5] = NewTrait("music", 0)
	skills[6] = NewTrait("repair", 0)
	skills[7] = NewTrait("security", 0)
	skills[8] = NewTrait("stealth", 0)
	skills[9] = NewTrait("survival", 0)

	knowledges[0] = NewTrait("bureaucracy", 0)
	knowledges[1] = NewTrait("computer", 0)
	knowledges[2] = NewTrait("finance", 0)
	knowledges[3] = NewTrait("investigation", 0)
	knowledges[4] = NewTrait("law", 0)
	knowledges[5] = NewTrait("linguistics", 0)
	knowledges[6] = NewTrait("medicine", 0)
	knowledges[7] = NewTrait("occult", 0)
	knowledges[8] = NewTrait("politics", 0)
	knowledges[9] = NewTrait("science", 0)


	physicalPoints := distribute(physical, attributePriority[0], 5, rule)
	socialPoints := distribute(social, attributePriority[1], 5, rule)
	mentalPoints := distribute(mental, attributePriority[2], 5, rule)

	talentsPoints := distribute(talents, abilityPriority[0], 5, rule)
	skillsPoints := distribute(skills, abilityPriority[1], 5, rule)
	knowledgesPoints := distribute(knowledges, abilityPriority[2], 5, rule)

	disciplinesPoints := distribute(getDisciplines(info), getDisciplinePoints(info), 5, rule)
	backgroundsPoints := getBackgrounds(10)
	virtuesPoints := getVirtues(7)

	attributes := Attributes{physicalPoints, socialPoints, mentalPoints}
	if info.Clan == "nosferatu" {
		attributes.Social[2].Value = 0
	}
	abilities := Abilities{talentsPoints, skillsPoints, knowledgesPoints}
	advantages := Advantages{disciplinesPoints, backgroundsPoints, virtuesPoints}

	if backgroundsPoints[5].Value > 0 {
		info.Generation = info.Generation - backgroundsPoints[5].Value
	}

	character := Character{info, attributes, abilities, advantages, getVitals(virtuesPoints, info)}
	character = freebies(15, character)

	jsonString, err := json.Marshal(character)
	if err != nil {
		fmt.Println(err)
	}

	return string(jsonString)
}

func NewTrait(name string, val int) Trait {
	return Trait{name, val}
}

func priority (primary, secondary, tertiary int) []int {
	a := []int{primary, secondary, tertiary}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	return a
}

func distribute (category []Trait, points, max int, rules Rule) []Trait {
	for points > 0 {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(len(category))

		if category[r].Name != "" && category[r].Value < max {
			if rules.Trait == "none" || (category[r].Name == rules.Trait && category[r].Value < rules.Value) {
				category[r].Value = category[r].Value + 1
				points = points - 1
			}
		}
	}

	return category
}

func selectClan() string {
	clans := make([]string, 8)
	clans[0] = "brujah"
	clans[1] = "gangrel"
	clans[2] = "malkavian"
	clans[3] = "nosferatu"
	clans[4] = "toreador"
	clans[5] = "tremere"
	clans[6] = "ventrue"
	clans[7] = "caitiff"

	rand.Seed(time.Now().UnixNano())
	return clans[rand.Intn(7)]
 }

func getDisciplinePoints(info Info) int {
	p := 3
	if info.Clan == "caitiff" {
		p = 1
	}
	return p
}

func getDisciplines(info Info) []Trait {
	d := make([]Trait, 10)

	switch info.Clan {
		case "brujah":
			d[0] = NewTrait("celerity", 0)
			d[1] = NewTrait("potence", 0)
			d[2] = NewTrait("presence", 0)
			break
		case "gangrel":
			d[0] = NewTrait("animalism", 0)
			d[1] = NewTrait("fortitude", 0)
			d[2] = NewTrait("protean", 0)
			break
		case "malkavian":
			d[0] = NewTrait("auspex", 0)
			d[1] = NewTrait("dominate", 0)
			d[2] = NewTrait("obfuscate", 0)
			break
		case "nosferatu":
			d[0] = NewTrait("animalism", 0)
			d[1] = NewTrait("obfuscate", 0)
			d[2] = NewTrait("potence", 0)
			break
		case "toreador":
			d[0] = NewTrait("auspex", 0)
			d[1] = NewTrait("celerity", 0)
			d[2] = NewTrait("presence", 0)
			break
		case "tremere":
			d[0] = NewTrait("auspex", 0)
			d[1] = NewTrait("dominate", 0)
			d[2] = NewTrait("thaumaturgy", 0)
			break
		case "ventrue":
			d[0] = NewTrait("dominate", 0)
			d[1] = NewTrait("fortitude", 0)
			d[2] = NewTrait("presence", 0)
			break
		default:
			d[0] = NewTrait("animalism", 0)
			d[1] = NewTrait("auspex", 0)
			d[2] = NewTrait("celerity", 0)
			d[3] = NewTrait("dominate", 0)
			d[4] = NewTrait("fortitude", 0)
			d[5] = NewTrait("obfuscate", 0)
			d[6] = NewTrait("potence", 0)
			d[7] = NewTrait("presence", 0)
			d[8] = NewTrait("protean", 0)
			d[9] = NewTrait("thaumaturgy", 0)
			break
	}

	return d
}

func getBackgrounds(points int) []Trait {
	b := make([]Trait, 10)
	b[0] = NewTrait("allies", 0)
	b[1] = NewTrait("contacts", 0)
	b[2] = NewTrait("fame", 0)
	b[3] = NewTrait("herd", 0)
	b[4] = NewTrait("influence", 0)
	b[5] = NewTrait("generation", 0)
	b[6] = NewTrait("mentor", 0)
	b[7] = NewTrait("resouces", 0)
	b[8] = NewTrait("retainer", 0)
	b[9] = NewTrait("status", 0)

	rules := Rule{"none", 10}

	return distribute(b, points, 10, rules)
}

func getVirtues(points int) []Trait {
	v := make([]Trait, 3)
	v[0] = NewTrait("conscience", 1)
	v[1] = NewTrait("self-control", 1)
	v[2] = NewTrait("courage", 1)

	rules := Rule{"none", 10}
	return distribute(v, points, 5, rules)
}

func getVitals(virtues []Trait, info Info) []Trait {
	v := make([]Trait, 3)
	v[0] = NewTrait("path", virtues[0].Value + virtues[1].Value)
	v[1] = NewTrait("willpower", virtues[2].Value)

	bp := make([]int, 14)
	bp[13] = 10
	bp[12] = 11
	bp[11] = 12
	bp[10] = 15
	bp[9] = 20
	bp[8] = 25
	bp[7] = 30
	bp[6] = 35
	bp[5] = 40
	bp[4] = 45
	bp[3] = 50
	bp[2] = 75
	bp[1] = 100
	bp[0] = 1000

	v[2] = NewTrait("bloodpool", bp[info.Generation])

	return v
}

func freebies(points int, character Character) Character {
	cost := make([]int, 7)
	cost[0] = 5 //attribute
	cost[1] = 2 //ability
	cost[2] = 7 //discipline
	cost[3] = 1 //background
	cost[4] = 2 //virtue
	cost[5] = 2 //path
	cost[6] = 1 //willpower
	incremented := false

	for points > 0 {
		incremented = false
		traitmax := getTraitMax(character)
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(7)
		if cost[r] <= points {
			switch r {
				case 0:
					character.Attributes, incremented = getRandomAttribute(character.Attributes, getRules(character.Info), traitmax)
					break
				case 1:
					character.Abilities, incremented = getRandomAbility(character.Abilities, getRules(character.Info), traitmax)
					break
				case 2:
					character.Advantages.Disciplines, incremented = getRandomAdvantage(character.Advantages.Disciplines, getRules(character.Info), traitmax)
					break
				case 3:
					character.Advantages.Backgrounds, incremented = getRandomAdvantage(character.Advantages.Backgrounds, getRules(character.Info), 10)
					break
				case 4:
					character.Advantages.Virtues, incremented = getRandomAdvantage(character.Advantages.Virtues, getRules(character.Info), 5)
					break
				case 5:
					if character.Vitals[0].Value < 10 {
						character.Vitals[0].Value = character.Vitals[0].Value + 1
						incremented = true
					}
					break
				case 6:
					if character.Vitals[1].Value < 10 {
						character.Vitals[1].Value = character.Vitals[1].Value + 1
						incremented = true
					}
 					break
			}
			if incremented == true {
				points = points - cost[r]
			}
		}
	}

	return character
}

func getRandomAttribute(attribute Attributes, rules Rule, traitmax int) (Attributes, bool) {
	r := Rand(3)
	ra := Rand(3)

	incremented := false

	switch r {
	case 0:
		att := attribute.Physical[ra]
		if att.Value < traitmax {
			if att.Name != rules.Trait || (att.Name == rules.Trait && att.Value < rules.Value) {
				attribute.Physical[ra].Value = attribute.Physical[ra].Value + 1
				incremented = true
			}
		}
		break
	case 1:
		att := attribute.Social[ra]
		if att.Value < traitmax {
			if att.Name != rules.Trait || (att.Name == rules.Trait && att.Value < rules.Value) {
				attribute.Social[ra].Value = attribute.Social[ra].Value + 1
				incremented = true
			}
		}
		break
	case 2:
		att := attribute.Mental[ra]
		if att.Value < traitmax {
			if att.Name != rules.Trait || (att.Name == rules.Trait && att.Value < rules.Value) {
				attribute.Mental[ra].Value = attribute.Mental[ra].Value + 1
				incremented = true
			}
		}
		break
	}

	return attribute, incremented
}

func getRandomAbility(ability Abilities, rules Rule, traitmax int) (Abilities, bool) {
	r := Rand(3)
	ra := Rand(10)
	incremented := false

	switch r {
		case 0:
			abi := ability.Talents[ra]
			if abi.Value < traitmax && abi.Name != "" {
				if abi.Name != rules.Trait || (abi.Name == rules.Trait && abi.Value < rules.Value) {
					ability.Talents[ra].Value = ability.Talents[ra].Value + 1
					incremented = true
				}

			}
			break
		case 1:
			abi := ability.Skills[ra]
			if abi.Value < traitmax && abi.Name != ""{
				if abi.Name != rules.Trait || (abi.Name == rules.Trait && abi.Value < rules.Value) {
					ability.Skills[ra].Value = ability.Skills[ra].Value + 1
					incremented = true
				}
			}
			break
		case 2:
			abi := ability.Knowledges[ra]
			if abi.Value < traitmax && abi.Name != "" {
				if abi.Name != rules.Trait || (abi.Name == rules.Trait && abi.Value < rules.Value) {
					ability.Knowledges[ra].Value = ability.Knowledges[ra].Value + 1
					incremented = true
				}
			}
			break
	}
	return ability, incremented
}

func getRandomAdvantage(advantage []Trait, rules Rule, traitmax int) ([]Trait, bool) {
	r := Rand(len(advantage))
	a := advantage[r]
	incremented := false

	if a.Value < traitmax && a.Name != ""{
		if a.Name != rules.Trait || (a.Name == rules.Trait && a.Value < rules.Value) {
			advantage[r].Value = advantage[r].Value + 1
			incremented = true
		}
	}
	return advantage, incremented
}

func Rand(rng int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(rng)
}

func getRules(info Info) Rule {
	rule := Rule{"none", 10}
	if info.Clan == "nosferatu" {
		rule = Rule{"appearance", 0}
	}
	return rule
}

func getTraitMax(character Character) int {
	tm := character.Advantages.Backgrounds[5].Value

	if tm < 5 {
		tm = 5
	}
	return tm
}