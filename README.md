# ВКР Бакалавра студента Киселева В.В.
Здесь хранится информация о выборе архитектур, обучении моделей для классификации кожных заболеваний, сегментации области поражения кожного заболевания, алгоритме подсчета пор на фотографии, алгоритме определения абсолютной площади поражения в (см^2), алгоритме определения диаметра области поражения.

Запуск:

* Классификация изображений:

Генерация файлов для работы с proto для бэк-енда:
(Back-end folder)
protoc -I ../GRPC/DiseaseClassification/ --go_out=../Back-end/clients/GRPC/DiseaseClassification --go-grpc_out=../Back-end/clients/GRPC/DiseaseClassification ../GRPC/DiseaseClassification/image.proto

Генерация файлов для работы с proto для микросервиса:
(DiseaseClassification folder)
python -m grpc_tools.protoc -I ../../GRPC/DiseaseClassification/ --python_out=../../GRPC/DiseaseClassification/ML --grpc_python_out=../../GRPC/DiseaseClassification/ML ../../GRPC/DiseaseClassification/image.proto

* Подсчет пор:

Генерация файлов для работы с proto для бэк-енда:
(Back-end folder)
protoc -I ../GRPC/CountPores/ --go_out=../Back-end/clients/GRPC/CountPores --go-grpc_out=../Back-end/clients/GRPC/CountPores ../GRPC/CountPores/countPores.proto

Генерация файлов для работы с proto для микросервиса:
(PoresSegmentation folder)
python -m grpc_tools.protoc -I ../../GRPC/CountPores/ --python_out=../../GRPC/CountPores/ML --grpc_python_out=../../GRPC/CountPores/ML ../../GRPC/CountPores/countPores.proto

* Сегментация области поражения заболевания

Генерация файлов для работы с proto для сервера:
(LesionSegmentation folder)
python -m grpc_tools.protoc -I ../../GRPC/LesionSegmentation/ --python_out=../../GRPC/LesionSegmentation/ML --grpc_python_out=../../GRPC/LesionSegmentation/ML ../../GRPC/LesionSegmentation/lesionSegmentation.proto
