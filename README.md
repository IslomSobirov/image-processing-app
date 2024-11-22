# image-processing-app
Склонируйте репозиторий
git clone https://github.com/IslomSobirov/image-processing-app.git
cd image-processing-app


Соберите Docker-образ
docker build -t image-processor .


Запустите контейнер
docker run -d -p 8085:8085 \
  -v $(pwd)/img_orig:/img_orig \
  -v $(pwd)/img_res:/img_res \
  image-processor \
  ./app -path-orig /img_orig -path-res /img_res -width 300 -height 300
