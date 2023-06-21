package telegram

const (
	msgHelp = `Я помогаю определять кожное за болевание по изображению!
	
	Чтобы получить результат распознавания просто отправьте мне фотографию.
	Если вы хотите получить больше информации о распознавании, то укажите информацию о себе после ввода команды: "`+ StartSurveyCmd + `"`
	msgHello = "Привет! 👾\n\n" + msgHelp
	msgSurveyAge = "Сейчас вам потребуется ввести информацию о себе, чтобы улучшить точность работы бота. \nВведите свой возраст числом 🕰"
	msgSurveyWrongAge = "Возраст должен являться числом"
	msgSurveyGender = "Выберите свой пол 👨 👩‍🦰\n'м' или 'ж' для мужчины и женщины, соответственно"
	msgSurveyWrongGender = "Для выбора пола следует ввести 'м' или 'ж' для мужчины и женщины, соответственно"
	msgSurveyNationality = "Введите название своей национальности 👨🏻 👳‍♀️ из данного списка:\n'Бразилец', 'Кавказец', 'Китаец', 'Индус', 'Японец', 'Другая'"
	msgSurveyWrongNationality = "Введенная национальность должна быть из данного списка:\n'Бразилец', 'Кавказец', 'Китаец', 'Индус', 'Японец', 'Другая'"
	msgSurveyEnd = "Информация сохранена, отлично!"
	msgWrongFileExtension = "Поддреживаемые разрешения изображений: 'jpg', 'png'"
	msgWrongResolution = "Разрешение изображения должно быть хотя бы 256x256"
)

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgNoSavedPages   = "You have no saved pages 🙊"
	msgSaved          = "Saved! 👌"
	msgAlreadyExists  = "You have already have this page in your list 🤗"
	msgDiagnosis   	  = "Предполагаемое заболевание: %s 🦠"
	msgLesionPatameters   = "\nПлощадь поражения: %.2f см^2;\nДиаметр области поражения: %.2f см."
)
