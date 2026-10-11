package repository

import (
	"fmt"
	"strings"
)

// Адрес бакета Minio, в котором лежат фото и видео участников
const minioURL = "http://localhost:9000/volleyreg/"

// Participant — участник соревнования (услуга)
type Participant struct {
	ID          int
	FullName    string
	Description string
	BirthDate   string
	Position    string
	HeightCm    int
	SportRank   string
	City        string
	ImageURL    string
	VideoURL    string
}

// TeamApplicationParticipant — участник в составе заявки (м-м) с игровым номером
type TeamApplicationParticipant struct {
	ParticipantID int
	JerseyNumber  int
}

// TeamApplication — заявка команды на участие в соревновании
type TeamApplication struct {
	ID              int
	Status          string
	CreatedAt       string
	TeamName        string
	Institution     string
	CoachName       string
	CompetitionDate string
	AverageAge      float64
	Participants    []TeamApplicationParticipant
}

type Repository struct {
	participants     []Participant
	teamApplications map[int]TeamApplication
}

func NewRepository() *Repository {
	return &Repository{
		participants:     participants,
		teamApplications: teamApplications,
	}
}

// GetParticipants возвращает участников, в ФИО которых встречается participantName
func (r *Repository) GetParticipants(participantName string) []Participant {
	if participantName == "" {
		return r.participants
	}
	var result []Participant
	for _, p := range r.participants {
		if strings.Contains(strings.ToLower(p.FullName), strings.ToLower(participantName)) {
			result = append(result, p)
		}
	}
	return result
}

func (r *Repository) GetParticipant(id int) (Participant, error) {
	for _, p := range r.participants {
		if p.ID == id {
			return p, nil
		}
	}
	return Participant{}, fmt.Errorf("участник %d не найден", id)
}

func (r *Repository) GetTeamApplication(id int) (TeamApplication, error) {
	app, ok := r.teamApplications[id]
	if !ok {
		return TeamApplication{}, fmt.Errorf("заявка %d не найдена", id)
	}
	return app, nil
}

// GetDraftTeamApplication возвращает текущий черновик заявки представителя команды
func (r *Repository) GetDraftTeamApplication() (TeamApplication, bool) {
	for _, app := range r.teamApplications {
		if app.Status == "Черновик" {
			return app, true
		}
	}
	return TeamApplication{}, false
}

