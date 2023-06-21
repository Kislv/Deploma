package telegram

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	PBCountPores "read-adviser-bot/clients/GRPC/CountPores/grpc"
	PBDiseaseClassification "read-adviser-bot/clients/GRPC/DiseaseClassification/grpc"
	"read-adviser-bot/clients/telegram"
	"read-adviser-bot/events"
	"read-adviser-bot/lib/e"

	"read-adviser-bot/utils/cast"
	database "read-adviser-bot/storage/PostgreSQL"
	"read-adviser-bot/storage/PostgreSQL/queries"
	"read-adviser-bot/storage/files"
)
const(
	STATUSBEGIN = iota
	STATUSAGE 
	STATUSGENDER 
	STATUSNATIONALITY
	STATUSEND
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	survey_status int
	db *database.DBManager
}


type Meta struct {
	ChatID   int
	Username string
	Photo []telegram.Photo
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

	func New(client *telegram.Client, database *database.DBManager) *Processor {
	return &Processor{
		tg:      client,
		survey_status: 0,
		db: database,
	}
}

func (p *Processor) Fetch(ctx context.Context, limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(ctx, p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, p.event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(ctx, event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}
func nationalityToTableIndex (nationality string) (uint8, error){
	nationalityMap := map[string]uint8 {"Бразилец":0, "Кавказец":1, "Китаец":2, "Индус":3, "Японец":4, "Другая": 1}
	var index uint8
	index, ok := nationalityMap[nationality]
	if  !ok {
		return 0, fmt.Errorf("%s", "cant do nationalityToTableIndex")
	}
	return index, nil
}

func ageToTableIndex (age uint64) (uint8){
	if age < 29 {
		return 0
	}
	if age < 39 {
		return 1
	}
	if age < 49 {
		return 2
	}
	if age < 59 {
		return 3
	}
	if age < 69 {
		return 4
	}
	return 5
}

func (p *Processor) UserParameters(chatId uint64) (age uint64, nationality string, err error) {
	defer func() { err = e.WrapIfErr("can't do command: save age to db", err) }()
	nationality = "Японец"
	age = uint64(20)
	query := queries.QuerySelectParameters
	resp, err := p.db.Query(query,
		chatId,
	)
	if err != nil {
		return 0, "", err
	}

	if len(resp) != 0 {
		age = cast.ToUint64(resp[0][0])
		nationality = cast.ToString(resp[0][1])
	}
	return age, nationality , nil
}


func (p *Processor) poresDensity(chatId uint64) (density uint64, err error){
	defer func() { err = e.WrapIfErr("can't get pores density", err) }()
	age, nationality, err :=  p.UserParameters(chatId)
	if err != nil {
		return 0, err
	}
	densities := [5][6]uint64 {
		{91,82,69,62,62,62},
		{64,62,64,60,61,61},
		{21,23,21,25,20,22},
		{83,80,77,78,74,71},
		{67,70,70,72,70,69},
	}
	nationalityIndex, err := nationalityToTableIndex(nationality)
	if err != nil {
		return 0, err
	}
	return densities[nationalityIndex][ageToTableIndex(age)], nil
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) (err error) {
	defer func() { err = e.WrapIfErr("can't process message", err) }()
	meta, err := meta(event)
	if err != nil {
		return err
	}
	fmt.Printf("len(meta.Photo): %v\n", len(meta.Photo))
	diagnosis := ""
	lesionParameters := events.LesionParameters{Area: 0, Diameter: 0}
	if len(meta.Photo) > 0{

		var wg sync.WaitGroup

		bestQualityPhoto := meta.Photo[len(meta.Photo)-1]
		imageData, extenstion, err := p.DownLoadImage(bestQualityPhoto.FileId)
		if err == files.ErrBadFileExtension{
			if err := p.tg.SendMessage(ctx, meta.ChatID, msgWrongFileExtension); err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
		wg.Add(1)
		var errDiagnose error
		go func() {
			defer wg.Done()
			diagnosis, errDiagnose = p.Diagnosis(bestQualityPhoto, extenstion, imageData)

		}()
		

		if bestQualityPhoto.Height < 256 ||bestQualityPhoto.Width < 256 {
			if err := p.tg.SendMessage(ctx, meta.ChatID, msgWrongResolution); err != nil {
				return err
			}
		}
		
		poresDensity, err := p.poresDensity(uint64(meta.ChatID))
		if err != nil {
			return err
		}
		imageParameters := events.ImageParameters{Extension: extenstion, Height: int32(bestQualityPhoto.Height), Width: int32(bestQualityPhoto.Width), Data: *imageData }
		var errLesionParameters error
		wg.Add(1)
		go func() {
			defer wg.Done()
			lesionParameters, errLesionParameters = p.areaAndDiameter(&imageParameters, poresDensity)
		}()
		wg.Wait()
		if errDiagnose != nil {
			return errDiagnose
		}
		if errLesionParameters != nil {
			return errLesionParameters
		}
	}
	if diagnosis != "Melanoma Skin Cancer Nevi and Moles"{
		lesionParameters = events.LesionParameters{Area: 0, Diameter: 0}
		
	}

	if err := p.doCmd(ctx, event.Text, meta.ChatID, meta.Username, diagnosis, lesionParameters); err != nil {
		return e.Wrap("can't process message", err)
	}


	return nil
}

func (p *Processor) DownLoadImage (fileId string) (*[]byte, string, error) {
	filePath, err := p.tg.File(context.Background(), fileId)

	if err != nil {
		return nil, "", err
	}
	
	data, err := p.tg.DownloadFile(context.Background(), filePath)
	if err != nil {
		return nil, "", err
	}
	fmt.Printf("filePath: %v\n", filePath)
	ext := filepath.Ext(filePath)[1:]
	fmt.Printf("ext: %v\n", ext)
	fmt.Printf("filepath.Ext(filePath): %v\n", filepath.Ext(filePath))
	fmt.Printf("len(data): %v\n", len(data))
	dataPointer := &data
	return dataPointer, ext, nil
}

func (p *Processor) Diagnosis (photo telegram.Photo, extension string, imageData *[]byte) (string, error) {
	conn, err := grpc.Dial("localhost:40041", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}

	client := PBDiseaseClassification.NewClassificateImageClient(conn)

	ImageClassificationResponse, err := client.Classificate(context.Background(), &PBDiseaseClassification.ImageClassificationRequest{Extenstion: extension,Height: int32(photo.Height), Width: int32(photo.Width), Image: *imageData})
	if err != nil {
		return "", err
	}
	fmt.Printf("ImageClassificationResponse.DiseaseName: %v\n", ImageClassificationResponse.DiseaseName)
	return ImageClassificationResponse.DiseaseName, nil
}

func (p *Processor) areaAndDiameter (ImageParameters *events.ImageParameters, poresDensity uint64) (events.LesionParameters, error) {
	conn, err := grpc.Dial("localhost:40042", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return events.LesionParameters{}, err
	}
	client := PBCountPores.NewCountPoresClient(conn)
	ImageClassificationResponse, err := client.Count(context.Background(), &PBCountPores.CountPoresRequest{Extenstion: ImageParameters.Extension, Height: ImageParameters.Height, Width: ImageParameters.Width, Image: ImageParameters.Data, Density: int32(poresDensity)})
	if err != nil {
		return events.LesionParameters{}, err
	}
	fmt.Printf("ImageClassificationResponse.Area: %v\n", ImageClassificationResponse.Area)
	fmt.Printf("ImageClassificationResponse.Diameter: %v\n", ImageClassificationResponse.Diameter)
	return events.LesionParameters{Area: ImageClassificationResponse.Area, Diameter: ImageClassificationResponse.Diameter}, nil

}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func (p *Processor) event (upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
			Photo: upd.Message.Photo,
		}
	}
	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
