package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"read-adviser-bot/events"
	"read-adviser-bot/lib/e"
	"read-adviser-bot/storage/PostgreSQL/queries"
	"strconv"
	"strings"
)

const (
	EmptyCmd = ""
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
	StartSurveyCmd = "/survey"
)



func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string, diagnosis string, lesionParameters events.LesionParameters) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if text == StartSurveyCmd || p.survey_status > STATUSBEGIN {
		return p.surveyRouter(ctx, chatID, text)
	}

	switch text {
	case EmptyCmd:
		return p.sendDiagnosis(ctx, chatID, diagnosis, lesionParameters)
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}


func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHello)
}

func (p *Processor) sendSurveyBegin(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgSurveyAge)
}

func (p *Processor) sendSurveyAge(ctx context.Context, chatID int, age uint64) error {
	fmt.Printf("age: %v\n", age)

	return p.tg.SendMessage(ctx, chatID, msgSurveyGender)
}

func (p *Processor) saveAge(chatId, age uint64) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save age to db", err) }()
	query := queries.QuerySaveAge
	_, err = p.db.Query(query,
		chatId,
		age,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) saveGender (chatId uint64, isMale bool) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save gender to db", err) }()
	query := queries.QuerySaveGender
	_, err = p.db.Query(query,
		chatId,
		isMale,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) saveNationality (chatId uint64, nationality string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save nationality to db", err) }()
	query := queries.QuerySaveNationality
	_, err = p.db.Query(query,
		chatId,
		nationality,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) surveyRouter (ctx context.Context, chatID int, text string) (err error) {
	defer func() {
		if err == nil{
			p.survey_status += 1
		}
		err = e.WrapIfErr("can't do command: save age to db", err) }()
	switch p.survey_status {
	case STATUSBEGIN:
		err = p.sendSurveyBegin(ctx, chatID)
		if err != nil {
			return err
		}
	case STATUSAGE:
		age, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			p.tg.SendMessage(ctx, chatID, msgSurveyWrongAge)
			return err
		}
		err = p.saveAge(uint64(chatID), age)
		if err != nil {
			return err
		}

		err = p.sendSurveyAge(ctx, chatID, age)
		if err != nil {
			return err
		}
	case STATUSGENDER:
		err = p.saveGender(uint64(chatID), text == "м")
		if err != nil {
			return err
		}
		err = p.sendSurveyGender(ctx, chatID, text)
		if err != nil {
			return err
		}
	case STATUSNATIONALITY:
		err = p.saveNationality(uint64(chatID), text)
		if err != nil {
			return err
		}
		err = p.sendSurveyNationality(ctx, chatID, text)
		if err != nil {
			return err
		}
		p.survey_status = STATUSBEGIN - 1
	} 
	return err
}

func (p *Processor) sendSurveyGender(ctx context.Context, chatID int, text string) error {
	fmt.Printf("text: %v\n", text)
	if text != "м" && text != "ж" {
		p.tg.SendMessage(ctx, chatID, msgSurveyWrongGender)
		return errors.New("Wrong gender")
	}
	fmt.Printf("text: %v\n", text)
	return p.tg.SendMessage(ctx, chatID, msgSurveyNationality)
}

func (p *Processor) sendSurveyNationality (ctx context.Context, chatID int, text string) error {
	fmt.Printf("text: %v\n", text)
	nationalities_list := []string{"Бразилец", "Кавказец", "Китаец", "Индус", "Японец", "Другая"}
	nationalities := make(map[string]struct{})
	for _, el := range nationalities_list {
		nationalities[el] = struct{}{}
	}
	if _, ok := nationalities[text]; !ok {
		p.tg.SendMessage(ctx, chatID, msgSurveyWrongNationality)
		return errors.New("Wrong nationality")
	}
	fmt.Printf("text: %v\n", text)
	return p.tg.SendMessage(ctx, chatID, msgSurveyEnd)
}

func (p *Processor) sendDiagnosis (ctx context.Context, chatID int, diagnosis string, lesionParameters events.LesionParameters) error {
	message := constructDiagnosisMessage(diagnosis, lesionParameters.Area, lesionParameters.Diameter)
	return p.tg.SendMessage(ctx, chatID, message)
}

func constructDiagnosisMessage (diagnosis string, area, diameter float32) string {
	message := fmt.Sprintf(msgDiagnosis, diagnosis)
	if area != 0 {
		lesionParametersMessage := fmt.Sprintf(msgLesionPatameters, area, diameter)
		message = fmt.Sprintf("%s%s", message, lesionParametersMessage)
	}
	return message
}
func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