// Коллекция участников (услуг)
var participants = []Participant{
	{ID: 1, FullName: "Ковалёва Анна Сергеевна", BirthDate: "14.03.2004", Position: "Связующая", HeightCm: 176, SportRank: "КМС", City: "Норильск",
		Description: "Связующая с точной передачей в первый темп. Хорошо читает блок соперника и быстро разводит нападающих по зонам. Капитан студенческой сборной ЗГУ в сезоне 2025/26.",
		ImageURL: minioURL + "participant_1.jpg", VideoURL: minioURL + "video_1.mp4"},
	{ID: 2, FullName: "Миронова Екатерина Олеговна", BirthDate: "02.07.2003", Position: "Доигровщица", HeightCm: 181, SportRank: "1 разряд", City: "Норильск",
		Description: "Доигровщица универсального плана: уверенный приём подачи и сильный удар из четвёртой зоны. Стабильно подаёт планирующую подачу в зону либеро соперника.",
		ImageURL: minioURL + "participant_2.jpg", VideoURL: minioURL + "video_2.mp4"},
	{ID: 3, FullName: "Белова Дарья Андреевна", BirthDate: "23.01.2005", Position: "Центральная блокирующая", HeightCm: 188, SportRank: "КМС", City: "Красноярск",
		Description: "Центральная блокирующая с высоким прыжком. Ставит групповой блок и атакует быстрым ударом с передачи в первый темп. Лучшая блокирующая краевого первенства 2025 года.",
		ImageURL: minioURL + "participant_3.jpg", VideoURL: minioURL + "video_3.mp4"},
	{ID: 4, FullName: "Сафина Алина Ринатовна", BirthDate: "09.11.2004", Position: "Либеро", HeightCm: 165, SportRank: "1 разряд", City: "Норильск",
		Description: "Либеро, отвечает за приём подачи и игру в защите. Отличная реакция на сложных мячах, подстраховка блока и точная доводка мяча на связующую.",
		ImageURL: minioURL + "participant_4.jpg", VideoURL: minioURL + "video_1.mp4"},
	{ID: 5, FullName: "Громова Полина Игоревна", BirthDate: "30.05.2002", Position: "Диагональная", HeightCm: 184, SportRank: "МС", City: "Москва",
		Description: "Диагональная нападающая, основная атакующая сила команды. Мощный удар из второй зоны и с задней линии, силовая подача в прыжке. Мастер спорта, выступала в высшей лиге «Б».",
		ImageURL: minioURL + "participant_5.jpg", VideoURL: minioURL + "video_2.mp4"},
	{ID: 6, FullName: "Тихонова Валерия Денисовна", BirthDate: "17.02.2006", Position: "Доигровщица", HeightCm: 178, SportRank: "2 разряд", City: "Дудинка",
		Description: "Молодая доигровщица, первый сезон в студенческом волейболе. Надёжный приём, техничный удар с обводкой блока, работает над подачей в прыжке.",
		ImageURL: minioURL + "participant_6.jpg", VideoURL: minioURL + "video_3.mp4"},
	{ID: 7, FullName: "Орлова Ксения Максимовна", BirthDate: "05.09.2003", Position: "Центральная блокирующая", HeightCm: 186, SportRank: "1 разряд", City: "Норильск",
		Description: "Центральная блокирующая, сильна в одиночном блоке и в атаке с короткой передачи. Хорошо перемещается вдоль сетки и закрывает направление удара соперника.",
		ImageURL: minioURL + "participant_7.jpg", VideoURL: minioURL + "video_1.mp4"},
	{ID: 8, FullName: "Иванова Софья Павловна", BirthDate: "21.06.2005", Position: "Связующая", HeightCm: 172, SportRank: "2 разряд", City: "Талнах",
		Description: "Связующая, делает ставку на быструю игру и передачи в вторую линию. Хорошая подача и игра в защите. Запасная связующая сборной ЗГУ.",
		ImageURL: minioURL + "participant_8.jpg", VideoURL: minioURL + "video_2.mp4"},
	{ID: 9, FullName: "Шарипова Амина Руслановна", BirthDate: "12.08.2004", Position: "Либеро", HeightCm: 163, SportRank: "1 разряд", City: "Кайеркан",
		Description: "Либеро с отличным приёмом силовой подачи. Уверенно играет в защите на задней линии, организует подстраховку и командует расстановкой при приёме.",
		ImageURL: minioURL + "participant_9.jpg", VideoURL: minioURL + "video_3.mp4"},
	{ID: 10, FullName: "Федорова Мария Алексеевна", BirthDate: "03.12.2001", Position: "Диагональная", HeightCm: 183, SportRank: "КМС", City: "Норильск",
		Description: "Опытная диагональная нападающая. Атакует из второй зоны и пайпом с задней линии, сильная подача в прыжке. Играет за команду ветеранов ЗГУ.",
		ImageURL: minioURL + "participant_10.jpg", VideoURL: minioURL + "video_1.mp4"},
}

// Словарь заявок команд: ключ — id заявки
var teamApplications = map[int]TeamApplication{
	1: {
		ID:              1,
		Status:          "Черновик",
		CreatedAt:       "10.10.2026 18:40",
		TeamName:        "Северное сияние",
		Institution:     "ЗГУ им. Н.М. Федоровского, факультет ЭЭиУ",
		CoachName:       "Петров Игорь Викторович",
		CompetitionDate: "20.11.2026",
		AverageAge:      22.7,
		Participants: []TeamApplicationParticipant{
			{ParticipantID: 1, JerseyNumber: 7},
			{ParticipantID: 4, JerseyNumber: 1},
			{ParticipantID: 5, JerseyNumber: 12},
		},
	},
}
